package parser

import (
	"bytes"
	"fmt"
	"strings"
	"sort"
	"strconv"
	"os"
)

const (
	UBUNTU_WELCOME = "Welcome to Ubuntu"
	GDB = "(gdb) "
	XGDB = "(xgdb) "
	EMPTY = ""
	SHELL = "~shell ** "

	KNRM = "\x1B[0m"
	KRED = "\x1B[31m"
	KGRN = "\x1B[32m"
	KYEL = "\x1B[33m"
	KBLU = "\x1B[34m"
	KMAG = "\x1B[35m"
	KCYN = "\x1B[36m"
	KWHT = "\x1B[37m"

	VERSION = "0.0.1"
)

var PROMPT = ""
var PrintPrompt = 1

type RawData struct {
	PID int
	Data string
}

func ParseStdoutStream(data []*RawData) {
	accumulateOutput := make(map[string][]int)

	for _, rawData := range data {
		accumulateOutput[rawData.Data] = append(accumulateOutput[rawData.Data], rawData.PID)	
	}

	for output, pid := range accumulateOutput {
		s, err := parseGeneric(output)
		if err != nil {
			panic(err)
		}

		if s != "" {
			fmt.Println("\n" + KGRN + "PID:" + KRED + "[" + formatPID(pid) + "]" + KNRM)
			fmt.Printf("%s\n", s)

			if PrintPrompt == 0 {
				fmt.Printf("%s", PROMPT)
			}
		}
		
	}
}

func ParseErrStream(data []*RawData) {
	accumulateOutput := make(map[string][]int)

	for _, rawData := range data {
		accumulateOutput[rawData.Data] = append(accumulateOutput[rawData.Data], rawData.PID)

		
	}

	for output, pid := range accumulateOutput {
		s, err := parseGeneric(output)
		if err != nil {
			panic(err)
		}

		if s != "" {
			fmt.Println("\n" + KGRN + "PID:" + KRED + "[" + formatPID(pid) + "]" + KNRM)
			fmt.Printf("%s\n", s)

			if PrintPrompt == 0 {
				fmt.Printf("%s", PROMPT)
			}
		}
		
	}
}

func parseGeneric(data string) (string, error) {
	if strings.Index(data, UBUNTU_WELCOME) > -1 {
		return "", nil
	}
	newData := strings.Replace(data, GDB, EMPTY, -1)
	return newData, nil
}

func formatPID(pidArray []int) string {

	if len(pidArray) == 1 {
		return strconv.Itoa(pidArray[0])
	}

	sort.Ints(pidArray)

	if len(pidArray) == 2 {
		return strconv.Itoa(pidArray[0]) + "..." + strconv.Itoa(pidArray[1])
	}

	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(pidArray[0]))

	for i := 1; i < len(pidArray); i++ {
		if (pidArray[i-1] + 1) != pidArray[i] {
			buffer.WriteString("..." + strconv.Itoa(pidArray[i]))
		}
	}

	return buffer.String()
}

func PreProcess(data []byte) bool {
	s := string(data[:len(data)-1])
	if strings.Compare(s, "version") == 0 {
		fmt.Printf("Version: %s\n\n%s", VERSION, PROMPT)
		return true
	} else if strings.Compare(s, "quit") == 0 {
		os.Exit(1)
	}
	return false
}
