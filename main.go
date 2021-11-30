package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"time"
)

const PUB_IP_PING = "https://digitalocean.andbrant.com"

func reportToHomeAssistant(device_name string, report_key string, state string, time_stamp string, file_location string) {

	access_token, err := os.ReadFile(file_location)
	check(err)

	url := "http://homeassistant:8123/api/states/" + device_name + "." + report_key + ""

	jsonStr := []byte(`{"state":"` + state + `","attributes":{"time_stamp":"` + time_stamp + `"}}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	req.Header.Set("Authorization", "Bearer "+string(access_token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println(string(jsonStr))

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

func main() {

	if len(os.Args) < 3 {
		panic("No device name to report or token file. Provide as arguments, in order, when running")
	}

	device_name := os.Args[1]
	time_stamp := fmt.Sprint(time.Now().Unix())
	file_location := os.Args[2]

	type IpList struct {
		Pub    string `json:"pub"`
		Lan    string `json:"lan"`
		Nebula string `json:"nebula"`
	}

	pub_ip := pingForIP(PUB_IP_PING)
	lan_ip := getSystemIP("192.168")
	nebula_ip := getSystemIP("10.10.10")

	reportToHomeAssistant(device_name, "pub_ip", pub_ip, time_stamp, file_location)
	reportToHomeAssistant(device_name, "lan_ip", lan_ip, time_stamp, file_location)
	reportToHomeAssistant(device_name, "nebula_ip", nebula_ip, time_stamp, file_location)

	ip_list := IpList{pub_ip, lan_ip, nebula_ip}
	ip_list_json, _ := json.Marshal(ip_list)

	fmt.Println(device_name)
	fmt.Println(string(ip_list_json))
}

func pingForIP(url_to_ping string) string {
	resp, err := http.Get(url_to_ping)
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
