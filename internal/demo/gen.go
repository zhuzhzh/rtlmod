// Package vdemo  provides ...
package vdemo

import (
	"io/ioutil"
	"path"

	"github.com/zhuzhzh/vmod/internal/helper"
)

func Write(filepath string, ctext string) error {

	dir := path.Dir(filepath)
	if err := helper.CreateOutputDir(dir); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath, []byte(ctext), 0644); err != nil {
		return err
	}

	return nil
}

func GenConfig(filepath string) error {
	ctext := `{
  "opcode": [
		  { "op": "replace", "begin": "primitive udp_dff", "end": "endprimitive", "src": "./tests/udp_dff.v"},
		  { "op": "replace", "begin": "primitive udp_sedfft", "end": "endprimitive", "src": "./tests/udp_sedfft.v"},
		  { "op": "dummy", "begin": "module and001", "end": "endmodule", "src": ""},
		  { "op": "remove", "begin": "module or001", "end": "endmodule", "src": ""},
		  { "op": "deleteline", "begin": "celldefine", "end": "", "src": ""}
  ]
}`
	return Write(filepath, ctext)
}

func GenLib(filepath string) error {
	ctext := `module test();
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

	return Write(filepath, ctext)
}

func GenDemo(dir string) error {
	configpath := path.Dir(dir) + "/config.json"
	libpath := path.Dir(dir) + "/lib.v"
	if err := GenConfig(configpath); err != nil {
		return err
	}
	if err := GenLib(libpath); err != nil {
		return err
	}
	return nil
}
