def main(ctx):
  return [
    {
      "kind": "pipeline",
      "type": "docker",
      "name": "check",
      "steps": [
        {
          "name": "gofmt",
          "image": "golang:1.13-buster",
          "commands": [
            "make check-gofmt"
          ]
        },
        {
          "name": "revive",
          "image": "golang:1.13-buster",
          "commands": [
            "go get -u github.com/mgechev/revive",
            "make check-revive"
          ]
        }
      ]
    },
    {
      "kind": "pipeline",
      "type": "docker",
      "name": "test",
      "steps": [
        {
          "name": "test",
          "image": "golang:1.13-buster",
          "commands": [
            "make test"
          ]
        },
        {
          "name": "coverage",
          "image": "golang:1.13-buster",
          "environment": {
            "COVERALLS_TOKEN": {
              "from_secret": "COVERALLS_TOKEN"
            }
          },
          "commands": [
            "go get github.com/mattn/goveralls",
            "goveralls -coverprofile=coverage.out -service=drone -repotoken $${COVERALLS_TOKEN}"
          ],
          "when": {
            "event": {
              "exclude": [
                "pull_request"
              ]
            }
          }
        }
      ],
      "depends_on": [
        "check"
      ]
    },
    {
      "kind": "pipeline",
      "type": "docker",
      "name": "docker",
      "steps": [
        {
          "name": "build",
          "image": "golang:1.13-alpine",
          "commands": [
            "apk --update add build-base git",
            "make build"
          ],
          "environment": {
            "GOOS": "linux",
            "GOARCH": "amd64",
            "WAIT4X_BUILD_OUTPUT": ".",
            "WAIT4X_BINARY_NAME": "wait4x"
          }
        },
        {
          "name": "docker",
          "image": "plugins/docker",
          "settings": {
            "repo": "atkrad/wait4x",
            "auto_tag": "true",
            "username": {
              "from_secret": "DOCKER_USERNAME"
            },
            "password": {
              "from_secret": "DOCKER_PASSWORD"
            }
          },
          "when": {
            "event": {
              "exclude": [
                "pull_request"
              ]
            }
          },
          "depends_on": [
            "build"
          ]
        },
        {
          "name": "dockerhub-description",
          "image": "peterevans/dockerhub-description:2.0.0",
          "environment": {
            "DOCKERHUB_REPOSITORY": "atkrad/wait4x",
            "DOCKERHUB_USERNAME": {
              "from_secret": "DOCKER_USERNAME"
            },
            "DOCKERHUB_PASSWORD": {
              "from_secret": "DOCKER_PASSWORD"
            }
          },
          "when": {
            "event": {
              "exclude": [
                "pull_request"
              ]
            }
          },
          "depends_on": [
            "docker"
          ]
        }
      ],
      "depends_on": [
        "test"
      ]
    },
    build_pipeline("linux", "amd64"),
    build_pipeline("windows", "amd64"),
    build_pipeline("darwin", "amd64"),
  ]

def build_pipeline(os, arch):
  return {
    "kind": "pipeline",
    "type": "docker",
    "name": "build-%s-%s" % (os, arch),
    "steps": [
      {
        "name": "fetch",
        "image": "alpine/git",
        "commands": [
          "git fetch --tags"
        ],
        "when": {
          "event": [
            "tag"
          ]
        }
      },
      {
        "name": "build",
        "image": "golang:1.13-buster",
        "commands": [
          "make build"
        ],
        "environment": {
          "GOOS": os,
          "GOARCH": arch,
          "WAIT4X_BUILD_OUTPUT": ".",
          "WAIT4X_BINARY_NAME": "wait4x-%s-%s" % (os, arch)
        },
        "depends_on": [
          "fetch"
        ]
      },
      {
        "name": "pre-release",
        "image": "plugins/github-release",
        "settings": {
          "api_key": {
            "from_secret": "github_token"
          },
          "prerelease": "true",
          "files": [
            "wait4x-%s-%s" % (os, arch)
          ],
          "checksum": [
            "sha256"
          ],
          "checksum_file": "wait4x-%s-%s.CHECKSUMsum" % (os, arch)
        },
        "when": {
          "event": [
            "tag"
          ],
          "ref": {
            "include": [
              "refs/tags/*rc*",
              "refs/tags/*alpha*",
              "refs/tags/*beta*"
            ]
          }
        },
        "depends_on": [
          "build"
        ]
      },
      {
        "name": "release",
        "image": "plugins/github-release",
        "settings": {
          "api_key": {
            "from_secret": "github_token"
          },
          "prerelease": "false",
          "files": [
            "wait4x-%s-%s" % (os, arch)
          ],
          "checksum": [
            "sha256"
          ],
          "checksum_file": "wait4x-%s-%s.CHECKSUMsum" % (os, arch)
        },
        "when": {
          "event": [
            "tag"
          ],
          "ref": {
            "include": [
              "refs/tags/*"
            ],
            "exclude": [
              "refs/tags/*rc*",
              "refs/tags/*alpha*",
              "refs/tags/*beta*"
            ]
          }
        },
        "depends_on": [
          "build"
        ]
      }
    ],
    "depends_on": [
      "test"
    ]
  }
