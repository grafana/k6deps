// Package file contains esbuild file plugin.
package file

import (
	"errors"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

const (
	pluginName = "file"
)

var (
	errNotFound = errors.New("module not found")

	loaderByExtension = map[string]api.Loader{ //nolint:gochecknoglobals
		".js":   api.LoaderJS,
		".json": api.LoaderJSON,
		".txt":  api.LoaderText,
		".ts":   api.LoaderTS,
	}
)

type plugin struct {
	fs      fs.FS
	resolve func(path string, options api.ResolveOptions) api.ResolveResult
	options *api.BuildOptions
}

// New creates new http plugin instance.
func New(fs fs.FS) api.Plugin {
	plugin := &plugin{
		fs: fs,
	}

	return api.Plugin{
		Name:  pluginName,
		Setup: plugin.setup,
	}
}

func (plugin *plugin) setup(build api.PluginBuild) {
	plugin.resolve = build.Resolve
	plugin.options = build.InitialOptions

	build.OnResolve(api.OnResolveOptions{Filter: ".*", Namespace: "file"}, plugin.onResolve)
	build.OnLoad(api.OnLoadOptions{Filter: ".*", Namespace: "file"}, plugin.onLoad)
}

func (plugin *plugin) load(path string) (*api.OnLoadResult, error) {
	path, err := filepath.Rel("/", path)
	if err != nil {
		return nil, err
	}

	file, err := plugin.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() //nolint:errcheck

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	contents := string(bytes)

	// FIXME: what happens if the extension is not supported?
	var loader api.Loader
	if ldr, ok := loaderByExtension[filepath.Ext(path)]; ok {
		loader = ldr
	}

	return &api.OnLoadResult{
			Contents:   &contents,
			PluginName: pluginName,
			Loader:     loader,
		},
		nil
}

func (plugin *plugin) onLoad(args api.OnLoadArgs) (api.OnLoadResult, error) {
	if res, ok := args.PluginData.(*api.OnLoadResult); ok {
		return *res, nil
	}

	return onLoadError(args, errNotFound)
}

func onResolveError(_ api.OnResolveArgs, err error) (api.OnResolveResult, error) {
	return api.OnResolveResult{}, err
}

func onLoadError(_ api.OnLoadArgs, err error) (api.OnLoadResult, error) {
	return api.OnLoadResult{}, err
}

func newOnResolveResult(path string, plugindata interface{}) api.OnResolveResult {
	return api.OnResolveResult{
		Namespace:  "file",
		Path:       path,
		PluginData: plugindata,
	}
}

func (plugin *plugin) onResolve(args api.OnResolveArgs) (api.OnResolveResult, error) {
	path := args.Path

	if !filepath.IsAbs(path) {
		path = filepath.Join(filepath.Dir(args.Importer), path)
	}

	loadResult, err := plugin.load(path)
	if err == nil {
		return newOnResolveResult(path, loadResult), nil
	}

	return onResolveError(args, err)
}
