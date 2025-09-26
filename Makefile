docs:
	swag init -g ./cmd/main.go -o ./docs --parseDependency --parseInternal
.PHONY: docs