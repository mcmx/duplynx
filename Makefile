.PHONY: lint test e2e perf tidy ci

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

.PHONY: ci
ci:
	@set -eu; \
	start=$$(date +%s); \
	$(MAKE) lint; \
	$(MAKE) test; \
	$(MAKE) perf; \
	end=$$(date +%s); \
	duration=$$((end - start)); \
	echo "CI suite completed in $$duration seconds"; \
	if [ $$duration -gt 480 ]; then \
		echo "CI suite exceeded 480 seconds (8 minute budget)"; \
		exit 1; \
	fi
