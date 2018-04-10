package hub

import (
	"fmt"
	"time"

	"github.com/floeit/floe/config"
	nt "github.com/floeit/floe/config/nodetype"
	"github.com/floeit/floe/event"
	"github.com/floeit/floe/log"
)

type timer struct {
	flow   config.FlowRef
	nodeID string
	period int // time between triggers in seconds
	next   time.Time
	opts   nt.Opts
}

type timers struct {
	list map[string]*timer
}

func newTimers(h *Hub) *timers {
	t := &timers{
		list: map[string]*timer{},
	}

	go func() {
		for now := range time.Tick(time.Second) {
			for id, tim := range t.list {
				if !now.After(tim.next) {
					continue
				}
				fmt.Println("---------------------- ping", id)
				tim.next = now.Add(time.Duration(tim.period) * time.Second)

				e := event.Event{
					RunRef: event.RunRef{
						FlowRef: tim.flow,
					},
					SourceNode: config.NodeRef{
						Class: "trigger",
						ID:    tim.nodeID,
					},
					Tag:  "timer", // to match the timer trigger
					Good: true,
					Opts: nt.Opts{
						"period": tim.period,
					},
				}

				ref, err := h.addToPending(tim.flow, h.hostID, e)
				if err != nil {
					log.Errorf("<%s> - from timer trigger did not add to pending: %s", e.RunRef, err)
				}
				log.Debugf("<%s> - from timer trigger added to pending", ref)
			}
		}
	}()
	return t
}

func (t *timers) register(flow config.FlowRef, nodeID string, opts nt.Opts) {
	t.list[flow.String()+"-"+nodeID] = &timer{
		flow:   flow,
		nodeID: nodeID,
		period: opts["period"].(int),
	}
}