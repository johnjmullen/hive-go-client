# hive-go-client
go client for hive-rest

Since this is a priate repo you need to run
`git config --global url."git@github.com:".insteadOf "https://github.com/"`

install it with: `go get github.com/hive-io/hive-go-client/hioctl`

```
 ~/go/bin/hioctl
hiveio cli

Usage:
  hioctl [command]

Available Commands:
  alert       alert operations
  cluster     cluster operations
  guest       guest operations
  help        Help about any command
  host        host operations
  pool        pool operations
  profile     profile operations
  realm       realm operations
  storage     Storage
  template    template operations

Flags:
      --config string     config file
      --format string     format (json/yaml) (default "json")
  -h, --help              help for hioctl
      --host string       Server to connect to
  -k, --insecure          ignore certificate errors
  -p, --password string   Admin user password
      --port uint         port (default 8443)
  -r, --realm string      Admin user realm (default "local")
  -u, --user string       Admin username (default "admin")

Use "hioctl [command] --help" for more information about a command.
```

