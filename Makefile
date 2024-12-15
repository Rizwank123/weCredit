setup:
	@echo "Generating env files..."
	cp sample.env .env
	cp sample.env test.env
	@echo ".env & test.env created. Now, update values in them"

# Migrate database
migrate:
	@echo "Running migrations..."
	sh script/migration.sh


# Generate API documentation
doc:
	@echo "Generating swagger docs..."
	swag fmt --exclude ./internal/domain
	swag init --parseDependency --parseInternal -g internal/http/api/we_credit_api.go -ot go,yaml -o internal/http/swagger

wire:
	cd internal/dependency/ && wire && cd ../..