package main

import (
	"log"
	"net/http"
)

func main() {

	// The endpoints in use check API doc for description of call structure
	getURL := "/retrieveData"
	createPassportURL := "/publishPassport"
	addRemanafactureEventURL := "/addEvent"
	retrieveEventURL := "/retrieveEvent"
	addMutableProductURL := "/addMutableProduct"
	retrieveMutableLogURL := "/retrieveMutableLog"
	generateQcodeURL := "/getQrCode"

	//os.Setenv("CA_route", "http://localhost:3000/api/v1/CA/")

	http.HandleFunc(getURL, getHandler)
	http.HandleFunc(createPassportURL, createPassportHandler)
	http.HandleFunc(addRemanafactureEventURL, addMutableData)
	http.HandleFunc(retrieveEventURL, retriveEvent)
	http.HandleFunc(addMutableProductURL, addMutableProduct)
	http.HandleFunc(retrieveMutableLogURL, retrieveMutableLog)
	http.HandleFunc(generateQcodeURL, generateQrCode)

	// Start the server on port 80
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Println("There was an error listening on port :80", err)
	}
}
