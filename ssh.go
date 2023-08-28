package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

//////////////////////////////////////////////////////////////////////////
//
// This section is all for SSH
//
//////////////////////////////////////////////////////////////////////////

// This will log ssh password attempts
func passwd(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	do_log(c.RemoteAddr().String(), "SSH", "Password", c.User(), string(pass))
	time.Sleep(time.Second)
	return nil, fmt.Errorf("Permission denied, please try again.")
}

// This will log that an ssh key was _attempted_ but won't log the key
func key(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	do_log(c.RemoteAddr().String(), "SSH", "Key", c.User(), key.Type())
	time.Sleep(time.Second)
	return nil, fmt.Errorf("Bad key", string(key.Type()))
}

// SSH server
func do_ssh(port string) {
	config := &ssh.ServerConfig{
		PasswordCallback:  passwd,
		PublicKeyCallback: key,
		ServerVersion:     "SSH-2.0-OpenSSH_7.4",
	}

	// You can generate keys with "make keys"
	keys := []string{"ssh_host_ecdsa_key", "ssh_host_ed25519_key", "ssh_host_rsa_key"}

	for _, key := range keys {
		privateBytes, err := ioutil.ReadFile(key)
		if err != nil {
			log.Fatal("SSH: Failed to load private key "+key+". Did you run `make keys`?")
		}

		private, err := ssh.ParsePrivateKey(privateBytes)
		if err != nil {
			log.Fatal("SSH: Failed to parse private key " + key)
		}

		config.AddHostKey(private)
	}

	port = ":"+port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("SSH: Failed to listen on %s (%s)", port, err)
	}

	do_log(port, "SSH", "Listening", "", "")
	for {
		conn, err := listener.Accept()
		if err != nil {
			do_log("", "SSH", "Accept", "Failed", err.Error())
			continue
		}
		do_log(conn.RemoteAddr().String(), "SSH", "Connection Established", "", "")

		// Now handle the SSH connection
		// If we do this in a goroutine then we can handle multiple
		// concurrent requests, but need to be careful about DoS
		// attempts
		// Doing it in the main routine means we can only process 1
		// request at a time, but that's probably good enough...

		// We expect this to fail because we reject all authn so we
		// do nothing and just drop any return values
		ssh.NewServerConn(conn, config)
	}
}
