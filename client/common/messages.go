package common

import (
	"bufio"
	"encoding/binary"
)

const (
	Query    byte = 0
	Response byte = 1
)

type Message struct {
	Type byte
}

type ContestantsQuery struct {
	Contestants []Contestant
}

type WinnerResponse struct {
	Winners []uint32
}

func ContestantToBytes(c *Contestant, writer *bufio.Writer) {
	binary.Write(writer, binary.BigEndian, c.Document)
	writer.WriteString(c.FirstName)
	writer.WriteByte(0)
	writer.WriteString(c.LastName)
	writer.WriteByte(0)
	writer.WriteString(c.BirthDate)
	writer.WriteByte(0)
}

func ContestantFromBytes(reader *bufio.Reader) *Contestant {
	c := &Contestant{}
	binary.Read(reader, binary.BigEndian, &c.Document)
	c.FirstName, _ = reader.ReadString(0)
	c.LastName, _ = reader.ReadString(0)
	c.BirthDate, _ = reader.ReadString(0)
	return c
}

func WriteContestantsMessage(msg *ContestantsQuery, writer *bufio.Writer) {
	writer.WriteByte(Query)
	binary.Write(writer, binary.BigEndian, int32(len(msg.Contestants)))

	for _, s := range msg.Contestants {
		ContestantToBytes(&s, writer)
	}
}

func ReadWinnerMessage(reader *bufio.Reader) *WinnerResponse {
	reader.ReadByte()

	msg := &WinnerResponse{}

	var len int32
	binary.Read(reader, binary.BigEndian, &len)

	msg.Winners = make([]uint32, len)

	var i int32
	for i = 0; i < len; i++ {
		var document uint32
		binary.Read(reader, binary.BigEndian, &document)
		msg.Winners[i] = document
	}

	return msg
}
