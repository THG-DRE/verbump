CLI_ARCHITECTURES=arm64 amd64
CLI_PLATFORMS=darwin windows linux
CLI_VERSION=0.0.2

default: build

.PHONY: build
build:
	$(foreach GOOS, $(CLI_PLATFORMS), \
		$(foreach GOARCH, $(CLI_ARCHITECTURES), \
			$(shell \
				export GOOS=$(GOOS); \
				export GOARCH=$(GOARCH); \
				go build -v -o ./bin/verbump-$(GOOS)-$(GOARCH)-$(CLI_VERSION)/verbump \
					-ldflags "-X verbump/cmd.Version=$(CLI_VERSION)" \
					.; \
				tar -czf ./bin/verbump-$(GOOS)-$(GOARCH)-$(CLI_VERSION)/verbump-$(GOOS)-$(GOARCH)-$(CLI_VERSION).tar.gz -C ./bin/verbump-$(GOOS)-$(GOARCH)-$(CLI_VERSION) verbump; \
			)\
		)\
	)

clean:
	rm -rf bin/