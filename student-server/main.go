package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"grpc/initial/database"
	"grpc/initial/server"
	"grpc/initial/studentpb"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	errorVar := godotenv.Load()
	if errorVar != nil {
		panic(errorVar)
	}

	list, err := net.Listen("tcp", ":5060")

	if err != nil {
		log.Fatalf("Error listening: %s", err.Error())
	}

	connStr := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SERVER"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	repo, err := database.NewPostgresRepository(connStr)

	server := server.NewStudentServer(repo)

	if err != nil {
		log.Fatalf("Error creating repository: %s", err.Error())
	}

	s := grpc.NewServer()
	studentpb.RegisterStudentServiceServer(s, server)

	reflection.Register(s)

	if err := s.Serve(list); err != nil {
		log.Fatalf("Error serving: %s", err.Error())
	}
}
