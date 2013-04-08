// A Go package that provides bindings to Balanced Payemnts API.
//
// See https://www.balancedpayments.com/
package balanced

import (
	"encoding/json"
	"flag"
	"github.com/stathat/jconfig"
	"log"
	"os"
)

const (
	testApiRoot = "https://api.balancedpayments.com"
)

// Basic information needed to connect to the balanced REST API.
var apiRoot, apiKey, marketplaceId string

func init() {
	var stage string

	// Get stage flag passed in. Default to test
	flag.StringVar(&stage, "stage", "test", "flag for deployment stage")
	flag.Parse()

	if stage == "test" {
		// apiRoot = testApiRoot
		// apiKey = "ec7e126692cd11e28157026ba7f8ec28"
		// marketplaceId = "TEST-MP7cnju2M0ojMCpPWS8geAPK"
		setupTestEnvironment()
	} else {
		// Retrieve config from balanced.conf
		config := jconfig.LoadConfig(stage + "/balanced.conf")
		apiRoot = config.GetString("balanced_api_root")
		apiKey = config.GetString("balanced_api_key")
		marketplaceId = config.GetString("balanced_marketplace_id")
	}
}

// Setup basic information needed to connect to balanced. Overrides any config
// set during init.
func SetupEnvironment(root, key, marketId string) {
	apiRoot = root
	apiKey = key
	marketplaceId = marketId
}

// Used when running test, or when no config file was specified.
func setupTestEnvironment() {
	apiRoot = testApiRoot

	// Get test api key from balanced
	resp, err := post(apiKeyUri, nil)
	if err != nil {
		log.Println("Unable to generate test key")
		os.Exit(1)
	}

	// Attempt to parse response into ApiKey
	key := ApiKey{}
	if err := json.Unmarshal(resp, &key); err != nil {
		log.Println("Unable to parse generated test key")
		os.Exit(1)
	}

	apiKey = key.Secret

	// Get test marketplace from balanced
	resp, err = post(marketplaceUri, nil)
	if err != nil {
		log.Println("Unable to generate test marketplace")
		os.Exit(1)
	}

	// Attempt to parse response into ApiKey
	marketplace := Marketplace{}
	if err := json.Unmarshal(resp, &marketplace); err != nil {
		log.Println("Unable to parse generated test marketplace")
		os.Exit(1)
	}

	marketplaceId = marketplace.Id
}
