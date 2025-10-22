linters:
	docker run --rm \
		-v $$(pwd):/client \
		-v $$(pwd)/golangci-lint/.cache/golangci-lint/v2.4.0:/root/.cache \
		-w /client \
		golangci/golangci-lint:v2.4.0 \
		golangci-lint run ./... -c .golangci.yml