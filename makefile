mocksDestination="_mocks/generated"

mocks:
	@echo "Generating mocks..."
	@mockgen -source=discord/types.go -package=mock_createmessage -destination=$(mocksDestination)/discord/mock_createmessage.go