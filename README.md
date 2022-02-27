# Wait4X

[![Build Status](https://cloud.drone.io/api/badges/atkrad/wait4x/status.svg)](https://cloud.drone.io/atkrad/wait4x) [![Coverage Status](https://coveralls.io/repos/github/atkrad/wait4x/badge.svg?branch=master)](https://coveralls.io/github/atkrad/wait4x?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/atkrad/wait4x)](https://goreportcard.com/report/github.com/atkrad/wait4x) [![Docker Pulls](https://img.shields.io/docker/pulls/atkrad/wait4x)](https://hub.docker.com/r/atkrad/wait4x) [![Go Reference](https://pkg.go.dev/badge/github.com/atkrad/wait4x.svg)](https://pkg.go.dev/github.com/atkrad/wait4x)

**Wait4X** allows you to wait for a port or a service to enter the requested state, with a customizable timeout and
interval time.

**Table of Contents**

- [Features](#features)
- [Installation](#installation)
    - [with Docker](#with-docker)
    - [From binary](#from-binary)
        - [Verify SHA256 Checksum](#verify-sha256-checksum)
    - [From package](#from-package)
        - [Alpine Linux](#on-alpine-linux)
        - [Arch Linux (AUR)](#on-arch-linux-aur)
- [Examples](#examples)
    - [TCP](#tcp)
    - [HTTP](#http)
    - [Redis](#redis)
    - [MySQL](#mysql)
    - [PostgreSQL](#postgresql)
    - [InfluxDB](#influxdb)
    - [MongoDB](#mongodb)

## Features

- **Supports various protocols:**
    - TCP
    - HTTP
- **Supports various services:**
    - Redis
    - MySQL
    - PostgreSQL
    - InfluxDB
    - MongoDB
- **Reverse Checking:** Invert the sense of checking to find a free port or non-ready services
- **CI/CD Friendly:** Well-suited to be part of a CI/CD pipeline step
- **Cross Platform:** One single pre-built binary for Linux, Mac OSX, and Windows
- **Importable:** Beside the CLI tool, Wait4X can be imported as a pkg in your Go app

## Installation

There are many different ways to install **Wait4X**

### with Docker

**Wait4X** provides automatically updated Docker images within Docker Hub. It is possible to always use the latest
stable tag.

Pull the image from the docker index.

```bash
docker pull atkrad/wait4x:latest
```

then you can launch the `wait4x` container.

```bash
docker run --rm --name='wait4x' \
    atkrad/wait4x:latest --help
```

### From binary

Choose the file matching the destination platform from the [release page](https://github.com/atkrad/wait4x/releases),
copy the URL and replace the URL within the commands below:

#### Linux

```bash
curl -L https://github.com/atkrad/wait4x/releases/latest/download/wait4x-linux-amd64 -o /usr/local/bin/wait4x
chmod +x /usr/local/bin/wait4x
```

#### Mac OSX

```bash
curl -L https://github.com/atkrad/wait4x/releases/latest/download/wait4x-darwin-amd64 -o /usr/local/bin/wait4x
chmod +x /usr/local/bin/wait4x
```

#### Windows

```bash
curl -L https://github.com/atkrad/wait4x/releases/latest/download/wait4x-windows-amd64 -o wait4x.exe
```

#### Verify SHA256 Checksum

**Wait4X** generates checksum for all binaries with **sha256sum** to prevent against unwanted modification of binaries.
To validate the binary, download the checksum file which ends in `.sha256sum` for the binary you downloaded and use
the `sha256sum` command line tool.

```bash
curl -SLO https://github.com/atkrad/wait4x/releases/latest/download/wait4x-linux-amd64.sha256sum
sha256sum --check wait4x-linux-amd64.sha256sum
```

### From package

You can find the **Wait4X** package in some Linux distributions.

[![Packaging status](https://repology.org/badge/vertical-allrepos/wait4x.svg)](https://repology.org/project/wait4x/versions)

#### On Alpine Linux

You can install the [wait4x](https://pkgs.alpinelinux.org/packages?name=wait4x) package from the official sources:

```shell
apk add wait4x
```

#### On Arch Linux (AUR)

You can install the [wait4x](https://aur.archlinux.org/packages/wait4x/) package from the Arch User Repository:

```shell
yay -S wait4x
```

## Examples

### TCP

```shell
# If you want checking just tcp connection
wait4x tcp 127.0.0.1:9090
```

### HTTP

```shell
# If you want checking just http connection
wait4x http https://ifconfig.co

# If you want checking http connection and expect specify http status code
wait4x http https://ifconfig.co --expect-status-code 200

# If you want checking http connection, status code and match the response body.
# Note: You can write any regex that compatible with Golang syntax (https://pkg.go.dev/regexp/syntax#hdr-Syntax)
wait4x http https://ifconfig.co/json --expect-status-code 200 --expect-body='"country":\s"Netherlands"'

# If you want to check a http response header
# NOTE: the value in the expected header is regex.
# Sample response header: Authorization Token 1234ABCD
# You can match it by these ways:

# Full key value:
wait4x http https://ifconfig.co --expect-header "Authorization=Token 1234ABCD"

# Value starts with:
wait4x http https://ifconfig.co --expect-header "Authorization=Token"

# Regex value:
wait4x http https://ifconfig.co --expect-header "Authorization=Token\s.+"

# Body JSON value:
# Note: the value in the expected JSON will be processed by gjson.
# Note: the complete response MUST be in JSON format to be processed. If one part of the response is JSON
# and the rest is something else like HTML, please use --expect-body instead.
wait4x http https://ifconfig.co/json --expect-body-json "user_agent.product"
To know more about JSON syntax https://github.com/tidwall/gjson/blob/master/SYNTAX.md
```

### Redis

```shell
# Checking Redis connection
wait4x redis redis://127.0.0.1:6379

# Specify username, password and db
wait4x redis redis://user:password@localhost:6379/1

# Checking Redis connection over unix socket
wait4x redis unix://user:password@/path/to/redis.sock?db=1

# Checking a key existence
wait4x redis redis://127.0.0.1:6379 --expect-key FOO

# Checking a key existence and matching the value
# Note: You can write any regex that compatible with Golang syntax (https://pkg.go.dev/regexp/syntax#hdr-Syntax)
wait4x redis redis://127.0.0.1:6379 --expect-key "FOO=^b[A-Z]r$"
```

### MySQL

```shell
# Checking MySQL TCP connection
wait4x mysql user:password@tcp(localhost:5555)/dbname?tls=skip-verify

# Checking MySQL UNIX Socket connection
wait4x mysql username:password@unix(/tmp/mysql.sock)/myDatabase
```

### PostgreSQL

```shell
# Checking PostgreSQL TCP connection
wait4x postgresql 'postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full'

# Checking PostgreSQL Unix socket connection
wait4x postgresql 'postgres://bob:secret@/mydb?host=/var/run/postgresql'
```

### InfluxDB

```shell
# Checking InfluxDB connection
wait4x influxdb http://localhost:8086
```

### MongoDB

```shell
# Checking MongoDB connection
wait4x mongodb 'mongodb://127.0.0.1:27017'

# Checking MongoDB connection with credentials and options
wait4x mongodb 'mongodb://user:pass@127.0.0.1:27017/?maxPoolSize=20&w=majority'
```
