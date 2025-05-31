.PHONY: proto dev

proto:
	mkdir -p pb/userpb
	protoc --proto_path=proto proto/*.proto \
		--go_out=pb/userpb --go_opt=paths=source_relative \
		--go-grpc_out=pb/userpb --go-grpc_opt=paths=source_relative


dev:
	docker compose build && docker compose up

