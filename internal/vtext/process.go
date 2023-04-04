package vtext

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/zhuzhzh/vmod/internal/helper"
)

type Config struct {
	Opcode []struct {
		Op    string `json:"op"`
		Begin string `json:"begin"`
		End   string `json:"end"`
		Src   string `json:"src"`
	} `json:"opcode"`
}

type pIndex struct {
	beginIndex int
	endIndex   int
}

type func SingleParamFunc func(files []string, kw string) (string, error)
type func TwoParamFunc func(files []string, bw string, ew string) (string, error)
type func ThreeParamFunc func(files []string, bw string, ew string, subst string) (string, error)

func findFirstBeginEnd(text, begin, end string) (int, int) {
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

// findAllBeginEnd searches for pairs of beginword and endword in the text.
// The function skips over line and block comments during the search.
// endword should be the first one following beginword tightly.
func findAllBeginEnd(text, beginword, endword string) []pIndex {
	var res []pIndex
	var i int
	for i < len(text) {
		if text[i] == '/' && i+1 < len(text) && text[i+1] == '/' {
			log.Debug("Skipping line comment")
			i = skipLineComment(text, i)
		} else if text[i] == '/' && i+1 < len(text) && text[i+1] == '*' {
			log.Debug("Skipping block comment")
			i = skipBlockComment(text, i)
		} else if strings.HasPrefix(text[i:], beginword) {
			log.Debugf("Found beginword at index %d", i)
			j := i + len(beginword)
			for j < len(text) {
				if text[j] == '/' && j+1 < len(text) && text[j+1] == '/' {
					log.Debug("Skipping line comment")
					j = skipLineComment(text, j)
				} else if text[j] == '/' && j+1 < len(text) && text[j+1] == '*' {
					log.Debug("Skipping block comment")
					j = skipBlockComment(text, j)
				} else if strings.HasPrefix(text[j:], endword) {
					log.Debugf("Found endword at index %d", j)
					break
				} else {
					j++
				}
			}
			if j < len(text) && strings.HasPrefix(text[j:], endword) {
				res = append(res, pIndex{i, j + len(endword)})
				i = j + len(endword)
			} else {
				i = j
			}
		} else {
			i++
		}
	}
	return res
}

// skipLineComment skips over a line comment starting at index i in the text.
func skipLineComment(text string, i int) int {
	for i < len(text) && text[i] != '\n' {
		i++
	}
	return i
}

// skipBlockComment skips over a block comment starting at index i in the text.
func skipBlockComment(text string, i int) int {
	for i < len(text)-1 && !(text[i] == '*' && text[i+1] == '/') {
		i++
	}
	return i + 2
}

func removeText(input string, keyword string, p []pIndex) (output string) {
	var start int
	for _, pair := range p {
		output += input[start:pair.beginIndex]
		output += ("// remove " + keyword + "\n")
		start = pair.endIndex
	}
	output += input[start:]
	return
}

func RemoveAction(fileContent string, begin string, end string) (string, error) {
	log.WithFields(log.Fields{
		"begin": begin,
		"end":   end,
	}).Debug("Removing content between begin and end indices")

	occurs := findAllBeginEnd(fileContent, begin, end)
	newContent := removeText(fileContent, begin+"..."+end, occurs)
	return newContent, nil
}

func dummyText(input string, keyword string, p []pIndex) (output string) {
	var start int
	for _, pair := range p {
		output += input[start:pair.beginIndex]

		moduleContent := input[pair.beginIndex:pair.endIndex]
		lines := strings.Split(moduleContent, "\n")
		newLines := []string{}
		for _, line := range lines {
			if strings.Contains(line, "module") || strings.Contains(line, "endmodule") || strings.Contains(line, "input") || strings.Contains(line, "output") || strings.Contains(line, "inout") {
				newLines = append(newLines, line)
			}
		}
		output += ("// dummy " + keyword + "\n")
		output += strings.Join(newLines, "\n")
		start = pair.endIndex
	}
	output += input[start:]
	return
}

func DummyAction(fileContent string, begin string, end string) (string, error) {
	log.WithFields(log.Fields{
		"begin": begin,
		"end":   end,
	}).Debug("Dummying content between begin and end indices")

	occurs := findAllBeginEnd(fileContent, begin, end)
	newContent := dummyText(fileContent, begin+"..."+end, occurs)
	return newContent, nil
}

func replaceText(input string, repl string, keyword string, p []pIndex) (output string) {
	var start int
	for _, pair := range p {
		output += input[start:pair.beginIndex]
		output += ("// replace " + keyword + "\n")
		output += repl
		start = pair.endIndex
	}
	output += input[start:]
	return
}

func ReplaceAction(fileContent string, replFile string, begin string, end string) (string, error) {
	log.WithFields(log.Fields{
		"begin": begin,
		"end":   end,
	}).Debug("Replacing content between begin and end indices")

	srcData, err := ioutil.ReadFile(replFile)
	if err != nil {
		panic(err)
	}
	occurs := findAllBeginEnd(fileContent, begin, end)
	newContent := replaceText(fileContent, string(srcData), begin, occurs)
	return newContent, nil
}

func DeletelineAction(fileContent string, begin string) (string, error) {
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
	return newLines, nil
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

func ActionHelper(files []string, bw string, ew string, repl string, outDir string, funcHelper interface{}, funcDesc string) {
	var (
		err error
		wg  sync.WaitGroup
	)
	log.WithFields(log.Fields{
		"outDir": outDir,
	}).Debug("Creating output directory")

	if err = helper.CreateOutputDir(outDir); err != nil {
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

			var newContent string

			switch f := funcHelper.(type) {
			case SingleParamFunc:
				newContent, err = f(fileContent, bw)
			case TwoParamFunc:
				newContent, err = f(fileContent, bw, ew)
			case ThreeParamFunc:
				newContent, err = f(fileContent, bw, ew, repl)
			default:
				log.Error("Unsupported function type")
				return
			}

			if err != nil {
				log.WithFields(log.Fields{
					"op":     op,
					"error":  err,
					"action": funcDesc,
				}).Error("Error processing content")
				continue
			}
			fileContent = newContent
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

func RemoveHelper(files []string, bw string, ew string, outDir string)  {
	ActionHelper(files, bw, ew, "", outDir, RemoveAction, "remove")
}

func DeleteLineHelper(files []string, kw string, outDir string)  {
	ActionHelper(files, kw, "", "", outDir, DeletelineAction, "delete line")
}

func DummyHelper(files []string, bw string, ew string, outDir string)  {
	ActionHelper(files, bw, ew, "", outDir, DummyAction, "dummy")
}

func ReplaceHelper(files []string, bw string, ew string, repl string, outDir string) {
	ActionHelper(files, bw, ew, repl, outDir, ReplaceAction, "replace")
}

func ProcessFiles(configFile string, files []string, outDir string) {
	var (
		config Config
		err    error
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
		"outDir": outDir,
	}).Debug("Creating output directory")

	if err = helper.CreateOutputDir(outDir); err != nil {
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
					newContent, err := ReplaceAction(fileContent, op.Src, op.Begin, op.End)
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
					newContent, err := DummyAction(fileContent, op.Begin, op.End)
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
					newContent, err := RemoveAction(fileContent, op.Begin, op.End)
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
					newContent := DeletelineAction(fileContent, op.Begin)
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
