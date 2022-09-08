package common

import (
	"encoding/csv"
	"io"
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

	dataset, err := os.Open("/dataset.csv")
	if err != nil {
		log.Fatal("Unable to read dataset file", err)
		return err
	}
	defer dataset.Close()

	csvReader := csv.NewReader(dataset)
	csvReader.ReuseRecord = true
	csvReader.FieldsPerRecord = 4

	slice := make([]Contestant, 0)

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err != io.EOF {
				log.Errorf("Error reading csv %s", err)
				return err
			}
			break
		}

		document, err := strconv.ParseUint(record[2], 10, 32)
		if err != nil {
			log.Errorf("Error parsing document field %s", err)
			return err
		}

		contestant := Contestant{
			Document:  uint32(document),
			FirstName: record[0],
			LastName:  record[1],
			BirthDate: record[3],
		}

		slice = append(slice, contestant)

		//log.Debugf("[CLIENT 1] Contestant: %d %s %s %s", contestant.Document, contestant.FirstName, contestant.LastName, contestant.BirthDate)
	}

	msg := &ContestantsQuery{
		Contestants: slice,
	}

	res := domain.client.queryWinners(msg)

	if len(slice) > 0 {
		percentage := float64(len(res.Winners)) / float64(len(slice))
		log.Infof("Winner percentage: %f", 100*percentage)
	}

	return nil
}
