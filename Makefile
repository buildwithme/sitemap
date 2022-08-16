tidy:
	go mod tidy
build:
	cd cmd && go build sitemap.go
test: build
	go test ./...
debug: build
	go run -race cmd/sitemap.go https://google.com/
run: build
	./cmd/sitemap https://google.com/