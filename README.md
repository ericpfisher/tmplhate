# tmplhate
Let the .tmpl hate flow through you...

This was a "reinvent the wheel" project.  Heavily inspired by the wonderful tool [Gomplate](https://docs.gomplate.ca).

## Usage

`tmplhate` uses [`cobra`](https://github.com/spf13/cobra) under the hood to provide CLI functionality.

As with most Go applications, run it with `go run main.go`:

```bash
❯ go run main.go --help
Tool to generate .tmpls from values.

Usage:
  tmplhate [flags]

Flags:
      --config string   config file (default is $HOME/.tmplhate.yaml)
  -h, --help            help for tmplhate
  -t, --tmpl string     .tmpl location
  -l, --values string   values location
  -v, --version         print version info
```

`tmplhate` can also read from `STDIN`:

```bash
❯ export NAME=GitHub
❯ echo 'Hi {{ .Env.NAME }}' | go run main.go
Hi GitHub
```

## [License](https://github.com/ericpfisher/tmplhate/blob/main/LICENSE)

Make sure you observe the licenses of libraries included in this project, too.
