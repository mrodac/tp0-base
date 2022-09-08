package common

import (
	"bufio"
	"encoding/binary"

	log "github.com/sirupsen/logrus"
)

const (
	WinnerQuery         byte = 0
	WinnerResponse      byte = 1
	TotalWinnerQuery    byte = 2
	TotalWinnerResponse byte = 3
)

type Message struct {
	Type byte
}

type ContestantsMessage struct {
	Contestants []Contestant
}

type WinnersMessage struct {
	Winners []uint32
}

type TotalWinnersMessage struct {
	TotalWinners uint32
	Pending      byte
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

func WriteContestantsMessage(msg *ContestantsMessage, writer *bufio.Writer) {
	log.Infof("Sending message: %d", WinnerQuery)

	writer.WriteByte(WinnerQuery)
	binary.Write(writer, binary.BigEndian, int32(len(msg.Contestants)))

	for _, s := range msg.Contestants {
		ContestantToBytes(&s, writer)
	}
}

func WriteTotalWinnerMessage(writer *bufio.Writer) {
	log.Infof("Sending message: %d", TotalWinnerQuery)
	writer.WriteByte(TotalWinnerQuery)
}

func ReadWinnerMessage(reader *bufio.Reader) *WinnersMessage {
	reader.ReadByte()

	msg := &WinnersMessage{}

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

func ReadTotalWinnersMessage(reader *bufio.Reader) *TotalWinnersMessage {
	reader.ReadByte()

	msg := &TotalWinnersMessage{}

	binary.Read(reader, binary.BigEndian, &msg.TotalWinners)
	b, _ := reader.ReadByte()
	msg.Pending = b

	return msg
}
