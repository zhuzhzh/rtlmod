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
endmodule