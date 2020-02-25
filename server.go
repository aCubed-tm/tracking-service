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

func (s server) GetAllObjects(_ context.Context, req *proto.GetAllObjectsRequest) (*proto.GetAllObjectsReply, error) {
	objects, err := getAllObjects()
	if err != nil {
		return nil, err
	}

	newObjects := make([]*proto.ObjectInfo, len(objects))
	for i, e := range objects {
		newObjects[i] = &proto.ObjectInfo{
			Uuid: e.Uuid,
			Name: e.Name,
			Note: e.Note,
			LastLocation: &proto.ObjectLocation{
				X:    e.LastLocation.X,
				Y:    e.LastLocation.Y,
				Z:    e.LastLocation.Z,
				Time: e.LastLocation.Time,
			},
		}
	}

	return &proto.GetAllObjectsReply{Objects: newObjects}, nil
}

func (s server) GetObject(_ context.Context, req *proto.GetObjectRequest) (*proto.GetObjectReply, error) {
	history, err := getObjectHistory(req.Uuid)
	if err != nil {
		return nil, err
	}

	newHistory := make([]*proto.ObjectLocation, len(history))
	for i, e := range history {
		newHistory[i] = &proto.ObjectLocation{
			X:    e.X,
			Y:    e.Y,
			Z:    e.Z,
			Time: e.Time,
		}
	}

	return &proto.GetObjectReply{Locations: newHistory}, nil
}
