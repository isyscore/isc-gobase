package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/listener"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestPublish(t *testing.T) {
	listener.AddListener("event1", Event1Lister1)
	listener.AddListener("event1", Event1Lister2)
	listener.AddListener("event1", Event1Lister3)

	listener.AddListener("event2", Event2Lister1)
	listener.AddListener("event2", Event2Lister2)

	e1 := Event1{Company: "公司"}
	listener.PublishEvent(e1)

	e2 := Event2{Address: "杭州"}
	listener.PublishEvent(e2)
}

func TestPublish1_1(t *testing.T) {
	listener.AddListener("event1", Event1Lister1)
	listener.AddListenerWithGroup("newGroup", "event1", Event1Lister1_1)

	e1 := Event1{Company: "公司"}
	listener.PublishEvent(e1)

	e2 := Event1_1{Company: "杭州"}
	listener.PublishEvent(e2)
}

func TestPublish2(t *testing.T) {
	listener.AddListener("event1", func(event listener.BaseEvent) {
		assert.Equal(t, "公司", event.(Event1).Company)
	})
	listener.AddListener("event1", func(event listener.BaseEvent) {
		assert.Equal(t, "公司", event.(Event1).Company)
	})
	listener.AddListener("event1", func(event listener.BaseEvent) {
		assert.Equal(t, "公司", event.(Event1).Company)
	})

	listener.PublishEvent(Event1{Company: "公司"})
}

func Event1Lister1(event listener.BaseEvent) {
	ev := event.(Event1)
	fmt.Println("Event1Lister1: " + ev.Company)
}

func Event1Lister1_1(event listener.BaseEvent) {
	ev := event.(Event1_1)
	fmt.Println("Event1Lister1_1: " + ev.Company)
}

func Event1Lister2(event listener.BaseEvent) {
	ev := event.(Event1)
	fmt.Println("Event1Lister2: " + ev.Company)
}

func Event1Lister3(event listener.BaseEvent) {
	ev := event.(Event1)
	fmt.Println("Event1Lister3: " + ev.Company)
}

func Event2Lister1(event listener.BaseEvent) {
	ev := event.(Event2)
	fmt.Println("Event2Lister1: " + ev.Address)
}

func Event2Lister2(event listener.BaseEvent) {
	ev := event.(Event2)
	fmt.Println("Event2Lister2: " + ev.Address)
}

type Event1 struct {
	Company string
}

func (e1 Event1) Name() string {
	return "event1"
}

func (e1 Event1) Group() string {
	return listener.DefaultGroup
}

type Event2 struct {
	Address string
}

func (e1 Event2) Name() string {
	return "event2"
}

func (e1 Event2) Group() string {
	return listener.DefaultGroup
}

type Event1_1 struct {
	Company string
}

func (e1 Event1_1) Name() string {
	return "event1"
}

func (e1 Event1_1) Group() string {
	return "newGroup"
}
