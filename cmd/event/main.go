package main

import (
	"encoding/json"
	"fmt"
	"io"
	"joplin-fuse/internal/joplin"
	"net/http"
	"strconv"
)

const (
	tokenLocation = "./token"
	host          = "http://localhost:41184"
)

// Goal: I want to catch Joplin events
// https://joplinapp.org/help/api/references/rest_api/#events

func main() {
	ses, err := joplin.Authenticate(host, tokenLocation)
	if err != nil {
		panic(err)
	}
	fmt.Printf("token <%s>\n", ses.Token)

	cursor, err := GetEventCursor(*ses)
	if err != nil {
		panic(err)
	}
	fmt.Println("cursor:", cursor)
	cursor = cursor - 10

	events, err := GetEvents(*ses, cursor)
	if err != nil {
		panic(err)
	}

	fmt.Println(events)

	for i, ev := range events {
		fmt.Printf("%d: %v\n", i, ev)
	}
}

// OPTIM: cursor as int instead of string
func GetEventCursor(ses joplin.Session) (cursor int64, err error) {

	req := fmt.Sprintf("%s/events?token=%s", ses.Host, ses.Token)
	response, err := http.Get(req)
	if err != nil {
		return 0, err
	}

	bs, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	fmt.Println("resp:", string(bs))

	var jPage EventPageResponse
	err = json.Unmarshal(bs, &jPage)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(jPage.Cursor, 10, 64)
}

func GetEvents(ses joplin.Session, cursor int64) (events []EventResponse, err error) {
	hasMore := true
	page := 0

	for hasMore {
		req := fmt.Sprintf("%s/events?token=%s&cursor=%d&page=%d", ses.Host, ses.Token, cursor, page)
		if page == 0 {
			req = fmt.Sprintf("%s/events?token=%s&cursor=%d", ses.Host, ses.Token, cursor)
		}
		response, err := http.Get(req)
		if err != nil {
			return events, err
		}

		bs, err := io.ReadAll(response.Body)
		if err != nil {
			return events, err
		}

		var jPage EventPageResponse
		err = json.Unmarshal(bs, &jPage)
		if err != nil {
			return events, err
		}

		hasMore = jPage.Has_more

		events = append(events, jPage.Items...)
		page++
	}

	return events, err
}

type EventPageResponse struct {
	Has_more bool
	Cursor   string
	Items    []EventResponse
}

type EventResponse struct {
	Id                 int
	Item_type          int         // The item type (see table above for the list of item types)
	Item_id            string      // The item ID
	Type               int         // The type of change - either 1 (created), 2 (updated) or 3 (deleted)
	Created_time       joplin.Time // When the event was generated
	Source             int         // Unused
	Before_change_item string      // Unused
}

/*
id	int
item_type	int	The item type (see table above for the list of item types)
item_id	text	The item ID
type	int	The type of change - either 1 (created), 2 (updated) or 3 (deleted)
created_time	int	When the event was generated
source	int	Unused
before_change_item	text	Unused
*/
