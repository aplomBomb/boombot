mocksDestination="_mocks/generated"

test: mocks runtests

mocks:
	@echo "Generating mocks..."
	@mockgen -source=discord/types.go -package=mock_sendmsg -destination=$(mocksDestination)/discord/sendmsg/mock_sendmsg.go
	@mockgen -source=discord/types.go -package=mock_newclient -destination=$(mocksDestination)/discord/newclient/mock_newclient.go

runtests:
	@echo "Running tests..."
	@go test ./...