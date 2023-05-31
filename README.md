# Wait4X

[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/atkrad/wait4x/ci.yaml?branch=main)](https://github.com/atkrad/wait4x/actions/workflows/ci.yaml) [![Coverage Status](https://coveralls.io/repos/github/atkrad/wait4x/badge.svg?branch=main)](https://coveralls.io/github/atkrad/wait4x?branch=main) [![Go Report Card](https://goreportcard.com/badge/github.com/atkrad/wait4x)](https://goreportcard.com/report/github.com/atkrad/wait4x) [![Docker Pulls](https://img.shields.io/docker/pulls/atkrad/wait4x?logo=docker)](https://hub.docker.com/r/atkrad/wait4x) [![GitHub all releases](https://img.shields.io/github/downloads/atkrad/wait4x/total?logo=github)](https://github.com/atkrad/wait4x/releases) [![Go Reference](https://pkg.go.dev/badge/github.com/atkrad/wait4x.svg)](https://pkg.go.dev/wait4x.dev/v2)

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
    - [RabbitMQ](#rabbitmq)
    - [Temporal](#temporal)

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
    - RabbitMQ
    - Temporal
- **Reverse Checking:** Invert the sense of checking to find a free port or non-ready services
- **Parallel Checking**: You can define multiple inputs to be checked
- **Exponential Backoff Checking**: Retry using an exponential backoff approach to improve efficiency and reduce errors.
- **CI/CD Friendly:** Well-suited to be part of a CI/CD pipeline step
- **Cross Platform:** One single pre-built binary for Linux, Mac OSX, and Windows
- **Importable:** Beside the CLI tool, Wait4X can be imported as a pkg in your Go app
- **Command execution:** Execute your desired command after a successful wait

## Installation

There are many ways to install **Wait4X**

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
curl -#LO https://github.com/atkrad/wait4x/releases/latest/download/wait4x-linux-amd64.tar.gz
tar --one-top-level -xvf wait4x-linux-amd64.tar.gz
cp ./wait4x-linux-amd64/wait4x /usr/local/bin/wait4x
```

#### Mac OSX

```bash
curl -#LO https://github.com/atkrad/wait4x/releases/latest/download/wait4x-darwin-amd64.tar.gz
tar --one-top-level -xvf wait4x-darwin-amd64.tar.gz
cp ./wait4x-darwin-amd64/wait4x /usr/local/bin/wait4x
```

#### Windows

```bash
curl -#LO https://github.com/atkrad/wait4x/releases/latest/download/wait4x-windows-amd64.tar.gz
tar --one-top-level -xvf wait4x-windows-amd64.tar.gz
```

#### Verify SHA256 Checksum

**Wait4X** generates checksum for all binaries with **sha256sum** to prevent against unwanted modification of binaries.
To validate the archive files, download the checksum file which ends in `.sha256sum` for the archive file that you downloaded and use
the `sha256sum` command line tool.

```bash
curl -#LO https://github.com/atkrad/wait4x/releases/latest/download/wait4x-linux-amd64.tar.gz.sha256sum
sha256sum --check wait4x-linux-amd64.tar.gz.sha256sum
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
# To know more about JSON syntax https://github.com/tidwall/gjson/blob/master/SYNTAX.md
wait4x http https://ifconfig.co/json --expect-body-json "user_agent.product"

# Body XPath
wait4x http https://www.kernel.org/ --expect-body-xpath "//*[@id="tux-gear"]"

# Request headers:
wait4x http https://ifconfig.co --request-header "Content-Type: application/json" --request-header "Authorization: Token 123"

# Enable exponential backoff retry
wait4x http https://ifconfig.co --expect-status-code 200 --backoff-policy exponential  --backoff-exponential-max-interval 120s --timeout 120s
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
wait4x mysql user:password@tcp(localhost:5555)/dbname

# Checking MySQL UNIX Socket connection
wait4x mysql username:password@unix(/tmp/mysql.sock)/myDatabase
```
Syntax for the database connection string: https://github.com/go-sql-driver/mysql#dsn-data-source-name

### PostgreSQL

```shell
# Checking PostgreSQL TCP connection
wait4x postgresql 'postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=disable'

# Checking PostgreSQL Unix socket connection
wait4x postgresql 'postgres://bob:secret@/mydb?host=/var/run/postgresql'
```
Syntax for the database URL: https://pkg.go.dev/github.com/lib/pq

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

### RabbitMQ

```shell
# Checking RabbitMQ connection
wait4x rabbitmq 'amqp://127.0.0.1:5672'

# Checking RabbitMQ connection with credentials and vhost
wait4x rabbitmq 'amqp://guest:guest@127.0.0.1:5672/vhost'
```

### Temporal

```shell
# Checking just Temporal server health check
wait4x temporal server 127.0.0.1:7233

# Checking insecure Temporal server (no TLS)
wait4x temporal server 127.0.0.1:7233 --insecure-transport

# Checking a task queue that has registered workers (pollers) or not
wait4x temporal worker 127.0.0.1:7233 --namespace __YOUR_NAMESPACE__ --task-queue __YOUR_TASK_QUEUE__

# Checking the specific a Temporal worker (pollers)
wait4x temporal worker 127.0.0.1:7233 --namespace __YOUR_NAMESPACE__ --task-queue __YOUR_TASK_QUEUE__ --expect-worker-identity-regex ".*@__HOSTNAME__@.*"
```

### Command Execution

We need to wait for something in order to execute something else. This feature is also supported by using `--` after a feature parameter.

Let's have some scenarios:

* As soon as the MySQL server becomes ready, run a Django migration:

```shell
wait4x mysql username:password@unix(/tmp/mysql.sock)/myDatabase -- django migrate
```

* As soon as the new version of the website becomes ready, send an email to the support team:

```shell
wait4x http https://www.kernel.org/ --expect-body-xpath "//*[@id="website-logo"]" -- mail -s "The new version of the website has just been released" support@company < /dev/null
```

* Chain type: when PostgreSQL becomes ready, then wait for RabbitMQ, when it becomes ready, then run the integration tests:

```shell
wait4x postgresql 'postgres://bob:secret@/mydb?host=/var/run/postgresql' -t 1m -- wait4x rabbitmq 'amqp://guest:guest@127.0.0.1:5672/vhost' -t 30s -- ./integration-tests.sh
```

* Using an environment variable in the command:

```shell
EMAIL="support@company" wait4x http https://www.kernel.org/ --expect-body-xpath "//*[@id="website-logo"]" -- mail -s "New version of the website has just released" $EMAIL < /dev/null

# Or

export EMAIL="support@company"
wait4x http https://www.kernel.org/ --expect-body-xpath "//*[@id="website-logo"]" -- mail -s "New version of the website has just released" $EMAIL < /dev/null
```
