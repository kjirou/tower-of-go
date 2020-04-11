run:
	go run main.go

run-with-debug-mode:
	go run main.go --debug-mode

test-all:
	go test -v ./...
