package main

import (
	"check-cdn/mycloudflare"
	"check-cdn/mytencent"
	"fmt"
	"log"
	"os"
)

func main() {

	cloudflare_token := os.Getenv("CLOUDFLARE_API_TOKEN")
	Tencent_1_SecretId := os.Getenv("TENCENT_1_SECRETID")
	Tencent_1_SecretKey := os.Getenv("TENCENT_1_SECRETKEY")

	filename := "target.txt"
	lines, err := mycloudflare.ReadFileToArray(filename)
	if err != nil {
		log.Printf("Failed to read file: %s\n", err)
		return
	}

	for _, purge_target := range lines {
		//--------------------------------------------------------

		clouflare_purge_result := mycloudflare.CloudflarePurgeResult(cloudflare_token, purge_target, "A")
		fmt.Println("++ " + clouflare_purge_result + " ++")
		if clouflare_purge_result == "DOMAIN_NOT_EXIST" {
			fmt.Println("========> DOMAIN_NOT_EXIST pleae check again.")
		} else if clouflare_purge_result == "SUBDOMAIN_NOT_EXIST" {
			fmt.Println("========> go other cloud CDN check(" + purge_target + ")")
			//tencent-1
			if mytencent.PurgeCdn(Tencent_1_SecretId, Tencent_1_SecretKey, purge_target) {
				continue
			}
			/*
								if(mytencent.PurgeCdn(tencent-1,x,x,c)) {continue}
								else if (tencent.PurgeCdn(tencent-2,x,x,c))
								{continue}
								else if (tencent.PurgeCdn(tencent-3,x,x,c))
								{continue}
				                else if (huawei.PurgeCdn(huawei-1,x,x,c))
								{continue}
				                else if (huawei.PurgeCdn(huawei-2,x,x,c))
								{continue}
				                else if (huawei.PurgeCdn(huawei-3,x,x,c))
								{continue}
								else{
									fmt.Println("Not found !")
								}
			*/

		} else {
			fmt.Println("========> Cache purge successful for hostname.")
		}

		//--------------------------------------------------------
	}
}
