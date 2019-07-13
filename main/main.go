package main

import (
	"fmt"
	"os"

	rest "bitbucket.org/johnmullen/hiveio-go-client/rest"
)

//TODO: use cobra for cli

func main() {
	client := &rest.Client{Host: "hive1", Port: 8443, AllowInsecure: true}
	client.Login("admin", "admin", "local")

	version, _ := client.HostVersion()
	fmt.Println(version.Version)

	//Create pool
	/*confDir := "/home/john1/work/hiveio/conf"
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/pools/", confDir))
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		jsonFile, err := os.Open(fmt.Sprintf("%s/pools/%s", confDir, f.Name()))
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		json, _ := ioutil.ReadAll(jsonFile)
		var pool rest.Pool
		pool.FromJson(json)
		pool.ProfileID = "c730dc7c-0892-43d3-a856-308f2525eeb3"
		fmt.Println(pool)
		msg, err := pool.Create(client)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(msg)
	}*/

	pools, err := client.ListStoragePools()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, pool := range pools {
		files, _ := pool.Browse(client)
		fmt.Println(files)
	}

	clusters, err := client.ListClusters()
	expiration, licenseType, err := clusters[0].GetLicenseInfo(client)
	fmt.Printf("%s, %s", expiration, licenseType)

	alerts, err := client.ListAlerts()
	for _, alert := range alerts {
		yaml, _ := alert.ToYaml()
		fmt.Println(string(yaml))
	}
	/*body, err := client.LoadTemplate("ubuntu-minimal", "disk")
	fmt.Println(string(body), err)*/
}
