package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

type Config struct {
	Opcode []struct {
		Op    string `json:"op"`
		Begin string `json:"begin"`
		End   string `json:"end"`
		Src   string `json:"src"`
	} `json:"opcode"`
}

func processFiles(configFile, fileList, outDir string) {
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		panic(err)
	}

	filesData, err := ioutil.ReadFile(fileList)
	if err != nil {
		panic(err)
	}

	files := strings.Split(string(filesData), "\n")

	// Create the output directory if it doesn't already exist
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err = os.MkdirAll(outDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	for _, file := range files {
		if file == "" {
			continue
		}

		fileData, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		fileContent := string(fileData)

		for _, op := range config.Opcode {
			switch op.Op {
			case "replace":
				srcData, err := ioutil.ReadFile(op.Src)
				if err != nil {
					panic(err)
				}

				beginIndex := strings.Index(fileContent, op.Begin)
				endIndex := strings.Index(fileContent[beginIndex:], op.End) + beginIndex + len(op.End)

				fileContent = fileContent[:beginIndex] + "// replace " + op.Begin + "\n" + string(srcData) + fileContent[endIndex:]
			case "dummy":
				beginIndex := strings.Index(fileContent, op.Begin)
				endIndex := strings.Index(fileContent[beginIndex:], op.End) + beginIndex + len(op.End)

				moduleContent := fileContent[beginIndex:endIndex]

				lines := strings.Split(moduleContent, "\n")
				newLines := []string{}
				for _, line := range lines {
					if strings.Contains(line, "module") || strings.Contains(line, "endmodule") || strings.Contains(line, "input") || strings.Contains(line, "output") || strings.Contains(line, "inout") {
						newLines = append(newLines, line)
					}
				}

				fileContent = fileContent[:beginIndex] + "// dummy " + op.Begin + "\n" + strings.Join(newLines, "\n") + fileContent[endIndex:]
			case "remove":
				beginIndex := strings.Index(fileContent, op.Begin)
				endIndex := strings.Index(fileContent[beginIndex:], op.End) + beginIndex + len(op.End)

				fileContent = fileContent[:beginIndex] + "// remove " + op.Begin + "\n" + fileContent[endIndex:]
			default:
				fmt.Printf("Unknown opcode: %s\n", op.Op)
			}
		}

		// Write the modified content to the output directory
		outPath := outDir + "/" + file[strings.LastIndex(file, "/")+1:]
		err = ioutil.WriteFile(outPath, []byte(fileContent), 0644)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	app := &cli.App{
		Name:  "program",
		Usage: "Usage: <program> -c <json> -f <file list> -o <out dir>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "c",
				Value:    "",
				Usage:    "config file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "f",
				Value:    "",
				Usage:    "file list",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "o",
				Value:    "",
				Usage:    "output directory",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			configFile := c.String("c")
			fileList := c.String("f")
			outDir := c.String("o")

			processFiles(configFile, fileList, outDir)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
