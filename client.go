package balanced

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	version   = "0.0.1"
	userAgent = "balanced-go/" + version
)

func get(path string, payload url.Values, out interface{}) error {
	return request("GET", path, payload, out)
}

func post(path string, payload url.Values, out interface{}) error {
	return request("POST", path, payload, out)
}

func put(path string, payload url.Values, out interface{}) error {
	return request("PUT", path, payload, out)
}

func delete(path string, payload url.Values, out interface{}) error {
	return request("DELETE", path, payload, out)
}

func request(method, path string, payload url.Values, out interface{}) error {
	// Build Uri
	var uri bytes.Buffer
	uri.WriteString(apiRoot)
	uri.WriteString(path)

	// Build Body
	var body io.Reader
	if payload != nil && len(payload) != 0 {
		if method == "GET" {
			// GET request encode payload in uri
			uri.WriteString("?")
			uri.WriteString(payload.Encode())
		} else {
			// Not a GET request, encode payload into body of request
			body = strings.NewReader(payload.Encode())
		}
	}

	// Build Request
	req, err := http.NewRequest(method, uri.String(), body)
	if err != nil {
		log.Printf("Balanced API: Error creating %v request message.", method)
		return err
	}

	// Add Headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	// Add Basic Authentication
	if len(apiKey) != 0 {
		// Balanced does not have a traditional username and password. Just a key
		// that's passed in as username, password is left empty.
		req.SetBasicAuth(apiKey, "")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Balanced API: Error sending %v request message.", method)
		return err
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response bytes.")
		return err
	}

	// Attempt to parse response as a balanced api error
	apiError := ApiError{}
	if err := json.Unmarshal(respBytes, &apiError); err == nil {
		// Check if api error is valid
		if len(apiError.Status) != 0 {
			return &apiError
		}
	}

	// Attempt to parse response into out
	if out != nil {
		if err := json.Unmarshal(respBytes, out); err != nil {
			log.Printf("Balanced API: Error parsing response message: %v ", err.Error())
			return err
		}
	}

	return nil
}

func addToPayload(payload url.Values, key, value string) {
	// Check if empty
	if len(value) != 0 {
		payload.Add(key, value)
	}
}

func defaultPayload(limit, offset int) (payload url.Values) {
	payload = url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}

	return
}
