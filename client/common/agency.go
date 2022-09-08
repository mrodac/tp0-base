package common

import (
	"encoding/csv"
	"io"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

type Contestant struct {
	Document  uint32
	FirstName string
	LastName  string
	BirthDate string
}

type Agency struct {
	client        *Client
	interrupt     chan os.Signal
	retryDuration time.Duration
}

func NewAgency(client *Client) *Agency {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM)

	dur, _ := time.ParseDuration("10s")

	agency := &Agency{
		client:        client,
		interrupt:     interrupt,
		retryDuration: dur,
	}
	return agency
}

func (agency *Agency) CheckWinners() error {

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
	}

	msg := &ContestantsMessage{
		Contestants: slice,
	}

	response := agency.client.queryWinners(msg)

	if len(slice) > 0 {
		percentage := float64(len(response.Winners)) / float64(len(slice))
		log.Infof("Winner percentage: %f %%", 100*percentage)
	}

	return nil
}

func (agency *Agency) TotalWinners() error {

	response := agency.client.queryTotalWinners()

loop:
	for response.Pending > 0 {
		log.Infof("Partial Winner count: %d. Processing %d agencies", response.TotalWinners, response.Pending)

		select {
		case <-agency.interrupt:
			log.Infof("Got SIGTERM")
			response = nil
			break loop
		case <-time.After(agency.retryDuration):
		}
		response = agency.client.queryTotalWinners()
	}
	if response != nil {
		log.Infof("Total winner count: %d", response.TotalWinners)
	}

	return nil
}
