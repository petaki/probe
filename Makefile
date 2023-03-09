VERSION := $(if $(RELEASE_VERSION),$(RELEASE_VERSION),"master")

all: pre_clean darwin darwin_arm64 linux linux_arm64 windows post_clean

pre_clean:
	rm -rf dist
	mkdir dist
	cp .env.example dist/.env

darwin:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o dist/probe .
	cd dist && zip probe_$(VERSION)_darwin_amd64.zip .env probe
	rm -f dist/probe

darwin_arm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o dist/probe .
	cd dist && zip probe_$(VERSION)_darwin_arm64.zip .env probe
	rm -f dist/probe

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/probe .
	cd dist && zip probe_$(VERSION)_linux_amd64.zip .env probe
	rm -f dist/probe

linux_arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o dist/probe .
	cd dist && zip probe_$(VERSION)_linux_arm64.zip .env probe
	rm -f dist/probe

windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o dist/probe.exe .
	cd dist && zip probe_$(VERSION)_windows_amd64.zip .env probe.exe
	rm -f dist/probe.exe

post_clean:
	rm -rf dist/.env
