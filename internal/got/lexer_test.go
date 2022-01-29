package got

import (
	"bufio"
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func TestSplitParams(t *testing.T) {
	type test struct {
		name  string
		input string
		want  []string
	}
	tests := []test{
		{"one", "test", []string{"test"}},
		{"two", "test, test2", []string{"test", "test2"}},
		{"three", "test, test2, test3", []string{"test", "test2", "test3"}},
		{"no space", "test,test2,test3", []string{"test", "test2", "test3"}},
		{"one quote", `"test"`, []string{"test"}},
		{"one quote with space", `"test test2"`, []string{"test test2"}},
		{"one quote with param", `"test", test2`, []string{"test", "test2"}},
		{"one quote with space and param", `"test test2", test3, test4`, []string{"test test2", "test3", "test4"}},
		{"quote with comma", `"test, test2"`, []string{"test, test2"}},
		{"quote with quote", `"test]\", test2"`, []string{`test]", test2`}},
		{"quote with 2 quotes", `"test]\", \"test2"`, []string{`test]", "test2`}},
		{"go-like call", `a.b.c("d"), test2`, []string{`a.b.c("d")`, "test2"}},
		{"go-like call with quote", `a.b.c("d\""), test2`, []string{`a.b.c("d\"")`, "test2"}},
		{"empty param", `test1,,test2`, []string{`test1`, "", "test2"}},
		{"empty space param", `test1, ,test2`, []string{`test1`, "", "test2"}},
		{"space param", `test1," " ,test2`, []string{`test1`, " ", "test2"}},
		{"3 empty param", `,,`, []string{"", "", ""}},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := splitParams(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitParamsError(t *testing.T) {
	type test struct {
		name  string
		input string
	}
	tests := []test{
		{"one quote", `"test`},
		{"one quote with quote", `"test \"`},
		{"one quote with comma", `"test ,`},
		{"two param, one quote with comma", `test1, "test ,`},
		{"three quote", `"test1"", test ,`},
		{"three quote no space", `"test1"",test,`},
		{"2nd no quote", `test1,"`},
		{"Only quote", `"`},
		{"Only quote 2", `,"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := splitParams(tt.input)
			assert.Error(t, err)
		})
	}
}

func Test_lexBlock(t *testing.T) {
	type args struct {
		blockName string
		content   string
	}
	tests := []struct {
		name string
		args args
		want *lexer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lexBlock(tt.args.blockName, tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lexBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lexer_currentLen(t *testing.T) {
	l := newTestLexer("")
	assert.Equal(t, 0, l.currentLen())

	l = newTestLexer("abc")
	assert.Equal(t, 0, l.currentLen())
	l.acceptRun()
	assert.Equal(t, 3, l.currentLen())
	l.backup()
	c := l.currentLen()
	assert.Equal(t, 2, c)
}

func Test_lexer_currentString(t *testing.T) {
	l := newTestLexer("")
	assert.Equal(t, "", l.currentString())

	l = newTestLexer("abc")
	assert.Equal(t, "", l.currentString())
	l.acceptRun()
	assert.Equal(t, "abc", l.currentString())
	l.backup()
	assert.Equal(t, "ab", l.currentString())
}

func Test_lexer_ignoreN(t *testing.T) {
	l := newTestLexer("")
	l.ignoreN(2) // make sure this does not panic

	l = newTestLexer("abc123")
	l.ignoreN(2) // make sure this does not panic
	l.acceptRun()
	assert.Equal(t, "c123", l.currentString())
}

func Test_lexer_ignoreCloseTag(t *testing.T) {
	l := newTestLexer("abc")
	l.ignoreCloseTag()
	l.acceptRun()
	assert.Equal(t, "abc", l.currentString())

	l = newTestLexer("}}abc")
	l.ignoreCloseTag()
	l.acceptRun()
	assert.Equal(t, "abc", l.currentString())

	l = newTestLexer(" }}abc")
	l.ignoreCloseTag()
	l.acceptRun()
	assert.Equal(t, "abc", l.currentString())
}

func Test_lexer_ignoreNewline(t *testing.T) {
	tests := []struct {
		name   string
		content string
		expected string
	}{
		{"empty", "", ""},
		{"no new line", "abc123", "abc123"},
		{"beginning new line", "\nabc123", "abc123"},
		{"win new line", "\r\nabc123", "abc123"},
		{"mid new line", "abc\n123", "abc\n123"},
		{"fake new line", "\r123", "\r123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newTestLexer(tt.content)
			l.ignoreNewline()
			l.acceptRun()
			assert.Equal(t, tt.expected, l.currentString())
		})
	}
}

func Test_lexer_ignoreOneSpace(t *testing.T) {
	tests := []struct {
		name   string
		content string
		expected string
	}{
		{"empty", "", ""},
		{"no space", "abc123", "abc123"},
		{"space", " abc123", "abc123"},
		{"newline", "\nabc123", "abc123"},
		{"win new line", "\r\nabc123", "abc123"},
		{"fake win new line", "\rabc123", "\rabc123"},
		{"mid space", "abc 123", "abc 123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newTestLexer(tt.content)
			l.ignoreOneSpace()
			l.acceptRun()
			assert.Equal(t, tt.expected, l.currentString())
		})
	}
}

func Test_lexer_ignoreWhiteSpace(t *testing.T) {
	tests := []struct {
		name   string
		content string
		expected string
	}{
		{"empty", "", ""},
		{"no space", "abc123", "abc123"},
		{"space begin", " abc123", "abc123"},
		{"newline", "\nabc123", "abc123"},
		{"win new line", "\r\nabc123", "abc123"},
		{"fake win new line", "\rabc123", "abc123"},
		{"mid space", "abc 123", "abc 123"},
		{"multi space", "   123", "123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newTestLexer(tt.content)
			l.ignoreWhiteSpace()
			l.acceptRun()
			assert.Equal(t, tt.expected, l.currentString())
		})
	}
}

func Test_lexer_ignoreSpace(t *testing.T) {
	tests := []struct {
		name   string
		content string
		expected string
	}{
		{"empty", "", ""},
		{"no space", "abc123", "abc123"},
		{"space begin", " abc123", "abc123"},
		{"newline", "\nabc123", "\nabc123"},
		{"win new line", "\r\nabc123", "\r\nabc123"},
		{"tab", "\tabc123", "abc123"},
		{"multi space", "   \t123", "123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newTestLexer(tt.content)
			l.ignoreSpace()
			l.acceptRun()
			assert.Equal(t, tt.expected, l.currentString())
		})
	}
}

func Test_lexer_putBackCurBuffer(t *testing.T) {
	l := newTestLexer("abc")
	l.acceptRun()
	l.putBackCurBuffer()
	l.acceptRun()
	assert.Equal(t, "abc", l.currentString())

	l = newTestLexer("abc{{123")
	l.acceptRun()
	l.ignoreN(2)
	l.acceptRun()
	l.putBackCurBuffer()
	l.acceptRun()
	assert.Equal(t, "123", l.currentString())
}

func Test_lexer_backup(t *testing.T) {
	// just test for panic. Mostly tested in other places.
	l := newTestLexer("")
	assert.Panics(t, func() {
		l.backup()
	})
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func newTestLexer(content string) *lexer {
	l := &lexer{
		input:    bufio.NewReader(strings.NewReader(content)),
		blockName: "test",
		items:    make(chan tokenItem),
	}
	return l
}


// A lexer that simply returns an error for testing of error responses
func lexError() *lexer {

	l := &lexer{
		input:    bufio.NewReader(errReader(0)),
		blockName: "test",
		items:    make(chan tokenItem),
	}

	return l
}

func Test_lexer_next(t *testing.T) {
	l := lexError()
	assert.Equal(t, errRune, l.next())

	l = newTestLexer("")
	assert.Equal(t, eof, l.next())
}

func Test_lexer_acceptTag(t *testing.T) {
	tests := []struct {
		name   string
		content string
		expected string
	}{
		{"short tag", "{{", "{{"},
		{"long tag", "{{abc ad", "{{abc"},
		{"block tag", "{{abc}}", "{{abc}}"},
		{"empty", "", ""},
		{"no tag", "abc123", ""},
		{"g tag", "{{g\nabc", "{{g"},
		{"block tag with trailer", "{{abc}}ert", "{{abc}}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newTestLexer(tt.content)
			s := l.acceptTag()
			assert.Equal(t, tt.expected, s)
		})
	}
}

func Test_lexer_acceptUntil(t *testing.T) {
	tests := []struct {
		name   string
		content string
		expected string
	}{
		{"start", "{{abc}}123", ""},
		{"end", "123{{abc}}", "123"},
		{"mid", "12{{abc}}34", "12"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newTestLexer(tt.content)
			l.acceptUntil("{{abc}}")
			assert.Equal(t, tt.expected, l.currentString())
		})
	}
}

func Test_lexer_acceptRun(t *testing.T) {
	tests := []struct {
		name   string
		content string
		expected string
	}{
		{"start", "{{abc}}123", ""},
		{"end", "123{{abc}}", "123"},
		{"mid", "12{{abc}}34", "12"},
		{"empty", "", ""},
		{"eof", "abc123", "abc123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newTestLexer(tt.content)
			l.acceptRun()
			assert.Equal(t, tt.expected, l.currentString())
		})
	}

	t.Run("error", func(t *testing.T) {
		l := lexError()
		l.acceptRun()
		assert.Equal(t, "", l.currentString())
	})

}

func Test_lexer_peekN(t *testing.T) {
	tests := []struct {
		name   string
		content string
		n int
		expected string
	}{
		{"start", "{{abc}}123", 2, "{{"},
		{"all", "123", 3, "123"},
		{"more", "123", 4, "123"},
		{"empty", "", 1, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newTestLexer(tt.content)
			s := l.peekN(tt.n)
			assert.Equal(t, tt.expected, s)
			assert.Equal(t, "", l.currentString()) // make sure nothing is in the current buffer
		})
	}

	t.Run("error", func(t *testing.T) {
		l := lexError()
		s := l.peekN(1)
		assert.Equal(t, "", s)
		assert.Equal(t, "", l.currentString())
	})
}

func Test_lexer_peek(t *testing.T) {
	tests := []struct {
		name   string
		content string
		expected rune
	}{
		{"start", "{{abc}}123", '{'},
		{"empty", "", eof},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newTestLexer(tt.content)
			s := l.peek()
			assert.Equal(t, tt.expected, s)
			assert.Equal(t, "", l.currentString()) // make sure nothing is in the current buffer
		})
	}

	t.Run("error", func(t *testing.T) {
		l := lexError()
		s := l.peek()
		assert.Equal(t, errRune, s)
		assert.Equal(t, "", l.currentString())
	})
}

func runBlockLexer(content string) (ret []tokenItem, l *lexer) {
	l = lexBlock("test", content)

	for tok := range l.items {
		ret = append(ret, tok)
	}
	return ret, l
}

type lexTypesTest struct {
	name   string
	content string
	want []tokenType
}
// a rough test to make sure the correct types of tokens are returned
// does not drill down into individual tokens
func Test_lexTypes(t *testing.T) {
	tests := []lexTypesTest{
		{"run", "abc123", []tokenType{itemRun}},
		{"run with text", "abc{{ 123 }}", []tokenType{itemRun, itemText, itemRun, itemEnd}},
		{"go token not terminated", "{{g abc", []tokenType{itemGo, itemRun}},
		{"go token terminated", "{{g abc}}", []tokenType{itemGo, itemRun, itemEnd}},
		{"go with text", "{{g abc {{ 123 }} }}", []tokenType{itemGo, itemRun, itemText, itemRun, itemEnd, itemEnd}},
		{"go with text", "{{g abc {{ 123 }} }}", []tokenType{itemGo, itemRun, itemText, itemRun, itemEnd, itemEnd}},
		{"go value", "{{abc}}", []tokenType{itemInterface, itemRun, itemEnd}},
		{"strict block", "{{begin abc}} 123 {{g }} {{end abc}}", []tokenType{itemStrictBlock}},
		{"strict block error 1", "{{begin abc {{sf}} }} 123 {{g }} {{end abc}}", []tokenType{itemError}},
		{"strict block error 2", "{{begin abc}} 123 {{g }} {{end abcd}}", []tokenType{itemError}},
		{"comment", "{{// adfaf }}abc", []tokenType{itemRun}},
		{"join", "{{join a,b }}c{{join}}", []tokenType{itemJoin, itemParam, itemParam, itemEnd, itemRun,itemEndBlock}},
		{"join error", "{{join a,\"b }}c{{join}}", []tokenType{itemJoin, itemError}},
		{"join error 2", "{{join a,b {{d}} }}c{{join}}", []tokenType{itemJoin, itemError}},
		{"if", "{{if a>b }}c{{if}}", []tokenType{itemIf, itemRun, itemEnd, itemRun,itemEndBlock}},
		{"else", "{{if a>b }}c{{else}}d{{if}}", []tokenType{itemIf, itemRun, itemEnd, itemRun,itemEndBlock, itemRun, itemEndBlock}},
		{"elseif", "{{if a>b }}c{{elseif c<d}}d{{if}}", []tokenType{itemIf, itemRun, itemEnd, itemRun,itemEndBlock, itemRun, itemEnd, itemRun, itemEndBlock}},
		{"for", "{{for _,i := range g.ar }}c{{for}}", []tokenType{itemFor, itemRun, itemEnd, itemRun,itemEndBlock}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items,_ := runBlockLexer(tt.content)
			if len(items) == len(tt.want) {
				for i := 0; i < len(items); i++ {
					assert.Equal(t, items[i].typ, tt.want[i])
				}
			} else {
				if len(items) > 0 && items[0].typ == itemError {
					t.Errorf("lexRun() error: %s", items[0].val)
				} else {
					t.Errorf("lexRun(): len = %d, want %d", len(items), len(tt.want))
				}
			}
		})
	}

}

func Test_lexBlocks(t *testing.T) {
	items, l := runBlockLexer("{{< abc}}123{{end abc}}")
	assert.Len(t, items, 0)
	assert.Equal(t, namedBlockEntry{"123", 0}, l.namedBlocks["abc"])

	items, l = runBlockLexer("{{< abc}}123{{end abc}} {{< abc}}123{{end abc}}")
	assert.Len(t, items, 1)
	assert.Equal(t, items[0].typ, itemError)

	items, l = runBlockLexer("{{< abc}}123{{end abc}} {{< def 2}}456{{end def}}")
	assert.Len(t, items, 0)
	assert.Equal(t, namedBlockEntry{"456", 2}, l.namedBlocks["def"])

	items, l = runBlockLexer("{{< abc}}123")
	assert.Len(t, items, 1)
	assert.Equal(t, items[0].typ, itemError)

	items, l = runBlockLexer("{{< abc {{xyz}} }}")
	assert.Len(t, items, 1)
	assert.Equal(t, items[0].typ, itemError)

	items, l = runBlockLexer("{{< abc def }}")
	assert.Len(t, items, 1)
	assert.Equal(t, items[0].typ, itemError)

	items, l = runBlockLexer("{{< abc def tre}}")
	assert.Len(t, items, 1)
	assert.Equal(t, items[0].typ, itemError)

	items, l = runBlockLexer("{{< abc\n}}")
	assert.Len(t, items, 1)
	assert.Equal(t, items[0].typ, itemError)

	items, l = runBlockLexer("{{< include}}")
	assert.Len(t, items, 1)
	assert.Equal(t, items[0].typ, itemError)

}

func Test_subsitute(t *testing.T) {
	t.Run("std block", func(t *testing.T) {
		items, _ := runBlockLexer("{{< abc}}123{{end abc}}{{> abc}}")
		assert.Equal(t, 1, len(items))
		assert.Equal(t, items[0].typ, itemRun)
		assert.Equal(t, items[0].val, "123")
	})

	t.Run("short block", func(t *testing.T) {
		items, _ := runBlockLexer("{{< abc}}123{{end abc}}{{abc}}")
		assert.Equal(t, 1, len(items))
		assert.Equal(t, items[0].typ, itemRun)
		assert.Equal(t, items[0].val, "123")
	})

	t.Run("error1", func(t *testing.T) {
		items, _ := runBlockLexer("{{< abc}}123{{end abc}}{{> {{abc}} }}")
		assert.Equal(t, 1, len(items))
		assert.Equal(t, items[0].typ, itemError)
	})

	t.Run("error2", func(t *testing.T) {
		items, _ := runBlockLexer("{{< abc}}123{{end abc}}{{> a {{abc}} }}")
		assert.Equal(t, 1, len(items))
		assert.Equal(t, items[0].typ, itemError)
	})


	t.Run("missing block error", func(t *testing.T) {
		items, _ := runBlockLexer("{{< abc}}123{{end abc}}{{> def }}")
		assert.Equal(t, 1, len(items))
		assert.Equal(t, items[0].typ, itemError)
	})

	t.Run("1 param", func(t *testing.T) {
		items, _ := runBlockLexer("{{< abc 1}}123$1{{end abc}}{{> abc d}}")
		assert.Equal(t, 1, len(items))
		assert.Equal(t, items[0].typ, itemRun)
		assert.Equal(t, items[0].val, "123d")
	})

	t.Run("2 param", func(t *testing.T) {
		items, _ := runBlockLexer("{{< abc 2}}123$1$2{{end abc}}{{> abc d,e}}")
		assert.Equal(t, 1, len(items))
		assert.Equal(t, items[0].typ, itemRun)
		assert.Equal(t, items[0].val, "123de")
	})

	t.Run("param error 1", func(t *testing.T) {
		items, _ := runBlockLexer("{{< abc 2}}123$1$2{{end abc}}{{> abc d e}}")
		assert.Equal(t, 1, len(items))
		assert.Equal(t, items[0].typ, itemError)
	})

	t.Run("param error 2", func(t *testing.T) {
		items, _ := runBlockLexer("{{< abc 2}}123$1$2{{end abc}}{{> abc d,\"e}}")
		assert.Equal(t, 1, len(items))
		assert.Equal(t, items[0].typ, itemError)
	})

	t.Run("error in block", func(t *testing.T) {
		items, _ := runBlockLexer("{{< abc}}{{# {{end abc}}{{> abc}}")
		assert.Equal(t, 2, len(items))
		assert.Equal(t, items[0].typ, itemError)
		assert.Equal(t, items[1].typ, itemError)
	})




}