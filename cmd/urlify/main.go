package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/WillChangeThisLater/urlify/pkg/urlify"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "S3 uploader",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "output",
				Value: "line",
				Usage: "format of the URLs output ('line', 'json', 'csv')",
			},
		},
		Action: func(c *cli.Context) error {

			// CLI variables
			region := "us-east-2"
			bucket := "urlify"
			prefix := "urlify"

			// make sure there are some arguments
			if c.Args().Len() == 0 {
				cli.ShowAppHelpAndExit(c, 0)
			}

			var urls []string

			for i := 0; i < c.Args().Len(); i++ {

				localFileName := c.Args().Get(i)
				file, err := os.Open(localFileName)
				if err != nil {
					fmt.Printf("failed to open file %q, %v\n", c.Args().Get(i), err)
					continue
				}
				defer file.Close()

				fileBytes, err := io.ReadAll(file)
				if err != nil {
					fmt.Printf("failed to open file %q, %v\n", localFileName, err)
					continue
				}

				// TODO: how do I import this?
				urlStr, err := urlify.Urlify(bucket, prefix, region, fileBytes)
				if err != nil {
					fmt.Printf("failed to urlify file %q, %v\n", localFileName, err)
					continue
				}
				urls = append(urls, urlStr)
			}

			switch c.String("output") {
			case "line":
				for _, url := range urls {
					fmt.Println(url)
				}
			case "json":
				jsonOut, err := json.Marshal(urls)
				if err != nil {
					fmt.Println("failed to serialize urls to JSON:", err)
				}
				fmt.Println(string(jsonOut))
			case "csv":
				fmt.Println(strings.Join(urls, ","))
			default:
				fmt.Println("unrecognized output format")
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
