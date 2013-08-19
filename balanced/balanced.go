// A Go package that provides bindings to Balanced Payemnts API.
//
// See https://www.balancedpayments.com/
package balanced

import (
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
// The api invoked by this function is not a public endpoint at balanced.
// May not work in the future.
func setupTestEnvironment() {
	apiRoot = testApiRoot

	// Get test api key from balanced
	key := ApiKey{}
	err := post(apiKeyUri, nil, &key)
	if err != nil {
		log.Println("Unable to generate test key")
		os.Exit(1)
	}

	apiKey = key.Secret

	// Get test marketplace from balanced
	marketplace := Marketplace{}
	err = post(marketplaceUri, nil, &marketplace)
	if err != nil {
		log.Println("Unable to generate test marketplace")
		os.Exit(1)
	}

	marketplaceId = marketplace.Id
}
