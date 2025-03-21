package main

import "google.golang.org/grpc"

func main() {
	conn, err := grpc.NewClient("0.0.0.0:8000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	
}
