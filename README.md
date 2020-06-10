# coffee-tracker
> Google App Engine app in Go 1.12 for tracking coffee with the purpose of knowing when to re-order


Environment Variables (you might want to use something like direnv)

```.env
# mimics for local dev what's set by GAE
# https://cloud.google.com/appengine/docs/standard/go/runtime
export GOOGLE_CLOUD_PROJECT=XXXXXXXXXXXXX

# configuration parameters from you Auth0 app
export AUTH0_AUD=XXXXXXXXXXXXXX
export AUTH0_ISS=XXXXXXXXXXXXXX
```

Links:

* https://www.freecodecamp.org/news/how-to-build-a-web-app-with-go-gin-and-react-cffdc473576/
* https://cloud.google.com/secret-manager/docs/quickstart
* https://cloud.google.com/appengine/docs/standard/go/runtime
* https://medium.com/tech-tajawal/deploying-react-app-to-google-app-engine-a6ea0d5af132
    