#!/bin/sh

APPNAME="quickyiitest"

env GOOS=linux GOARCH=amd64 go build -o "./distr/"$APPNAME"_amd64"
