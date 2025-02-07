build:
	go build -o bin/$(shell basename $(PWD)) ./main.go
