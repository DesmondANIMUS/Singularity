package main

import (
	"fmt"
	"net/http"
)

func main() {
	//TODO: Write API's
	http.HandleFunc("/regLogUp", regLogUp)
	http.HandleFunc("/you", userProfile)
	http.HandleFunc("/checkMD5", checkMD5)

	// above API's working, tested, ok
	http.HandleFunc("/listFiles", listFiles)
	http.HandleFunc("/uploadFile", uploadFile)

	// Try & write different paths for mobile apps &
	// web app, but keep them as optimized as possible
	fmt.Println("Server is listening at port 8888")
	http.ListenAndServe(":8888", nil)
}
