# jqplay-cli

`jqplay-cli` is a fork of the wonderful
[jqplay](https://github.com/owenthereal/jqplay) with a focus on local command
line usage. Play with your JSON data without sending your private data to a
third party server.

![web-based](./demo/web.gif)

You can use the default web-based UI or a [complete terminal-based UI](./demo/term.gif),
all without hitting the internet.

## Quickstart

```sh
brew install jittering/kegs/jqplay
curl -s https://jsonplaceholder.typicode.com/todos/1 | jqplay
```

## Installation

via homebrew (mac or linux):

```sh
brew install jittering/kegs/jqplay
```

or manually:

Download a [pre-built binary](https://github.com/jittering/jqplay/releases) or
build it from source:

```sh
go get github.com/jittering/jqplay/cmd/jqplay
```

## Usage

```text
Usage of jqplay:

  -cli
    	CLI mode
  -no-open
    	Do not open browser on startup
  -verbose
    	Verbose output
  -web
    	Web mode (default true)
```

## License

jqplay is released under the MIT license. See [LICENSE.md](https://github.com/jingweno/jqplay/blob/master/LICENSE.md).
