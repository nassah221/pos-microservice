.PHONY: proto

proto:
	protoc --go_out=./contract --go_opt=module=/cashier \
	--go-grpc_out=./contract --go-grpc_opt=module=/cashier \
	contract/cashier.proto

test:
	go test -v --cover pos-microservices/cashier...