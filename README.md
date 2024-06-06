<!-- The badges are prepared, it is worth displaying them when the repo becomes public.
[![Go Reference](https://pkg.go.dev/badge/github.com/grafana/k6deps.svg)](https://pkg.go.dev/github.com/grafana/k6deps)
[![GitHub Release](https://img.shields.io/github/v/release/grafana/k6deps)](https://github.com/grafana/k6deps/releases/)
[![Go Report Card](https://goreportcard.com/badge/github.com/grafana/k6deps)](https://goreportcard.com/report/github.com/grafana/k6deps)
[![GitHub Actions](https://github.com/grafana/k6deps/actions/workflows/test.yml/badge.svg)](https://github.com/grafana/k6deps/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/grafana/k6deps/graph/badge.svg?token=PCRNQE9LGQ)](https://codecov.io/gh/grafana/k6deps)
![GitHub Downloads](https://img.shields.io/github/downloads/grafana/k6deps/total)
-->

<h1 name="title">k6deps</h1>

**Dependency analysis for k6 tests**

The goal of k6deps is to extract dependencies from k6 test scripts. For this purpose, k6deps analyzes the k6 test scripts and the modules imported from it in a recursive manner.

k6deps is primarily used as a go library for [k6](https://github.com/grafana/k6) and [xk6](https://github.com/grafana/xk6). In addition, it also contains a command-line tool, which is suitable for listing the dependencies of k6 test scripts.

The command line tool can be integrated into other command line tools as a subcommand. For this purpose, the library also contains the functionality of the command line tool as a factrory function that returns [cobra.Command](https://pkg.go.dev/github.com/spf13/cobra#Command).

## Install

Precompiled binaries can be downloaded and installed from the [Releases](https://github.com/grafana/k6deps/releases) page.

If you have a go development environment, the installation can also be done with the following command:

```
go install github.com/grafana/k6deps/cmd/k6deps@latest
```

## Usage

<!-- #region cli -->
## k6deps

Extension dependency detection for k6.

### Synopsis

Analyze the k6 test script and extract the extensions that the script depends on.

**Sources**

Dependencies can come from three sources: k6 test script, manifest file, `K6_DEPENDENCIES` environment variable.

Primarily, the k6 test script is the source of dependencies. The test script and the local and remote JavaScript modules it uses are recursively analyzed. The extensions used by the test script are collected. In addition to the require function and import expression, the `"use k6 ..."` directive can be used to specify additional extension dependencies. If necessary, the `"use k6 ..."` directive can also be used to specify version constraints.

       "use k6>0.49";
       "use k6 with k6/x/faker>=0.2.0";
       "use k6 with k6/x/toml>v0.1.0";
       "use k6 with xk6-dashboard*";

Dependencies and version constraints can also be specified in the so-called manifest file. The default name of the manifest file is `package.json` and it is automatically searched from the directory containing the test script up to the root directory. The `dependencies` property of the manifest file contains the dependencies in JSON format.

    {"dependencies":{"k6":">0.49","k6/x/faker":">=0.2.0","k6/x/toml":>v0.1.0","xk6-dashboard":"*"}}

Dependencies and version constraints can also be specified in the `K6_DEPENDENCIES` environment variable. The value of the variable is a list of dependencies in a one-line text format.

       k6>0.49;k6/x/faker>=0.2.0;k6/x/toml>v0.1.0;xk6-dashboard*

**Format**

By default, dependencies are written as a JSON object. The property name is the name of the dependency and the property value is the version constraints of the dependency.

    {"k6":">0.49","k6/x/faker":">=0.2.0","k6/x/toml":>v0.1.0","xk6-dashboard":"*"}

Additional output formats:

 * `text` - One line text format. A semicolon-separated sequence of the text format of each dependency. The first element of the series is `k6` (if there is one), the following elements follow each other in lexically increasing order based on the name.

       k6>0.49;k6/x/faker>=0.2.0;k6/x/toml>v0.1.0;xk6-dashboard*

 * `js` - A consecutive, one-line JavaScript string directives. The first element of the series is `k6` (if there is one), the following elements follow each other in lexically increasing order based on the name.

       "use k6>0.49";
       "use k6 with k6/x/faker>=0.2.0";
       "use k6 with k6/x/toml>v0.1.0";
       "use k6 with xk6-dashboard*";

**Output**

By default, dependencies are written to standard output. By using the `-o/--output` flag, the dependencies can be written to a file.

**Arguments**

The only (optional) argument of the command is the name of the k6 test script file.

```
k6deps [flags] [script-file]
```

### Flags

```
      --format string     output format, possible values: json,env,script (default "json")
  -h, --help              help for k6deps
      --ignore-manifest   disable package.json detection and processing
      --ignore-script     disable script processing
      --ingnore-env       ignore K6_DEPENDENCIES environment variable processing
      --manifest string   manifest file to analyze (default 'package.json' nearest to script-file)
  -o, --output string     write output to file (default stdout)
```

<!-- #endregion cli -->

## Development

### Tasks

This section contains a description of the tasks performed during development. If you have the [xc (Markdown defined task runner)](https://github.com/joerdav/xc) command-line tool, individual tasks can be executed simply by using the `xc task-name` command.

<details><summary>Click to expand</summary>

#### readme

Update documentation in README.md.

```
go run ./tools/gendoc README.md
```

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

</details>
