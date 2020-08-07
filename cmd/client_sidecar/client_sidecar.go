package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type responseJson struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	ExpiresOn   string `json:"expires_on"`
}

func getAccessToken(resource string) responseJson {
	// Create HTTP request for a managed services for Azure resources token to access Azure Resource Manager
	var msi_endpoint *url.URL
	msi_endpoint, err := url.Parse("http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01")
	if err != nil {
		log.Panicln(err)
	}
	msi_parameters := url.Values{}
	msi_parameters.Add("resource", "https://management.azure.com/")
	msi_endpoint.RawQuery = msi_parameters.Encode()
	req, err := http.NewRequest("GET", msi_endpoint.String(), nil)
	if err != nil {
		log.Panicln(err)
	}
	req.Header.Add("Metadata", "true")

	// Call managed services for Azure resources token endpoint
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error calling token endpoint: ", err)
	}

	// Extract response body
	responseBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("Error reading the response body: ", err)
	}

	// Unmarshall response body into struct
	var r responseJson
	err = json.Unmarshal(responseBytes, &r)
	if err != nil {
		log.Println("Error unmarshalling the response: ", err)
	}
	return r
}

func main() {

	resp := getAccessToken("https://management.azure.com/")
	// fmt.Println("Response status: ", resp.Status)
	fmt.Println("Access Token: ", resp.AccessToken)
	fmt.Println("Expires In: ", resp.ExpiresIn)
	fmt.Println("Expires On: ", resp.ExpiresOn)
}
