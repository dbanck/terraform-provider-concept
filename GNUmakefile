default: build

build:
	go build -v .

install: build
	go install -v .

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir .

.PHONY: build install docs
