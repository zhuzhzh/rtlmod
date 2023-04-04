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
		  { "op": "replace", "begin": "primitive udp_dff", "end": "endprimitive", "src": "./src/udp_dff.v"},
		  { "op": "replace", "begin": "primitive udp_sedfft", "end": "endprimitive", "src": "./src/udp_sedfft.v"},
		  { "op": "dummy", "begin": "module and001", "end": "endmodule", "src": ""},
		  { "op": "remove", "begin": "module or001", "end": "endmodule", "src": ""},
		  { "op": "deleteline", "begin": "celldefine", "end": "", "src": ""}
  ]
}`
	return Write(filepath, ctext)
}

func GenDff(filepath string) error {
	srctext := `module udp_dff();
<anything>
endmodule
`
	return Write(filepath, srctext)
}

func GenSeDfft(filepath string) error {
	srctext := `
module udp_sedfft();
<anything>
endmodule
`
	return Write(filepath, srctext)
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
endmodule`

	return Write(filepath, ctext)
}

func GenMk(mkpath string) error {
	mktext := `
run:
	vmod chain -c config.json -f filelist -o out 
`
	return Write(mkpath, mktext)
}

func GenFilelist(fpath string) error {
	ftext := `
./lib.v
`
	return Write(fpath, ftext)
}

func GenDemo(dir string) error {
	configpath := path.Dir(dir) + "/config.json"
	libpath := path.Dir(dir) + "/lib.v"
	fpath := path.Dir(dir) + "/filelist"
	mkpath := path.Dir(dir) + "/Makefile"
	dffpath := path.Dir(dir) + "/src/udp_dff.v"
	sedfftpath := path.Dir(dir) + "/src/udp_sedfft.v"
	if err := GenDff(dffpath); err != nil {
		return err
	}
	if err := GenSeDfft(sedfftpath); err != nil {
		return err
	}
	if err := GenConfig(configpath); err != nil {
		return err
	}
	if err := GenLib(libpath); err != nil {
		return err
	}
	if err := GenFilelist(fpath); err != nil {
		return err
	}
	if err := GenMk(mkpath); err != nil {
		return err
	}
	return nil
}
