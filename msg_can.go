package main

import (
	"fmt"
	"strconv"
)

type MsgCan struct {
	Id int
	Length int
	Data [8]byte
}

func ParseMsgCan(cmd string) (*MsgCan, error) {
	msg := MsgCan{}

	id, err := strconv.ParseUint(cmd[1:4], 16, 16)
	if err != nil {
		return nil, err
	}
	msg.Id = int(id)

	length, err := strconv.ParseUint(cmd[4:5], 10, 8)
	if err != nil {
		return nil, err
	}
	msg.Length = int(length)

	for i := 0; i < msg.Length; i++ {
		d, err := strconv.ParseUint(cmd[5+2*i:5+2*i+2], 16, 8)
		if err != nil {
			return nil, err
		}
		msg.Data[i] = byte(d)
	}

	return &msg, nil
}

func (m *MsgCan) String() string {
	str := fmt.Sprintf("%x  [%d] ", m.Id, m.Length)

	for i := 0; i < m.Length; i++ {
		str += fmt.Sprintf("%02x ", m.Data[i])
	}

	return str
}