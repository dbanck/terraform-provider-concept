default: build

build:
	go build -v .

install: build
	go install -v .

test:
	go test -count=1 ./...

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir .

.PHONY: build install test docs
