build:
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go
	zip links-app.zip bootstrap

test:
	go test -race -coverprofile=coverage.txt ./...