package template

import (
	"github.com/valyala/fasttemplate"
	"strings"
)

// Data is a wrapper
type Data map[string]interface{}

// Parse ...
func Parse(template string, data Data) string {
	for strings.Contains(template, "{{ ") {
		template = strings.Replace(template, "{{ ", "{{", -1)
	}
	for strings.Contains(template, " }}") {
		template = strings.Replace(template, " }}", "}}", -1)
	}
	t, err := fasttemplate.NewTemplate(template, "{{", "}}")
	if err != nil {
		panic(err)
	}
	return t.ExecuteString(data)
}
