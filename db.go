package main

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

const databaseUrl = "bolt://neo4j-public.default:7687"
const username = "neo4j"
const password = ""

var driver neo4j.Driver = nil

func newSession(accessMode neo4j.AccessMode) neo4j.Session {
	var (
		err     error
		session neo4j.Session
	)

	if driver == nil {
		driver, err = neo4j.NewDriver(databaseUrl, neo4j.BasicAuth(username, password, ""))
		if err != nil {
			panic(err.Error())
		}
	}

	session, err = driver.Session(accessMode)
	if err != nil {
		panic(err.Error())
	}

	return session
}

// should be able to return whether added or not :(
func insertCapture(objUuid, camUuid string, time int32, posX, posY float32) error {
	query := `MATCH (o:TracableObject{uuid:{objUuid}})
	          MATCH (c:Camera{uuid:{camUuid}})
	          CREATE (c)<-[:CAPTURED_BY]-(x:Capture{time:{time}, cameraX:{posX}, cameraY:{posY}})-[:CAPTURE_OF]->(o)`
	variables := map[string]interface{}{"objUuid": objUuid, "camUuid": camUuid, "time": time, "posX": posX, "posY": posY}
	return Write(query, variables)
}

type ObjectWithLastLocation struct {
	Uuid, Name, Note string
	LastLocation     Location
}

type Location struct {
	X, Y, Z float32
	Time    int32
}

func getAllObjects() ([]ObjectWithLastLocation, error) {
	query := `MATCH (p:TrackingPoint)-[:TRACKS]->(t:TracableObject) WITH p, t
	          ORDER BY p.time DESC WITH collect(p)[0] AS p, t
	          RETURN t.uuid, t.name, t.note, p.x, p.y, p.z, p.time`
	ret, err := Fetch(query, map[string]interface{}{}, func(result neo4j.Result) (interface{}, error) {
		var ret []ObjectWithLastLocation
		for result.Next() {
			rec := result.Record()
			ret = append(ret, ObjectWithLastLocation{
				Uuid: rec.GetByIndex(0).(string),
				Name: rec.GetByIndex(1).(string),
				Note: rec.GetByIndex(2).(string),
				LastLocation: Location{
					X:    rec.GetByIndex(3).(float32),
					Y:    rec.GetByIndex(4).(float32),
					Z:    rec.GetByIndex(5).(float32),
					Time: rec.GetByIndex(6).(int32),
				},
			})
		}

		return ret, nil
	})
	if err != nil {
		return nil, err
	}

	return ret.([]ObjectWithLastLocation), nil
}

func getObjectHistory(uuid string) ([]Location, error) {
	query := `MATCH (p:TrackingPoint)-[:TRACKS]->(:TracableObject{uuid:{uuid}}) WITH p
	          ORDER BY p.time DESC
	          RETURN p.x, p.y, p.z, p.time`
	params := map[string]interface{}{"uuid": uuid}
	ret, err := Fetch(query, params, func(result neo4j.Result) (interface{}, error) {
		var ret []Location
		for result.Next() {
			rec := result.Record()
			ret = append(ret, Location{
				X:    rec.GetByIndex(0).(float32),
				Y:    rec.GetByIndex(1).(float32),
				Z:    rec.GetByIndex(2).(float32),
				Time: rec.GetByIndex(3).(int32),
			})
		}

		return ret, nil
	})
	if err != nil {
		return nil, err
	}

	return ret.([]Location), nil
}

type CaptureWithId struct {
	CameraUuid                string
	X, Y                      float32
	Time                      int32
	Id                        int64
	CamResX, CamResY          int32
	Fov                       float32
	CamPosX, CamPosY, CamPosZ float32
	CamYaw, CamPitch          float32
}

func getObjectTraces(uuid string) ([]CaptureWithId, error) {
	query := `MATCH (p:TrackingPoint)-[:TRACKS]->(:TracableObject{uuid:{uuid}}) WITH p
	          ORDER BY p.time DESC
	          RETURN p.x, p.y, p.z, p.time, ID(p),
	              cam.resolutionX, cam.resolutionY, cam.fieldOfView
	              cam.locationX, cam.locationY, cam.locationZ, cam.yaw, cam.pitch`
	params := map[string]interface{}{"uuid": uuid}
	ret, err := Fetch(query, params, func(result neo4j.Result) (interface{}, error) {
		var ret []CaptureWithId
		for result.Next() {
			rec := result.Record()
			ret = append(ret, CaptureWithId{
				X:          rec.GetByIndex(0).(float32),
				Y:          rec.GetByIndex(1).(float32),
				Time:       rec.GetByIndex(2).(int32),
				CameraUuid: rec.GetByIndex(3).(string),
				Id:         rec.GetByIndex(4).(int64),
				CamResX:    rec.GetByIndex(5).(int32),
				CamResY:    rec.GetByIndex(6).(int32),
				Fov:        rec.GetByIndex(7).(float32),
				CamPosX:    rec.GetByIndex(8).(float32),
				CamPosY:    rec.GetByIndex(9).(float32),
				CamPosZ:    rec.GetByIndex(10).(float32),
				CamYaw:     rec.GetByIndex(11).(float32),
				CamPitch:   rec.GetByIndex(12).(float32),
			})
		}

		return ret, nil
	})
	if err != nil {
		return nil, err
	}

	return ret.([]CaptureWithId), nil
}

func getAllTraceableObjects() ([]string, error) {
	return FetchStringArray("MATCH (o:TracableObject) RETURN o.uuid", map[string]interface{}{})
}

func insertTrackingPoint(uuid string, time int32, v Vector3) error {
	query := `MATCH (o:TracableObject{uuid:{uuid}})
	          CREATE (:TrackingPoint {time:{time}, x:{x}, y:{y}, z:{z}})-[:TRACKS]->(o)`
	params := map[string]interface{}{"uuid": uuid, "time": time, "x": v.x, "y": v.y, "z": v.z}
	return Write(query, params)
}

func deleteNode(id int64) error {
	query := "MATCH (n) where ID(n)={id} DETACH DELETE n"
	params := map[string]interface{}{"id": id}
	return Write(query, params)
}
