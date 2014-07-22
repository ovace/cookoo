/* The Cookoo `fmt` package provides utility wrappers for formatting text.
*/
package fmt

import (
	"github.com/Masterminds/cookoo"
	"text/template"
	"crypto/md5"
	"bytes"
	"fmt"
)

// Template is a template-based text formatter.
//
// This uses the core `text/template` to process a given string template.
//
// Params
// 	- template (string): A template string.
// 	- template.Context (bool): If true, the context will be placed into the
// 		template renderer as 'Cxt', and can be used as `{{.Cxt.Foo}}`. False
// 		by default.
// 	- ... (interface{}): Values passed into the template.
//
// Conventionally, template variables should start with an initial capital.
//
// Returns a formatted string.
func Template(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	format := cookoo.GetString("template", "", p)
	withCxt := cookoo.GetBool("template.Context", false, p)

	name := fmt.Sprintf("%x", md5.Sum([]byte(format)))

	//c.Logf("debug", "Template %s is '%s'\n", name, format)

	tpl, err := template.New(name).Parse(format)
	if err != nil {
		return "", err
	}

	data := p.AsMap()
	if withCxt {
		//c.Logf("debug", "Adding context.")
		data["Cxt"] = c.AsMap()
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, data); err != nil {
		return "", err
	}

	return out.String(), nil
}

// Printf formats a string and then sends it to standard out.
//
// This is a command wrapper for the core `fmt.Printf` function, but tooled
// to work the Cookoo way.
//
// The following prints 'Hello' to standard out.
//
// 	//...
// 	Does(Printf, "out").
// 	Using("format").WithDefault("%s\n")
// 	Using("0").WithDefault("Hello")
//
// Params:
// 	- format (string): The format string
// 	- "0"... (string): String representation of an integer ascending from 0.
// 	  These are treated as positional.
func Printf(c cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	msg := params.Get("format", "").(string)

	maxP := len(params.AsMap())
	vals := make([]interface{}, 0, maxP - 1)

	var istr string
	var i = 0
	for i < maxP {
		istr = fmt.Sprintf("%d", i) // FIXME
		if v, ok := params.Has(istr); ok {
			//fmt.Printf("%d: Found %v\n", i, v)
			vals = append(vals, v)
			i++
		} else {
			break
		}
	}
	fmt.Printf(msg, vals...)
	return true, nil
}
// Sprintf formats a string and then returns the result to the context.
//
// This is a command wrapper for the core `fmt.Sprintf` function, but tooled
// to work the Cookoo way.
//
// The following returns 'Hello World' to the context.
//
// 	//...
// 	Does(Sprintf, "out").
// 	Using("format").WithDefault("%s %s\n")
// 	Using("0").WithDefault("Hello")
// 	Using("1").WithDefault("World")
//
// Params:
// 	- format (string): The format string
// 	- "0"... (string): String representation of an integer ascending from 0.
// 	  These are treated as positional.
func Sprintf(c cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	msg := params.Get("format", "").(string)

	maxP := len(params.AsMap())
	vals := make([]interface{}, 0, maxP - 1)

	var istr string
	var i = 0
	for i < maxP {
		istr = fmt.Sprintf("%d", i) // FIXME
		if v, ok := params.Has(istr); ok {
			//fmt.Printf("%d: Found %v\n", i, v)
			vals = append(vals, v)
			i++
		} else {
			break
		}
	}

	return fmt.Sprintf(msg, vals...), nil
}

