package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	shell "github.com/ipfs/go-ipfs-api"
)

type TestClaim struct {
	CID string `json: CID`
}

func getHandler(writer http.ResponseWriter, request *http.Request) {
	keyD := "hej"
	writer.Header().Set("Content-Type", "application/json")
	var response []byte
	//Check that messages is GET
	if request.Method != http.MethodGet {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var testClaim TestClaim
	fmt.Println(request.Body)

	// Decode JSON from the request body into the Message struct
	err := json.NewDecoder(request.Body).Decode(&testClaim)
	if err != nil {
		fmt.Println("Get request failed", err)
		return
	}
	fmt.Println("--->CID ", testClaim.CID)

	if testClaim.CID[0] == 107 { // checks if the first char is k
		//fmt.Println("This is a public key ", key)
		key := "/ipns/" + testClaim.CID
		output := lsIPNS(key)
		content, contentLenght := splitListContent(output)
		stringindex := catContent(content, contentLenght)

		// for output := 0; output < len(stringindex); output++ {
		// SAMUELS SKIT
		// }

		for output := range stringindex {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write(decryptIt([]byte(stringindex[output]), "hej"))
		}

	}
	if testClaim.CID[0] == 81 { // checks if the first char is Q
		//fmt.Println("This is a CID ", key)
		Dpp := getPassport(testClaim.CID, keyD)
		if err != nil {
			fmt.Println("Wrong CID", err)
			return
		}
		response, err = json.Marshal(Dpp)
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write(response)

	}

}

type APICheck struct {
	APIKey string `json: api_key`
}

func createPassportHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	//Check that messages is Put
	if request.Method != http.MethodPut {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Println("Error reading body", err)
		http.Error(writer, "Error reading body", http.StatusMethodNotAllowed)
		return
	}

	//var apiCheck APICheck

	fmt.Println("Body", string(body))
	sh := shell.NewShell("localhost:5001")

	cid, err := addFile(sh, string(body))
	fmt.Println("cid----", cid)

	response, err := json.Marshal(cid)

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(response)

}

func addRemanafactureEventHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	//Check that messages is Put
	if request.Method != http.MethodPut {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Println("Error reading body", err)
		http.Error(writer, "Error reading body", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Body", string(body))
}
