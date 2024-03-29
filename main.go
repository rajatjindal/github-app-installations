package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-github/v45/github"
	"github.com/pkg/errors"
)

const (
	//KeyHeader that contains private key for app
	KeyHeader = "X-Private-Key"

	//AppIDHeader contains the app id
	AppIDHeader = "X-App-Id"

	//FormatHeader specifies the return format
	FormatHeader = "X-Resp-Format"
)

type Repo struct {
	Name    string `json:"name"`
	HtmlURL string `json:"htmlURL"`
}

type resp struct {
	GithubLogin         string  `json:"ghLogin"`
	OrgUserURL          string  `json:"orgUserURL"`
	RepositorySelection string  `json:"repositorySelection"`
	Repositories        []*Repo `json:"repositories,omitempty"`
}

func main() {
	http.HandleFunc("/", Handle)
	http.ListenAndServe(":8080", nil)
}

func Handle(w http.ResponseWriter, req *http.Request) {
	key := req.Header.Get(KeyHeader)
	appID := req.Header.Get(AppIDHeader)
	format := req.Header.Get(FormatHeader)

	if key == "" || appID == "" {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("%s and %s headers are mandatory", AppIDHeader, KeyHeader)
		w.Write([]byte(msg))
		return
	}

	iss, err := strconv.Atoi(appID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("invalid value of %s header [%s]. it should be an integer", AppIDHeader, appID)))
		return
	}

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("invalid value of %s. failed to decode", KeyHeader)))
		return
	}

	transport := &authenticator{
		transport: http.DefaultTransport,
		iss:       iss,
		key:       decodedKey,
	}

	client := &http.Client{Transport: transport}
	c := github.NewClient(client)
	l := &github.ListOptions{
		PerPage: 100,
	}
	i, _, err := c.Apps.ListInstallations(context.Background(), l)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	channel := make(chan *resp, len(i))
	for _, installation := range i {
		go func(ii *github.Installation, respChan chan *resp) {
			repositorySelection := stringValue(ii.RepositorySelection)
			e := &resp{
				GithubLogin:         stringValue(ii.Account.Login),
				OrgUserURL:          stringValue(ii.Account.HTMLURL),
				RepositorySelection: repositorySelection,
			}

			if repositorySelection == "all" {
				respChan <- e
				return
			}

			installationToken, _, err := c.Apps.CreateInstallationToken(context.Background(), ii.GetID(), &github.InstallationTokenOptions{})
			if err != nil {
				respChan <- e
				return
			}

			repos, err := getRepositoriesForInstallation(installationToken.GetToken())
			if err != nil {
				respChan <- e
				return
			}

			e.Repositories = []*Repo{}
			for _, r := range repos.Repositories {
				e.Repositories = append(e.Repositories, &Repo{
					Name:    r.GetName(),
					HtmlURL: stringValue(r.HTMLURL),
				})
			}

			respChan <- e
		}(installation, channel)
	}

	out := []*resp{}
	for loop := 0; loop < len(i); loop++ {
		installation := <-channel
		out = append(out, installation)
	}

	if format == "readme" {
		outString := []string{"| Org/User | Repository |"}
		outString = append(outString, "| ------ | ------ |")

		for _, i := range out {
			if i.RepositorySelection == "all" {
				outString = append(outString, fmt.Sprintf("| [%s](%s) | [All](%s) |", i.GithubLogin, i.OrgUserURL, i.OrgUserURL))
				continue
			}

			listVar := ""
			if len(i.Repositories) > 1 {
				listVar = "- "
			}

			repos := []string{}
			max := 5
			for index, r := range i.Repositories {
				repos = append(repos, fmt.Sprintf("%s[%s](%s)", listVar, r.Name, r.HtmlURL))
				done := index + 1
				remaining := len(i.Repositories) - done
				if done >= max && remaining > 0 {
					repos = append(repos, fmt.Sprintf("%sand %d more...", listVar, remaining))
					break
				}
			}
			outString = append(outString, fmt.Sprintf("| [%s](%s) | %s |", i.GithubLogin, i.OrgUserURL, strings.Join(repos, "<br/>")))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strings.Join(outString, "\n")))
		return
	}

	if format == "readme-logos" {
		outString := []string{"Note: The images below are the profile images of orgs/users who have enabled the `goodfirstissue` bot on one or more repository.\n\n"}
		for _, i := range out {
			outString = append(outString, fmt.Sprintf(`<a href="%s"><img src="%s.png" width="100"></a><span width="10px">&nbsp;</span>`, i.OrgUserURL, i.OrgUserURL))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strings.Join(outString, "\n")))
		return
	}

	outBytes, err := json.Marshal(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to marshal"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(outBytes)
}

type installationAuth struct {
	transport http.RoundTripper
	token     string
}

func (a *installationAuth) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("bearer %s", a.token))
	resp, err := a.transport.RoundTrip(r)
	return resp, err
}

func getRepositoriesForInstallation(token string) (*github.ListRepositories, error) {
	transport := &installationAuth{
		transport: http.DefaultTransport,
		token:     token,
	}

	client := &http.Client{
		Transport: transport,
	}

	ghClient := github.NewClient(client)

	repos, _, err := ghClient.Apps.ListRepos(context.Background(), &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

type authenticator struct {
	transport http.RoundTripper
	iss       int
	key       []byte
}

func (a *authenticator) GetToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": int32(time.Now().Unix()),
		"exp": int32(time.Now().Add(1 * time.Minute).Unix()),
		"iss": a.iss,
	})

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(a.key)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse secretr into RSA private key")
	}

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign token")
	}

	return tokenString, nil
}

func (a *authenticator) RoundTrip(r *http.Request) (*http.Response, error) {
	token, err := a.GetToken()
	if err != nil {
		return nil, err
	}

	r.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	resp, err := a.transport.RoundTrip(r)
	return resp, err
}
