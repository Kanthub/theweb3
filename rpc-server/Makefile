rpc-server:
	go mod tidy
	env GO111MODULE=on go build -v $(LDFLAGS) ./cmd/rpc-server

proto:
	sh ./script/compile.sh


.PHONY: proto rpc-server


