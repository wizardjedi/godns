package main

import (
	"net"
	"os"
	"fmt"
	"encoding/hex"
	"bytes"
)

type DNSMSG interface{}

type DNSMessage struct {
	messageId uint16
	flags uint16
	qcount uint16
	acount uint16
	nscount uint16
	arcount uint16
}

func main() {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":10053")

	if err != nil {
		fmt.Println("Error on resolve udp address:", err)
		os.Exit(1)
	}

	ServerConn, err := net.ListenUDP("udp", ServerAddr)

	if err != nil {
		fmt.Println("Error on listen:", err)
		os.Exit(1)
	}

	defer ServerConn.Close()

	buf := make([]byte, 1024);

	for {
		fmt.Println("Ready to use")

		n, addr, err := ServerConn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("Error: ",err)
		}

		buffer := bytes.NewBuffer(buf[0:n])

		fmt.Println("Received ",n, " bytes from ",addr)

		go processRequest(*ServerConn, buffer)
	}
}

func processRequest(ServerConnection net.UDPConn, buffer *bytes.Buffer) {
	fmt.Println(hex.Dump(buffer.Bytes()))

	ServerConnection.Write([]byte{0x00})

	decode(buffer)
}

func decode(buffer *bytes.Buffer) DNSMSG {
	dn := DNSMessage{};

	var id uint16

	b1,_ := buffer.ReadByte()
	b2,_ := buffer.ReadByte()

	id = (uint16)(b1 << 8 + b2)

	dn.messageId = id

	fmt.Println("DNSMessage:", dn)

	return dn
}

func (dn DNSMessage ) String() string {
	return fmt.
		Sprintf(
			"id:%04x flags:%04x qcount:%d nscount:%d acount:%d arcount:%d",
			dn.messageId,
			dn.flags,
			dn.qcount,
			dn.nscount,
			dn.acount,
			dn.arcount)
}