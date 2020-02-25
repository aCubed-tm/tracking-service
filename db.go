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
func insertCapture(objUuid, camUuid, time string, posX, posY float32) error {
	query := `MATCH (o:TracableObject{uuid:{objUuid}})
	          MATCH (c:Camera{uuid:{camUuid}})
	          CREATE (c)<-[:CAPTURED_BY]-(x:Capture{time:{time}, cameraX:{posX}, cameraY:{posY}})-[:CAPTURE_OF]->(o)`
	variables := map[string]interface{}{"objUuid": objUuid, "camUuid": camUuid, "time": time, "posX": posX, "posY": posY}
	return Write(query, variables)
}
