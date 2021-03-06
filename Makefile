VERSION := $(if $(RELEASE_VERSION),$(RELEASE_VERSION),"master")

all: pre_clean dotenv darwin linux windows post_clean

pre_clean:
	rm -rf dist

dotenv:
	mkdir dist
	cp .env.example dist/.env

darwin:
	GOOS=darwin GOARCH=amd64 go build -o dist/probe .
	cd dist && zip probe_$(VERSION)_darwin_amd64.zip .env probe
	rm -f dist/probe

linux:
	GOOS=linux GOARCH=amd64 go build -o dist/probe .
	cd dist && zip probe_$(VERSION)_linux_amd64.zip .env probe
	rm -f dist/probe

windows:
	GOOS=windows GOARCH=amd64 go build -o dist/probe.exe .
	cd dist && zip probe_$(VERSION)_windows_amd64.zip .env probe.exe
	rm -f dist/probe.exe

post_clean:
	rm -rf dist/.env
