package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const imds = "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01"

type responseJson struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

func getAccessToken(resource string, clientId string) responseJson {
	// Create HTTP request for a managed services for Azure resources token to access Azure Resource Manager
	var msi_endpoint *url.URL
	msi_endpoint, err := url.Parse(imds)
	if err != nil {
		log.Panicln(err)
	}
	msi_parameters := url.Values{}
	msi_parameters.Add("resource", resource)
	msi_parameters.Add("client_id", clientId)
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

func shareToken(token string, path string) error {
	// Share the token via a memory based volume. Check your k8s manifest
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	bytes, err := w.WriteString(token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Wrote %d bytes to %s\n", bytes, path)
	w.Flush()
	return err
}

func main() {

	resp := getAccessToken("https://ossrdbms-aad.database.windows.net", os.Getenv("client_id"))
	shareToken(resp.AccessToken, "/token/.token")
	duration, err := strconv.Atoi(resp.ExpiresIn)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("Sleeping for %s seconds ...\n", resp.ExpiresIn)
	time.Sleep(time.Duration(duration) * time.Second)

}
