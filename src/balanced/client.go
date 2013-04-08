package balanced

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	version   = "0.0.1"
	userAgent = "balanced-go/" + version
)

func get(path string, payload url.Values) ([]byte, error) {
	req, err := newRequest("GET", path, payload)
	if err != nil {
		log.Println("Balanced API: Error creating Get request message.")
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Balanced API: Error sending Get request message.")
		return nil, err
	}

	return responseReader(resp)
}

func post(path string, payload url.Values) ([]byte, error) {
	req, err := newRequest("POST", path, payload)
	if err != nil {
		log.Println("Balanced API: Error creating Post request message.")
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Balanced API: Error sending Post request message.")
		return nil, err
	}
	return responseReader(resp)
}

func put(path string, payload url.Values) ([]byte, error) {
	req, err := newRequest("PUT", path, payload)
	if err != nil {
		log.Println("Balanced API: Error creating Put request message.")
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Balanced API: Error sending Put request message.")
		return nil, err
	}

	return responseReader(resp)
}

func delete(path string, payload url.Values) ([]byte, error) {
	req, err := newRequest("DELETE", path, payload)
	if err != nil {
		log.Println("Balanced API: Error creating Delete request.")
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Balanced API: Error sending Delete request message.")
		return nil, err
	}

	return responseReader(resp)
}

func newRequest(method, path string, payload url.Values) (req *http.Request, err error) {
	// Build Uri
	var uri bytes.Buffer
	uri.WriteString(apiRoot)
	uri.WriteString(path)

	// Build Request
	if payload != nil && len(payload) != 0 {
		// Request with Payload
		if method == "GET" {
			// GET request encode payload in uri
			uri.WriteString("?")
			uri.WriteString(payload.Encode())
			req, err = http.NewRequest(method, uri.String(), nil)
		} else {
			// Not a GET request, encode payload into body of request
			req, err = http.NewRequest(method, uri.String(),
				strings.NewReader(payload.Encode()))
		}
	} else {
		// Request w/o Payload
		req, err = http.NewRequest(method, uri.String(), nil)
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

	return
}

func addToPayload(payload url.Values, key, value string) {
	// Check if empty
	if len(value) != 0 {
		payload.Add(key, value)
	}
}

func responseReader(resp *http.Response) ([]byte, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response bytes.")
		return nil, err
	}

	// Attempt to parse response as a balanced api error
	apiError := ApiError{}
	if err := json.Unmarshal(respBytes, &apiError); err == nil {
		// Check if api error is valid
		if len(apiError.Status) != 0 {
			return respBytes, &apiError
		}
	}

	return respBytes, nil
}
