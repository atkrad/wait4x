BUILD_OUTPUT ?= bin/wait4x
WAIT4X_FLAGS ?=
COMMIT_REF_SLUG = $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
COMMIT_SHORT_SHA = $(shell git rev-parse --verify --short=8 HEAD)
COMMIT_DATETIME = $(shell git log -1 --format="%at" | TZ=UTC xargs -I{} date -d @{} '+%FT%TZ')

help:
	@echo " __      __        .__  __     _________  ___"
	@echo "/  \    /  \_____  |__|/  |_  /  |  \   \/  /"
	@echo "\   \/\/   /\__  \ |  \   __\/   |  |\     / "
	@echo " \        /  / __ \|  ||  | /    ^   /     \ "
	@echo "  \__/\  /  (____  /__||__| \____   /___/\  \\"
	@echo "       \/        \/              |__|     \_/"
	@echo ""
	@echo ""
	@echo "build"
	@echo "  Build Wait4X."
	@echo ""
	@echo "run"
	@echo '  Run Wait4X.'
	@echo '  You can pass subcommand and arguements with "Wait4X" e.g. "make run WAIT4X_FLAGS.'
	@echo ""

build:
	go build -v \
	-ldflags "-X github.com/atkrad/wait4x/cmd.AppVersion=$(COMMIT_REF_SLUG) -X github.com/atkrad/wait4x/cmd.GitCommit=$(COMMIT_SHORT_SHA) -X github.com/atkrad/wait4x/cmd.BuildTime=$(COMMIT_DATETIME)" \
	-o $(BUILD_OUTPUT)

run:
	go run \
	-ldflags "-X github.com/atkrad/wait4x/cmd.AppVersion=$(COMMIT_REF_SLUG) -X github.com/atkrad/wait4x/cmd.GitCommit=$(COMMIT_SHORT_SHA) -X github.com/atkrad/wait4x/cmd.BuildTime=$(COMMIT_DATETIME)" \
	main.go $(WAIT4X_FLAGS)
