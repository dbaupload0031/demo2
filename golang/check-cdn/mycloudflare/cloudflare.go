//package main

package mycloudflare

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

func ChkUser() string {
	// Create a new Cloudflare API client
	//api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	api, err := cloudflare.NewWithAPIToken("gTUuzv4crjAtBH7aUPbR_YQxiXDRkDmrFUAAYP7u")
	if err != nil {
		log.Fatal(err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Fetch user details on the account
	u, err := api.UserDetails(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(u.ID)
	return u.Email
}

// Function to check if a domain exists in a Cloudflare account
func DomainExists(api *cloudflare.API, ctx context.Context, dnsName string) (bool, error) {
	// List all zones in the Cloudflare account
	zones, err := api.ListZones(ctx)
	if err != nil {
		return false, err
	}

	// Check if any zone has a matching domain name
	for _, zone := range zones {
		if zone.Name == dnsName {
			zoneID, err := api.ZoneIDByName(dnsName)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("zoneID: " + zoneID)
			return true, nil
		}
	}
	return false, nil
}

func SimpleChkZoneId(api *cloudflare.API, ctx context.Context, dnsName string) string {
	zoneID, err := api.ZoneIDByName(dnsName)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("zoneID: " + zoneID)
	return zoneID
}

func GetSubdomainAndDomain(domainName string) (string, string) {
	parts := strings.Split(domainName, ".")
	if len(parts) > 2 {
		subdomain := strings.Join(parts[:len(parts)-2], ".")
		domain := strings.Join(parts[len(parts)-2:], ".")
		return subdomain, domain
	} else if len(parts) == 2 {
		return "", domainName
	}
	return "", ""
}

//https://pkg.go.dev/github.com/cloudflare/cloudflare-go@v0.70.0#API.ListDNSRecords
// func check A type

func ChkDnsRecord(api *cloudflare.API, ctx context.Context, zoneID string, dnstype string, fullName string) bool {
	recs, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{Type: dnstype})
	if err != nil {
		fmt.Println(err)

	}
	//fmt.Println(recs)
	for _, r := range recs {
		if r.Name == fullName {
			fmt.Println(r.Name, r.Content)
			return true // Return true as soon as a matching record is found
		}
	}
	return false
}

func PurgeCache(api *cloudflare.API, ctx context.Context, zoneID string, hostname string) error {
	req := cloudflare.PurgeCacheRequest{
		//Hosts: []string{hostname},//Only enterprise zones can purge by host
		Files: []string{"https://jenkins.sre360.store/job/test/38 "},
		//Everything: true, //ok
	}

	resp, err := api.PurgeCache(ctx, zoneID, req)
	if err != nil {
		log.Printf("Cache purge failed for hostname %s: %s\n", hostname, err)
		return err
	}

	log.Printf("Cache purge successful for hostname: %s\n", hostname)
	log.Printf("Response: %+v\n", resp)

	return nil
}

func CloudflarePurgeResult(token string, fullname string, sub_dns_type string) string {
	//init
	//api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		fmt.Println(err)
		//return
	}
	ctx := context.Background()
	fullName := fullname
	subdomainname, domain := GetSubdomainAndDomain(fullName)
	dns_type := sub_dns_type
	fmt.Println("domain: " + domain)
	fmt.Println("subdomainname " + subdomainname)

	//check domain
	exists, err := DomainExists(api, ctx, domain)
	if err != nil {
		log.Fatal(err)
	}
	//check sub domain
	if exists {
		zoneid_res := SimpleChkZoneId(api, ctx, domain)
		//check subdomain A record
		dns_res := ChkDnsRecord(api, ctx, zoneid_res, dns_type, fullName)

		if dns_res {
			fmt.Println("purge target:" + fullName)
			PurgeCache(api, ctx, zoneid_res, fullName)
			return "TARGET_PURGE"
		} else {
			fmt.Println("subdomain (" + subdomainname + ") does not exist in the domain (" + domain + ") with A record ")
			return "SUBDOMAIN_NOT_EXIST"
		}

	} else {
		fmt.Printf("Domain (%s) does not exist in current Cloudflare\n", domain)
		return "DOMAIN_NOT_EXIST"

	}
}

func ReadFileToArray(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Failed to open file: %s\n", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file: %s\n", err)
		return nil, err
	}

	return lines, nil
}

/*
//pending to test,模組化需將main(),提到外面處理-----------
func main() {
	//ChkUser()
	var a = ChkUser()
	fmt.Println("++ " + a + " ++")

	//init
	api, err := cloudflare.NewWithAPIToken(".......")
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	fullName := "jenkins.devopslife.shop"
	subdomainname, dnsName := GetSubdomainAndDomain(fullName)
	dns_type := "A"
	fmt.Println("dnsName: " + dnsName)
	fmt.Println("subdomainname " + subdomainname)

	//call
	exists, err := DomainExists(api, ctx, dnsName)
	if err != nil {
		log.Fatal(err)
	}
	if exists {
		zoneid_res := SimpleChkZoneId(api, ctx, dnsName)
		dns_res := ChkDnsRecord(api, ctx, zoneid_res, dns_type, fullName)

		if dns_res {
			fmt.Println("purge target:" + fullName)
			PurgeCache(api, ctx, zoneid_res, fullName)
		} else {
			fmt.Println("Not found subdomain " + subdomainname + " in the domain " + dnsName + " with A record ")
		}

	} else {
		fmt.Printf("Domain %s does not exist in Cloudflare\n", dnsName)
	}

}

//-------------------------------------------------------
*/
