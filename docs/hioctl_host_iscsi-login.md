## hioctl host iscsi-login

login to an iscsi target

```
hioctl host iscsi-login [flags]
```

### Options

```
  -h, --help                    help for iscsi-login
  -i, --id string               hostid
      --ip string               host ip address
      --iscsi-password string   iscsi password
      --iscsi-username string   iscsi username
  -n, --name string             hostname
      --portal string           portal
      --target string           target
```

### Options inherited from parent commands

```
      --config string     config file
      --format string     format (json/yaml) (default "json")
      --host string       Hostname or ip address
  -k, --insecure          ignore certificate errors
  -p, --password string   Admin user password
      --port uint         port (default 8443)
      --profile string    Load a profile from the config file
  -r, --realm string      Admin user realm (default "local")
  -u, --user string       Admin username (default "admin")
```

### SEE ALSO

* [hioctl host](hioctl_host.md)	 - host operations

###### Auto generated by spf13/cobra on 26-Jun-2025
