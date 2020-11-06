mocksDestination="_mocks/generated"

test: mocks runtests

mocks:
	@echo "Generating mocks..."
	@mockgen -destination=$(mocksDestination)/discord/mock_client.go -package=mock_disgordclient github.com/aplombomb/boombot/discord/ifaces DisgordClientAPI
	@mockgen -destination=$(mocksDestination)/youtube/mock_client.go -package=mock_youtubeclient github.com/aplombomb/boombot/youtube/ifaces YoutubeClientAPI


runtests:
	@echo "Running tests..."
	@go test ./...