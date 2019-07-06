package main

import (
	"fmt"

	rest "bitbucket.org/johnmullen/hiveio-go-client/rest"
)

//TODO: use cobra for cli

func main() {
	client := rest.Client{Host: "hive1", Port: 8443}
	client.Login("admin", "admin", "local")

	version, _ := client.HostVersion()
	fmt.Println(version.Version)
	body, err := client.LoadTemplate("ubuntu-minimal", "disk")
	fmt.Println(string(body), err)
}
