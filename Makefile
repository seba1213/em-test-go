.PHONY: swagger
swagger:
	go tool swag init -g main.go -o docs --parseInternal
	@echo "Regenerated docs/swagger.json and docs/docs.go"
