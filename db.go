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
func insertCapture(objUuid, camUuid string, time int64, posX, posY float32) error {
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
	X, Y, Z float64
	Time    int64
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
					X:    rec.GetByIndex(3).(float64),
					Y:    rec.GetByIndex(4).(float64),
					Z:    rec.GetByIndex(5).(float64),
					Time: rec.GetByIndex(6).(int64),
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
				X:    rec.GetByIndex(0).(float64),
				Y:    rec.GetByIndex(1).(float64),
				Z:    rec.GetByIndex(2).(float64),
				Time: rec.GetByIndex(3).(int64),
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
	X, Y                      int64
	Time, Id                  int64
	CamResX, CamResY          int64
	Fov                       int64
	CamPosX, CamPosY, CamPosZ float64
	CamYaw, CamPitch          float64
}

func getObjectCaptures(uuid string) ([]CaptureWithId, error) {
	query := `MATCH (cam:Camera)<-[:CAPTURED_BY]-(p:Capture)-[:CAPTURE_OF]->(:TracableObject{uuid:{uuid}}) WITH p, cam
	          ORDER BY p.time DESC
	          RETURN p.cameraX, p.cameraY, p.time, ID(p), cam.uuid,
	              cam.resolutionX, cam.resolutionY, cam.fieldOfView,
	              cam.locationX, cam.locationY, cam.locationZ, cam.yaw, cam.pitch`
	params := map[string]interface{}{"uuid": uuid}
	ret, err := Fetch(query, params, func(result neo4j.Result) (interface{}, error) {
		var ret []CaptureWithId
		for result.Next() {
			rec := result.Record()
			ret = append(ret, CaptureWithId{
				X:          rec.GetByIndex(0).(int64),
				Y:          rec.GetByIndex(1).(int64),
				Time:       rec.GetByIndex(2).(int64),
				Id:         rec.GetByIndex(3).(int64),
				CameraUuid: rec.GetByIndex(4).(string),
				CamResX:    rec.GetByIndex(5).(int64),
				CamResY:    rec.GetByIndex(6).(int64),
				Fov:        rec.GetByIndex(7).(int64),
				CamPosX:    rec.GetByIndex(8).(float64),
				CamPosY:    rec.GetByIndex(9).(float64),
				CamPosZ:    rec.GetByIndex(10).(float64),
				CamYaw:     rec.GetByIndex(11).(float64),
				CamPitch:   rec.GetByIndex(12).(float64),
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

func insertTrackingPoint(uuid string, time int64, v Vector3) error {
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
