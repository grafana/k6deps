// Package cmd contains k6deps cobra command factory function.
package cmd

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/k6deps"
	"github.com/spf13/cobra"
)

//go:generate go run github.com/dmarkham/enumer@v1.5.9 -type=format -transform=lower -trimprefix format -output format_gen.go

type format int

func (f *format) Set(v string) error {
	var err error
	*f, err = formatString(v)
	return err
}

func (f format) Type() string {
	return strings.Join(formatStrings(), "|")
}

const (
	formatJSON format = iota
	formatText
	formatJS
)

type options struct {
	k6deps.Options
	input  string
	format format
	output string
}

//go:embed help.md
var help string

// New creates new cobra command for deps command.
func New() *cobra.Command {
	opts := new(options)

	cmd := &cobra.Command{
		Use:   "k6deps [flags] [script-file]",
		Short: "Extension dependency detection for k6.",
		Long:  help,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return deps(opts, args)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	flags := cmd.Flags()

	flags.StringVar(&opts.Manifest.Name, "manifest", "",
		"manifest file to analyze (default 'package.json' nearest to script-file)")
	flags.Var(&opts.format, "format", "output format, possible values: json,env,script")
	flags.StringVarP(&opts.output, "output", "o", "", "write output to file (default stdout)")
	flags.BoolVar(&opts.Env.Ignore, "ingnore-env", false,
		"ignore "+k6deps.EnvDependencies+" environment variable processing")
	flags.BoolVar(&opts.Manifest.Ignore, "ignore-manifest", false, "disable package.json detection and processing")
	flags.BoolVar(&opts.Script.Ignore, "ignore-script", false, "disable script processing")
	flags.StringVarP(&opts.input, "input", "i", "", "input format ('js', 'ts' or 'tar' for archives)")
	return cmd
}

func deps(opts *options, args []string) error {
	var ignoreStdin bool

	if len(args) > 0 {
		filename := args[0]
		switch filepath.Ext(filename) {
		case ".js", ".ts":
			opts.Script.Name = filename
		case ".tar":
			opts.Archive.Name = filename
		default:
			return fmt.Errorf("unsupported file extension: %s", filepath.Ext(filename))
		}
		ignoreStdin = true
	}

	if opts.input != "" && !ignoreStdin {
		switch opts.input {
		case "js", "ts":
			buffer := &bytes.Buffer{}
			buffer.ReadFrom(os.Stdin) //nolint:errcheck,forbidigo,gosec
			opts.Script.Name = "stdin"
			opts.Script.Contents = buffer.Bytes()
		case "tar":
			opts.Archive.Name = "stdin"
			opts.Archive.Reader = os.Stdin //nolint:forbidigo
		default:
			return fmt.Errorf("unsupported input format: %s", opts.input)
		}
	}

	var out io.Writer

	if len(opts.output) == 0 {
		out = os.Stdout //nolint:forbidigo
	} else {
		file, err := os.Create(filepath.Clean(opts.output)) //nolint:forbidigo
		if err != nil {
			return err
		}

		defer file.Close() //nolint:errcheck

		out = file
	}

	deps, err := k6deps.Analyze(&opts.Options)
	if err != nil {
		return err
	}

	return printDependencies(deps, out, opts.format)
}

func printDependencies(deps k6deps.Dependencies, out io.Writer, outFormat format) error {
	switch outFormat {
	case formatText:
		text, err := deps.MarshalText()
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(out, string(text))
		return err
	case formatJS:
		text, err := deps.MarshalJS()
		if err != nil {
			return err
		}

		_, err = out.Write(text)
		return err
	case formatJSON:
		fallthrough
	default:
		encoder := json.NewEncoder(out)
		encoder.SetEscapeHTML(false)

		return encoder.Encode(deps)
	}
}
