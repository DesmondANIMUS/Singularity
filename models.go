package main

type fileData struct {
	FileName  string `bson:"file_name"` //we'll generate random name for file
	FileDesc  string `bson:"file_desc"` // file name we get will go in file_desc
	MD5Hash   string `bson:"md5_hash"`  // we'll get it
	FileURL   string `bson:"file_URL"`  // the location/address of file on server
	TimeStamp string `bson:"time"`      // we'll generate the timestamp
}

type personData struct {
	Name      string `bson:"name"`
	UID       string `bson:"uid"`
	TimeStamp string `bson:"time"`
}
