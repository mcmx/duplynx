.PHONY: lint test e2e perf tidy

lint:
	cd backend && golangci-lint run ./...

test:
	cd backend && go test ./...

e2e:
	npx playwright test

perf:
	cd backend && go test -run=^$ -bench=. ./...

tidy:
	cd backend && go mod tidy
