# github-app-installations

github-app-installations queries github to find all orgs/users who have installed this app. It needs your app's private-key and installation-id for generating the token and making request to github api. You can find this information from `https://github.com/settings/apps/<your-app-name>`

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
    "repositorySelection": "All Repositories"
  },
  {
    "ghLogin": "reactiverse",
    "repositorySelection": "Selected Repositories"
  },
  {
    "ghLogin": "pmlopes",
    "repositorySelection": "Selected Repositories"
  },
  {
    "ghLogin": "asyncapi",
    "repositorySelection": "All Repositories"
  },
  {
    "ghLogin": "google",
    "repositorySelection": "Selected Repositories"
  },
  {
    "ghLogin": "asyncy",
    "repositorySelection": "All Repositories"
  },
  {
    "ghLogin": "Ewocker",
    "repositorySelection": "Selected Repositories"
  },
  {
    "ghLogin": "jetstack",
    "repositorySelection": "Selected Repositories"
  },
  {
    "ghLogin": "openfaas",
    "repositorySelection": "All Repositories"
  },
  {
    "ghLogin": "alexellis",
    "repositorySelection": "Selected Repositories"
  },
  {
    "ghLogin": "rajatjindal",
    "repositorySelection": "All Repositories"
  }
]
```

# credits
This is a [openfaas](https://github.com/openfaas/faas) function which is built and deployed on [openfaas-cloud-community-cluster](http://github.com/apps/openfaas-cloud-community-cluster).