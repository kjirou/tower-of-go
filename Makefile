get:
	go get -v -t -d ./...

run:
	go run main.go

run-with-debug-mode:
	go run main.go -debug

test:
	go test -v ./...
