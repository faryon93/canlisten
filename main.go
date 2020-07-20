package main

import (
	"bytes"
	"flag"
	"github.com/tarm/serial"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	Device string
)

func read(s *serial.Port) (string, error) {
	buf := make([]byte, 1)
	str := bytes.NewBufferString("")

	for {
		n, err := s.Read(buf)
		if err != nil {
			return "", err
		}

		if n < 1 {
			continue
		}

		err = str.WriteByte(buf[0])
		if err != nil {
			return "", err
		}

		if buf[0] == '\r' || buf[0] == '\a' {
			return str.String(), nil
		}
	}
}

func main() {
	flag.StringVar(&Device, "dev", "/dev/ttyACM0", "")
	flag.Parse()

	c := &serial.Config{Name: Device, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.Write([]byte("C\rS8\rO\r"))
	if err != nil {
		log.Println("failed to open serial port")
	}

	stdout := log.New(os.Stdout, "", log.LstdFlags | log.Lmicroseconds)
	for {
		cmd, err := read(s)
		if err != nil {
			log.Println("read failed:", err.Error())
			continue
		}

		if cmd[0] != 't' {
			continue
		}

		can, err := ParseMsgCan(cmd)
		if err != nil {
			cmd = strings.Replace(cmd, "\r", "!", -1)
			cmd = strings.Replace(cmd, "\a", "?", -1)
			log.Printf("failed to parse message \"%s\": %s", cmd, err.Error())
			continue
		}

		stdout.Printf("%s  %s", filepath.Base(Device), can.String())
	}
}
