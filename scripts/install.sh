#! /bin/bash
set -xe

# export variables
export $(grep -v '^#' ./.env | xargs)

if [[ -n "$CHANNEL" ]]
then
    echo "channel found"
else
    exit 1
fi

CHART_NAME="slack-twitch-slash"

helm install slack-slash ./charts/slack-twitch-slash --create-namespace -n $CHART_NAME \
    --set chart."name=${CHART_NAME}" \
    --set app."bearertoken=${BEARERTOKEN}" \
    --set app."callbackurl=${CALLBACKURL}" \
    --set app."client_id=${CLIENT_ID}" \
    --set app."client_secret=${CLIENT_SECRET}" \
    --set app."slackapikey=${SLACKAPIKEY}" \
    --set app."channel=${CHANNEL}" \
    --set imageCredentials."username=${DOCKER_USERNAME}" \
    --set imageCredentials."password=${DOCKER_PASSWORD}" \
    --set imageCredentials."email=${DOCKER_EMAIL}"
