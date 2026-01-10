tools:
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/swaggo/swag/cmd/swag@latest

check:
	govulncheck ./...
	gosec ./...

lint:
	golangci-lint -c=.golangci.yml run ./...

test:
	go test -cover -v ./...

swagger:
	swag fmt && swag init -g ./main.go -o ./docs --parseInternal=true --parseDependency=true