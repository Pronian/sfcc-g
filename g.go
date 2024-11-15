package main

import (
	"bytes"
	"fmt"
	"os"
	"sfcc/g/kv"
	"sfcc/g/log"
	"sfcc/g/sfcc"
	"sfcc/g/util"

	env "github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var apiID string
var apiSecret string
var constantPath = true

func getConstantPath(constantPath bool, fileName string) string {
	if constantPath {
		return util.GetFilePathInExecutableDirectory(fileName)
	}
	return fileName
}

func loadEnv() {
	path := getConstantPath(constantPath, ".env")

	err := env.Load(path)
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	apiID = os.Getenv("SF_API_ID")
	apiSecret = os.Getenv("SF_API_SECRET")

	if apiID == "" || apiSecret == "" {
		log.Fatal("SF_API_ID and SF_API_SECRET must be set in .env file")
	}
}

func main() {
	log.Start(log.LevelInfo, getConstantPath(constantPath, "sfcc-g.log"))
	defer log.Stop()
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

						summary.WriteString(fmt.Sprintf("%s🔹%s%s%s", sb.ID, sb.HostName, stateEmoji, sb.State))

						if sb != sandboxList[len(sandboxList)-1] {
							summary.WriteString("\n")
						}
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
