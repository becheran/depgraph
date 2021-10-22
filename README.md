# depgraph

[![Go Report Card][go-report-image]][go-report-url]
[![PRs Welcome][pr-welcome-image]][pr-welcome-url]
[![License][license-image]][license-url]
[![Go Reference](https://pkg.go.dev/badge/github.com/becheran/depgraph.svg)](https://pkg.go.dev/github.com/becheran/depgraph)

[license-url]: https://github.com/becheran/depgraph/blob/main/LICENSE
[license-image]: https://img.shields.io/badge/License-MIT-brightgreen.svg
[go-report-image]: https://goreportcard.com/badge/github.com/becheran/depgraph
[go-report-url]: https://goreportcard.com/report/github.com/becheran/depgraph
[pr-welcome-image]: https://img.shields.io/badge/PRs-welcome-brightgreen.svg
[pr-welcome-url]: https://github.com/becheran/depgraph/blob/main/CONTRIBUTING.md

Interactive dependency graph visualization tool for golang using the awesome [cytoscape](https://cytoscape.org/) graph visualizer.

![screenshot](./doc/screenshot-0.PNG)

## Install

Install via:

``` sh
go install github.com/becheran/depgraph@latest
```

## Quick Start

Run the *depgraph* command line tool next to you *go.mod* file of your project.

The first required argument is the path to the file or module which shall be observed. For example:

``` sh
depgraph ./main.go
```

or:

``` sh
depgraph github.com/becheran/depgraph 
```

Per default the frontend will be started on **http://localhost:3001**. The address can be changed via the *host* flag:

``` sh
depgraph -host=localhost8080 github.com/becheran/depgraph 
```
