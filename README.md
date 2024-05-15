<!-- The badges are prepared, it is worth displaying them when the repo becomes public.
[![Go Reference](https://pkg.go.dev/badge/github.com/grafana/k6deps.svg)](https://pkg.go.dev/github.com/grafana/k6deps)
[![GitHub Release](https://img.shields.io/github/v/release/grafana/k6deps)](https://github.com/grafana/k6deps/releases/)
[![Go Report Card](https://goreportcard.com/badge/github.com/grafana/k6deps)](https://goreportcard.com/report/github.com/grafana/k6deps)
[![GitHub Actions](https://github.com/grafana/k6deps/actions/workflows/test.yml/badge.svg)](https://github.com/grafana/k6deps/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/grafana/k6deps/graph/badge.svg?token=PCRNQE9LGQ)](https://codecov.io/gh/grafana/k6deps)
![GitHub Downloads](https://img.shields.io/github/downloads/grafana/k6deps/total)
-->

# k6deps

**Dependency analysis for k6 tests**

The goal of k6deps is to extract dependencies from k6 test scripts. For this purpose, k6deps analyzes the k6 test scripts and the modules imported from it in a recursive manner.

k6deps is primarily used as a go library for [k6](https://github.com/grafana/k6) and [xk6](https://github.com/grafana/xk6). In addition, it also contains a command-line tool, which is suitable for listing the dependencies of k6 test scripts.

The command line tool can be integrated into other command line tools (for example k6, xk6) as a subcommand. For this purpose, the library also contains the functionality of the command line tool as a factrory function that returns [cobra.Command](https://pkg.go.dev/github.com/spf13/cobra#Command).

## How It Works

The first step is to collect the import paths used by the k6 test script. This is followed by the normalization of the import paths. Those import paths that do not refer to k6 extensions after normalization are traversed recursively.

### Import path normalization

For various reasons, it may be necessary to change the import path, so the normalization process may overwrite the original import path. For example, import path as URL can optionally contain version constraints in the hash part of the URL. This is removed by the normalization process.

The normalized import path is used for recursive traversal of modules, if the normalization result does not refer to k6 extension.

As a side effect of normalization, the following information will be available:

- kind (JavaScript module, golang extension), based on the pattern matching of the module name
- version constraints (optional), parsed from the hash part of the import path URL

Parsing of k6 test scripts and recursive traversal of module imports is done using the esbuild library (bundling feature).

### Resolvers

Normalized import paths are resolved using different resolvers. The resolvers are selected by pattern matching on the import path.

#### HTTPS resolver

An import path starting with `https://` is a remote JavaScript module reference and the normalized import path also specifies the location of the module.

#### Local resolver

An import path starting with the character `.` is a local JavaScript module reference, in which case the normalized import path also specifies the location of the module.

#### Other resolvers

TBD (@scope/name, name, github:user/repo, gitlab:user/repo, bitbucket:user/repo, etc)

## Install

Precompiled binaries can be downloaded and installed from the [Releases](https://github.com/grafana/k6deps/releases) page.

If you have a go development environment, the installation can also be done with the following command:

```
go install github.com/grafana/k6deps/cmd/k6deps@latest
```

## Usage


## Development

### Tasks

This section contains a description of the tasks performed during development. If you have the [xc (Markdown defined task runner)](https://github.com/joerdav/xc) command-line tool, individual tasks can be executed simply by using the `xc task-name` command.

<details><summary>Click to expand</summary>

#### lint

Run the static analyzer.

```
golangci-lint run
```

#### test

Run the tests.

```
go test -count 1 -race -coverprofile=build/coverage.txt ./...
```

#### coverage

View the test coverage report.

```
go tool cover -html=build/coverage.txt
```

#### build

Build the executable binary.

This is the easiest way to create an executable binary (although the release process uses the goreleaser tool to create release versions).

```
go build -ldflags="-w -s" -o build/k6deps ./cmd/k6deps
```

#### snapshot

Creating an executable binary with a snapshot version.

The goreleaser command-line tool is used during the release process. During development, it is advisable to create binaries with the same tool from time to time.

```
goreleaser build --snapshot --clean --single-target -o build/k6deps
```

#### clean

Delete the build directory.

```
rm -rf build
```

#### all

Run all tasks.

Requires: lint,test,build,snapshot

</details>
