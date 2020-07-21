package can

// ---------------------------------------------------------------------------------------
//  imports
// ---------------------------------------------------------------------------------------

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

type Frame struct {
	Id     int
	Length int
	Data   [8]byte

	Timestamp time.Time
}

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

func ParseFrame(cmd string) (*Frame, error) {
	msg := Frame{}

	if len(cmd) < 5 {
		return nil, errors.New("cmd not long enough")
	}

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

	if len(cmd) < 5+msg.Length*2 {
		return nil, errors.New("cmd not long enough")
	}

	for i := 0; i < msg.Length; i++ {
		d, err := strconv.ParseUint(cmd[5+2*i:5+2*i+2], 16, 8)
		if err != nil {
			return nil, err
		}
		msg.Data[i] = byte(d)
	}

	return &msg, nil
}

// ---------------------------------------------------------------------------------------
//  public members
// ---------------------------------------------------------------------------------------

func (f *Frame) ToUint64() uint64 {
	u := uint64(0)

	for i := 0; i < f.Length; i++ {

		u |= uint64(f.Data[i]) << (i * 8)
	}

	return u
}

func (f *Frame) String() string {
	str := fmt.Sprintf("%x  [%d] ", f.Id, f.Length)

	for i := 0; i < f.Length; i++ {
		str += fmt.Sprintf("%02x ", f.Data[i])
	}

	return str
}
