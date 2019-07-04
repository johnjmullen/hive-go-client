package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

//TODO: use cobra for cli

func main() {
	client := Client{Host: "hive1", Port: 8443}
	client.Login("admin", "admin", "local")

	version, _ := client.HostVersion()
	fmt.Println(version.Version)

	//storage pools
	confDir := "/home/john1/work/hiveio/conf"
	pools := []string{"vms", "iso", "uvs"}
	for _, pool := range pools {
		jsonFile, err := os.Open(fmt.Sprintf("%s/pools/%s.conf", confDir, pool))
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		json, _ := ioutil.ReadAll(jsonFile)
		var sp StoragePool
		sp.FromJson(json)

		msg, err := client.CreateStoragePool(&sp)
		fmt.Println(msg)
		if err != nil {
			log.Fatal(err)
		}
	}

	//realm
	home := Realm{Name: "HOME", FQDN: "home.john-mullen.net", Verified: true}

	client.CreateRealm(&home)
	clusters, _ := client.ListClusters()
	for _, cluster := range clusters {
		fmt.Println(cluster)
	}

	//Default Profile
	profile := Profile{Name: "default", BypassBroker: true}
	client.CreateProfile(&profile)

	//templates
	vms, err := client.GetStoragePoolByName("vms")
	if err != nil {
		fmt.Println(err)
		return
	}
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/templates/", confDir))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		jsonFile, err := os.Open(fmt.Sprintf("%s/templates/%s", confDir, f.Name()))
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		json, _ := ioutil.ReadAll(jsonFile)
		var tmpl Template
		tmpl.FromJson(json)
		for i := range tmpl.Disks {
			tmpl.Disks[i].StorageID = vms.ID
		}
		msg, err := client.CreateTemplate(&tmpl)
		fmt.Println(msg)
		if err != nil {
			log.Fatal(err)
		}
	}

	extraHosts := []string{"hive2", "hive3"}
	for _, host := range extraHosts {
		ips, err := net.LookupIP(host)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to find ip for %s", host)
			os.Exit(1)
		}
		err = client.JoinHost("admin", "admin", ips[0].String())
		fmt.Println(err)
	}
}
