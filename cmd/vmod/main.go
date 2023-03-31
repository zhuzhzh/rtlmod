package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/zhuzhzh/vmod/internal/vtext"
)

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
				if err = vtext.CreateOutputDir("log"); err != nil {
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
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
