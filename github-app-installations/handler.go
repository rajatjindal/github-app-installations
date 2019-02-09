package function

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type payload struct {
	iat int32
	exp int32
	iss int
}

const (
	//KeyHeader that contains private key for app
	KeyHeader = "X-Private-Key"

	//InstallationIDHeader contains the installation id of app
	InstallationIDHeader = "X-Installation-Id"
)

type resp struct {
	GithubLogin         string `json:"ghLogin"`
	RepositorySelection string `json:"repositorySelection"`
}

// func main() {
// 	http.HandleFunc("/", Handle)
// 	http.ListenAndServe(":8081", nil)
// }

//Handle handles the function call
func Handle(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(KeyHeader)
	installationID := r.Header.Get(InstallationIDHeader)

	if key == "" || installationID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("X-Private-Key and X-Installation-Id headers are mandatory"))
		return
	}

	iss, err := strconv.Atoi(installationID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("invalid value of X-Installation-Id [%s]. it should be an integer", installationID)))
		return
	}

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("invalid value of X-Private-Key. failed to decode")))
		return
	}

	transport := &authenticator{
		transport: http.DefaultTransport,
		iss:       iss,
		key:       decodedKey,
	}

	client := &http.Client{Transport: transport}
	c := github.NewClient(client)
	l := &github.ListOptions{}
	i, _, err := c.Apps.ListInstallations(context.Background(), l)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	out := []*resp{}
	for _, ii := range i {
		var respositorySelection string
		if stringValue(ii.RepositorySelection) == "all" {
			respositorySelection = "All Repositories"
		}

		if stringValue(ii.RepositorySelection) == "selected" {
			respositorySelection = "Selected Repositories"
		}

		out = append(out,
			&resp{
				GithubLogin:         stringValue(ii.Account.Login),
				RepositorySelection: respositorySelection,
			})
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

func newAuth() *authenticator {
	return &authenticator{
		transport: http.DefaultTransport,
	}
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
