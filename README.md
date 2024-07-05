# Wait4X

[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/atkrad/wait4x/ci.yaml?branch=main)](https://github.com/atkrad/wait4x/actions/workflows/ci.yaml) [![Coverage Status](https://coveralls.io/repos/github/atkrad/wait4x/badge.svg?branch=main)](https://coveralls.io/github/atkrad/wait4x?branch=main) [![Go Report Card](https://goreportcard.com/badge/github.com/atkrad/wait4x)](https://goreportcard.com/report/github.com/atkrad/wait4x) [![Docker Pulls](https://img.shields.io/docker/pulls/atkrad/wait4x?logo=docker)](https://hub.docker.com/r/atkrad/wait4x) [![GitHub all releases](https://img.shields.io/github/downloads/atkrad/wait4x/total?logo=github)](https://github.com/atkrad/wait4x/releases) [![Packaging status](https://repology.org/badge/tiny-repos/wait4x.svg)](https://repology.org/project/wait4x/versions) [![Go Reference](https://pkg.go.dev/badge/github.com/atkrad/wait4x.svg)](https://pkg.go.dev/wait4x.dev/v2)

**Wait4X** allows you to wait for a port or a service to enter the requested state, with a customizable timeout and interval time.

**Table of Contents**

- [Features](#features)
- [Installation](#installation)
  - [With Docker](#with-docker)
  - [From Binary](#from-binary)
    - [Verify SHA256 Checksum](#verify-sha256-checksum)
  - [From Package](#from-package)
    - [Alpine Linux](#on-alpine-linux)
    - [Arch Linux (AUR)](#on-arch-linux-aur)
    - [NixOS](#on-nixos)
    - [Scoop (Windows)](#on-scoop-windows)
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
- [Command Execution](#command-execution)

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
- **Parallel Checking:** You can define multiple inputs to be checked
- **Exponential Backoff Checking:** Retry using an exponential backoff approach to improve efficiency and reduce errors
- **CI/CD Friendly:** Well-suited to be part of a CI/CD pipeline step
- **Cross Platform:** One single pre-built binary for Linux, Mac OSX, and Windows
- **Importable:** Beside the CLI tool, Wait4X can be imported as a pkg in your Go app
- **Command Execution:** Execute your desired command after a successful wait

## Installation

There are several ways to install **Wait4X**.

### With Docker

**Wait4X** provides automatically updated Docker images within Docker Hub. It is possible to always use the latest stable tag.

Pull the image from Docker Hub:

```bash
docker pull atkrad/wait4x:latest
```

Then you can launch the `wait4x` container:

```bash
docker run --rm --name='wait4x' \
    atkrad/wait4x:latest --help
```

### From Binary

Choose the file matching the destination platform from the [release page](https://github.com/atkrad/wait4x/releases), copy the URL and replace it within the commands below:

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

**Wait4X** generates checksums for all binaries with **sha256sum** to prevent against unwanted modification of binaries. To validate the archive files, download the checksum file which ends in `.sha256sum` for the archive file that you downloaded and use the `sha256sum` command line tool.

```bash
curl -#LO https://github.com/atkrad/wait4x/releases/latest/download/wait4x-linux-amd64.tar.gz.sha256sum
sha256sum --check wait4x-linux-amd64.tar.gz.sha256sum
```

### From Package

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

#### On NixOS

You can install **Wait4X** using Nix:

```shell
nix-env -iA nixpkgs.wait4x
```

#### On Scoop (Windows)

You can install **Wait4X** using Scoop:

```shell
scoop install wait4x
```

## Examples

### TCP

```shell
# Check a TCP connection
wait4x tcp 127.0.0.1:9090
```

### HTTP

```shell
# Check an HTTP connection
wait4x http https://ifconfig.co

# Check HTTP connection and expect a specific status code
wait4x http https://ifconfig.co --expect-status-code 200

# Check HTTP connection, status code, and match the response body (using regex)
wait4x http https://ifconfig.co/json --expect-status-code 200 --expect-body='"country":\s"Netherlands"'

# Check an HTTP response header (value in expected header is regex)
wait4x http https://ifconfig.co --expect-header "Authorization=Token 1234ABCD"
wait4x http https://ifconfig.co --expect-header "Authorization=Token"
wait4x http https://ifconfig.co --expect-header "Authorization=Token\s.+"

# Check a body JSON value (value in expected JSON will be processed by gjson)
wait4x http https://ifconfig.co/json --expect-body-json "user_agent.product"

# Check body XPath
wait4x http https://www.kernel.org/ --expect-body-xpath "//*[@id='tux-gear']"

# Set request headers
wait4x http https://ifconfig.co --request-header "Content-Type: application/json" --request-header "Authorization: Token 123"

# Enable exponential backoff retry
wait4x http https://ifconfig.co --expect-status-code 200 --backoff-policy exponential --backoff-exponential-max-interval 120s --timeout 120s
```

### Redis

```shell
# Check Redis connection
wait4x redis redis://127.0.0.1:6379

# Specify username, password, and db
wait4x redis redis://user:password@localhost:6379/1

# Check Redis connection over Unix socket
wait4x redis unix://user:password@/path/to/redis.sock?db=1

# Check a key existence
wait4x redis redis://127.0.0.1:6379 --expect-key FOO

# Check a key existence and match the value (using regex)
wait4x redis redis://127.0.0.1:6379 --expect-key "FOO=^b[A-Z]r$"
```

### MySQL

```shell
# Check MySQL TCP connection
wait4x mysql user:password@tcp(localhost:5555)/dbname

# Check My

SQL Unix socket connection
wait4x mysql username:password@unix(/tmp/mysql.sock)/myDatabase
```
Syntax for the database connection string: [DSN Data Source Name](https://github.com/go-sql-driver/mysql#dsn-data-source-name).

### PostgreSQL

```shell
# Check PostgreSQL TCP connection
wait4x postgresql 'postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=disable'

# Check PostgreSQL Unix socket connection
wait4x postgresql 'postgres://bob:secret@/mydb?host=/var/run/postgresql'
```
Syntax for the database URL: [lib/pq](https://pkg.go.dev/github.com/lib/pq).

### InfluxDB

```shell
# Check InfluxDB connection
wait4x influxdb http://localhost:8086
```

### MongoDB

```shell
# Check MongoDB connection
wait4x mongodb 'mongodb://127.0.0.1:27017'

# Check MongoDB connection with credentials and options
wait4x mongodb 'mongodb://user:pass@127.0.0.1:27017/?maxPoolSize=20&w=majority'
```

### RabbitMQ

```shell
# Check RabbitMQ connection
wait4x rabbitmq 'amqp://127.0.0.1:5672'

# Check RabbitMQ connection with credentials and vhost
wait4x rabbitmq 'amqp://guest:guest@127.0.0.1:5672/vhost'
```

### Temporal

```shell
# Check Temporal server health check
wait4x temporal server 127.0.0.1:7233

# Check insecure Temporal server (no TLS)
wait4x temporal server 127.0.0.1:7233 --insecure-transport

# Check a task queue that has registered workers (pollers)
wait4x temporal worker 127.0.0.1:7233 --namespace __YOUR_NAMESPACE__ --task-queue __YOUR_TASK_QUEUE__

# Check a specific Temporal worker (pollers)
wait4x temporal worker 127.0.0.1:7233 --namespace __YOUR_NAMESPACE__ --task-queue __YOUR_TASK_QUEUE__ --expect-worker-identity-regex ".*@__HOSTNAME__@.*"
```

## Command Execution

You can wait for something to execute a command afterward. Use `--` after a feature parameter.

### Scenarios

* Run a Django migration as soon as the MySQL server is ready:

```shell
wait4x mysql username:password@unix(/tmp/mysql.sock)/myDatabase -- django migrate
```

* Send an email to the support team as soon as the new version of the website is ready:

```shell
wait4x http https://www.kernel.org/ --expect-body-xpath "//*[@id='website-logo']" -- mail -s "The new version of the website has just been released" support@company < /dev/null
```

* Chain checks: when PostgreSQL is ready, wait for RabbitMQ, then run the integration tests:

```shell
wait4x postgresql 'postgres://bob:secret@/mydb?host=/var/run/postgresql' -t 1m -- wait4x rabbitmq 'amqp://guest:guest@127.0.0.1:5672/vhost' -t 30s -- ./integration-tests.sh
```

* Use an environment variable in the command:

```shell
EMAIL="support@company" wait4x http https://www.kernel.org/ --expect-body-xpath "//*[@id='website-logo']" -- mail -s "New version of the website has just released" $EMAIL < /dev/null

# Or

export EMAIL="support@company"
wait4x http https://www.kernel.org/ --expect-body-xpath "//*[@id='website-logo']" -- mail -s "New version of the website has just released" $EMAIL < /dev/null
```
