package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/zhuzhzh/vmod/internal/helper"
	"github.com/zhuzhzh/vmod/internal/vtext"
)

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
		Description: "<[replace module]|[delete module]|[empty module]|[delete line]> in the verilog file",
		Copyright:   "(c) MIT",
		Commands: []*cli.Command{
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

					vtext.ProcessFiles(configFile, fileList, outDir)
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
