mocksDestination="_mocks/generated"

test: mocks runtests

mocks:
	@echo "Generating mocks..."
	@mockgen -destination=$(mocksDestination)/discord/mock_client.go -package=mock_client github.com/aplombomb/boombot/discord/ifaces DisgordClientAPI

runtests:
	@echo "Running tests..."
	@go test ./...