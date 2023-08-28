package main

import (
	"log"
	"net"
	"time"
)

//////////////////////////////////////////////////////////////////////////
//
// This section is all for telnet
//
//////////////////////////////////////////////////////////////////////////

// This will just skip any telnet negotiation stuff.  We read one byte
// at a time and decide whether to add it to the buffer
// This isn't perfect 'cos if the client sends DO/DONT requests we
// ignore them, when we should send WONT responses, but good enough... maybe?

func telnet_read_str(prompt string, conn net.Conn) string {
	buff := make([]byte, 0, 1024)
	b := make([]byte, 1)
	skip := false

	conn.Write([]byte(prompt))

	for {
		n, _ := conn.Read(b)
		if n == 0 || b[0] == 10 {
			break
		}
		if skip {
			// 251->254 are two byte responses
			if b[0] < 251 || b[0] > 254 {
				skip = false
			}
		} else if b[0] == 255 {
			skip = true
		} else if b[0] != 13 && b[0] != 10 {
			buff = append(buff, b[0])
		}
	}
	return string(buff)
}

func do_telnet(port string) {
	port = ":"+port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("TELNET Failed to listen on %s (%s)", port, err)
	}

	do_log(port, "TELNET", "Listening", "", "")

	for {
		conn, err := listener.Accept()
		if err != nil {
			do_log("", "TELNET", "Accept", "Failed", err.Error())
			continue
		}

		do_log(conn.RemoteAddr().String(), "TELNET", "Connection Established", "", "")
		for i := 0; i < 3; i++ {
			user := telnet_read_str("login: ", conn)

			// This escape sequence tells the client to disable echo
			pswd := telnet_read_str("\xff\xfb\x01Password: ", conn)

			time.Sleep(time.Second)

			// re-enable echo
			conn.Write([]byte("\xff\xfc\x01\r\nLogin incorrect\r\n"))
			if user != "" {
				do_log(conn.RemoteAddr().String(), "TELNET", "Password", user, pswd)
			}
		}
		conn.Close()
	}

}

