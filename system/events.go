package system

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/misthios/skypod"
)

/*
Describes something that can generate events (container,network,volume)
*/
type Actor struct{
	ID string
	Attributes map[string]string
}

/*
Contains the relevant/non deprecated parts of the event
*/
type Event struct{
	Actor	Actor	
	Action  string	
	Type    string
	Time    int64 
}

/*
Streams the events from the api into a gochannel
returns error
*/
func StreamEvents(ctx context.Context, eventchan chan Event) error{
	client, err := skypod.GetClient(ctx)
	if err !=nil {return err}

	response, err := client.Request(ctx,http.MethodGet,"/events",nil,nil,nil)
	if err !=nil {return err}

	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)

	for err == nil {
		var e Event
		err = dec.Decode(&e)
		if err == nil {
			eventchan <- e
		}
	}

	close(eventchan)
	return fmt.Errorf("Failed to stream events: %w ",err)
}
