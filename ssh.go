package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

func connectSSH(d Device) (*ssh.Client, error) {
	var auth []ssh.AuthMethod
	auth = append(auth, ssh.Password(d.Password))
	config := &ssh.ClientConfig{
		Config: ssh.Config{
			KeyExchanges: sshkkeysAlgo,
			Ciphers:      sshChippers,
		},
		User:            d.User,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(d.SshTimeout) * time.Second,
	}
	str := fmt.Sprintf("%s:%d", d.Address, d.SshPort)
	return ssh.Dial("tcp", str, config)

}
