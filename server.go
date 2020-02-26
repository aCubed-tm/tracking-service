package main

import (
	"context"
	proto "github.com/acubed-tm/tracking-service/protofiles"
	"log"
	"math"
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

func (s server) UpdatePositions(_ context.Context, req *proto.UpdatePositionsRequest) (*proto.UpdatePositionsReply, error) {
	var toHandle []string
	if req.Uuid != "" {
		toHandle = []string{req.Uuid}
	} else {
		objs, err := getAllTraceableObjects()
		if err != nil {
			return nil, err
		}
		toHandle = objs
	}

	for _, e := range toHandle {
		// TODO: should we only take records that are older than 5s?
		history, err := getObjectCaptures(e)
		if err != nil {
			return nil, err
		}

		// could group into camera, then for each capture in every group find a mate within 5s?
		// otherwise after group, find shortest intervals between 2 different cameras, then remove those from list?
		// need to take into account that last 5s shouldn't be used? maybe try use 3 povs somehow?
		groupedEntries := make(map[string][]CaptureWithId)
		for _, e := range history {
			groupedEntries[e.CameraUuid] = append(groupedEntries[e.CameraUuid], e)
		}

		// TODO: remove these pairs from db!!
		pairs := ExtractCapturePairs(groupedEntries)
		for _, pair := range pairs {
			a, b := pair.A, pair.B
			capA := MakeCaptureInfo(float64(a.X), float64(a.Y), CameraInfo{
				pos:   Vec3(a.CamPosX, a.CamPosY, a.CamPosZ),
				pitch: a.CamPitch,
				yaw:   a.CamYaw,
				resX:  a.CamResX,
				resY:  a.CamResY,
				fov:   a.Fov,
			})
			capB := MakeCaptureInfo(float64(b.X), float64(b.Y), CameraInfo{
				pos:   Vec3(b.CamPosX, b.CamPosY, b.CamPosZ),
				pitch: b.CamPitch,
				yaw:   b.CamYaw,
				resX:  b.CamResX,
				resY:  b.CamResY,
				fov:   b.Fov,
			})
			intersection := CalculateIntersection(capA, capB)
			err := insertTrackingPoint(e, (a.Time+b.Time)/2, intersection)
			if err != nil {
				// just log, don't exit
				log.Printf("Error during insertion of new tracking point: %v", err)
			} else {
				// if no error, remove old points
				log.Printf("Inserted new tracking point: %v (from prespectives %v and %v)", intersection, capA, capB)
				// TODO: ignoring errors for now
				_ = deleteNode(a.Id)
				_ = deleteNode(b.Id)
			}
		}
	}

	return &proto.UpdatePositionsReply{}, nil
}

type CapturePair struct {
	A, B CaptureWithId
}

func ExtractCapturePairs(groups map[string][]CaptureWithId) []CapturePair {
	pairs := make([]CapturePair, 0)

	// get all keys so we can index by int
	keys := make([]string, len(groups))
	i := 0
	for k := range groups {
		keys[i] = k
		i++
	}

	// find closest pair until we can find no more
	for {
		var lowestTimeDelta int32 = 2147483647
		firstCamIndex, secondCamIndex, firstGroupIndex, secondGroupIndex := -1, -1, -1, -1

		// loop over every group except for the very last
		for groupIndex, groupKey := range keys {
			group := groups[groupKey]

			// skip very last group
			if groupIndex == len(groups)-1 {
				continue
			}

			// loop over each capture
			for captureIndex, capture := range group {
				// for this capture, compare to every capture in all groups after this
				for otherGroupIndex := groupIndex + 1; otherGroupIndex < len(groups); otherGroupIndex++ {
					otherGroup := groups[keys[otherGroupIndex]]

					// compare to every capture in this group
					for otherCaptureIndex, otherCapture := range otherGroup {
						timeDelta := int32(math.Abs(float64(capture.Time - otherCapture.Time)))
						if timeDelta < lowestTimeDelta {
							lowestTimeDelta = timeDelta
							firstCamIndex, secondCamIndex = captureIndex, otherCaptureIndex
							firstGroupIndex, secondGroupIndex = groupIndex, otherGroupIndex
						}
					}
				}
			}
		}

		// if we didnt find anything, we can break out of the loop
		if firstCamIndex == -1 {
			break
		}

		// add new pair
		pairs = append(pairs, CapturePair{
			A: groups[keys[firstGroupIndex]][firstCamIndex],
			B: groups[keys[secondGroupIndex]][secondCamIndex],
		})

		// remove from collection
		groups[keys[firstGroupIndex]] = Slice(groups[keys[firstGroupIndex]], firstCamIndex)
		groups[keys[secondGroupIndex]] = Slice(groups[keys[secondGroupIndex]], secondCamIndex)
	}

	return pairs
}

func Slice(a []CaptureWithId, i int) []CaptureWithId {
	// taken from https://yourbasic.org/golang/delete-element-slice/
	copy(a[i:], a[i+1:])          // shift everything before index right
	a[len(a)-1] = CaptureWithId{} // likely not needed, but let's keep this for now
	return a[:len(a)-1]           // slice off first item
}
