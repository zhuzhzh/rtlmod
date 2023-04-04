package vtext

import (
	"testing"
)

func TestFindBeginEnd(t *testing.T) {
	text := `This is some text with // a line <begin> comment
// a line <end> comment
and /* a block <begin> comment 
aaa bb
<end> */.
Here is the <begin> keyword.
// <end> keyword
/* aaa
<end> bbb
*/
And here is the <end> keyword.
Here is the <begin> keyword.
// <end> keyword
/* aaa
<end> bbb
*/
And here is the <end> keyword.`
	t.Logf("text = \n%s\n", text)

	beginIndex := 0
	endIndex := 0

	for (beginIndex != -1) && (endIndex != -1) {
		beginIndex, endIndex = findBeginEnd(text[endIndex:], "<begin>", "<end>")
		t.Logf("beginIndex = %d, endIndex = %d\n", beginIndex, endIndex)
	}

	expectBegin := 134
	expectEnd := 209

	if beginIndex != expectBegin {
		t.Errorf("Expected beginIndex to be %d, but got %d", expectBegin, beginIndex)
	}

	if endIndex != expectEnd {
		t.Errorf("Expected endIndex to be %d, but got %d", expectEnd, endIndex)
	}
	t.Logf("begin = \n[%s]\n", text[beginIndex:])
	t.Logf("end = \n[%s]\n", text[endIndex:])
}
