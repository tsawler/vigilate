#!/bin/zsh

# This is the bare minimum to run in development. For full list of flags,
# run ./vigilate -help

go build -o vigilate cmd/web/*.go && ./vigilate \
-dbuser='someuser' \
-pusherHost='pusher.com' \
-pusherSecret='somesecret' \
-pusherKey='somekey' \
-pusherApp="1"
