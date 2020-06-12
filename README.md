# coffee-tracker
> Google App Engine app in Go 1.12 for tracking coffee with the purpose of knowing when to re-order


Environment Variables (you might want to use something like direnv)

```.env
# mimics for local dev what's set by GAE
# https://cloud.google.com/appengine/docs/standard/go/runtime
export GOOGLE_CLOUD_PROJECT=XXXXXXXXXXXXX

# configuration parameters from you Auth0 app
export OAUTH_AUDIENCE=XXXXXXXXXXXXXX
export OAUTH_ISSUER=XXXXXXXXXXXXXX

# react vars
export REACT_APP_API_DOMAIN_PROD=<GAE PROJECT BASE URL>
export REACT_APP_API_DOMAIN_DEV=http://localhost:8080
```

Running:

in one terminal run the datastore
```
gcloud beta emulators datastore start
```

in another terminal
```
$(gcloud beta emulators datastore env-init)
go run cmd/server/main.go
```
in a 3rd terminal run the frontend
```
cd frontend
yarn start
```

Links:

* https://www.freecodecamp.org/news/how-to-build-a-web-app-with-go-gin-and-react-cffdc473576/
* https://cloud.google.com/secret-manager/docs/quickstart
* https://cloud.google.com/appengine/docs/standard/go/runtime
* https://medium.com/tech-tajawal/deploying-react-app-to-google-app-engine-a6ea0d5af132
    