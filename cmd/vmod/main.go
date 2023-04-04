package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	docgen "github.com/zhuzhzh/vmod/internal/demo"
	"github.com/zhuzhzh/vmod/internal/helper"
	"github.com/zhuzhzh/vmod/internal/vtext"
)

func main() {
	app := &cli.App{
		Name:     "program",
		Version:  "v1.0.0",
		Compiled: time.Now(),
		Usage:    "Usage: <program> <subcommands> <options>",
		Authors: []*cli.Author{
			{
				Name:  "Harris Zhu",
				Email: "zhuzhzh@163.com",
			},
		},
		Description: "<chain|demo> <options>",
		Copyright:   "(c) MIT",
		Commands: []*cli.Command{
			{
				Name:  "demo",
				Usage: "Usage: <program> demo",
				Action: func(c *cli.Context) error {
					err := docgen.GenDemo("./vmod_demo/")
					if err == nil {
						fmt.Print("The demo is generated in ./vmod_demo\n")
						fmt.Print("To run demo:  `cd ./vmod_demo && make run`\n")
					}
					return err
				},
			},
			{
				// add command dummy
				// flag : -f <file list>
				// flag : -o <output directory>
				// flag : -bw <begin word>
				// flag : -ew <end word>
				Name:  "replace",
				Usage: "Usage: <program> replace -f <file list> -o <out dir> -bw <begin word> -ew <end word> [--verbose <level>] [-tofile] <files1> <file2> ...",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "bw",
						Value:    "begin",
						Usage:    "begin word,
						Required: true,
					},
					&cli.StringFlag{
						Name:     "ew",
						Value:    "end",
						Usage:    "end word,
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
					bw := c.String("bw")
					ew := c.String("ew")
					fileList := c.String("f")
					outDir := c.String("o")
					tofile := c.Bool("tofile")
					files := c.Args().Slice()

					// Parse the log level from the command-line flag
					level, err := log.ParseLevel(c.String("verbose"))
					if err != nil {
						return err
					}

					if tofile {
						// create log directory
						if err = helper.CreateOutputDir("log"); err != nil {
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

					if fileFromLists, err := helper.ReadFiles(fileList); err != nil {
						fmt.Printf("can not open %s\n", fileList)
					} else {
						files = append(files, fileFromLists...)
					}

					vtext.ReplaceHelper(files, bw, ew, repl, outDir)
					return nil
				}
			}，
			{
				// add command dummy
				// flag : -f <file list>
				// flag : -o <output directory>
				// flag : -bw <begin word>
				// flag : -ew <end word>
				Name:  "dummy",
				Usage: "Usage: <program> dummy -f <file list> -o <out dir> -bw <begin word> -ew <end word> [--verbose <level>] [-tofile] <files1> <file2> ...",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "bw",
						Value:    "begin",
						Usage:    "begin word,
						Required: true,
					},
					&cli.StringFlag{
						Name:     "ew",
						Value:    "end",
						Usage:    "end word,
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
					bw := c.String("bw")
					ew := c.String("ew")
					fileList := c.String("f")
					outDir := c.String("o")
					tofile := c.Bool("tofile")
					files := c.Args().Slice()

					// Parse the log level from the command-line flag
					level, err := log.ParseLevel(c.String("verbose"))
					if err != nil {
						return err
					}

					if tofile {
						// create log directory
						if err = helper.CreateOutputDir("log"); err != nil {
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

					if fileFromLists, err := helper.ReadFiles(fileList); err != nil {
						fmt.Printf("can not open %s\n", fileList)
					} else {
						files = append(files, fileFromLists...)
					}

					vtext.DummyHelper(files, bw, ew, outDir)
					return nil

			}，
			{
				// add command remove
				// flag : -f <file list>
				// flag : -o <output directory>
				// flag : -bw <begin word>
				// flag : -ew <end word>
				// flag : -r <replacement file>
				Name:  "remove",
				Usage: "Usage: <program> remove -f <file list> -o <out dir> -bw <begin word> -ew <end word> [--verbose <level>] [-tofile] <files1> <file2> ...",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "bw",
						Value:    "begin",
						Usage:    "begin word,
						Required: true,
					},
					&cli.StringFlag{
						Name:     "ew",
						Value:    "end",
						Usage:    "end word,
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
					bw := c.String("bw")
					ew := c.String("ew")
					fileList := c.String("f")
					outDir := c.String("o")
					tofile := c.Bool("tofile")
					files := c.Args().Slice()

					// Parse the log level from the command-line flag
					level, err := log.ParseLevel(c.String("verbose"))
					if err != nil {
						return err
					}

					if tofile {
						// create log directory
						if err = helper.CreateOutputDir("log"); err != nil {
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

					if fileFromLists, err := helper.ReadFiles(fileList); err != nil {
						fmt.Printf("can not open %s\n", fileList)
					} else {
						files = append(files, fileFromLists...)
					}

					vtext.RemoveHelper(files, bw, ew, outDir)
					return nil

			}，
			{
				// add command deleteline
				// flag : -f <file list>
				// flag : -o <output directory>
				// flag : -bw <begin word>
				// flag : -ew <end word>
				// flag : -r <replacement file>
				Name:  "deleteline",
				Usage: "Usage: <program> deleteline -f <file list> -o <out dir> -kw <key word> [--verbose <level>] [-tofile] <files1> <file2> ...",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "kw",
						Value:    "keyword",
						Usage:    "key word",
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
					kw := c.String("kw")
					fileList := c.String("f")
					outDir := c.String("o")
					tofile := c.Bool("tofile")
					files := c.Args().Slice()

					// Parse the log level from the command-line flag
					level, err := log.ParseLevel(c.String("verbose"))
					if err != nil {
						return err
					}

					if tofile {
						// create log directory
						if err = helper.CreateOutputDir("log"); err != nil {
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

					if fileFromLists, err := helper.ReadFiles(fileList); err != nil {
						fmt.Printf("can not open %s\n", fileList)
					} else {
						files = append(files, fileFromLists...)
					}

					vtext.DeletelineHelper(files, kw, outDir)
					return nil

			}，
			{
				Name:  "chain",
				Usage: "Usage: <program> chain -c <json> -f <file list> -o <out dir> [--verbose <level>] [-tofile]",
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
					files := c.Args().Slice()

					// Parse the log level from the command-line flag
					level, err := log.ParseLevel(c.String("verbose"))
					if err != nil {
						return err
					}

					if tofile {
						// create log directory
						if err = helper.CreateOutputDir("log"); err != nil {
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

					if fileFromLists, err := helper.ReadFiles(fileList); err != nil {
						fmt.Printf("can not open %s\n", fileList)
					} else {
						files = append(files, fileFromLists...)
					}

					vtext.ProcessFiles(configFile, files, outDir)
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			cli.ShowAppHelp(c)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
