# hive-go-client

hive-go-client/rest is a wrapper around the hiveio rest api

hive-go-client/hioctl is a command line interface for the api
install it with: `go get -u github.com/hive-io/hive-go-client/hioctl`

```
hioctl
hiveio cli

Usage:
  hioctl [command]

Available Commands:
  alert       alert operations
  cluster     cluster operations
  export      export data
  guest       guest operations
  help        Help about any command
  host        host operations
  pool        pool operations
  profile     profile operations
  realm       realm operations
  storage     storage operations
  task        task operations
  template    template operations

Flags:
      --config string     config file
      --format string     format (json/yaml) (default "json")
  -h, --help              help for hioctl
      --host string       Hostname or ip address
  -k, --insecure          ignore certificate errors
  -p, --password string   Admin user password
      --port uint         port (default 8443)
  -r, --realm string      Admin user realm (default "local")
  -u, --user string       Admin username (default "admin")

Use "hioctl [command] --help" for more information about a command.
```
[hioctl Documentation](docs/hioctl.md)

You can add save your default settings in a json or yaml config file in ~/.hiveio/ or /etc/hive
~/.hiveio/hioctl.yaml
```
host: hive1
user: admin
password: admin
realm: local
insecure: true
```

