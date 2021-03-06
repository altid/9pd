# 9pd

9pd is an Altid server, used to connect to clients over the 9p protocol.

`go get github.com/altid/server/cmd/9pd`

![Tests](https://github.com/altid/9pd/workflows/Tests/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/altid/9pd)](https://goreportcard.com/report/github.com/altid/9pd) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)

## Usage

`9pd [-t] [-d <dir>] [-c certfile] [-k keyfile] [-u username]`

 - `-t` enables TLS use
 - `-c <certfile>` certificate file for use with TLS connections (Default /etc/ssl/certs/altid.pem)
 - `-k <keyfile>` key file for use with TLS connections (Not required for systems with factotum, default /etc/ssl/private/altid.pem)
 - `-d <dir>` directory to watch (Default /tmp/altid)
 - `-u <username>` Run as user (Default is current user)

## Configuration

```
# altid/config - place this in your operating system's default config directory

service=foo
	#listen_address=192.168.1.144:12345
```
 - listen_address is a more advanced topic, explained here: [Using listen_address](https://altid.github.io/using-listen-address.html)

## Plan9

 - On Plan9, the default certfile is set to $home/lib/altid.pem
 - You must run all services in the same namespace that 9pd is running
