package main

import (
	"context"
	proto "github.com/acubed-tm/tracking-service/protofiles"
)

type server struct{}

func (s server) AddCapture(_ context.Context, req *proto.AddCaptureRequest) (*proto.AddCaptureReply, error) {
	err := insertCapture(req.ObjectUuid, req.CameraUuid, req.Time, req.CaptureX, req.CaptureY)
	if err != nil {
		return nil, err
	}
	return &proto.AddCaptureReply{}, nil
}
