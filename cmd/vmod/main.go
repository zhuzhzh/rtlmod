package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
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

func findBeginEnd(text, begin, end string) (int, int) {
	log.WithFields(log.Fields{
		"text":  text,
		"begin": begin,
		"end":   end,
	}).Debug("Finding begin and end indices")

	inBlockComment := false
	inLineComment := false
	for i := 0; i < len(text); i++ {
		if !inBlockComment && !inLineComment && text[i] == '/' && i+1 < len(text) && text[i+1] == '*' {
			inBlockComment = true
			i++
		} else if inBlockComment && text[i] == '*' && i+1 < len(text) && text[i+1] == '/' {
			inBlockComment = false
			i++
		} else if !inBlockComment && !inLineComment && text[i] == '/' && i+1 < len(text) && text[i+1] == '/' {
			inLineComment = true
			i++
		} else if inLineComment && text[i] == '\n' {
			inLineComment = false
		} else if !inBlockComment && !inLineComment {
			if strings.HasPrefix(text[i:], begin) {
				beginIndex := i
				for j := i + len(begin); j < len(text); j++ {
					if !inBlockComment && !inLineComment && text[j] == '/' && j+1 < len(text) && text[j+1] == '*' {
						inBlockComment = true
						j++
					} else if inBlockComment && text[j] == '*' && j+1 < len(text) && text[j+1] == '/' {
						inBlockComment = false
						j++
					} else if !inBlockComment && !inLineComment && text[j] == '/' && j+1 < len(text) && text[j+1] == '/' {
						inLineComment = true
						j++
					} else if inLineComment && text[j] == '\n' {
						inLineComment = false
					} else if !inBlockComment && !inLineComment {
						if strings.HasPrefix(text[j:], end) {
							log.WithFields(log.Fields{
								"text":       text,
								"begin":      begin,
								"end":        end,
								"beginIndex": beginIndex,
								"endIndex":   j + len(end),
							}).Debug("can't find the begin")
							return beginIndex, j + len(end)
						}
					}
				}
				log.WithFields(log.Fields{
					"text":       text,
					"begin":      begin,
					"end":        end,
					"beginIndex": beginIndex,
				}).Error("can't find the end")
				return beginIndex, -1
			}
		}
	}
	log.WithFields(log.Fields{
		"text":  text,
		"begin": begin,
		"end":   end,
	}).Error("can't find the begin")
	return -1, -1
}

func removeAction(fileContent string, begin string, end string) (string, error) {
	log.WithFields(log.Fields{
		"begin": begin,
		"end":   end,
	}).Debug("Removing content between begin and end indices")

	beginIndex, endIndex := findBeginEnd(fileContent, begin, end)
	var newContent string
	if beginIndex != -1 && endIndex != -1 {
		newContent = fileContent[:beginIndex] + "// remove " + begin + "\n" + fileContent[endIndex:]
	} else {
		newContent = fileContent
	}
	return newContent, nil
}

func dummyAction(fileContent string, begin string, end string) (string, error) {
	log.WithFields(log.Fields{
		"begin": begin,
		"end":   end,
	}).Debug("Dummying content between begin and end indices")

	beginIndex, endIndex := findBeginEnd(fileContent, begin, end)
	var newContent string
	if beginIndex != -1 && endIndex != -1 {
		moduleContent := fileContent[beginIndex:endIndex]
		lines := strings.Split(moduleContent, "\n")
		newLines := []string{}
		for _, line := range lines {
			if strings.Contains(line, "module") || strings.Contains(line, "endmodule") || strings.Contains(line, "input") || strings.Contains(line, "output") || strings.Contains(line, "inout") {
				newLines = append(newLines, line)
			}
		}
		newContent = fileContent[:beginIndex] + "// dummy " + begin + "\n" + strings.Join(newLines, "\n") + fileContent[endIndex:]
	} else {
		newContent = fileContent
	}
	return newContent, nil
}

func replaceAction(fileContent string, src string, begin string, end string) (string, error) {
	log.WithFields(log.Fields{
		"begin": begin,
		"end":   end,
	}).Debug("Replacing content between begin and end indices")

	srcData, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}
	var newContent string
	beginIndex, endIndex := findBeginEnd(fileContent, begin, end)
	if beginIndex != -1 && endIndex != -1 {
		newContent = fileContent[:beginIndex] + "// replace " + begin + "\n" + string(srcData) + fileContent[endIndex:]
	} else {
		newContent = fileContent
	}
	return newContent, nil
}

func deletelineAction(fileContent string, begin string) string {
	log.WithFields(log.Fields{
		"begin": begin,
	}).Debug("deleting the line containing the key word")

	lines := strings.Split(fileContent, "\n")
	newLines := ""
	for _, line := range lines {
		if !strings.Contains(line, begin) {
			newLines += (line + "\n")
		} else {
			newLines += fmt.Sprintf("// remove the line %s\n", begin)
		}
	}
	return newLines
}

func readConfig(configFile string) (Config, error) {
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func readFiles(fileList string) ([]string, error) {
	filesData, err := ioutil.ReadFile(fileList)
	if err != nil {
		return nil, err
	}

	files := strings.Split(string(filesData), "\n")
	return files, nil
}

func createOutputDir(outDir string) error {
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err = os.MkdirAll(outDir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func processFiles(configFile, fileList, outDir string) {
	var (
		config Config
		err    error
		files  []string
		wg     sync.WaitGroup
	)

	log.WithFields(log.Fields{
		"configFile": configFile,
	}).Debug("Reading config file")

	if config, err = readConfig(configFile); err != nil {
		log.WithFields(log.Fields{
			"configFile": configFile,
			"error":      err,
		}).Error("Error reading config file")
		return
	}

	log.WithFields(log.Fields{
		"fileList": fileList,
	}).Debug("Reading file list")

	if files, err = readFiles(fileList); err != nil {
		log.WithFields(log.Fields{
			"fileList": fileList,
			"error":    err,
		}).Error("Error reading file list")
		return
	}

	log.WithFields(log.Fields{
		"outDir": outDir,
	}).Debug("Creating output directory")

	if err = createOutputDir(outDir); err != nil {
		log.WithFields(log.Fields{
			"outDir": outDir,
			"error":  err,
		}).Error("Error creating output directory")
		return
	}

	for _, file := range files {
		if file == "" {
			continue
		}

		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			log.WithFields(log.Fields{
				"file": file,
			}).Debug("Processing file")
			fileData, err := ioutil.ReadFile(file)
			if err != nil {
				log.WithFields(log.Fields{
					"file":  file,
					"error": err,
				}).Error("Error reading file")
				return
			}

			fileContent := string(fileData)

			for _, op := range config.Opcode {
				switch op.Op {
				case "replace":
					newContent, err := replaceAction(fileContent, op.Src, op.Begin, op.End)
					if err != nil {
						log.WithFields(log.Fields{
							"op":     op,
							"error":  err,
							"action": "replace",
						}).Error("Error replacing content")
						continue
					}
					fileContent = newContent
				case "dummy":
					newContent, err := dummyAction(fileContent, op.Begin, op.End)
					if err != nil {
						log.WithFields(log.Fields{
							"op":     op,
							"error":  err,
							"action": "dummy",
						}).Error("Error dummying content")
						continue
					}
					fileContent = newContent
				case "remove":
					newContent, err := removeAction(fileContent, op.Begin, op.End)
					if err != nil {
						log.WithFields(log.Fields{
							"op":     op,
							"error":  err,
							"action": "remove",
						}).Error("Error removing content")
						continue
					}
					fileContent = newContent
				case "deleteline":
					newContent := deletelineAction(fileContent, op.Begin)
					fileContent = newContent
				default:
					fmt.Printf("Unknown opcode: %s\n", op.Op)
				}
			}

			outPath := outDir + "/" + file[strings.LastIndex(file, "/")+1:]
			log.WithFields(log.Fields{
				"outPath": outPath,
			}).Debug("Writing modified content to output directory")
			err = ioutil.WriteFile(outPath, []byte(fileContent), 0644)
			if err != nil {
				log.WithFields(log.Fields{
					"outPath": outPath,
					"error":   err,
				}).Error("Error writing modified content to output directory")
				return
			}
		}(file)
	}
	wg.Wait()
}

var helptext string = `[config file]:
{
  "opcode": [
		  { "op": "replace", "begin": "primitive udp_dff", "end": "endprimitive", "src": "./tests/udp_dff.v"},
		  { "op": "replace", "begin": "primitive udp_sedfft", "end": "endprimitive", "src": "./tests/udp_sedfft.v"},
		  { "op": "dummy", "begin": "module and001", "end": "endmodule", "src": ""},
		  { "op": "remove", "begin": "module or001", "end": "endmodule", "src": ""},
		  { "op": "deleteline", "begin": "celldefine", "end": "", "src": ""}
  ]
}
[input file]:
module test();
reg a;
endmodule

module and001(a, b, c);
input a;
input b;
output c;
wire wa;
reg wb;
assign c = a & b;
endmodule

/*
primitive udp_dff (out, in, clk, clr_, set_, NOTIFIER);
	output out;
	input in, clk, clr_, set_, NOTIFIER;
	reg out;
	table
		0 r ? 1 ? : ?:0;
	endtable
endprimitive
*/

primitive udp_dff (out, in, clk, clr_, set_, NOTIFIER);
	output out;
	input in, clk, clr_, set_, NOTIFIER;
	reg out;
	table
		0 r ? 1 ? : ?:0;
	endtable
endprimitive

primitive udp_sedfft (out, in, clk, clr_, set_, NOTIFIER);
	output out;
	input in, clk, clr_, set_, NOTIFIER;
	reg out;
	table
		0 r ? 1 ? : ?:0;
	endtable
endprimitive

module or001();
reg b;
endmodule

module done();
reg c;
endmodule

[output file]:
module test();
reg a;
endmodule

// dummy module and001
module and001(a, b, c);
input a;
input b;
output c;
endmodule

/*
primitive udp_dff (out, in, clk, clr_, set_, NOTIFIER);
	output out;
	input in, clk, clr_, set_, NOTIFIER;
	reg out;
	table
		0 r ? 1 ? : ?:0;
	endtable
endprimitive
*/

// replace primitive udp_dff
module udp_dff();
reg a;
endmodule


// replace primitive udp_sedfft
module udp_sedfft();
reg a;
endmodule


// remove module or001


module done();
reg c;
endmodule`

func main() {
	app := &cli.App{
		Name:     "program",
		Version:  "v1.0.0",
		Compiled: time.Now(),
		Usage:    "Usage: <program> -c <json> -f <file list> -o <out dir> [--verbose <level>] [-tofile]",
		Authors: []*cli.Author{
			{
				Name:  "Harris Zhu",
				Email: "zhuzhzh@163.com",
			},
		},
		Description: helptext,
		Copyright:   "(c) MIT",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "c",
				Value:    "vmod.json",
				Usage:    "config file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "f",
				Value:    "filelist",
				Usage:    "file list",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "o",
				Value:    "newout",
				Usage:    "output directory",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "verbose",
				Value: "info",
				Usage: "set the log level (debug, info, warn, error, fatal, panic)",
			},
			&cli.BoolFlag{
				Name:  "tofile",
				Value: false,
				Usage: "redirect the log into the file log/vmod.log",
			},
		},
		Action: func(c *cli.Context) error {
			configFile := c.String("c")
			fileList := c.String("f")
			outDir := c.String("o")
			tofile := c.Bool("tofile")

			// Parse the log level from the command-line flag
			level, err := log.ParseLevel(c.String("verbose"))
			if err != nil {
				return err
			}

			if tofile {
				// create log directory
				if err = createOutputDir("log"); err != nil {
					panic(err)
				}

				// Open the log file
				logfile, err := os.OpenFile("log/vmod.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					log.Fatal(err)
				}
				defer logfile.Close()

				// Set the logger output to the log file
				log.SetOutput(logfile)
			}

			// Set the log level
			log.SetLevel(level)

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
