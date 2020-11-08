mocksDestination="_mocks/generated"

test: mocks runtests

mocks:
	@echo "Generating mocks..."
	@mockgen -destination=$(mocksDestination)/discordclient/mock_client.go -package=mock_disgordclient github.com/aplombomb/boombot/discord/ifaces DisgordClientAPI
	@mockgen -destination=$(mocksDestination)/discorduser/mock_user.go -package=mock_disgorduser github.com/aplombomb/boombot/discord/ifaces DisgordUserAPI
	@mockgen -destination=$(mocksDestination)/youtube/mock_searchservice.go -package=mock_youtubesearchservice github.com/aplombomb/boombot/youtube/ifaces YoutubeSearchServiceAPI


runtests:
	@echo "Running tests..."
	@go test ./...