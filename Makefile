lint:
	golangci-lint run --enable-all --disable gochecknoglobals,funlen,gomnd,gocognit ./...

test:
	go test ./... -race -cover
