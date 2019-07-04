package main

import (
	"fmt"

	"bitbucket.org/johnmullen/hiveio-go-client/rest"
)

//TODO: use cobra for cli

func main() {
	client := rest.Client{Host: "hive1", Port: 8443}
	client.Login("admin", "admin", "local")

	version, _ := client.HostVersion()
	fmt.Println(version.Version)

}
