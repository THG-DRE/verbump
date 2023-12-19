ARCHITECTURES=arm64 amd64
PLATFORMS=darwin windows linux
VERSION=`cat version`

default: build

bump-and-build: clean
	hack/bump-and-build.sh

.PHONY: build
build:
	$(foreach GOOS, $(PLATFORMS), \
		$(foreach GOARCH, $(ARCHITECTURES), \
			$(shell \
				export GOOS=$(GOOS); \
				export GOARCH=$(GOARCH); \
				go build -v -o ./bin/verbump-$(GOOS)-$(GOARCH)-$(VERSION)/verbump \
					-ldflags "-X verbump/cmd.Version=$(VERSION)" \
					.; \
				tar -czf ./bin/verbump-$(GOOS)-$(GOARCH)-$(VERSION)/verbump-$(GOOS)-$(GOARCH)-$(VERSION).tar.gz -C ./bin/verbump-$(GOOS)-$(GOARCH)-$(VERSION) verbump; \
			)\
		)\
	)

clean:
	rm -rf bin/