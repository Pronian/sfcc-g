package main

import (
	"fmt"
	"os"
	"sfcc/g/kv"
	"sfcc/g/sfcc"
	"sfcc/g/util"

	env "github.com/joho/godotenv"
)

var apiID string
var apiSecret string
var constantPath = false

func loadEnv() {
	var path string
	if constantPath {
		path = util.GetFilePathInExecutableDirectory(".env")
	} else {
		path = "./.env"
	}

	err := env.Load(path)
	if err != nil {
		fmt.Errorf("Error loading .env file", err)
		os.Exit(1)
	}

	apiID = os.Getenv("SF_API_ID")
	apiSecret = os.Getenv("SF_API_SECRET")

	if apiID == "" || apiSecret == "" {
		fmt.Errorf("SF_API_ID and SF_API_SECRET must be set in .env file")
		os.Exit(1)
	}
}

func main() {
	loadEnv()

	kv.Init(constantPath)
	defer kv.Close()

	value := sfcc.GetSfccAuthToken(apiID, apiSecret)
	fmt.Printf("SFCC Auth Token: %s\n", value)
}
