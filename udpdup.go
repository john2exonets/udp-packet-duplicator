package main

//
//  udpdup.go  --  Duplicate all incomming UPD packets to one or more UDP servers.
//    Useful for duplicating Syslog packets to multiple Syslog record scanners.
//
//  Example Output (with Debug = 10):
// <134>May 24 14:43:05 Testing1 a10logd: [ACOS]<6> Port 31615 type TCP on server 44.147.45.220 is deleted
// Sent  104 bytes [::]:5514 -> 127.0.0.1:8514
// Sent  104 bytes [::]:5514 -> 10.1.1.59:6514
// Sent  104 bytes [::]:5514 -> [::1]:7514
// <134>May 24 14:43:05 Testing1 a10logd: [ACOS]<6> Server 44.147.45.220 is deleted
// Sent  81 bytes [::]:5514 -> 127.0.0.1:8514
// Sent  81 bytes [::]:5514 -> [::1]:7514
// Sent  81 bytes [::]:5514 -> 10.1.1.59:6514
// <133>May 24 14:43:05 Testing1 a10logd: [SYSTEM]<5> Session ID 190 is now closed.
// Sent  81 bytes [::]:5514 -> 127.0.0.1:8514
// Sent  81 bytes [::]:5514 -> [::1]:7514
// Sent  81 bytes [::]:5514 -> 10.1.1.59:6514
//
//  John D. Allen
//  Global Solutions Architect - Cloud, A10 Networks
//  Apache 2.0 License Applies
//  May, 2021
//

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

type Dest struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type Configuration struct {
	Debug  int    `json:"debug"`
	Port   int    `json:"port"`
	MaxBuf int    `json:"maxbuf"`
	Dests  []Dest `json:"dests"`
}

var config Configuration

//
// Get the config
func getConfig(fn string) (Configuration, error) {
	jsonFile, err := os.Open(fn)
	if err != nil {
		return Configuration{}, errors.New("Unable to open Config File!")
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var c Configuration
	json.Unmarshal(byteValue, &c)

	return c, nil
}

//
// Send packet to UDP Server
func sendUDPpkt(conn *net.UDPConn, addr *net.UDPAddr, pkt string, debug int) {
	//fmt.Println(pkt)
	n, err := conn.WriteToUDP([]byte(pkt), addr)
	if err != nil {
		fmt.Printf("--> Error: Unable to send packet to %v", addr)
	}
	if debug > 8 {
		fmt.Println("Sent ", n, "bytes", conn.LocalAddr(), "->", addr)
	}

}

//
// MAIN Function
func main() {
	//
	// Get Config info
	config, err := getConfig("./config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//
	// Setup Input UDP port
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: config.Port,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	fmt.Printf("server listening %s\n", conn.LocalAddr().String())

	//-------------------[   MAIN   ]---------------------
	for {
		message := make([]byte, config.MaxBuf)
		rlen, _, err := conn.ReadFromUDP(message[:])
		if err != nil {
			fmt.Printf("<-- Error: %v", err)
		}

		//fmt.Printf("m>>%v\n", message)
		data := strings.TrimSpace(string(message[:rlen]))
		if config.Debug > 9 { // Print out received Packet
			fmt.Println(data)
		}

		for _, dst := range config.Dests { // Loop through destinations
			addr := net.UDPAddr{
				Port: dst.Port,
				IP:   net.ParseIP((dst.IP)),
			}
			go sendUDPpkt(conn, &addr, data, config.Debug)
		}

	}
}
