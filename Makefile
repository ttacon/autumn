# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help toml2cli cli

default: cli ## Run the default action (see cli).

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


tomlcli: ## Rebuild the CLI interface.
	toml2cli -in-file=templates/autumn-cli.toml -out-file=cmd/main.go

cli: tomlcli ## Make the CLI binary for local testing.
	CGO_ENABLED=0 go build -o autumn ./cmd

binaries:
	env GOOS=linux  GOARCH=amd64 go build -ldflags="-X 'main.AutumnVersion=${GITHUB_REF/refs\/tags\//}'" -o osnp.linux.amd64
    env GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.AutumnVersion=${GITHUB_REF/refs\/tags\//}'" -o osnp.darwin.amd64
    env GOOS=darwin GOARCH=arm64 go build -ldflags="-X 'main.AutumnVersion=${GITHUB_REF/refs\/tags\//}'" -o osnp.darwin.arm64
