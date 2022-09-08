package common

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Contestant struct {
	Document  uint32
	FirstName string
	LastName  string
	BirthDate string
}

type Agency struct {
	client *Client
}

func NewAgency(client *Client) *Agency {
	domain := &Agency{
		client: client,
	}
	return domain
}

func (domain *Agency) CheckWinners() error {
	dni, _ := strconv.ParseUint(os.Getenv("DOCUMENT"), 10, 32)
	contestant := Contestant{
		Document:  uint32(dni),
		FirstName: os.Getenv("FIRST_NAME"),
		LastName:  os.Getenv("LAST_NAME"),
		BirthDate: os.Getenv("BIRTH_DATE"),
	}

	log.Debugf("[CLIENT 1] Contestant: %d %s %s %s", contestant.Document, contestant.FirstName, contestant.LastName, contestant.BirthDate)

	slice := make([]Contestant, 1)
	slice[0] = contestant

	msg := &ContestantsQuery{
		Contestants: slice,
	}

	res := domain.client.queryWinners(msg)

	if len(res.Winners) > 0 {
		log.Infof("Es ganador")
	} else {
		log.Infof("No es ganador")
	}
	return nil
}
