package main

import (
	"flag"
	"log"
	"strings"

	"github.com/cpu/goacmedns"
)

func main() {
	apiBase := flag.String("api", "", "ACME-DNS server API URL")
	domain := flag.String("domain", "", "Domain to register an account for")
	storagePath := flag.String("storage", "", "Path to the JSON storage file to create/update")
	allowFrom := flag.String("allowFrom", "", "List of comma separated CIDR notation networks the account is allowed to be used from")
	username := flag.String("username", "", "Username to use for HTTP Basic Auth")
	password := flag.String("password", "", "Password to use for HTTP Basic Auth")
	flag.Parse()

	if *apiBase == "" {
		log.Fatal("You must provide a non-empty -api flag")
	}

	if *domain == "" {
		log.Fatal("You must provide a non-empty -domain flag")
	}

	if *storagePath == "" {
		log.Fatal("You must provide a non-empty -storage flag")
	}

	var allowedNetworks []string
	if *allowFrom != "" {
		allowedNetworks = strings.Split(*allowFrom, ",")
	}

	client := goacmedns.NewClient(*apiBase)
	storage := goacmedns.NewFileStorage(*storagePath, 0600)

	var auth goacmedns.AuthHandler = nil
	if *username != "" {
		auth = goacmedns.NewHttpBasicAuth(*username, *password)
	}
	newAcct, err := client.RegisterAccountWithAuth(allowedNetworks, auth)
	if err != nil {
		log.Fatal(err)
	}
	// Save it
	err = storage.Put(*domain, newAcct)
	if err != nil {
		log.Fatalf("Failed to put account in storage: %v", err)
	}

	err = storage.Save()
	if err != nil {
		log.Fatalf("Failed to save storage: %v", err)
	}

	log.Printf(
		"new account created for %q. "+
			"To complete setup for %q you must provision the following CNAME in your DNS zone:\n"+
			"%s CNAME %s.\n",
		*domain, *domain, "_acme-challenge."+*domain, newAcct.FullDomain)
}
