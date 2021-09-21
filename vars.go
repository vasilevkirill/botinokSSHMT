package main

import "errors"

var (
	errorDeviceFileNotFound  = errors.New("File Device list not found")
	errorImportFileNotFound  = errors.New("File botinok.rsc not found")
	errorStringNotTrueFormat = errors.New("Bad content in file devices")
	sshkkeysAlgo             = []string{
		"diffie-hellman-group1-sha1",
		"diffie-hellman-group14-sha1",
		"ecdh-sha2-nistp256",
		"ecdh-sha2-nistp384",
		"ecdh-sha2-nistp521",
		"diffie-hellman-group-exchange-sha1",
		"diffie-hellman-group-exchange-sha256",
		"curve25519-sha256@libssh.org"}
	sshChippers = []string{
		"aes128-gcm@openssh.com",
		"cast128-cbc",
		"aes128-ctr", "aes192-ctr", "aes256-ctr",
		"3des-cbc", "blowfish-cbc", "twofish-cbc", "twofish256-cbc", "twofish192-cbc", "twofish128-cbc", "aes256-cbc", "aes192-cbc", "aes128-cbc", "arcfour",
	}
)
