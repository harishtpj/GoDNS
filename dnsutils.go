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
}

func parseDNSQuery(msg []byte) (DNSHeader, DNSQuestion){
	var header DNSHeader
	buf := bytes.NewReader(msg[:12])
	binary.Read(buf, binary.BigEndian, &header)

	question, _ := parseQuestion(msg[12:])
	log.Printf("Query: ID=%d, Name=%s, Type=%d, Class=%d\n", header.ID, question.Name, question.Type, question.Class)

	return header, question
}

func parseQuestion(msg []byte) (DNSQuestion, int) {
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

	return q, offset
}

func buildResponse(header DNSHeader, question DNSQuestion, ip net.IP) []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, header.ID)
	binary.Write(&buf, binary.BigEndian, uint16(0x8180)) // std response
	binary.Write(&buf, binary.BigEndian, uint16(1))
	binary.Write(&buf, binary.BigEndian, uint16(1))
	binary.Write(&buf, binary.BigEndian, uint16(0))
	binary.Write(&buf, binary.BigEndian, uint16(0))

	labels := strings.Split(question.Name, ".")
	for _, label := range labels {
		buf.WriteByte(byte(len(label)))
		buf.WriteString(label)
	}
	buf.WriteByte(0)
	binary.Write(&buf, binary.BigEndian, question.Type)
	binary.Write(&buf, binary.BigEndian, question.Class)

	buf.WriteByte(0xC0)
	buf.WriteByte(0x0C)
	binary.Write(&buf, binary.BigEndian, uint16(1)) // Type A
	binary.Write(&buf, binary.BigEndian, uint16(1)) // Class IN
	binary.Write(&buf, binary.BigEndian, uint32(60)) // TTL
	binary.Write(&buf, binary.BigEndian, uint16(4)) // Rdlength
	buf.Write(ip)

	return buf.Bytes()
}
