//** This file was code generated by got. ***

package testOut

import (
	"bytes"
)

func anIncludeTest(buf bytes.Buffer) {

	buf.WriteString(`	Print me.
`)

}

func TestInclude(buf *bytes.Buffer) {

	buf.WriteString(`The end.
`)

	buf.WriteString(`    This is a named block.
`)

	// a test of substituting a name
	{
		var smallBlock string
		_ = smallBlock
	}

	buf.WriteString(`	<p>
	`)

	buf.WriteString(`Escaped html &lt;`)

	buf.WriteString(`
	</p>

`)

	buf.WriteString(`<p>This is text<br>
that is both escaped and has<br>
html paragraphs and breaks inserted.<br>
</p>
`)

	buf.WriteString(`
`)

	buf.WriteString(``)

	buf.WriteString(`
`)

	buf.WriteString(`<html>
<head>

</head>
<body>
<p>
    A typical html document
</p>
</body>
</html>`)

	buf.WriteString(`
`)

}
