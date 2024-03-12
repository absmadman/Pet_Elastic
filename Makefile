all: build

build:
	go build app/main.go

clean:
	rm -rf main

rebuild:
	rm -rf main
	go build app/main.go
