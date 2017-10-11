package main

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func fileFinder(hash string, c *mgo.Collection) (fileData, error) {
	result := fileData{}
	if err := c.Find(bson.M{checkHash: hash}).One(&result); err != nil {
		return result, err
	}

	return result, nil
}
