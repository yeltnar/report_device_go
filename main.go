package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
)

const PUB_IP_PING = "https://digitalocean.andbrant.com"

func main() {

	pub_ip := getPubIP()
	lan_ip := getLanIP("192.168")
	nebula_ip := getLanIP("10.10.10")

	fmt.Println(pub_ip)
	fmt.Println(lan_ip)
	fmt.Println(nebula_ip)
	log.Printf("done")
}

func getPubIP() string {
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

func getLanIP(regex_str string) string {
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
