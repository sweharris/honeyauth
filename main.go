package main

import (
	"flag"
	"log"
	"strconv"
	"strings"
)

//////////////////////////////////////////////////////////////////////////
//
// Common code
//
//////////////////////////////////////////////////////////////////////////

func do_log(addr string, proto string, method string, user string, pass string) {
	user = strings.ReplaceAll(user, ",", "%2C")
	pass = strings.ReplaceAll(pass, ",", "%2C")
	log.Printf("%s,%s,%s,%s,%s\n", addr, proto, method, user, pass)
}

func main() {
	telnet := flag.Int("t", 0, "Run TELNET listener on specified port")
	ssh := flag.Int("s", 0, "Run SSH listener on specified port")
	flag.Parse()

	// Run listeners in go-routines so it's in parallel
	if *ssh != 0 {
		go do_ssh(strconv.Itoa(*ssh))
	}

	if *telnet != 0 {
		go do_telnet(strconv.Itoa(*telnet))
	}

	// Have main block forever; all the work is in goroutines
	<-(chan int)(nil)
}
