VERSION := $(if $(RELEASE_VERSION),$(RELEASE_VERSION),"master")

all: clean darwin linux windows

clean:
	rm -rf dist

darwin:
	GOOS=darwin GOARCH=amd64 go build -o dist/probe .
	cd dist && zip probe_$(VERSION)_darwin_amd64.zip probe
	rm -f dist/probe

linux:
	GOOS=linux GOARCH=amd64 go build -o dist/probe .
	cd dist && zip probe_$(VERSION)_linux_amd64.zip probe
	rm -f dist/probe

windows:
	GOOS=windows GOARCH=amd64 go build -o dist/probe.exe .
	cd dist && zip probe_$(VERSION)_windows_amd64.zip probe.exe
	rm -f dist/probe.exe
