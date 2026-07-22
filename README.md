# freeswitch_module_golang_sample
Sample module for [FreeSWITCH](https://github.com/signalwire/freeswitch) using golang


## tl; dr

```bash
git clone https://github.com/iuridiniz/freeswitch_module_golang_sample.git mod_hello_world
cd mod_hello_world
make && make install

fs_cli -x 'load mod_hello_world'
fs_cli -x 'hello my friend'
```

## Requirements

Working gcc, golang, make and freeswitch with dev files.

## Compiling

Just call `make`

```bash
make
```

Makefile will use a freeswitch compiled and installed in `/usr/local/freeswitch`, you can change by passing `FREESWITCH_DIR=/path/to/your/freeswitch` to `make`:

```bash
make FREESWITCH_DIR="/opt/freeswitch"
``` 

Also, this program will try to use `go` tool from your PATH, but you can change this by passing `GO_BINARY=/path/to/your/go` to `make`:

```bash
make GO_BINARY="/host/home/iuri/.local/opt/go-1.17.2.linux-amd64/bin/go"
```

## Install

```bash
make install
```

## Test

On fs_cli, call:

```cli
freeswitch@localhost> load mod_hello_world
freeswitch@localhost> hello golang
```

## File Descriptions

- **glue/glue.c**: C glue — `SWITCH_MODULE_DEFINITION` and the wrapper functions that call into Go.
- **glue/glue.h**: Header file for `glue.c`, declaring the C helper functions used by the glue layer.
- **glue/glue.go**: The only Go file using cgo. Contains all FreeSWITCH↔Go plumbing: the `//export`ed entry points, the `Module` interface, `Register()`, and the `Stream`/`Session`/`Log` wrappers passed to your module.
- **mod_hello_world.go**: Pure module logic — no cgo. Implements the `glue.Module` interface (`Load`/`Runtime`/`Shutdown`/`ApiHandler`) and registers itself via `glue.Register()`. Edit this file to build your own module.
