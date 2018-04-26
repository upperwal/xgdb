package ssh

import (
	"golang.org/x/crypto/ssh"
)

type SSHConnection struct {
	Hostname string
	ClientConfig *ssh.ClientConfig
	Session *ssh.Session
}

func NewSSHConnection(hostname string) (*SSHConnection, error) {
	sc := &SSHConnection {
		Hostname: hostname,
	}

	sc.ClientConfig = &ssh.ClientConfig {
		User: "marslab",
		Auth: []ssh.AuthMethod {
			ssh.Password("admin"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	c, err := ssh.Dial("tcp", hostname + ":22", sc.ClientConfig)
	if err != nil {
		return nil, err
	}

	sc.Session, err = c.NewSession()
	if err != nil {
		return nil, err
	}

	return sc, nil
}

func (sshConn *SSHConnection) Shell() {
	sshConn.Session.Shell()
}
