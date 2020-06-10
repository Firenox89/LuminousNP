#!/usr/bin/env bash
env GOOS=linux GOARCH=arm GOARM=6 go build -o nodemcu-controller main.go
