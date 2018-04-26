package main

import (
	//"bytes"
	"fmt"
	"io"
	"os"

	//"golang.org/x/crypto/ssh"

	"github.com/upperwal/xgdb/gdb"
)

const (
	PROMT = "(xgdb)"
)

type Shell struct {
	StdErr io.Reader
	Stdout io.Reader
	Stdin io.WriteCloser
}

func NewShell(r io.Reader, w io.WriteCloser, e io.Reader) Shell {
	return Shell{e,r,w}
}

func (s Shell) reader() (string, error) {
	var buf [100]byte

	n, err := s.Stdout.Read(buf[:])

	return string(buf[:n]), err
}

func main() {

	gdb, err := gdb.NewGDB("main")
	if err != nil {
		panic(err)
	}

	gdb.Start()

	/*sshConfig := &ssh.ClientConfig {
		User: "marslab",
		Auth: []ssh.AuthMethod {
			ssh.Password("admin"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	c, err := ssh.Dial("tcp", "localhost:22", sshConfig)
	if err != nil {
		panic(err)
	}

	s, err := c.NewSession()
	if err != nil {
		panic(err)
	}

	defer s.Close()

	//s.Stdin = os.Stdin

	in, _ := s.StdinPipe()
	out, _ := s.StdoutPipe()
	er, _ := s.StderrPipe()

	s.Shell()

	sh := NewShell(out, in, er)*/






	/*_, err = ip.Write([]byte("/usr/bin/whoami\n"))

	if err != nil {
		panic(err)
	}


	op.Read(buf[:])

	fmt.Println(string(buf[:]))*/

	/*go readRoutine(sh)
	go writeRoutine(sh)
	go errRoutine(sh)*/


	/*_, err = sh.Stdin.Write([]byte("export ARR=/home/\n"))
	fmt.Println("hkhkh")
	if err != nil {
		panic(err)
	}

	_, err = sh.Stdin.Write([]byte("echo $ARR\n"))

	_, err = sh.Stdin.Write([]byte("whoami\n"))

	if err != nil {
		panic(err)
	}*/


	/*var buffer bytes.Buffer
	s.Stdout = &buffer

	if err := s.Run("/usr/bin/whoami"); err != nil {
		panic(err)
	}
	fmt.Printf("%s", buffer.String())*/

	select {}
}

func writeRoutine(s Shell) {
	var buf [100]byte

	fmt.Printf("%s ", PROMT)

	for {
		n, _ := os.Stdin.Read(buf[:])
		s.Stdin.Write(buf[:n])
	}



}

func errRoutine(s Shell) {

	var buf [100]byte
	for {
		da, _ := s.StdErr.Read(buf[:])
		fmt.Printf("%s", string(buf[:da]))
		fmt.Printf("%s ", PROMT)
	}
}

func readRoutine(s Shell) {
	for {
		da, _ := s.reader()
		fmt.Printf("%s", da)
		fmt.Printf("%s ", PROMT)
	}
}


