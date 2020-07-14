all:
	GOOS=linux GOARCH=amd64 go build -ldflags '-X main.key=$(key)' -o main main.go decrypt.go
	zip main.zip main
clean:
	rm kewei main
