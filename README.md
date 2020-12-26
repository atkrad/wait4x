# Wait4X 
[![Build Status](https://cloud.drone.io/api/badges/atkrad/wait4x/status.svg)](https://cloud.drone.io/atkrad/wait4x) [![codecov](https://img.shields.io/codecov/c/github/atkrad/wait4x)](https://codecov.io/gh/atkrad/wait4x) [![Go Report Card](https://goreportcard.com/badge/github.com/atkrad/wait4x)](https://goreportcard.com/report/github.com/atkrad/wait4x) [![Docker Pulls](https://img.shields.io/docker/pulls/atkrad/wait4x)](https://hub.docker.com/r/atkrad/wait4x) [![Go Reference](https://pkg.go.dev/badge/github.com/atkrad/wait4x.svg)](https://pkg.go.dev/github.com/atkrad/wait4x)

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Wait4X](#wait4x)
    - [Introduction](#introduction)
    - [Installation](#installation)
        - [with Docker](#with-docker)
        - [From binary](#from-binary)
            - [Verify SHA256 Checksum](#verify-sha256-checksum)
        - [From package](#from-package)
            - [Alpine Linux](#alpine-linux)

<!-- markdown-toc end -->

## Introduction
**Wait4X** is a cli tool to wait for everything! It can be wait for a port to open or enter to rquested state.

**Features:**
- Checking TCP connection
- Checking HTTP connection and also status code and response body
- Checking Redis connection (Ping request) and key existence
- Checking MySQL connection
- Checking PostgreSQL connection

## Installation
There are many different methods to install **Wait4X**

### with Docker
**Wait4X** provides automatically updated Docker images within Docker Hub. It is possible to always use the latest stable tag.

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
Choose the file matching the destination platform from the [release page](https://github.com/atkrad/wait4x/releases), copy the URL and replace the URL within the commands below:
```bash
curl -L https://github.com/atkrad/wait4x/releases/download/v0.4.0/wait4x-linux-amd64 -o /usr/local/bin/wait4x
chmod +x /usr/local/bin/wait4x
```

#### Verify SHA256 Checksum
**Wait4X** generates checksum for all binaries with **sha256sum** to prevent against unwanted modification of binaries. To validate the binary, download the checksum file which ends in `.sha256sum` for the binary you downloaded and use the `sha256sum` command line tool.
```bash
curl -SLO https://github.com/atkrad/wait4x/releases/download/v0.4.0/wait4x-linux-amd64.sha256sum
sha256sum --check wait4x-linux-amd64.sha256sum
```

### From package
You can install **Wait4X** from package, currentlly we only support **Alpine Linux**.

#### Alpine Linux
Currently you can find it in [edge](https://pkgs.alpinelinux.org/packages?name=wait4x&branch=edge) repository so you have to ensure the repository has enabled. ([Repository pinning](https://wiki.alpinelinux.org/wiki/Alpine_Linux_package_management#Repository_pinning))
```bash
apk add wait4x@edgecommunity
```
