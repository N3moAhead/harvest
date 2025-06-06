.DEFAULT_GOAL := start

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet: fmt
	go vet ./...

.PHONY: test
test:
	go test ./...

.PHONY: build
build: vet
	go build -o harvest-game ./cmd/harvest/main.go

.PHONY: start
start: build
	./harvest-game

commit: test
	git commit


.PHONY: dev
dev:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "./harvest-game" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go,tpl,tmpl,html,css,scss,js,ts,sql,jpeg,jpg,gif,png,bmp,svg,webp,ico" \
		--misc.clean_on_exit "true"

.PHONY: clean
clean:
	rm -rf ./harvest-game
