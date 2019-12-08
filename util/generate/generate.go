package generate

import (
	"fmt"
	"time"
)

// ParameterGenerate ...
type ParameterGenerate struct {
	OutFile  string
	PkgName  string
	FuncName string
}

func generator() string { return "gdk" }
func timestamp() string { return time.Now().Format(time.RFC3339) }
func die(err error) {
	if err != nil {
		panic(err)
	}
}
func log(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
