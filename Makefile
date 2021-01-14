fmt:
	go fmt github.com/loilo-inc/logos/...
test:
	go test -race -cover -coverprofile=coverage.out -covermode=atomic \
		github.com/loilo-inc/logos/... -count 1
gen:
	go run set/gen/main.go
