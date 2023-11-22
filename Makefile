DBCONN = postgresql://urls:urls@localhost:5432/urls?sslmode=disable
PORT = 8080
JAEGER_GRPC_ADDR = replaceme #-jaeger-grpc-address=$(JAEGER_GRPC_ADDR)

migrate:
	bin/migrate/migrate -path=internal/adapters/migrations -database=$(DBCONN) up

start:
	go run ./cmd/api/main.go -sqlconn=$(DBCONN) -port=$(PORT)

testreport:
	go test ./... -coverpkg=./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html