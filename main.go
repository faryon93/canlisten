package main

// ---------------------------------------------------------------------------------------
//  imports
// ---------------------------------------------------------------------------------------

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"canlisten/can"
	"canlisten/filter"

	"github.com/tarm/serial"
)

// ---------------------------------------------------------------------------------------
//  constants
// ---------------------------------------------------------------------------------------

const (
	CmdCanTransmit = 't'
)

// ---------------------------------------------------------------------------------------
//  global variables
// ---------------------------------------------------------------------------------------

var (
	Device     string
	FilterExpr string
	Running    = true
)

// ---------------------------------------------------------------------------------------
//  helper functions
// ---------------------------------------------------------------------------------------

func ReadNextCommand(s *serial.Port) (string, error) {
	buf := make([]byte, 1)
	str := bytes.NewBufferString("")

	for Running {
		n, err := s.Read(buf)
		if err == io.EOF {
			continue
		} else if err != nil {
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

	return "", errors.New("empty response")
}

func GetPrintableCmd(cmd string) string {
	cmd = strings.Replace(cmd, "\r", "!", -1)
	cmd = strings.Replace(cmd, "\a", "?", -1)

	return cmd
}

// ---------------------------------------------------------------------------------------
//  application entry
// ---------------------------------------------------------------------------------------

func main() {
	flag.StringVar(&Device, "dev", "/dev/ttyACM0", "")
	flag.StringVar(&FilterExpr, "filter", "", "")
	flag.Parse()

	c := &serial.Config{Name: Device, Baud: 115200, ReadTimeout: 50 * time.Millisecond}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// Open CAN interface with 1MBit/s
	_, err = s.Write([]byte("C\rS8\rO\r"))
	if err != nil {
		log.Println("failed to open serial port")
	}

	// gracefull shutdown on sigint/sigterm
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		Running = false
		log.Println("closing application")
	}()

	log.Printf("Using filter: \"%s\"", FilterExpr)
	filterVm := filter.NewFilter(FilterExpr)

	// print all can messages to stdout
	stdout := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	canCounter := 0
	for Running {
		cmd, err := ReadNextCommand(s)
		if err != nil {
			log.Println("ReadNextCommand failed:", err.Error())
			continue
		}

		// ignore all non can messages
		if cmd[0] != CmdCanTransmit {
			log.Printf("ignoring command \"%s\": unexpected opcode", GetPrintableCmd(cmd))
			continue
		}

		frame, err := can.ParseFrame(cmd)
		if err != nil {
			log.Printf("failed to parse message \"%s\": %s", GetPrintableCmd(cmd), err.Error())
			continue
		}
		frame.Timestamp = time.Now()

		// apply the filterVm
		show, err := filterVm.Eval(frame)
		if err != nil {
			log.Println("failed to evaluate filterVm:", err.Error())
			continue
		}

		if !show {
			continue
		}

		// count message and print to stdout
		canCounter++
		stdout.Printf("%s  %s", filepath.Base(Device), frame.String())
	}

	log.Printf("statistics: rx=%d", canCounter)
}
