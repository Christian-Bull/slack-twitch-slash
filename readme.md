![build](https://github.com/Christian-Bull/slack-twitch-slash/actions/workflows/docker-image.yml/badge.svg)


# slack-twitch-slash

Simple PoC slackbot that allows subscription to twitch events (i.e. channel goes live)

Primarily used to test k8s deploys and scaling

Export all env vars:  
`export $(grep -v '^#' .env | xargs)`