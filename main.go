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

// func hardcodedResponse(id uint16) []byte {
// 	base := []byte{
// 		0x00, 0x00, // ID placeholder
// 		0x81, 0x80, // Flags: response + recursion
// 		0x00, 0x01, // QDCOUNT = 1
// 		0x00, 0x01, // ANCOUNT = 1
// 		0x00, 0x00, // NSCOUNT = 0
// 		0x00, 0x00, // ARCOUNT = 0
//
// 		// Question: harish.com A IN
// 		0x06, 'h', 'a', 'r', 'i', 's', 'h',
// 		0x03, 'c', 'o', 'm',
// 		0x00,
// 		0x00, 0x01, // Type A
// 		0x00, 0x01, // Class IN
//
// 		// Answer section
// 		0xC0, 0x0C, // Pointer to offset 12 (name)
// 		0x00, 0x01, // Type A
// 		0x00, 0x01, // Class IN
// 		0x00, 0x00, 0x00, 0x3C, // TTL = 60
// 		0x00, 0x04,             // RDLENGTH = 4
// 		0x7F, 0x00, 0x00, 0x01, // 127.0.0.1
// 	}
//
// 	// Patch in the request ID
// 	base[0] = byte(id >> 8)
// 	base[1] = byte(id & 0xFF)
//
// 	return base
// }

func handleDNSQuery(conn *net.UDPConn, clientAddr *net.UDPAddr, msg []byte) {
	header, question := parseDNSQuery(msg)
	ipStr, found := records[question.Name]
	if !found {
		log.Println("Record doesn't exist", question.Name)
		return
	}
	ip := net.ParseIP(ipStr).To4()
	response := buildResponse(header, question, ip)
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
			"harish.com": "127.0.0.1",
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
