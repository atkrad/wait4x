# Wait4X

[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/atkrad/wait4x/ci.yaml?branch=main&style=flat-square)](https://github.com/atkrad/wait4x/actions/workflows/ci.yaml) [![Coverage Status](https://img.shields.io/coverallsCoverage/github/atkrad/wait4x?branch=main&style=flat-square
)](https://coveralls.io/github/atkrad/wait4x?branch=main) [![Go Report Card](https://goreportcard.com/badge/wait4x.dev/v3?style=flat-square)](https://goreportcard.com/report/wait4x.dev/v3) [![Docker Pulls](https://img.shields.io/docker/pulls/atkrad/wait4x?logo=docker&style=flat-square)](https://hub.docker.com/r/atkrad/wait4x) [![GitHub all releases](https://img.shields.io/github/downloads/atkrad/wait4x/total?logo=github&style=flat-square)](https://github.com/atkrad/wait4x/releases) [![Packaging status](https://img.shields.io/repology/repositories/wait4x?style=flat-square
)](https://repology.org/project/wait4x/versions) [![Go Reference](https://img.shields.io/badge/reference-007D9C.svg?style=flat-square&logo=go&logoColor=white&labelColor=5C5C5C)](https://pkg.go.dev/wait4x.dev/v3)

## Introduction

**Wait4X** is a versatile command-line tool designed to wait for various ports or services to reach a specified state. It supports multiple protocols and services, making it an essential tool for CI/CD pipelines, automated testing, and deployment processes.

## Features

- **Supports various protocols:**
    - TCP
    - HTTP
    - DNS
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
- **Importable:** Wait4X can be imported as a package in your Go application
- **Command Execution:** Execute your desired command after a successful wait

## Installation

There are several ways to install **Wait4X**.

### With Docker

**Wait4X** provides automatically updated Docker images within Docker Hub. It is possible to always use the latest stable tag.

Pull the image from Docker Hub:

```bash
docker pull atkrad/wait4x:latest
```

Then launch the `wait4x` container:

```bash
docker run --rm --name='wait4x' atkrad/wait4x:latest --help
```

### From Binary

Choose the file matching your platform from the [release page](https://github.com/atkrad/wait4x/releases), then run the following commands:

#### Linux

```bash
curl -#LO https://github.com/atkrad/wait4x/releases/latest/download/wait4x-linux-amd64.tar.gz
tar --one-top-level -xvf wait4x-linux-amd64.tar.gz
sudo cp ./wait4x-linux-amd64/wait4x /usr/local/bin/wait4x
```

#### Mac OSX

```bash
curl -#LO https://github.com/atkrad/wait4x/releases/latest/download/wait4x-darwin-amd64.tar.gz
tar --one-top-level -xvf wait4x-darwin-amd64.tar.gz
sudo cp ./wait4x-darwin-amd64/wait4x /usr/local/bin/wait4x
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

[![Packaging status](https://repology.org/badge/vertical-allrepos/wait4x.svg?exclude_unsupported=1)](https://repology.org/project/wait4x/versions)

#### On Alpine Linux

Wait4X is available in the Alpine Linux [community](https://pkgs.alpinelinux.org/packages?name=wait4x) repository.

```shell
apk add wait4x
```

#### On Arch Linux (AUR)

Wait4X is available in the Arch User Repository ([AUR](https://aur.archlinux.org/packages/wait4x/)).

```shell
yay -S wait4x-bin
```

#### On NixOS

Wait4X is available in the NixOS repository.

```shell
nix-env -iA nixpkgs.wait4x
```

#### On Scoop (Windows)

Wait4X is available in the Scoop bucket.

```bash
scoop install wait4x
```

## Examples

### TCP

Check TCP connection:

```bash
wait4x tcp 127.0.0.1:9090
```
This command waits until the TCP port `9090` on `127.0.0.1` is available.

### DNS

```shell
# Check A records existence
wait4x dns A wait4x.dev

# Check A records with expected ips
wait4x dns A wait4x.dev --expected-ip 172.67.154.180

# Check A records by defined nameserver
wait4x dns A wait4x.dev --expected-ip 172.67.154.180 -n gordon.ns.cloudflare.com

# Check AAAA records existence
wait4x dns AAAA wait4x.dev

# Check AAAA records with expected ips
wait4x dns AAAA wait4x.dev --expected-ip '2606:4700:3033::ac43:9ab4'

# Check AAAA records by defined nameserver
wait4x dns AAAA wait4x.dev --expected-ip '2606:4700:3033::ac43:9ab4' -n gordon.ns.cloudflare.com

# Check CNAME record existence
wait4x dns CNAME 172.67.154.180

# Check CNAME records with expected ips
wait4x dns CNAME 172.67.154.180 --expected-domain wait4x.dev

# Check CNAME record by defined nameserver
wait4x dns CNAME 172.67.154.180 --expected-domain wait4x.dev -n gordon.ns.cloudflare.com

# Check MX records existence
wait4x dns MX wait4x.dev

# Check MX records with expected domains
wait4x dns MX wait4x.dev --expected-domain 'route1.mx.cloudflare.net'

# Check MX records by defined nameserver
wait4x dns MX wait4x.dev --expected-domain 'route1.mx.cloudflare.net.' -n gordon.ns.cloudflare.com

# Check NS records existence
wait4x dns NS wait4x.dev

# Check NS records with expected nameservers
wait4x dns NS wait4x.dev --expected-nameserver 'emma.ns.cloudflare.com'

# Check NS records by defined nameserver
wait4x dns NS wait4x.dev --expected-nameserver 'emma.ns.cloudflare.com' -n gordon.ns.cloudflare.com

# Check TXT records existence
wait4x dns TXT wait4x.dev

# Check TXT records with expected values
wait4x dns TXT wait4x.dev --expected-value 'include:_spf.mx.cloudflare.net'

# Check TXT records by defined nameserver
wait4x dns TXT wait4x.dev --expected-value 'include:_spf.mx.cloudflare.net' -n gordon.ns.cloudflare.com
```

### HTTP

Check HTTP connection and expect a specific status code:

```shell
wait4x http https://ifconfig.co --expect-status-code 200
```
This command waits until the URL `https://ifconfig.co` returns an HTTP status code of `200`.

Check HTTP connection, status code, and match the response body:

```shell
wait4x http https://ifconfig.co/json --expect-status-code 200 --expect-body-regex='"country":\s"Netherlands"'
```

Check an HTTP response header value:
```shell
wait4x http https://ifconfig.co --expect-header "Authorization=Token\s.+"
```
This command waits until the URL `https://ifconfig.co` returns an HTTP status code of `200` and the response header matches the provided regex pattern.

Check a body JSON value (value in expected JSON will be processed by gjson):
```shell
wait4x http https://ifconfig.co/json --expect-body-json "user_agent.product"
```
This command waits until the URL `https://ifconfig.co/json` returns an HTTP status code of `200` and the response body matches the provided [GJSON](https://github.com/tidwall/gjson?tab=readme-ov-file#path-syntax) path.

Check body XPath value:
```shell
wait4x http https://www.kernel.org/ --expect-body-xpath "//*[@id='tux-gear']"
```
This command waits until the URL `https://www.kernel.org/` returns an HTTP status code of `200` and the response body matches the provided XPath path.

Set request headers:

```shell
wait4x http https://ifconfig.co --request-header "Content-Type: application/json" --request-header "Authorization: Token 123"
```
This command sets the `Content-Type` and `Authorization` HTTP request headers and waits until the URL `https://ifconfig.co` returns an HTTP status code of `200`.

### Redis

Check Redis connection:
```shell
wait4x redis redis://127.0.0.1:6379
```
This command waits until the Redis server on `127.0.0.1:6379` is ready.

Check Redis connection (with database and credentials):
```shell
wait4x redis redis://user:password@localhost:6379/1
```
This command waits until the Redis server on `localhost:6379` is ready to accept connections to the `1` database.

Check Redis connection (Unix socket):
```shell
wait4x redis unix://user:password@/path/to/redis.sock?db=1
```
This command waits until the Redis server on `/path/to/redis.sock` is ready to accept connections to the `1` database.

Check Redis connection and match a key:
```shell
wait4x redis redis://127.0.0.1:6379 --expect-key FOO
```
This command waits until the Redis server on `127.0.0.1:6379` is ready and the key `FOO` exists.

Check Redis connection and match a pair of key and value:
```shell
wait4x redis redis://127.0.0.1:6379 --expect-key "FOO=^b[A-Z]r$"
```
This command waits until the Redis server on `127.0.0.1:6379` is ready and the key `FOO` exists and the value matches the provided regex pattern.

### MySQL

Check MySQL connection (TCP):

```shell
wait4x mysql 'user:password@tcp(localhost:5555)/dbname'
```
This command waits until the MySQL server on `127.0.0.1:3306` is ready to accept connections to the `dbname` database.

Check MySQL connection (Unix socket):
```shell
wait4x mysql 'username:password@unix(/tmp/mysql.sock)/myDatabase'
```
This command waits until the MySQL server on `/tmp/mysql.sock` is ready to accept connections to the `myDatabase` database.

**Note:** Syntax for the database connection string: [DSN Data Source Name](https://github.com/go-sql-driver/mysql#dsn-data-source-name).

### PostgreSQL

Check PostgreSQL connection (TCP):

```shell
wait4x postgresql 'postgres://bob:secret@1.2.3.4:5432/mydatabase?sslmode=disable'
```
This command waits until the PostgreSQL server on `127.0.0.1:5432` is ready to accept connections to the `mydatabase` database.

Check PostgreSQL connection (Unix socket):
```shell
wait4x postgresql 'postgres://bob:secret@/mydb?host=/var/run/postgresql'
```
This command waits until the PostgreSQL server on `/var/run/postgresql` is ready to accept connections to the `mydb` database.

*Note:* Syntax for the database DSN: [lib/pq](https://pkg.go.dev/github.com/lib/pq).

### InfluxDB

Check InfluxDB connection:
```shell
wait4x influxdb http://localhost:8086
```
This command waits until the InfluxDB server on `localhost:8086` is ready.

### MongoDB

Check MongoDB connection (with credentials and options):

```shell
wait4x mongodb 'mongodb://user:pass@127.0.0.1:27017/?maxPoolSize=20&w=majority'
```
This command waits until the MongoDB server on `127.0.0.1:27017` is ready.

### RabbitMQ

Check RabbitMQ connection (with credentials and vhost):

```shell
wait4x rabbitmq 'amqp://guest:guest@127.0.0.1:5672/vhost'
```
This command waits until the RabbitMQ server on `localhost:5672` is ready.

### Temporal

Check Temporal server connection:

```shell
wait4x temporal server 127.0.0.1:7233
```
This command waits until the Temporal server on `127.0.0.1:7233` is ready.

Check insecure Temporal server (no TLS):

```shell
wait4x temporal server 127.0.0.1:7233 --insecure-transport
```

Check a task queue that has registered workers (pollers):
```shell
wait4x temporal worker 127.0.0.1:7233 --namespace __YOUR_NAMESPACE__ --task-queue __YOUR_TASK_QUEUE__
```
This command waits until the Temporal server on `127.0.0.1:7233` is ready and the task queue `__YOUR_TASK_QUEUE__` has registered workers (pollers).

#Check a specific Temporal worker (pollers):
```shell
wait4x temporal worker 127.0.0.1:7233 --namespace __YOUR_NAMESPACE__ --task-queue __YOUR_TASK_QUEUE__ --expect-worker-identity-regex ".*@__HOSTNAME__@.*"
```
This command waits until the Temporal server on `127.0.0.1:7233` is ready and the task queue `__YOUR_TASK_QUEUE__` has a worker (poller) with an identity matching the provided regex pattern.

## Advanced Features

### Exponential Backoff

Enable exponential backoff retry:

```shell
wait4x http https://ifconfig.co --expect-status-code 200 --backoff-policy exponential --backoff-exponential-max-interval 120s --timeout 120s
```
This command retries the HTTP connection with exponential backoff until the status code `200` is returned or the timeout of `120s` is reached.

### Reverse Checking

Check for a free port:

```shell
wait4x tcp 127.0.0.1:9090 --reverse
```
This command waits until the TCP port `9090` on `127.0.0.1` is free.

### Parallel Checking

Check multiple services simultaneously:

```bash
wait4x tcp 127.0.0.1:9090 127.0.0.1:8080 127.0.0.1:9050
```
This command waits for the TCP ports `9090`, `8080` and `9050` on `127.0.0.1` to be available.

## Command Execution

You can execute a command after a successful wait. Use the `--` separator to separate the wait4x command from the command to execute.

Example:

```bash
wait4x tcp 127.0.0.1:9090 -- echo "Service is up!"
```
This command will echo "Service is up!" after the TCP port `9090` on `127.0.0.1` is available.

## Using Wait4X as an Importable Package

Besides being a standalone CLI tool, Wait4X can be used as an importable package in your Go applications. This allows you to integrate Wait4X's powerful waiting and service checking capabilities directly into your code.

### Adding Wait4X to Your Go Project

Add Wait4X as a dependency to your Go project:

```bash
go get wait4x.dev/v3
```

Then import the packages you need:

```go
import (
    "wait4x.dev/v3/checker/tcp"      // TCP checker
    "wait4x.dev/v3/checker/http"     // HTTP checker
    "wait4x.dev/v3/checker/redis"    // Redis checker
    // ...other checkers
    "wait4x.dev/v3/waiter"           // Waiter functionality
)
```

### Examples

Here are several examples of how to use Wait4X in your Go applications. Find the full examples in the [examples/pkg](examples/pkg) directory.

#### Basic TCP Check

```go
// Create a TCP checker for localhost:6379 with a 5-second connection timeout
tcpChecker := tcp.New("localhost:6379", tcp.WithTimeout(5*time.Second))

// Wait for the TCP port to be available
err := waiter.WaitContext(
    ctx,
    tcpChecker,
    waiter.WithTimeout(time.Minute),        // Total wait timeout
    waiter.WithInterval(2*time.Second),     // Time between retry attempts
    waiter.WithBackoffPolicy("exponential"), // Use exponential backoff
)
```

#### Advanced HTTP Check

```go
// Create custom HTTP headers
headers := http.Header{}
headers.Add("Authorization", "Bearer my-token")
headers.Add("Content-Type", "application/json")

// Create an HTTP checker with multiple validations
checker := httpChecker.New(
    "https://api.example.com/health",
    httpChecker.WithTimeout(5*time.Second),
    httpChecker.WithExpectStatusCode(200),
    httpChecker.WithExpectBodyJSON("status"),             // Check that 'status' field exists in JSON
    httpChecker.WithExpectBodyRegex(`"healthy":\s*true`), // Regex to check response
    httpChecker.WithExpectHeader("Content-Type=application/json"),
    httpChecker.WithRequestHeaders(headers),
)
```

#### Parallel Checking of Multiple Services

```go
// Create checkers for different services
checkers := []checker.Checker{
    // Redis checker
    redis.New(
        "redis://localhost:6379",
        redis.WithTimeout(5*time.Second),
    ),
    
    // PostgreSQL checker
    postgresql.New(
        "postgres://postgres:password@localhost:5432/app_db?sslmode=disable",
        postgresql.WithTimeout(5*time.Second),
    ),
    
    // HTTP API checker
    http.New(
        "http://localhost:8080/health",
        http.WithTimeout(3*time.Second),
        http.WithExpectStatusCode(200),
    ),
}

// Wait for all services in parallel
err := waiter.WaitParallelContext(ctx, checkers, waitOptions...)
```

#### Reverse Checking (Waiting for a Port to Become Free)

```go
// Create a TCP checker for the port
tcpChecker := tcp.New(port, tcp.WithTimeout(2*time.Second))

// Use invert check to wait until the TCP connection fails (port is free)
err := waiter.WaitContext(
    ctx,
    tcpChecker,
    waiter.WithTimeout(45*time.Second),
    waiter.WithInterval(3*time.Second),
    waiter.WithInvertCheck(true), // Key option to invert success condition
)
```

#### Creating Custom Checkers

You can create your own custom checkers by implementing the `checker.Checker` interface:

```go
// Implement the Checker interface
type MyCustomChecker struct {
    // Your fields here
}

// Identity returns a string identifying this checker
func (c *MyCustomChecker) Identity() (string, error) {
    return "my-custom-checker", nil
}

// Check performs the actual check
func (c *MyCustomChecker) Check(ctx context.Context) error {
    // Your checking logic here
    // Return nil if check passes, or an error if it fails
}
```

### Best Practices

1. Always use contexts with timeouts to prevent indefinite waiting
2. Consider using exponential backoff for services that might take a while to start
3. Use parallel checking when waiting for multiple independent services
4. Handle errors appropriately - distinguish between timeout errors and other errors
5. Add logging where appropriate to understand what's happening during waiting

### Core Components

1. **Checkers**: Implements the `checker.Checker` interface
2. **Waiter**: Provides waiting functionality with options like timeout, interval, backoff, etc.
3. **Context Usage**: All checkers and waiters support context for cancellation and timeouts

For more detailed examples, see the [examples/pkg](examples/pkg) directory in the repository. Each example is in its own directory with a complete, runnable `main.go` file.

## Contributing

### Reporting Issues

If you encounter any issues, please report them [here](https://github.com/atkrad/wait4x/issues).

### Submitting Pull Requests

1. Fork the repository
2. Create a new branch (`git checkout -b feature-branch`)
3. Make your changes
4. Commit your changes (`git commit -am 'Add new feature'`)
5. Push to the branch (`git push origin feature-branch`)
6. Create a new Pull Request

## License

This project is licensed under the Apache-2.0 license - see the [LICENSE](LICENSE) file for details.
```
Copyright 2019-2025 The Wait4X Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
