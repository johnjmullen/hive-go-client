package main

import "fmt"

//"log"

func main() {
	client := Client{Host: "hive1", Port: 8443}
	client.Login("admin", "admin", "local")

	pools, _ := client.ListStoragePools()
	for _, pool := range pools {
		fmt.Println(pool)
	}

	// realm, _ := client.GetRealm("default")
	// fmt.Println(realm)

	// template, _ := client.GetTemplate("kubuntu")
	// fmt.Println(template)

	// profiles, _ := client.ListProfiles()
	// for _, profile := range profiles {
	// 	fmt.Println(profile)
	// }
	// hosts, _ := client.ListHosts()
	// for _, host := range hosts {
	// 	fmt.Printf("%v", host)
	// }

	/*ryzen := StoragePool{ Name: "ryzen",
	                        Type: "nfs",
	                        Server: "ryzen",
	                        Path: "/mnt/files/vms",
	                        Roles: []string{"template", "guest"},
	                        MountOptions: []string{"vers=4.1"},
	                      }
	  msg, err := client.CreateStoragePool(&ryzen)
	  fmt.Println(msg)
	  if err != nil {
	      log.Fatal(err)
	  }
	  pools, err = client.ListStoragePools()
	  var ryzenId string
	  for _, pool := range pools {
	      fmt.Println(pool)
	      if pool.Name == "ryzen" {
	          ryzenId = pool.Id
	      }
	  }
	  if ryzenId != "" {
	      ryzen_copy, err := client.GetStoragePool(ryzenId)
	      fmt.Printf("ryzen_copy: %v\n", ryzen_copy)

	      fail, err := client.GetStoragePool("test123")
	      fmt.Println(fail, err)

	      err = client.DeleteStoragePool(ryzenId)
	      if err != nil {
	          log.Fatal("Failed to delete ryzen storage pool", err)
	      }
	  } else {
	      log.Print("ryzen pool not found")
	  }
	  pools, err = client.ListStoragePools()
	  for _, pool := range pools {
	      fmt.Println(pool)
	  }

	  if err != nil {
	      log.Fatal(err)
	  }*/
}
