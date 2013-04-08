package balanced

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

const (
	eventsUri = "/v1/events"
)

type Event struct {
	CallbackStatuses CallbackStatuses `json:"callback_statuses, omitempty"`
	CallbackUri      string           `json:"callback_uri, omitempty"`
	Entity           BankAccount      `json:"entity, omitempty"`
	Id               string           `json:"id, omitempty"`
	OccurredAt       time.Time        `json:"occurred_at, omitempty"`
	Type             string           `json:"type, omitempty"`
	Uri              string           `json:"uri, omitempty"`
}

type CallbackStatuses struct {
	Failed    int `json:"failed, omitempty"`
	Pending   int `json:"pending, omitempty"`
	Retrying  int `json:"retrying, omitempty"`
	Succeeded int `json:"succeeded, omitempty"`
}

type ListOfEvents struct {
	FirstUri    string  `json:"first_uri, omitempty"`
	Items       []Event `json:"items, omitempty"`
	LastUri     string  `json:"last_uri, omitempty"`
	Limit       int     `json:"limit, omitempty"`
	NextUri     string  `json:"next_uri, omitempty"`
	Offset      int     `json:"offset, omitempty"`
	PreviousUri string  `json:"previous_uri, omitempty"`
	Total       int     `json:"total, omitempty"`
	Uri         string  `json:"uri, omitempty"`
}

// Retrieves the details of an event that was previously created. Use the uri
// that was previously returned, and the corresponding event information will be
// returned.
func RetrieveEvent(uri string, limit, offset int) (*Event, error) {
	payload := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}

	resp, err := get(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Event
	event := Event{}
	if err := json.Unmarshal(resp, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

func ListAllEvents(limit, offset int) (*ListOfEvents, error) {
	payload := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}

	resp, err := get(eventsUri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into ListOfEvents
	listOfEvents := ListOfEvents{}
	if err := json.Unmarshal(resp, &listOfEvents); err != nil {
		return nil, err
	}

	return &listOfEvents, nil
}
