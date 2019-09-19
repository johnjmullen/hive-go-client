# hive-go-client

hive-go-client/rest is a go wrapper around the hiveio api

hive-go-client/hioctl is a command line interface for the api

install it with: `go get -u github.com/hive-io/hive-go-client/hioctl`

Since this is a private repo you need to run
`git config --global url."git@github.com:".insteadOf "https://github.com/"`

```
 ~/go/bin/hioctl
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
  storage     Storage                                                                                                                      
  template    template operations                                                                                                          
  util        hioctl utilities

Use "hioctl [command] --help" for more information about a command.
```
[hioctl Documentation](docs/hioctl.md)

You can add a json or yaml config file in ~/.hiveio/ for default settings
~/.hiveio/hioctl.yaml
```
host: hive1
user: admin
password: admin
realm: local
insecure: true
```

