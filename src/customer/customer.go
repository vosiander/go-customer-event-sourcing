package customer

import (
	"eventsourcing"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

type Customer struct {
	eventStream *eventsourcing.EventStream

	customerId uuid.UUID
	firstname  string
	lastname   string
	createdAt  time.Time
	pain       string
}

func (customer *Customer) Replay(eventStream *eventsourcing.EventStream) {
	customer.eventStream = eventStream

	customer.mutate()
}

func (customer *Customer) apply(event eventsourcing.Event) {
	customer.eventStream.Add(event)

	customer.mutate()
}

func (customer *Customer) mutate() {
	stream := customer.eventStream.Stream()
	var err error

	for _, e := range stream {
		switch true {
		case e.Name() == "Event.CreateId":
			customer.customerId, _ = uuid.FromString(e.Payload()["customerId"].(string))
			customer.createdAt, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", e.Payload()["createdAt"].(string))
			if err != nil {
				log.Fatal(err)
			}
			break
		case e.Name() == "Event.ChangeName":
			customer.firstname = e.Payload()["firstname"].(string)
			customer.lastname = e.Payload()["lastname"].(string)
			break
		case e.Name() == "Event.ExperiencePain":
			customer.pain = e.Payload()["pain"].(string)
			break
		}
	}
}

func (customer *Customer) Stream() []eventsourcing.Event {
	return customer.eventStream.Stream()
}

func (customer *Customer) CreateId(customerId uuid.UUID) {
	customer.apply(eventsourcing.NewEvent(uuid.NewV4(), map[string]interface{}{
		"name":       "Event.CreateId",
		"customerId": customerId.String(),
		"createdAt":  customer.createdAt.String(),
	}))
}

func (customer *Customer) ChangeName(firstname string, lastname string) {
	customer.apply(eventsourcing.NewEvent(uuid.NewV4(), map[string]interface{}{
		"name":      "Event.ChangeName",
		"firstname": customer.firstname,
		"lastname":  customer.lastname,
	}))
}

func (customer *Customer) ExperiencePain(pain string) {
	customer.apply(eventsourcing.NewEvent(uuid.NewV4(), map[string]interface{}{
		"name": "Event.ExperiencePain",
		"pain": pain,
	}))
}