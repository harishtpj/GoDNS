package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	_ "github.com/harishtpj/klassy/String"
)

var records map[string]string
var logger = slog.Default()

func handleIfError(err error) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func main() {

	switch len(os.Args) {
	case 1:
		records = make(map[string]string)
	case 2:
		ftxt, err := os.ReadFile(os.Args[1])
		handleIfError(err)
		json.Unmarshal(ftxt, &records)
		logger.Info(fmt.Sprintf("Read %d records from %s", len(records), os.Args[1]))
	default:
		logger.Error("Invalid no. of arguments passed")
		os.Exit(1)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", ":53")
	handleIfError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	handleIfError(err)
	defer conn.Close()
	logger.Info(fmt.Sprint("Server started and running at", conn.LocalAddr().String()))
}
