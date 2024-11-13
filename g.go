package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sfcc/g/kv"
	"sfcc/g/sfcc"
	"sfcc/g/util"

	env "github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
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
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	apiID = os.Getenv("SF_API_ID")
	apiSecret = os.Getenv("SF_API_SECRET")

	if apiID == "" || apiSecret == "" {
		fmt.Fprintf(os.Stderr, "SF_API_ID and SF_API_SECRET must be set in .env file")
		os.Exit(1)
	}
}

func main() {
	loadEnv()

	kv.Init(constantPath)
	defer kv.Close()

	app := &cli.App{
		Name:  "sfcc-g",
		Usage: "Perform sfcc actions with a Go CLI",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List fresh data for all sandboxes",
				Action: func(cCtx *cli.Context) error {
					sfcc.GetSfccAuthToken(apiID, apiSecret)
					sandboxList := sfcc.GetSandboxList(true)

					var summary bytes.Buffer
					for _, sb := range sandboxList {
						var stateEmoji string
						if sb.State == "started" {
							stateEmoji = "🟢"
						} else if sb.State == "stopped" {
							stateEmoji = "🔴"
						} else {
							stateEmoji = "🟡"
						}

						summary.WriteString(fmt.Sprintf("%s🔹%s%s%s\n", sb.ID, sb.HostName, stateEmoji, sb.State))
					}

					fmt.Println(summary.String())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
