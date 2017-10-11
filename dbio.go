package main

import (
	"encoding/json"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func basicDataDb(udata personData, session *mgo.Session) error {
	c := session.DB(databaseString).C(personCollectionString)
	err := c.Insert(udata)

	return err
}
func getProfileInfo(uid string, session *mgo.Session) ([]byte, error) {
	c := session.DB(databaseString).C(personCollectionString)
	result := personData{}

	err := c.Find(bson.M{"uid": uid}).One(&result)
	if err != nil {
		return nil, err
	}

	response, err := json.MarshalIndent(result, "", " ")

	return response, err
}
func fileUserData(fileInfo fileData, uid string, session *mgo.Session) error {
	userKol := userFileCollectionString + uid
	c := session.DB(databaseString).C(userKol)
	err := c.Insert(fileInfo)

	return err
}
func fileMasterData(fileInfo fileData, session *mgo.Session) error {
	c := session.DB(databaseString).C(fileCollectionString)
	err := c.Insert(fileInfo)

	return err
}
