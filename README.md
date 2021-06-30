# climan

:smile: CLI manager library, heavily inspired by [ffcli](https://github.com/peterbourgon/ff).

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/moul.io/climan)
[![License](https://img.shields.io/badge/license-Apache--2.0%20%2F%20MIT-%2397ca00.svg)](https://github.com/moul/climan/blob/main/COPYRIGHT)
[![GitHub release](https://img.shields.io/github/release/moul/climan.svg)](https://github.com/moul/climan/releases)
[![Docker Metrics](https://images.microbadger.com/badges/image/moul/climan.svg)](https://microbadger.com/images/moul/climan)
[![Made by Manfred Touron](https://img.shields.io/badge/made%20by-Manfred%20Touron-blue.svg?style=flat)](https://manfred.life/)

[![Go](https://github.com/moul/climan/workflows/Go/badge.svg)](https://github.com/moul/climan/actions?query=workflow%3AGo)
[![Release](https://github.com/moul/climan/workflows/Release/badge.svg)](https://github.com/moul/climan/actions?query=workflow%3ARelease)
[![PR](https://github.com/moul/climan/workflows/PR/badge.svg)](https://github.com/moul/climan/actions?query=workflow%3APR)
[![GolangCI](https://golangci.com/badges/github.com/moul/climan.svg)](https://golangci.com/r/github.com/moul/climan)
[![codecov](https://codecov.io/gh/moul/climan/branch/main/graph/badge.svg)](https://codecov.io/gh/moul/climan)
[![Go Report Card](https://goreportcard.com/badge/moul.io/climan)](https://goreportcard.com/report/moul.io/climan)
[![CodeFactor](https://www.codefactor.io/repository/github/moul/climan/badge)](https://www.codefactor.io/repository/github/moul/climan)

[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/moul/climan)

This package is originally based on [peterbourgon's `ff` package](https://github.com/peterbourgon/ff) (Apache2 License).

It implements small changes that don't fit with the original's author Goals and Non-goals.

---

Changes include:

* Adding an optional `Command.FlagSetBuilder` callback to configure commands and subcommands dynamically and support sharing the same flag targets.
* Using `flag.ContinueOnError` by default instead of `flag.ExitOnError`.
* Printing usage instead of an returning an error if a command does not implements an `Exec` func.
* Use a different `DefaultUsageFunc`.

## Example

[embedmd]:# (example_test.go /import\ / $)
```go
import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"

    "moul.io/climan"
)

func Example() {
    var opts struct {
        Debug bool
    }

    root := &climan.Command{
        Name:       "example",
        ShortUsage: "example [global flags] <subcommand> [flags] [args...]",
        ShortHelp:  "example's short help",
        LongHelp:   "example's longer help.\nwith more details.",
        FlagsBuilder: func(fs *flag.FlagSet) {
            fs.BoolVar(&opts.Debug, "debug", opts.Debug, "debug mode")
        },
        Exec: func(ctx context.Context, args []string) error {
            fmt.Println("args", args)
            return nil
        },
        Subcommands: []*climan.Command{
            &climan.Command{
                Name: "sub",
            },
        },
        // Options: []climan.Option{climan.WithEnvVarPrefix("EXAMPLE")},
    }
    if err := root.Parse(os.Args[1:]); err != nil {
        log.Fatal(fmt.Errorf("parse error: %w", err))
    }

    if err := root.Run(context.Background()); err != nil {
        log.Fatal(fmt.Errorf("run error: %w", err))
    }
}
```

## Usage

[embedmd]:# (.tmp/godoc.txt txt /TYPES/ $)
```txt
TYPES

type Command struct {
    Name         string
    Exec         func(context.Context, []string) error
    FlagsBuilder func(fs *flag.FlagSet)
    Subcommands  []*Command
    ShortUsage   string
    ShortHelp    string
    LongHelp     string

    // Has unexported fields.
}

func (c *Command) Parse(args []string) error

func (c *Command) Run(ctx context.Context) error

```

## Install

### Using go

```sh
go get moul.io/climan
```

### Releases

See https://github.com/moul/climan/releases

## Contribute

![Contribute <3](https://raw.githubusercontent.com/moul/moul/main/contribute.gif)

I really welcome contributions.
Your input is the most precious material.
I'm well aware of that and I thank you in advance.
Everyone is encouraged to look at what they can do on their own scale;
no effort is too small.

Everything on contribution is sum up here: [CONTRIBUTING.md](./.github/CONTRIBUTING.md)

### Dev helpers

Pre-commit script for install: https://pre-commit.com

### Contributors ✨

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-2-orange.svg)](#contributors)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="http://manfred.life"><img src="https://avatars1.githubusercontent.com/u/94029?v=4" width="100px;" alt=""/><br /><sub><b>Manfred Touron</b></sub></a><br /><a href="#maintenance-moul" title="Maintenance">🚧</a> <a href="https://github.com/moul/climan/commits?author=moul" title="Documentation">📖</a> <a href="https://github.com/moul/climan/commits?author=moul" title="Tests">⚠️</a> <a href="https://github.com/moul/climan/commits?author=moul" title="Code">💻</a></td>
    <td align="center"><a href="https://manfred.life/moul-bot"><img src="https://avatars1.githubusercontent.com/u/41326314?v=4" width="100px;" alt=""/><br /><sub><b>moul-bot</b></sub></a><br /><a href="#maintenance-moul-bot" title="Maintenance">🚧</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors)
specification. Contributions of any kind welcome!

### Stargazers over time

[![Stargazers over time](https://starchart.cc/moul/climan.svg)](https://starchart.cc/moul/climan)

## License

© 2021   [Manfred Touron](https://manfred.life)

Licensed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)
([`LICENSE-APACHE`](LICENSE-APACHE)) or the [MIT license](https://opensource.org/licenses/MIT)
([`LICENSE-MIT`](LICENSE-MIT)), at your option.
See the [`COPYRIGHT`](COPYRIGHT) file for more details.

`SPDX-License-Identifier: (Apache-2.0 OR MIT)`
