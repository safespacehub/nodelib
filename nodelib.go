package nodelib

import (
	"fmt"
    "io/ioutil"
    "net/http"
	"encoding/json"
	"bytes"
)

type Node struct {
	address string
	https bool
}

// NewSelf is a short function to connect to self using localhost:8080
func NewSelf() Node {	
return NewNode("127.0.0.1:8080", false)
}

// NewNode takes a connection and a bool whether to use HTTPS, and returns a node.
func NewNode(conn string, httpsFlag bool) Node{
	url := "http://"
	if httpsFlag {
		url = "https://"
	}

	url += conn

	n := Node{
		address: url,
		https: httpsFlag,
	}

	return n
}

// Get will get the specified key or return an error.
func (n Node) Get(key string) (error, string) {
	// The URL to send the GET request to
	url := fmt.Sprintf("%v/kv/%v", n.address, key)

	// Make the GET request
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Could not GET request", err), ""
	}
	defer response.Body.Close() // Ensure the response body is closed

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Could not read body", err), ""
	}

	return nil, string(body)

}

// Set will do a POST request with the given key and value or return an error.
func (n Node) Set(key string, val string) error {
	 // The URL to send the POST request to
	 url := n.address + "/kv"

	 fmt.Println(url)

	 // Data to be sent in the POST request body
	 postData := map[string]string{
		 "key": key,
		 "value": val,
	 }
 
	 // Convert the data to JSON format
	 jsonData, err := json.Marshal(postData)
	 if err != nil {
		 return fmt.Errorf("Could not encode JSON", err)
	 }
 
	 // Create a new POST request with the JSON body
	 response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	 if err != nil {
		 return fmt.Errorf("Error making the POST request:", err)
	 }
	 defer response.Body.Close()
 
	 // Read the response body
	 /*
	 body, err := ioutil.ReadAll(response.Body)
	 if err != nil {
		 return fmt.Errorf("Error reading the response body:", err)
	 }
		
	*/

	return nil
}

// Delete will create a DELETE request for the given key or return an error.
func (n Node) Delete(key string) error {
	 // The URL to send the DELETE request
	 url := n.address +"/kv/" + key

	 // Create a new HTTP request with the DELETE method
	 req, err := http.NewRequest("DELETE", url, nil)
	 if err != nil {
		 return fmt.Errorf("Error creating the DELETE request:", err)
	 }
 
	 // Send the DELETE request using the default HTTP client
	 client := &http.Client{}
	 response, err := client.Do(req)
	 if err != nil {
		 return fmt.Errorf("Error sending the DELETE request:", err)
	 }
	 defer response.Body.Close()
 
	return nil
}

// GetALL will run a get request for all keys on a node or return a node.
func (n Node) GetAll() (error, map[string]string) {
	// The URL to send the GET request to
	url := n.address + "/kv"

	// Make the GET request
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Could not GET request", err), nil
	}
	defer response.Body.Close() // Ensure the response body is closed

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Could not read body", err), nil
	}

	// Create a map to hold the parsed data
	var result map[string]string

	// Parse JSON string into the map
	err = json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("Error parsing JSON:", err), nil
	}

	return nil, result

}