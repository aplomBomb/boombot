mocksDestination="_mocks/generated"

test: mocks runtests

mocks:
	@echo "Generating mocks..."
	@mockgen -destination=$(mocksDestination)/discord/mock_client.go -package=mock_disgordclient github.com/aplombomb/boombot/discord/ifaces DisgordClientAPI
	@mockgen -destination=$(mocksDestination)/discord/mock_session.go -package=mock_disgordsession github.com/aplombomb/boombot/discord/ifaces DisgordSessionAPI
	@mockgen -destination=$(mocksDestination)/youtube/mock_searchservice.go -package=mock_youtubesearchservice github.com/aplombomb/boombot/youtube/ifaces YoutubeSearchServiceAPI


runtests:
	@echo "Running tests..."
	@go test ./...