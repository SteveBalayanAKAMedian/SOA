package server

import (
    "database/sql"
    "log"
    "net"

    "google.golang.org/grpc"
    "task-service/proto"
)

func RunGRPCServer(db *sql.DB) {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    proto.RegisterTaskServiceServer(s, NewServer(db))
    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
