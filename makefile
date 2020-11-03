mocksDestination="_mocks/generated"

mocks:
	@echo "Generating mocks..."
	@mockgen -source=discord/types.go -package=mock_sendmsg -destination=$(mocksDestination)/discord/mock_sendmsg.go