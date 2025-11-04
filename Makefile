.PHONY: lint test e2e perf tidy smoke-demo ci

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

.PHONY: smoke-demo
smoke-demo:
	@set -eu; \
	if [ ! -f go.work ]; then \
		go work init ./backend ./tests; \
	else \
		go work use ./backend ./tests; \
	fi; \
	go work sync; \
	start=$$(date +%s); \
	go test ./tests/smoke -count=1; \
	end=$$(date +%s); \
	duration=$$((end - start)); \
	echo "Smoke demo verified in $$duration seconds"; \
	if [ $$duration -gt 300 ]; then \
		echo "Smoke demo exceeded 300 second budget"; \
		exit 1; \
	fi

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
