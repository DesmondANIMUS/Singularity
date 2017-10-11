package main

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func checkIfRegistered(uid string, session *mgo.Session) error {
	result := personData{}

	c := session.DB(databaseString).C(personCollectionString)
	err := c.Find(bson.M{checkUID: uid}).One(&result)

	return err
}
func checkAndUpdate(udata personData, session *mgo.Session) string {
	result := personData{}

	c := session.DB(databaseString).C(personCollectionString)
	err := c.Find(bson.M{checkUID: udata.UID, checkName: udata.Name}).One(&result)
	if err != nil {
		colQuerier := bson.M{checkUID: udata.UID}
		err = c.Update(colQuerier, udata)

		return "Profile was updated"
	}

	return "No updates"
}

func testUserCol(hash, uid string, session *mgo.Session) error {
	userKol := userFileCollectionString + uid
	c := session.DB(databaseString).C(userKol)

	_, err := fileFinder(hash, c)

	return err
}
func testMasterCol(hash string, session *mgo.Session) (fileData, error) {
	c := session.DB(databaseString).C(fileCollectionString)
	result, err := fileFinder(hash, c)

	return result, err
}
