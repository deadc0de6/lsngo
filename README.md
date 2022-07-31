[![Tests Status](https://github.com/deadc0de6/lsngo/workflows/tests/badge.svg)](https://github.com/deadc0de6/lsngo/actions)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)

[![Donate](https://img.shields.io/badge/donate-KoFi-blue.svg)](https://ko-fi.com/deadc0de6)

# lsngo

[lsngo](https://github.com/deadc0de6/lsngo) is a terminal user interface to replace
the use of repetitive `ls ... cd ...`.

![](/resources/screenshot.png?raw=true "lsngo")

Install by picking up a binary from the [latest release](https://github.com/deadc0de6/checkah/releases) and adding it to your path.

Then add an alias:
```bash
alias lg=lsngo
```

# Usage

```bash
Usage of lg:
  -a	Show hidden files
  -debug
    	Debug mode
  -editor string
    	File editor
  -help
    	Show this help
  -l	Long format
  -version
    	Show version
```

# Build

```bash
## create a binary for your current host
go mod tidy
make
./lg --help

## create all architecture binaries
go mod tidy
make build-all
ls lg-*
```

# Shortcuts

* `j`: down (or arrow down)
* `k`: up (or arrow up)
* `h`: go to parent directory (or arrow left)
* `l`: open file/directory (or arrow right)
* `q`: exit
* `esc`: exit
* `H`: toggle hidden files
* `L`: toggle long format
* `enter`: open file/directory
* `?`: show help

# tview

This tools uses a modified version of [tview](https://github.com/rivo/tview)
with [this PR](https://github.com/rivo/tview/pull/745) applied.

# Thank you

If you like lsngo, [buy me a coffee](https://ko-fi.com/deadc0de6).

# License

This project is licensed under the terms of the GPLv3 license.
