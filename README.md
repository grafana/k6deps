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

## Install

Precompiled binaries can be downloaded and installed from the [Releases](https://github.com/grafana/k6deps/releases) page.

If you have a go development environment, the installation can also be done with the following command:

```
go install github.com/grafana/k6deps/cmd/k6deps@latest
```

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
