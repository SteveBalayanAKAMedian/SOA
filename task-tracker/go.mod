module tracker

go 1.22.4

require (
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.24.0
)

require (
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.64.0
	task-service v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace task-service/proto => ./task-service/proto

replace task-service => ./task-service
