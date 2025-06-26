package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

var records map[string]string

func handleIfError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func handleDNSQuery(conn *net.UDPConn, clientAddr *net.UDPAddr, msg []byte) {
	var response []byte
	header, question := parseDNSQuery(msg)

	ipStr, found := records[question.Name]
	if !found {
		log.Printf("Record doesn't exists: %s. Sending NXDOMAIN\n", question.Name)
		response = buildFailResponse(header, question)
	} else {
		ip := net.ParseIP(ipStr).To4()
		response = buildResponse(header, question, ip)
	}

	log.Println("Sending response of", len(response), "bytes")
	_, err := conn.WriteToUDP(response, clientAddr)
	if err != nil {
		log.Println("Error sending response:", err)
	}
}

func main() {
	switch len(os.Args) {
	case 1:
		records = map[string]string {
			"example.com": "127.0.0.1",
		}
	case 2:
		ftxt, err := os.ReadFile(os.Args[1])
		handleIfError(err)
		json.Unmarshal(ftxt, &records)
		log.Printf("Read %d records from %s\n", len(records), os.Args[1])
	default:
		log.Fatalln("Invalid no. of arguments passed")
	}

	udpAddr, err := net.ResolveUDPAddr("udp", ":53")
	handleIfError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	handleIfError(err)
	defer conn.Close()
	log.Println("Server started and running at", conn.LocalAddr().String())

	buffer := make([]byte, 512)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error reading:", err)
			continue
		}
		log.Printf("Received packet from CLIENT:%s\n", clientAddr)
		go handleDNSQuery(conn, clientAddr, buffer[:n])
	}
}
