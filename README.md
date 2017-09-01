## Install

go get github.com/jcnnghm/cmdtrack

## Deploying

First, download the google cloud sdk from https://cloud.google.com/sdk/.  Unzip it and move the directory to `~/google-cloud-sdk`.   Then run `~/google-cloud-sdk/install.sh` to install it - it's already setup in `shared.rc`.  Open a new shell so gcloud gets added to the path.

From there, these are the commands I used:

```
cd go/src/github.com/jcnnghm/cmdtrack
gcloud config set project cmdtrack-1127
go get ./...
cd server
go get github.com/gorilla/context
gcloud app deploy
```
