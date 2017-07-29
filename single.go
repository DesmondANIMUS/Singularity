package main

//--------------------------------------------------//
// Package Imports
import (
	"fmt"
	"log"
	"net/http"
	"time"

	"encoding/json"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//--------------------------------------------------//

//--------------------------------------------------//
// Data Models
type personData struct {
	Name      string `bson:"name"`
	UID       string `bson:"uid"`
	TimeStamp string `bson:"time"`
}

type fileData struct {
	FileName  string `bson:"filename"`
	FileDesc  string `bson:"filedesc"`
	MD5Hash   string `bson:"md5hash"`
	FileURL   string `bson:"fileURL"`
	TimeStamp string `bson:"time"`
}

type logData struct {
	Route     string `bson:"route"`
	TimeStamp string `bson:"time"`
	Result    string `bson:"result"`
}

//--------------------------------------------------//

//--------------------------------------------------//
// Constants
const (
	connectionString       = "mongodb://localhost/"
	databaseString         = "SingularityDB"
	fileCollectionString   = "masterFileCollection"
	personCollectionString = "personCollection"
	logCollectionString    = "logCollection"
)

//--------------------------------------------------//

//--------------------------------------------------//
func main() {
	//TODO: Write API's

	http.HandleFunc("/regLogUp", regLogUp)
	http.HandleFunc("/you", userProfile)
	http.HandleFunc("/checkMD5", checkMD5)
	http.HandleFunc("/file", file)

	// Try & write different paths for mobile apps &
	// web app, but keep them as optimized as possible
	fmt.Println("Server is listening at port 8888")
	http.ListenAndServe(":8888", nil)
}

//--------------------------------------------------//

//--------------------------------------------------//
// Handlers
func regLogUp(w http.ResponseWriter, r *http.Request) {
	var logInfo logData

	if r.Method == http.MethodPost {
		var user personData
		user.Name = r.FormValue("name")
		user.UID = r.FormValue("uid")
		user.TimeStamp = time.Now().Format(time.RFC850)

		err := checkIfRegistered(user.UID)
		if err == nil {
			up := checkAndUpdate(user)
			log.Println(up)
			log.Println("Log-in Success")

			logInfo.Result = "Log-in Success"
			logInfo.TimeStamp = time.Now().Format(time.RFC850)
			logDb(logInfo)

			fmt.Fprintf(w, `{"response":"Success"}`)

		} else {

			err = basicDataDb(user)

			if err != nil {
				log.Println("Failed :(")

				logInfo.Result = "Failed :("
				logInfo.TimeStamp = time.Now().Format(time.RFC850)
				logDb(logInfo)

				fmt.Fprintf(w, `{"response":"Failed :("}`)
			} else {
				log.Println("Sign-up Success")

				logInfo.Result = "Sign-up Success"
				logInfo.TimeStamp = time.Now().Format(time.RFC850)

				fmt.Fprintf(w, `{"response":"Success"}`)
			}
		}
	}

	logInfo.Route = r.URL.Path
	logInfo.TimeStamp = time.Now().Format(time.RFC850)
	log.Println(logDb(logInfo))
}
func userProfile(w http.ResponseWriter, r *http.Request) {
	var logInfo logData

	if r.Method == http.MethodPost {
		uid := r.FormValue("uid")

		response, err := getProfileInfo(uid)
		if err != nil {
			log.Println(err)
			logInfo.Result = err.Error()
			logInfo.Route = r.URL.Path
			logInfo.TimeStamp = time.Now().Format(time.RFC850)
			log.Println(logDb(logInfo))
		} else {
			logInfo.TimeStamp = time.Now().Format(time.RFC850)
			logInfo.Result = "Success"
			logInfo.Route = r.URL.Path
			log.Println("Success")
			log.Println(logDb(logInfo))

			fmt.Fprintf(w, string(response))
		}
	}

	log.Println(r.URL.Path)
}
func checkMD5(w http.ResponseWriter, r *http.Request) {
	var logger logData

	if r.Method == http.MethodPost {
		hash := r.FormValue("hash")
		uid := r.FormValue("uid")

		err := testUserCol(hash, uid)
		if err != nil {
			// Data not found in UserCollection
			// Test in Master Collection

			err = testMasterCol(hash)
			if err != nil {
				// Data NOT found in MasterCollection

				logger.Route = r.URL.Path
				logger.Result = err.Error()
				logger.TimeStamp = time.Now().Format(time.RFC850)
				log.Println(logDb(logger))

				fmt.Fprintf(w, `{"response":"500"}`)
			} else {
				// Data found in MasterCollection
				// SET data in UserCollection
			}
		} else {
			// Data found in UserCollection
			logger.Route = r.URL.Path
			logger.Result = "File Found"
			logger.TimeStamp = time.Now().Format(time.RFC850)
			log.Println(logDb(logger))

			fmt.Fprintf(w, `{"response":"200"}`)
		}
	}

	log.Println(logDb(logger))
}
func file(w http.ResponseWriter, r *http.Request) {}

//--------------------------------------------------//

//--------------------------------------------------//
// Helpers
func logDb(logInfo logData) string {
	session, err := mgo.Dial(connectionString)

	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(databaseString).C(logCollectionString)
	err = c.Insert(logInfo)

	return logInfo.Route
}
func uploadFile(w http.ResponseWriter, r *http.Request) {}

//--------------------------------------------------//

//--------------------------------------------------//
// Basic Database Related Functions

func basicDataDb(udata personData) error {
	session, err := mgo.Dial(connectionString)

	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(databaseString).C(personCollectionString)
	err = c.Insert(udata)

	return err
}
func getProfileInfo(uid string) ([]byte, error) {
	session, err := mgo.Dial(connectionString)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(databaseString).C(personCollectionString)
	result := personData{}

	err = c.Find(bson.M{"uid": uid}).One(&result)
	if err != nil {
		return nil, err
	}

	response, err := json.MarshalIndent(result, "", " ")

	return response, err
}
func fileUserData(fileInfo fileData, uid string) error {
	return nil
}
func fileMasterData(fileInfo fileData) error {
	return nil
}

//--------------------------------------------------//

//--------------------------------------------------//
// Check Database Related Functions
func checkIfRegistered(uid string) error {
	session, err := mgo.Dial(connectionString)

	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	result := personData{}

	c := session.DB(databaseString).C(personCollectionString)
	err = c.Find(bson.M{"uid": uid}).One(&result)

	return err
}
func checkAndUpdate(udata personData) string {
	session, err := mgo.Dial(connectionString)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	result := personData{}

	c := session.DB(databaseString).C(personCollectionString)
	err = c.Find(bson.M{"uid": udata.UID, "name": udata.Name}).One(&result)
	if err != nil {
		colQuerier := bson.M{"uid": udata.UID}
		err = c.Update(colQuerier, udata)

		return "Profile was updated"
	}

	return "No updates"
}
func testUserCol(hash, uid string) error {
	return nil
}
func testMasterCol(hash string) error {
	return nil
}

//--------------------------------------------------//

//--------------------------------------------------//
// Other Helper Functions
//--------------------------------------------------//
