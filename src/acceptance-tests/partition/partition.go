package partition

import (
	"bytes"
	"fmt"
	"log"

	"code.google.com/p/go.crypto/ssh"
)

func On(ip string, ports []int) {
	const password = "c1oudc0w"
	config := &ssh.ClientConfig{
		User: "vcap",
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ip), config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	for _, port := range ports {
		remoteCommands := fmt.Sprintf(`echo %[1]s | sudo -S true && \
    sudo iptables -A INPUT  ! -s 127.0.0.1 -p tcp ! --destination-port 22 -j DROP && \
    sudo iptables -A OUTPUT   -s 127.0.0.1 -p tcp ! --source-port 22 -j DROP && \
    sudo iptables -A INPUT  ! -s %[2]s     -p tcp ! --destination-port 22 -j DROP && \
    sudo iptables -A OUTPUT   -s %[2]s     -p tcp ! --source-port 22 -j DROP
    `, password, ip, port)

		println(remoteCommands)

		err = session.Run(remoteCommands)
		if err != nil {
			log.Fatal("Failed to run: " + err.Error())
		}
		fmt.Println(b.String())
	}
}

func Off(ip string) {
	const password = "c1oudc0w"
	config := &ssh.ClientConfig{
		User: "vcap",
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ip), config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	err = session.Run(fmt.Sprintf("echo %[1]s | sudo -S iptables -F", password))
	if err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
}
