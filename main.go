package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
)

const PUB_IP_PING = "https://digitalocean.andbrant.com"

func main() {

	type IpList struct {
		Pub    string `json:pub`
		Lan    string `json:lan`
		Nebula string `json:nebula`
	}

	pub_ip := pingForIP()
	lan_ip := getSystemIP("192.168")
	nebula_ip := getSystemIP("10.10.10")

	ip_list := IpList{pub_ip, lan_ip, nebula_ip}
	ip_list_json, _ := json.Marshal(ip_list)

	fmt.Println(string(ip_list_json))
}

func pingForIP() string {
	resp, err := http.Get(PUB_IP_PING)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//Convert the body to type string
	sb := string(body)

	return sb
}

func getSystemIP(regex_str string) string {
	ifaces, _ := net.Interfaces()
	// ifaces, err := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		// addrs, err := i.Addrs()
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// process IP address
			// fmt.Println(ip)

			ip_str := ip.String()

			match, _ := regexp.MatchString(regex_str, ip_str)
			// fmt.Println(match, ip)

			if match {
				return ip_str
			}
		}
	}
	return ""
}
