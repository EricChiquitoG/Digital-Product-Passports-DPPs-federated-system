package main

import (
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	// The endpoints in use check API doc for description of call structure
	getURL := "/retrieveData"
	createPassportURL := "/publishPassport"
	addRemanafactureEventURL := "/addEvent"
	retrieveEventURL := "/retrieveEvent"
	addMutableProductURL := "/addMutableProduct"
	retrieveMutableLogURL := "/retrieveMutableLog"
	generateQcodeURL := "/getQrCode"

	//os.Setenv("CA_route", "http://localhost:3000/api/v1/CA/")

	mux.HandleFunc(getURL, getHandler)
	mux.HandleFunc(createPassportURL, createPassportHandler)
	mux.HandleFunc(addRemanafactureEventURL, addMutableData)
	mux.HandleFunc(retrieveEventURL, retriveEvent)
	mux.HandleFunc(addMutableProductURL, addMutableProduct)
	mux.HandleFunc(retrieveMutableLogURL, retrieveMutableLog)
	mux.HandleFunc(generateQcodeURL, generateQrCode)

	// Start the server on port 80
	handler := cors.Default().Handler(mux)

	err := http.ListenAndServe(":8000", handler)
	if err != nil {
		log.Println("There was an error listening on port :80", err)
	}
}
