# LSNGO

[![Tests Status](https://github.com/deadc0de6/lsngo/workflows/tests/badge.svg)](https://github.com/deadc0de6/lsngo/actions)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)

[![Donate](https://img.shields.io/badge/donate-KoFi-blue.svg)](https://ko-fi.com/deadc0de6)

[lsngo](https://github.com/deadc0de6/lsngo) is a terminal user interface to replace
the use of repetitive `ls ... cd ...`.

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

# shortcuts

* `j`: down
* `k`: up
* `h`: go to parent directory
* `l`: open file/directory
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
