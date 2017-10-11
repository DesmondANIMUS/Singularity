package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"
)

func regLogUp(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(connectionString)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	if r.Method == http.MethodPost {
		var user personData
		user.Name = r.FormValue("name")
		user.UID = r.FormValue("uid")
		user.TimeStamp = time.Now().Format(time.RFC850)

		err := checkIfRegistered(user.UID, session)
		if err == nil {
			up := checkAndUpdate(user, session)
			log.Println(up)
			log.Println("Log-in Success")
			fmt.Fprintf(w, `{"response":"200"}`)
		} else {

			err = basicDataDb(user, session)

			if err != nil {
				log.Println("Failed :(")
				fmt.Fprintf(w, `{"response":"500"}`)
			} else {
				log.Println("Sign-up Success")
				fmt.Fprintf(w, `{"response":"200"}`)
			}
		}
	}

	log.Println(r.URL.Path)
}
func userProfile(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(connectionString)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	if r.Method == http.MethodPost {
		uid := r.FormValue("uid")

		response, err := getProfileInfo(uid, session)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, `{"response":"500"}`)
		} else {
			fmt.Fprintf(w, string(response))
		}
	}

	log.Println(r.URL.Path)
}
func checkMD5(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(connectionString)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	if r.Method == http.MethodPost {
		hash := r.FormValue("hash")
		uid := r.FormValue("uid")

		err := testUserCol(hash, uid, session)
		if err != nil {
			log.Println("not found in user")
			result, err := testMasterCol(hash, session)
			if err != nil {
				// Data NOT found in MasterCollection
				log.Println("not found in master")
				fmt.Fprintf(w, `{"response":"500"}`)
			} else {
				// Data found in MasterCollection
				// SET data in UserCollection
				if err := fileUserData(result, uid, session); err != nil {
					log.Println(err.Error())
					fmt.Fprintf(w, `{"response":"500"}`)
				}

				log.Println("found in master, populating user")
				fmt.Fprintf(w, `{"response":"200"}`)
			}
		} else {
			log.Println("found in user")
			fmt.Fprintf(w, `{"response":"200"}`)
		}
	}

	log.Println(r.URL.Path)
}
func listFiles(w http.ResponseWriter, r *http.Request)  {}
func uploadFile(w http.ResponseWriter, r *http.Request) {}
