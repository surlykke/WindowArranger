package sway

import (
	"encoding/binary"
	"encoding/json"
	"net"
	"os"
)

type Response struct {
	Success     bool
	Parse_error bool
	Error       string
}

var conn net.Conn = nil

func Execute(cmd string) (responses []Response) {
	execute(cmd, 0, &responses)
	return
}

type Output struct {
	Id     int
	Name   string
	Make   string
	Model  string
	Serial string
}

func GetOutputs() (outputs []Output) {
	execute("", 3, &outputs)
	return
}

func execute(cmd string, cmdType uint32, ret any) {
	if conn == nil {
		var err error
		if swaysock := os.Getenv("SWAYSOCK"); swaysock == "" {
			panic("SWAYSOCK not set")
		} else if conn, err = net.Dial("unix", swaysock); err != nil {
			panic(err)
		}
	}

	var (
		msg          = []byte("i3-ipc")
		respFromSway = make([]byte, 20000)
	)
	msg = binary.NativeEndian.AppendUint32(msg, uint32(len([]byte(cmd))))
	msg = binary.NativeEndian.AppendUint32(msg, cmdType)
	msg = append(msg, []byte(cmd)...)
	if _, err := conn.Write(msg); err != nil {
		panic(err)
	} else if n, err := conn.Read(respFromSway); err != nil {
		panic(err)
	} else {
		if err := json.Unmarshal(respFromSway[14:n], &ret); err != nil {
			panic(err)
		}
	}

}
