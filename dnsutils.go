package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"strings"
)

type DNSHeader struct {
	ID uint16
	Flags uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

type DNSQuestion struct {
	Name string
	Type uint16
	Class uint16
	Data []byte
}

func parseDNSQuery(msg []byte) (DNSHeader, DNSQuestion){
	var header DNSHeader
	buf := bytes.NewReader(msg[:12])
	binary.Read(buf, binary.BigEndian, &header)

	question := parseQuestion(msg[12:])
	log.Printf("Query: ID=%d, Name=%s, Type=%d, Class=%d\n", header.ID, question.Name, question.Type, question.Class)

	return header, question
}

func parseQuestion(msg []byte) DNSQuestion {
	var q DNSQuestion
	var labels []string
	offset := 0

	for {
		length := int(msg[offset])
		if length == 0 {
			offset++
			break
		}
		offset++
		labels = append(labels, string(msg[offset:offset+length]))
		offset += length
	}

	q.Name = strings.Join(labels, ".")
	q.Type = binary.BigEndian.Uint16(msg[offset:offset+2])
	q.Class = binary.BigEndian.Uint16(msg[offset+2:offset+4])
	offset += 4
	q.Data = msg[:offset]

	return q
}

func buildFailResponse(header DNSHeader, question DNSQuestion) []byte {
	var buf bytes.Buffer
	respHeader := DNSHeader{
		ID: header.ID,
		Flags: uint16(0x8183),
		QDCount: uint16(1),
		ANCount: uint16(0),
		NSCount: uint16(0),
		ARCount: uint16(0),
	}
	binary.Write(&buf, binary.BigEndian, respHeader)
	buf.Write(question.Data)

	return buf.Bytes()
}

func buildResponse(header DNSHeader, question DNSQuestion, ip net.IP) []byte {
	var buf bytes.Buffer
	respHeader := DNSHeader{
		ID: header.ID,
		Flags: uint16(0x8180),
		QDCount: uint16(1),
		ANCount: uint16(1),
		NSCount: uint16(0),
		ARCount: uint16(0),
	}
	binary.Write(&buf, binary.BigEndian, respHeader)
	buf.Write(question.Data)

	// Answer Section
	buf.WriteByte(0xC0) // Pointer to
	buf.WriteByte(0x0C) // 0x0C => 12
	binary.Write(&buf, binary.BigEndian, uint16(1)) // Type A
	binary.Write(&buf, binary.BigEndian, uint16(1)) // Class IN
	binary.Write(&buf, binary.BigEndian, uint32(60)) // TTL
	binary.Write(&buf, binary.BigEndian, uint16(4)) // Rdlength
	buf.Write(ip)

	return buf.Bytes()
}
