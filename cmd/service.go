// Copyright © 2017 Douglas Chimento <dchimento@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
)

var labels = make(map[string]prometheus.Labels)

var recordCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "passwdpot",
	Name:      "record",
	Help:      "count of requests",
	Subsystem: "total",
}, []string{"origin"})

type EventService interface {
	Record(ctx context.Context, event Event) (int64, error)
}

type eventService struct {
	eventRecorder
}

func NewEventService(er eventRecorder) EventService {
	return &eventService{er}
}

func (es *eventService) Record(ctx context.Context, event Event) (int64, error) {
	var l prometheus.Labels
	var ok bool
	if l, ok = labels[event.OriginAddr]; !ok {
		l = prometheus.Labels{"origin": event.OriginAddr}
	}
	recordCounter.With(l).Inc()
	return es.eventRecorder.recordEvent(event)
}

var (
	ErrNotFound = errors.New("not found")
)

func init() {
	prometheus.MustRegister(recordCounter)
}
