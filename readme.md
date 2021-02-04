
# Vigilate

This is the source code for the second project in the Udemy course Working with Websockets in Go (Golang).

A dead simple monitoring service, intended to replace things like Nagios.

## Build

Build in the normal way:

~~~
go build -o vigilate cmd/web/*.go
~~~

Or for a particular platform:

~~~
env GOOS=linux GOARCH=amd64 go build -o vigilate cmd/web/*.go
~~~

## Requirements

Vigilate requires:
- Postgres 11 or later (db is set up as a repository, so other databases are possible)
- An account with [Pusher](https://pusher.com/), or a Pusher alternative 
(like [ipe](https://github.com/dimiro1/ipe))

## Run

Run with flags:

~~~
./vigilate \
-dbuser='tcs' \
-pusherHost='some.host.com' \
-pusherSecret='somesecret' \
-pusherKey='somekey' 
~~~~

## All Flags

~~~~
tcs@grendel vigilate-udemy % ./vigilate -help
Usage of ./vigilate:
  -db string
        database name (default "vigilate")
  -dbhost string
        database host (default "localhost")
  -dbport string
        database port (default "5432")
  -dbssl string
        database ssl setting (default "disable")
  -dbuser string
        database user
  -domain string
        domain name (e.g. example.com) (default "localhost")
  -identifier string
        unique identifier (default "vigilate")
  -port string
        port to listen on (default ":4000")
  -production
        application is in production
  -pusherApp string
        pusher app id (default "9")
  -pusherHost string
        pusher host
  -pusherKey string
        pusher key
  -pusherPort string
        pusher port (default "443")
  -pusherSecret string
        pusher secret
~~~~


## Sample supervisor script

~~~
[program:vigilate]
command=/var/www/sites/vigilate/vigilate -port=':4000' -domain='example.com' -production=true -dbuser='postgres' -pusherHost='some.pusher.host.com' -pusherSecret='somescret' -pusherKey='somekey'
directory=/var/www/sites/vigilate
autostart=true
autorestart=true

stdout_logfile=/var/www/sites/vigilate/logs/vigilate.log
~~~