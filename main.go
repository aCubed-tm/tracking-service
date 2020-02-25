package main

import (
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/acubed-tm/tracking-service/protofiles"
)

type CameraInfo struct {
	pos        Vector3
	resX, resY int16
	pitch, yaw float32
	fov        float32
}

type CaptureInfo struct {
	origin    Vector3
	direction Vector3
}

const port = ":50551"

func main() {
	log.Print("Starting server")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTrackingServiceServer(s, &server{})
	for {
		if err := s.Serve(lis); err != nil {
			log.Printf("Failed to serve with error: %v", err)
		}
	}
}
