package main

import (
	"context"
	proto "github.com/acubed-tm/tracking-service/protofiles"
)

type server struct{}

func (s server) AddCapture(_ context.Context, req *proto.AddCaptureRequest) (*proto.AddCaptureReply, error) {
	panic("implement me")
}

