package k6deps_test

import (
	"fmt"

	"github.com/grafana/k6deps"
)

const script_without_pragma = `
import { Faker } from "k6/x/faker";
import sql from "k6/x/sql";
import driver from "k6/x/sql/driver/ramsql";

export default function() {
}
`

func ExampleAnalyze_without_pragma() {
	deps, _ := k6deps.Analyze(&k6deps.Options{
		Script: k6deps.Source{
			Name:     "script.js",
			Contents: []byte(script_without_pragma),
		},
		// disable automatic source detection
		Manifest: k6deps.Source{Ignore: true},
		Env:      k6deps.Source{Ignore: true},
	})

	fmt.Println(deps.String())

	out, _ := deps.MarshalJSON()
	fmt.Println(string(out))
	// Output:
	// k6/x/faker*;k6/x/sql*;k6/x/sql/driver/ramsql*
	// {"k6/x/faker":"*","k6/x/sql":"*","k6/x/sql/driver/ramsql":"*"}
}
