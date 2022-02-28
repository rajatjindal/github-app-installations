# github-app-installations

github-app-installations queries github to find all orgs/users who have installed a particular app. It needs the app's private-key and installation-id for generating the token and making request to github api. You can find this information from `https://github.com/settings/apps/<your-app-name>`

### Security disclaimer

You should NEVER submit private keys used for sensible personal data or production. Be aware that by giving away your app's private key, it could theoretically be abused to attack its users. A professional developer should be able to do this without sharing your key with a 3rd party. Even GitHub would never ask you to share private keys!

A sample request is as follows:

## generate base64 value of your private key
```
cat /Users/rajatjindal/.ssh/your-app.private-key.pem | base64
```

## make request to github-app-installations to find the details
```
curl -vXGET 'https://rajatjindal.o6s.io/github-app-installations' \
    -H'X-App-Id: <app-id>' \
    -H'X-Private-Key: <base64-value-generated-above>
```

and you will get output like follows:

```
[
  {
    "ghLogin": "helm",
    "orgUserURL": "https://github.com/helm",
    "repositorySelection": "all"
  },
  {
    "ghLogin": "reactiverse",
    "orgUserURL": "https://github.com/reactiverse",
    "repositorySelection": "selected",
    "repositories": [
      {
        "name": "es4x",
        "HtmlURL": "https://github.com/reactiverse/es4x"
      }
    ]
  },
  {
    "ghLogin": "pmlopes",
    "orgUserURL": "https://github.com/pmlopes",
    "repositorySelection": "selected",
    "repositories": [
      {
        "name": "vertx-starter",
        "HtmlURL": "https://github.com/pmlopes/vertx-starter"
      }
    ]
  },
  {
    "ghLogin": "asyncapi",
    "orgUserURL": "https://github.com/asyncapi",
    "repositorySelection": "all"
  },
  {
    "ghLogin": "google",
    "orgUserURL": "https://github.com/google",
    "repositorySelection": "selected",
    "repositories": [
      {
        "name": "go-github",
        "HtmlURL": "https://github.com/google/go-github"
      }
    ]
  },
  {
    "ghLogin": "asyncy",
    "orgUserURL": "https://github.com/asyncy",
    "repositorySelection": "all"
  },
  {
    "ghLogin": "Ewocker",
    "orgUserURL": "https://github.com/Ewocker",
    "repositorySelection": "selected",
    "repositories": [
      {
        "name": "vue-lodash",
        "HtmlURL": "https://github.com/Ewocker/vue-lodash"
      }
    ]
  },
  {
    "ghLogin": "jetstack",
    "orgUserURL": "https://github.com/jetstack",
    "repositorySelection": "selected",
    "repositories": [
      {
        "name": "cert-manager",
        "HtmlURL": "https://github.com/jetstack/cert-manager"
      }
    ]
  },
  {
    "ghLogin": "openfaas",
    "orgUserURL": "https://github.com/openfaas",
    "repositorySelection": "all"
  },
  {
    "ghLogin": "alexellis",
    "orgUserURL": "https://github.com/alexellis",
    "repositorySelection": "selected",
    "repositories": [
      {
        "name": "ubiquitous-octo-guacamole",
        "HtmlURL": "https://github.com/alexellis/ubiquitous-octo-guacamole"
      }
    ]
  },
  {
    "ghLogin": "rajatjindal",
    "orgUserURL": "https://github.com/rajatjindal",
    "repositorySelection": "all"
  }
]
```

It also supports a readme table format. If you add `X-Resp-Format: readme` header to the request, the output is returned in readme table format as follows:

```
| Org/User | Repository |
| ------ | ------ |
| [helm](https://github.com/helm) | [All](https://github.com/helm) |
| [reactiverse](https://github.com/reactiverse) | [es4x](https://github.com/reactiverse/es4x) |
| [pmlopes](https://github.com/pmlopes) | [vertx-starter](https://github.com/pmlopes/vertx-starter) |
| [asyncapi](https://github.com/asyncapi) | [All](https://github.com/asyncapi) |
| [google](https://github.com/google) | [go-github](https://github.com/google/go-github) |
| [asyncy](https://github.com/asyncy) | [All](https://github.com/asyncy) |
| [Ewocker](https://github.com/Ewocker) | [vue-lodash](https://github.com/Ewocker/vue-lodash) |
| [jetstack](https://github.com/jetstack) | [cert-manager](https://github.com/jetstack/cert-manager) |
| [openfaas](https://github.com/openfaas) | [All](https://github.com/openfaas) |
| [alexellis](https://github.com/alexellis) | [ubiquitous-octo-guacamole](https://github.com/alexellis/ubiquitous-octo-guacamole) |
| [rajatjindal](https://github.com/rajatjindal) | [All](https://github.com/rajatjindal) |
```

# credits
This is a [openfaas](https://github.com/openfaas/faas) function which is built and deployed on [openfaas-cloud-community-cluster](http://github.com/apps/openfaas-cloud-community-cluster).
