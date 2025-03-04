[![Go Reference](https://pkg.go.dev/badge/github.com/goradd/got.svg)](https://pkg.go.dev/github.com/goradd/got)
![Build Status](https://img.shields.io/github/actions/workflow/status/goradd/got/go.yml?branch=main)
[![codecov](https://codecov.io/gh/goradd/got/branch/main/graph/badge.svg?token=FHU0NR2N1V)](https://codecov.io/gh/goradd/got)
[![Go Report Card](https://goreportcard.com/badge/github.com/goradd/got)](https://goreportcard.com/report/github.com/goradd/got)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go#template-engines)
# GoT

GoT (short for Go Templates) is a flexible template engine that generates Go code.

This approach creates extremely fast templates. It also gives you much more freedom than Go's template
engine, since at any time you can just switch to Go code to do what you want.

GoT has:
- Mustache-like syntax similar to Go's built-in template engine
- The ability to define new tags, so you can create your own template syntax
- Include files, so you can create a library of tags and templates 

The GoT syntax is easy to learn. Get started in minutes.

- [Features](#features)
- [Installation](#installation)
- [Command Line Usage](#command-line-usage)
- [Basic Syntax](#basic-syntax)
- [Template Syntax](#template-syntax)
- [License](#license)
- [Acknowledgements](#acknowledgements)

## Features

- **High performance**. Since the resulting template is go code, your template will be compiled to fast
machine code.
- **Easy to use**. The templates themselves are embedded into your go code. The template language is pretty 
simple and you can do a lot with only a few tags. You can switch into and out of go code at will. Tags are 
Mustache-like, so similar to Go's template engine.
- **Flexible**. The template language makes very few assumptions about the go environment it is in. Most other
template engines require you to call the template with a specific function signature. **GoT** gives you the
freedom to call your templates how you want.
- **Translation Support**. You can specify that you want to send your strings to a translator before 
output.
- **Error Support**. You can call into go code that returns errors, and have the template stop at that
point and return an error to your wrapper function. The template will output its text up to that point,
allowing you to easily see where in the template the error occurred.
- **Include Files**. Templates can include other templates. You can specify
a list of search directories for the include files, allowing you to put include files in a variety of
locations, and have include files in one directory that override another directory.
- **Custom Tags**. You can define named fragments that you can include at will, 
and you can define these fragments in include files. To use these fragments, you just
use the name of the fragment as the tag. This essentially gives you the ability to create your own template
language. When you use a custom tag, you can also include parameters that will 
replace placeholders in your fragment, giving you even more power to your custom tags. 
- **Error Reporting**. Errors in your template files are identified by line and character
number. No need to guess where the error is.

Using other go libraries, you can have your templates compile when they are changed, 
use buffer pools to increase performance, and more. Since the
templates become go code, you can do what you want.


## Installation

```shell
go install github.com/goradd/got/got@latest
```

We also recommend installing `goimports` 
and passing GoT the -i flag on the command line. That will format your Go code and 
fix up the import lines of any generated go files.
```shell
go install golang.org/x/tools/cmd/goimports@latest
```

## Command Line Usage

```shell
got [options] [files]

options:
	- o: The output directory. If not specified, files will be output at the same 
	     location as the corresponding template. If using the -r option, files will
	     be output in subdirectories matching the directory names of the input files.
	- t  fileType: If set, will process all files in the current directory with this suffix. 
	     If not set, you must specify the files at the end of the command line.
	- i: Run `goimports` on the output files, rather than `go fmt`
	- I  directories and/or files:  A list of directories and/or files. 
	     If a directory, it is used as the search path for include files. 
	     If a file, it is automatically added to the front of every file that is processed.  
	     Directories are searched in the order specified and first matching file will be used. 
	     It will always look in the current directory last unless the current directory 
	     is specified in the list in another location. Relative paths must start with a dot (.) 
	     or double-dot (..).  Directories can start with a module name, and based on the 
	     current directory, the correct go.mod file will be searched to know where to look 
	     for include files.
	- d  directory: When using the -t option, will specify a directory to search.
	- v  verbose: Prints information about files being processed
	- r  recursive: Recursively processes directories. Used with the -t option and possibly -d.
	- f  force: Output files are normally not over-written if they are newer than the input file.
	     This option will force all input files to over-write the output files.
```
If a path described above starts with a module path, the actual disk location 
will be substituted.

examples:
```shell
	got -t got -i -o ../templates
	got -I .;../tmpl;example.com/projectTemplates file1.tmpl file2.tmpl
```

## Basic Syntax
Template tags start with `{{` and end with `}}`.

A template starts in go mode. To send simple text or html to output, surround the text with `{{` and `}}` tags
with a space or newline separating the tags from the surrounding text. Inside the brackets you will be in 
text mode.

In the resulting Go code, text will get written to output by calling:

 ```_, err = io.WriteString(_w, <text>)``` 

GoT assumes that the `_w` variable
is available and satisfies the io.Writer interface
and optionally the io.StringWriter interface.
Usually you would do this by declaring a function at the top of your template that looks like this:

``` func f(_w io.Writer) (err error) ```

After compiling the template output together with your program, you call
this function to get the template output. 

At a minimum, you will need to import the "io" package into the file with your template function.
Depending on what tags you use, you might need to add 
additional items to your import list. Those are mentioned below with each tag.

### Example
Here is how you might create a very basic template. For purposes of this example, we will call the file
`example.got` and put it in the `template` package, but you can name the file and package whatever you want.
```
package template

import "io"

func OutTemplate(_w io.Writer) (err error) {
	var world string = "World"
{{
<p>
    Hello {{world}}!
</p>
}}
  return // make sure the error gets returned
}
```

To compile this template, call GoT:
```shell
got example.got
```

This will create an `example.go` file, which you include in your program. You then can call the function
you declared:
```go
package main

import (
	"io"
	"os"
	"mypath/template"
)

func main() {
	_ = template.OutTemplate(os.Stdout)
}
```

This simple example shows a mix of go code and template syntax in the same file. Using GoT's include files,
you can separate your go code from template code.

## Template Syntax

The following describes how the various open tags work. Most tags end with a ` }}`, unless otherwise indicated.
Many tags have a short and a long form. Using the long form does not impact performance, its just there
to help your templates have some human readable context to them.

### Static Text
    {{<space or newline>   Begin to output text as written.
    {{! or {{esc           Html escape the text. Html reserved characters, like < or > are 
                           turned into html entities first. This happens when the template is 
                           compiled, so that when the template runs, the string will already 
                           be escaped. 
    {{h or {{html          Html escape and html format double-newlines into <p> tags.
    {{t or {{translate     Send the text to a translator


The `{{!` tag requires you to import the standard html package. `{{h` requires both the html and strings packages.

`{{t` will wrap the static text with a call to t.Translate(). Its up to you to define this object and make it available to the template. The
translation will happen during runtime of your program. We hope that a future implementation of GoT could
have an option to send these strings to an i18n file to make it easy to send these to a translation service.

#### Example
In this example file, note that we start in Go mode, copying the text verbatim to the template file.
```
package test

import (
	"io"
	"fmt"
)

type Translater interface {
	Translate(string) string
}

func staticTest(_w io.Writer) {
{{
<p>
{{! Escaped html < }}
</p>
}}

{{h
	This is text that is both escaped.
	 
	And has html paragraphs inserted.
}}

}

func translateTest(t Translater, buf *bytes.Buffer) {

{{t Translate me to some language }}

{{!t Translate & escape me > }}

}
```

### Switching Between Go Mode and Template Mode
From within any static text context described above you can switch into go context by using:

    {{g or {{go     Change to straight go code.
    
Go code is copied verbatim to the final template. Use it to set up loops, call special processing functions, etc.
End go mode using the `}}` closing tag. You can also include any other GoT tag inside of Go mode,
meaning you can nest Go mode and all the other template tags.

#### Example
```
// We start a template in Go mode. The next tag switches to text mode, and then nests
// switching back to go mode.
{{ Here 
{{go 
io.WriteString(_w, "is") 
}} some code wrapping text escaping to go. }}
```

### Dynamic Text
The following tags are designed to surround go code that returns a go value. The value will be 
converted to a string and sent to the buf. The go code could be a static value, or a function
that returns a value.

    Tag                       Description                     Example
    
    {{=, {{s, or  {{string    A go string                     {{= fmt.Sprintf("%9.2f", fVar) }}
    {{i or {{int              An int                          {{ The value is: {{i iVar }} }}
    {{u or {{uint             A uint                          {{ The value is: {{u uVar }} }}
    {{f or {{float            A floating point number         {{ The value is: {{f fVar }} }}
    {{b or {{bool             A boolean ("true" or "false")   {{ The value is: {{b bVar }} }}
    {{w or {{bytes            A byte slice                    {{ The value is: {{w byteSliceVar }} }}
    {{v or {{stringer or      Any value that implements       {{ The value is: {{objVar}} }}
       {{goIdentifier}}       the Stringer interface.

The last tag can be slower than the other tags since they use fmt.Sprint() internally, 
so if this is a heavily used template, avoid it. Usually you will not notice a speed difference though,
and the third option can be very convenient. This third option is simply any go variable surrounded by mustaches 
with no spaces.

The i, u, and f tags use the strconv package, so be sure to include that in your template.

#### Escaping Dynamic Text

Some value types potentially could produce html reserved characters. The following tags will html escape
the output.

    {{!=, {{!s or {{!string    HTML escape a go string
    {{!w or {{!bytes           HTML escape a byte slice
    {{!v or {{!stringer        HTML escape a Stringer
    {{!h                       Escape a go string and html format breaks and newlines

These tags require you to import the "html" package. The `{{!h` tag also requires the "strings" package.

#### Capturing Errors

These tags will receive two results, the first a value to send to output, and the second an error
type. If the error is not nil, processing will stop and the error will be returned by the template function. Therefore, these
tags expect to be included in a function that returns an error. Any template text
processed so far will still be sent to the output buffer.

    {{=e, {{se, {{string,err      Output a go string, capturing an error
    {{!=e, {{!se, {{!string,err   HTML escape a go string and capture an error
    {{ie or {{int,err             Output a go int and capture an error
    {{ue or {{uint,err            Output a go uint and capture an error
    {{fe or {{float,err           Output a go float64 and capture an error
    {{be or {{bool,err            Output a bool ("true" or "false") and capture an error
    {{we, {{bytes,err             Output a byte slice and capture an error
    {{!we or {{!bytes,err         HTML escape a byte slice and capture an error
    {{ve, {{stringer,err          Output a Stringer and capture an error
    {{!ve or {{!stringer,err      HTML escape a Stringer and capture an error
    {{e, or {{err                 Execute go code that returns an error, and stop if the error is not nil

##### Example
```go
func Tester(s string) (out string, err error) {
	if s == "bad" {
		err = errors.New("This is bad.")
	}
	return s
}

func OutTemplate(toPrint string, buf bytes.Buffer) error {
{{=e Tester(toPrint) }}
}
```

### Generating Go Code

These tags are useful if your templates will be generating Go code.

    {{L      The value as a Go literal (uses %#v from fmt package)
    {{T      The Go type of the value without package name (uses %T from the fmt package)
    {{PT     The Go type of the value WITH the package name (uses %T from the fmt package)

The {{L tag will quote strings, but not quote numbers.

The {{PT tag will produce the Go type with a package name if given. Built in types will not
have a package name, but standard library types will. For example, a time.Time type will 
result in the string "time.Time". The {{T tag will strip off the package name if there is one
in the value's type.

### Include Files
#### Include a GoT source file

    {{: "fileName" }} or 
    {{include "fileName" }}   Inserts the given file name into the template.

 The included file will start in whatever mode the receiving template is in, as if the text was inserted
 at that spot, so if the include tags are  put inside of go code, the included file will start in go mode. 
 The file will then be processed like any other GoT file. Include files can refer to other include files,
 and so are recursive.
 
 Include files are searched for in the current directory, and in the list of include directories provided
 on the command line by the -I option.

Example: `{{: "myTemplate.inc" }}`

#### Include a text file

    {{:! "fileName" }} or             Inserts the given file name into the template
    {{includeEscaped "fileName" }}    and html escapes it.
    {{:h "fileName" }} or             Inserts the given file name into the template,
    {{includeAsHtml "fileName" }}     html escapes it, and converts newlines into html breaks.

Use `{{:!` to include a file that you surround with a `<pre>` tag to include a text file
and have it appear in an html document looking the same. Use `{{:h` to include a file
without the `<pre>` tags, but if the file uses extra spaces for indent, those spaces will
not indent in the html. These kinds of include files will not be searched for GoT commands.
 
### Defined Fragments

Defined fragments start a block of text that can be included later in a template. The included text will
be sent as is, and then processed in whatever mode the template processor is in, as if that text was simply
inserted into the template at that spot. You can include the `{{` or `{{g` tags inside of the fragment to
force the processor into the text or go modes if needed. The fragment can be defined
any time before it is included, including being defined in other include files. 

You can add optional parameters
to a fragment that will be substituted for placeholders when the fragment is used. You can have up to 9
placeholders ($1 - $9). Parameters should be separated by commas, and can be surrounded by quotes if needed
to have a parameter that has a quote or comma in it.

    {{< fragName }} or {{define fragName }}   Start a block called "fragName".
    {{< fragName <count>}} or                 Start a block called "fragName" that will 
       {{define fragName <count>}}            have <count> parameters.
    {{> fragName param1,param2,...}} or       Substitute this tag for the given defined fragment.
      {{put fragName param1,param2,...}} or   
      {{fragName param1,param2,...}}
    {{>? fragName param1,param2,...}} or      Substitute this tag for the given defined fragment, 
      {{put? fragName param1,param2,...}}     but if the fragment is not defined, leave blank.

 
If you attempt to use a fragment that was not previously defined, GoT will panic and stop compiling,
unless you use {{>? or {{put? to include the fragment.

param1, param2, ... are optional parameters that will be substituted for $1, $2, ... in the defined fragment.
If a parameter is not included when using a fragment, an empty value will be substituted for the parameter in the fragment.

The fragment name is NOT surrounded by quotes, and cannot contain any whitespace in the name. Blocks are ended with a
`{{end fragName}}` tag. The end tag must be just like that, with no spaces after the fragName.

The following fragments are predefined:
* `{{templatePath}}` is the full path of the template file being processed
* `{{templateName}}` is the base name of the template file being processed, including any extensions
* `{{templateRoot}}` is the base name of the template file being processed without any extensions
* `{{templateParent}}` is the directory name of the template file being processed, without the preceeding path
* `{{includePath}}` is the full path of the file that the includePath fragment appears in.
* `{{includeName}}` is the base name of the file that the includeName fragment appears in.
* `{{includeRoot}}` is the base name of the file that the includeRoot fragment appears in without any extensions.
* `{{includeParent}}` is the directory name of the file that the includeParent fragment appears in without any extensions.
* `{{outPath}}` is the full path of the output file being written
* `{{outName}}` is the base name of the output file being written, including any extensions
* `{{outRoot}}` is the base name of the output file being written without any extensions
* `{{outParent}}` is the directory name of the output file being written, without the preceeding path
* `{{importPath}}` produces the import path of the output file
* `{{importParent}}` produces the parent of the import path of the output file

The template* fragments refer to the top level file, the one sent to the got command to process,
while the include* fragments refer to the file being processed currently. For example, if the
"a.tmpl" file included the "b.tmpl" file, then {{includeName}} in the b.tmpl file would
produce "b.tmpl", while {{templateName}} in the b.tmpl file would produce "a.tmpl".

importPath and importParent are used to generate import paths to files that are relative
to the output file, since go does not support relative import paths.

#### Example
```

{{< hFrag }}
<p>
This is my html body.
</p>
{{end hFrag}}

{{< writeMe 2}}
{{// The g tag here forces us to process the text as go code, 
     no matter where the fragment is included }}
{{g 
if "$2" != "" {
	io.WriteString(_w, "$1")
}
}}
{{end writeMe}}


func OutTemplate(_w io.Writer) (err error) {
{{
	<html>
		<body>
			{{> hFrag }}
		</body>
	</html>
}}

{{writeMe "Help Me!", a}}
{{writeMe "Help Me!", }}
 return
}
```

### Comment Tags

    {{# or {{//       Comment the template. This is removed from the compiled template.

These tags and anything enclosed in them is removed from the compiled template.

### Go Block Tags
    
    {{if <go condition>}}<text block>{{if}}                      
                           A convenience tag for surrounding text with a go "if" statement.
    {{if <go condition>}}<text block>{{else}}<text block>{{if}}  
                           Go "if" and "else" statement.
    {{if <go condition>}}<text block>{{elseif <go condition>}}<text block>{{if}}    
                           Go "if" and "else if" statement.
    {{for <go condition and optional range statement>}}<text block>{{for}}                                   
                           A convenience tag for surrounding text with a go "for" statement.

These tags are substitutes for switching into GO mode and using a `for` or `if` statements. 
`<text block>` will be in text mode to begin with, so that whatever you put there
will be output, but you can switch to go mode if needed.

####Example

```
{{
{{for num,item := range items }}
<p>Item {{num}} is {{item}}</p>
{{for}}
}}
```
### Join Tags

    {{join <slice>, <string>}}<text block>{{join}}    Joins the items of a slice with a string.

Join will output the `<text block>` for each item of `<slice>`. Within `<text block>` the variable
`_i` will be an integer representing the index of the slice item, and `_j` will be the
item itself. `<text block>` starts in text mode, but you can put GoT commands in it. `<string>` will be output
between the output of each item, creating an effect similar to joining a slice of strings.


####Example

```
{{g
  items := []string{"a", "b", "c"}
}}
{{join items,", "}}
{{ {{_i}} = {{_j}} }}
{{join}}
```


### Strict Text Block Tag

From within most of the GoT tags, you can insert another GoT tag. GoT will be looking for these
as it processes text. If you would like to turn off GoT's processing to output text 
that is not processed by GoT, you can use:

    {{begin *endTag*}}   Starts a strict text block and turns off the GoT parser. 
    
One thing this is useful for is to use GoT to generate GoT code.
End the block with a `{{end *endTag*}}` tag, where `*endTag*` is whatever you specified in the begin tag. 
There can be no space between the endTag and the final brackets
The following example will output the entire second line of code with no changes, 
including all brackets:

```
{{begin mystrict}}
{{! This is verbatim code }}
{{< all included}}
{{end mystrict}}
```

## Bigger Example

In this example, we will combine multiple files. One, a traditional html template with a place to fill in
some body text. Another, a go function declaration that we will use when we want to draw the template.
The function will use a traditional web server io.Writer pattern, including the use of a context parameter.
Finally, there is an example main.go illustrating how our template function would be called.

### index.html

```html
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
    </head>

    <body>
{{# The tag below declares that we want to substitute a named 
    fragment that is declared elsewhere }}
{{put body }}
    </body>
</html>
```

### template.got

```
package main

import {
	"context"
	"bytes"
}


func writeTemplate(ctx context.Context, buf *bytes.Buffer) {

{{# Define the body that will be inserted into the template }}
{{< body }}
<p>
The caller is: {{=s ctx.Value("caller") }}
</p>
{{end body}}

{{# include the html template. Since the template is html, 
    we need to enter static text mode first }}
{{ 
{{include "index.html"}}
}}

}
```

### main.go

```html
type myHandler struct {}

func (h myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	ctx :=  context.WithValue(r.Context(), "caller", r.Referer())
	writeTemplate(ctx, w)	// call the GoT template
}


func main() {
	var r myHandler
	var err error

	*local = "0.0.0.0:8000"

	err = http.ListenAndServe(*local, r)

	if err != nil {
		fmt.Println(err)
	}
}
```

To compile the template:

```shell
got template.got
```

Build your application and go to `http://localhost:8000` in your browser, to see your results

## License
GoT is licensed under the MIT License.

## Acknowldgements

GoT was influenced by:

- [hero](https://github.com/shiyanhui/hero)
- [fasttemplate](https://github.com/valyala/fasttemplate)
- [Rob Pike's Lexing/Parsing Talk](https://www.youtube.com/watch?v=HxaD_trXwRE)

## Syntax Changes

###v0.10.0
This was a major rewrite with the following changes:
- defined fragments end with {{end fragName}} tags, rather than {{end}} tags
- {{else if ...}} is now {{elseif ...}}
- {{join }} tag will join items with a string
- The backup tag {{- has been removed
- Reorganized the lexer and parser to be easier to debug
- Added many more unit tests. Code coverage > 90%.
- The output is sent to an io.Writer called _w. This allows more flexible use of the templates, and the ability to wrap them with middleware