# rtlmod

## Overview

This program is used to handle the RTL files.

Most actions will be taken on the text between <begin_word> and the <end_word>. This is very useful on RTL code.
For example, I want to add one line on module foo, then you could define <begin_word> as module foo, and <end_word> as endmodule.

Here is the supported actions.

- **replace**: replace one block starting from the keyword <begin_word> to the keyword <end_word> with another file
- **remove**: remove one block starting from the keyword <begin_word> to the keyword <end_word>
- **dummy**: dummy one block starting from the keyword <begin_word> to the keyword <end_word>
- **delete line**: delete the lines containing one keyword

## Usage

```shell
rtlmod replace -f <filelist> -o <output dir> -bw <bw> -ew <ew> -r <file to replace> <files>...
rtlmod remove -f <filelist> -o <output dir> -bw <bw> -ew <ew> -r <file to replace> <files>...
rtlmod dummy -f <filelist> -o <output dir> -bw <bw> -ew <ew> -r <file to replace> <files>...
rtlmod deleteline -f <filelist> -o <output dir> -kw <kw><files>...
```

## chain mode

All these actions can be applied on the files in the chain mode.
The chain mode needs one config file in json format. Here is one example.

```json
{
  "opcode": [
		  { "op": "replace", "begin": "primitive udp_dff", "end": "endprimitive", "src": "./test/udp_dff.v"},
		  { "op": "replace", "begin": "primitive udp_sedfft", "end": "endprimitive", "src": "./test/udp_sedfft.v"},
		  { "op": "dummy", "begin": "module and001", "end": "endmodule", "src": ""},
		  { "op": "remove", "begin": "module or001", "end": "endmodule", "src": ""},
		  { "op": "deleteline", "begin": "celldefine", "end": "", "src": ""}
  ]
}
```

