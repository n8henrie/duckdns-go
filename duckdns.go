package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var version = "undefined"

type configuration struct {
	domains []string
	token   string
}

var config configuration

func printVersion() {
	fmt.Println("duckdns version:", version)
	os.Exit(0)
}

// Documentation for main function
func main() {
	showVersion := flag.Bool("version", false, "Print version")
	configFile := flag.String("configFile", "config.json", "Path to config file")
	flag.Parse()

	if *showVersion {
		printVersion()
	}

	file, err := os.Open(*configFile)
	if err != nil {
		log.Fatal("Unable to open config file", err)
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatal("Unable to decode config", err)
	}

	domains := strings.Join(config.domains, ",")
	u, err := url.Parse("https://www.duckdns.org/update")
	if err != nil {
		log.Fatal("Error parsing url", err)
	}
	v := u.Query()
	v.Set("token", config.token)
	v.Set("domains", domains)
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal("Error sending GET request", err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}
