// Code generated by goa v3.19.1, DO NOT EDIT.
//
// Code Generator
//
// Command:
// goa

package main

import (
	_ "admin-panel/design"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/generator"
	"goa.design/goa/v3/eval"
	goa "goa.design/goa/v3/pkg"
)

func main() {
	var (
		out     = flag.String("output", "", "")
		version = flag.String("version", "", "")
		cmdl    = flag.String("cmd", "", "")
		ver     int
	)
	{
		flag.Parse()
		if *out == "" {
			fail("missing output flag")
		}
		if *version == "" {
			fail("missing version flag")
		}
		if *cmdl == "" {
			fail("missing cmd flag")
		}
		v, err := strconv.Atoi(*version)
		if err != nil {
			fail("invalid version %s", *version)
		}
		ver = v
	}

	if ver > goa.Major {
		fail("cannot run goa %s on design using goa v%s\n", goa.Version(), *version)
	}
	if err := eval.Context.Errors; err != nil {
		fail(err.Error())
	}
	if err := eval.RunDSL(); err != nil {
		fail(err.Error())
	}
	codegen.DesignVersion = ver
	outputs, err := generator.Generate(*out, "gen")
	if err != nil {
		fail(err.Error())
	}

	fmt.Println(strings.Join(outputs, "\n"))
}

func fail(msg string, vals ...any) {
	fmt.Fprintf(os.Stderr, msg, vals...)
	os.Exit(1)
}
