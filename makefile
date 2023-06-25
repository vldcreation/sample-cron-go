APP_ENV=prod
STORAGE_BUCKET=dci-auth-revamp
run:
	@go run main.go ${APP_ENV} ${STORAGE_BUCKET}
install:
	@echo "Installing depdencies"
	@go mod tidy && go mod download
	
.PHONY: run, install