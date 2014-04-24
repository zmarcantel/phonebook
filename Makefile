PROGNAME=phonebook

all: deps
	go build -o bin/$(PROGNAME)

deps:
	go list -f "{{ range .Deps }}{{ . }} {{ end }}" ./ | tr ' ' '\n' | awk '!/^.\//' | xargs go get

test:
	go test ./dns/record

run: all
	sudo bin/phonebook

.PHONY: test
