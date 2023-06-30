                       package main

                        import (
                            "fmt"
                            "github.com/cloudflare/cloudflare-go"
                            "log"
                            "os"
                            "context"
                        )

                        func main() {
                            // Create a new Cloudflare API client
                            api, err := cloudflare.NewWithAPIToken(<apiKey>)
                            if err != nil {
                                log.Fatal(err)
                            }

                            // Set up the domain name to check
                            domainName := "example.com"

                            // Most API calls require a Context
                            ctx := context.Background()

                            // Fetch user details on the account
                            u, err := api.UserDetails(ctx)
                            if err != nil {
                                log.Fatal(err)
                            }

                            // Print user details
                            fmt.Println(u)
                            fmt.Println("xxxxxxx")
                        }
