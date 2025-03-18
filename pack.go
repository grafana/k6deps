package k6deps

import (
	"encoding/json"
	"path/filepath"
	"text/scanner"
	"time"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/grafana/k6deps/internal/pack/plugins/http"
	"github.com/grafana/k6deps/internal/pack/plugins/k6"
)

type metafile struct {
	K6 *k6.Metadata `json:"k6,omitempty"`
}

// Metadata holds k6 related metadata, emitted under "k6" key of Metafile.
type Metadata struct {
	// Imports contains a list of k6 imports (core modules and extensions).
	Imports []string `json:"imports,omitempty"`
}

type packError struct {
	messages []api.Message
}

// PackOptions used to specify transform/build options.
type PackOptions struct {
	Directory  string
	Filename   string
	Timeout    time.Duration
	TypeScript bool
	Externals  []string
	SourceRoot string
}

func (o *PackOptions) setDefaults() *PackOptions {
	if !o.TypeScript {
		o.TypeScript = filepath.Ext(o.Filename) == ".ts"
	}

	return o
}

func (o *PackOptions) stdinOptions(contents string) *api.StdinOptions {
	dir := filepath.Dir(o.Filename)
	base := filepath.Base(o.Filename)
	if base == "." { // empty filename = stdin
		base = ""
	}

	return &api.StdinOptions{
		Contents:   contents,
		Sourcefile: base,
		Loader:     o.loaderType(),
		ResolveDir: dir,
	}
}

func (o *PackOptions) loaderType() api.Loader {
	if o.TypeScript {
		return api.LoaderTS
	}

	return api.LoaderJS
}

// Pack gathers dependencies and transforms TypeScript/JavaScript sources into single k6 compatible JavaScript test
// script.
func Pack(source string, opts *PackOptions) ([]byte, *Metadata, error) {
	opts.setDefaults()

	result := api.Build(api.BuildOptions{ //nolint:exhaustruct
		Stdin:      opts.stdinOptions(source),
		Bundle:     true,
		LogLevel:   api.LogLevelSilent,
		Sourcemap:  api.SourceMapNone,
		SourceRoot: opts.SourceRoot,
		Plugins:    []api.Plugin{http.New(), k6.New()},
		External:   opts.Externals,
		Metafile:   true,
	})

	if has, err := checkError(&result); has {
		return nil, nil, err
	}

	var meta metafile

	err := json.Unmarshal([]byte(result.Metafile), &meta)
	if err != nil {
		return nil, nil, wrapError(err)
	}

	metadata := &Metadata{
		Imports: meta.K6.Imports,
	}

	return result.OutputFiles[0].Contents, metadata, nil
}

func checkError(result *api.BuildResult) (bool, error) {
	if len(result.Errors) == 0 {
		return false, nil
	}

	return true, &packError{messages: result.Errors}
}

func wrapError(err error) error {
	return &packError{[]api.Message{{Text: err.Error()}}}
}

func (e *packError) Error() string {
	msg := e.messages[0]

	if msg.Location == nil {
		return msg.Text
	}

	pos := scanner.Position{
		Filename: msg.Location.File,
		Line:     msg.Location.Line,
		Column:   msg.Location.Column,
	}

	return pos.String() + " " + msg.Text
}
