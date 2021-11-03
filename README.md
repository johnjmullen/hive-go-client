# hive-go-client

## rest
hive-go-client/rest is a wrapper around the hiveio rest api

[Documentation](https://pkg.go.dev/github.com/hive-io/hive-go-client/rest)


---
## hioctl
hive-go-client/hioctl is a command line interface for the api
install it with: `go get -u github.com/hive-io/hive-go-client/hioctl`

```
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
  user        user operations

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

Save connection settings in hioctl.yaml in ~/.hiveio/, /etc/hive, or the current directory

```
host: hive-hostname
user: admin
password: password
realm: local
insecure: true
```