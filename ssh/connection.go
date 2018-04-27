package ssh

import (
	"io/ioutil"
	"os/user"

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

	/*sc.ClientConfig = &ssh.ClientConfig {
		User: "",
		Auth: []ssh.AuthMethod {
			ssh.Password(""),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}*/

	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	sc.ClientConfig = &ssh.ClientConfig{
		User: currentUser.Username,
		Auth: []ssh.AuthMethod{
			sc.publicKeyFile(currentUser.HomeDir + "/.ssh/id_rsa"),
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

func (sshConn *SSHConnection) publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func (sshConn *SSHConnection) Shell() {
	sshConn.Session.Shell()
}
