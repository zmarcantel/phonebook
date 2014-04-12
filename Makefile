PROGNAME=phonebook

all: deps
	go build -o bin/$(PROGNAME)

deps:
	go list -f "{{ range .Deps }}{{ . }} {{ end }}" ./ | tr ' ' '\n' | awk '!/^.\//' | xargs go get

run: all
	sudo bin/phonebook
