# test-tcp

This my simple implemetation for tcp connection with Golang.
Program can be run in two modes: client and server. Parameters for connect used from TOML config file.

## Build

    $ make build

## Config

Configuration should be writen in TOML format. Examples places in `data` directory.

Location for config file can be handled from command line with `-c` key.

Example for TOML-file:

    host = '127.0.0.1'
    port = 1514

## Using

Help message with all flags:

    $ ./test-tcp -h
    usage: test-tcp [<flags>]

    Flags:
      -h, --help                  Show context-sensitive help (also try --help-long and --help-man).
      -c, --config="config.toml"  Config file in TOML format. Used only ip and port variables.
      -l, --log="logfile"         Log file name.
      -m, --mode="server"         Client/Server mode. Example: -m "server"
      -v, --version               Show application version.

By default program used `config.toml` file from current directory.

Server mode used by default. For run in client mode, please use:

    $ ./test-tcp -m client
