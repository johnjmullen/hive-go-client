# hive-go-client

## rest
hive-go-client/rest is a wrapper around the hiveio rest api

[Documentation](https://pkg.go.dev/github.com/hive-io/hive-go-client/rest)


---
## hioctl
hive-go-client/hioctl is a command line interface for the api
install it with: `go get -u github.com/hive-io/hive-go-client/hioctl`

```
hive fabric rest api client

Usage:
  hioctl [command]

Available Commands:
  alert       alert operations
  cluster     cluster operations
  completion  Generate the autocompletion script for the specified shell
  export      export cluster configuration
  guest       guest operations
  help        Help about any command
  host        host operations
  import      import configuration from an export file
  metric      metrics operations
  pool        pool operations
  profile     profile operations
  realm       realm operations
  storage     storage operations
  task        task operations
  template    template operations
  user        user operations
  version     hioctl version information

Flags:
      --config string     config file
      --format string     format (json/yaml) (default "json")
  -h, --help              help for hioctl
      --host string       Hostname or ip address
  -k, --insecure          ignore certificate errors
  -p, --password string   Admin user password
      --port uint         port (default 8443)
      --profile string    Load a profile from the config file
  -r, --realm string      Admin user realm (default "local")
  -u, --user string       Admin username (default "admin")
  -v, --version           version for hioctl

Use "hioctl [command] --help" for more information about a command.
```
[hioctl Documentation](docs/hioctl.md)

Save default settings in hioctl.yaml in ~/.hiveio/, /etc/hive, or the current directory.

For multiple profiles, use `hioctl --profile cluster1` to select a profile

```
host: hive-hostname
user: admin
password: password
realm: local
insecure: true
profiles:
  cluster1:
    host: hive-hostname
    user: user1
    password: my-password
    realm: my-realm
  cluster2:
    host: hive-fabric-ip-address
    user: admin
    password: my-password
    realm: local
```