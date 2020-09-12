# drone-slack

[![Build Status](http://cloud.drone.io/api/badges/drone-plugins/drone-slack/status.svg)](http://cloud.drone.io/drone-plugins/drone-slack)
[![Gitter chat](https://badges.gitter.im/drone/drone.png)](https://gitter.im/drone/drone)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![](https://images.microbadger.com/badges/image/plugins/slack.svg)](https://microbadger.com/images/plugins/slack "Get your own image badge on microbadger.com")
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-slack?status.svg)](http://godoc.org/github.com/drone-plugins/drone-slack)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-slack)](https://goreportcard.com/report/github.com/drone-plugins/drone-slack)

Drone plugin for sending Slack notifications. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-slack/).

## Build

Build the binary with the following commands:

```
go build
```

## Docker

Build the Docker image with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-slack
docker build -f docker/Dockerfile.linux.amd64 -t plugins/slack .
```

## Slack App Creation

The functionality to send messages as PMs to a Slack user using their email as the recipient value requires the use of a Slack Token instead of simple webhook. 
To use a Slack Token you need to create a Slack App with the appropriate OAuth scope.

1. Go to https://api.slack.com/apps/new
2. In the New Slack App modal form give the app a name (eg Drone Bot) and select the Slack Workspace you want to use it in
3. Go to the _OAuth and Permissions_ page for your app and add the following OAuth Scopes in the _Bot Token_ section:
  * chat:write - Used to send messages
  * users:read - Used to find recipient Slack User IDs
  * users:read.email - Used to look up the recipient by a email address
4. (Optional) Go to the _Basic Information_ page and give the app a App Icon (./logo.png in repo), background colour (eg. #2C2D30), Description (eg "Drone is container native CICD platform."). Save your App changes.
5. Go to the _Basic Information_ page and _Install your app to your workspace_
6. Go to the _OAuth and Permissions_ page and use the _Bot User OAuth Access Token_ for the *SLACK_TOKEN* plugin setting in your pipeline.

## Usage

Execute from the working directory:

```
docker run --rm \
  -e SLACK_TOKEN=https://hooks.slack.com/services/... \
  -e PLUGIN_CHANNEL=foo \
  -e PLUGIN_USERNAME=drone \
  -e DRONE_REPO_OWNER=octocat \
  -e DRONE_REPO_NAME=hello-world \
  -e DRONE_COMMIT_SHA=7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=octocat \
  -e DRONE_COMMIT_AUTHOR_EMAIL=octocat@github.com \
  -e DRONE_COMMIT_AUTHOR_AVATAR="https://avatars0.githubusercontent.com/u/583231?s=460&v=4" \
  -e DRONE_COMMIT_AUTHOR_NAME="The Octocat" \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_STATUS=success \
  -e DRONE_BUILD_LINK=http://github.com/octocat/hello-world \
  -e DRONE_TAG=1.0.0 \
  plugins/slack
```
