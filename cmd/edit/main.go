package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"joplin-fuse/internal/joplin"
	"net/http"
)

const (
	tokenLocation = "./token"
	host          = "http://localhost:41184"
)

func main() {
	ses, err := joplin.Authenticate(host, tokenLocation)
	if err != nil {
		panic(err)
	}
	fmt.Printf("token <%s>\n", ses.Token)

	noteId := "cee90765acbd4297a09c31657c105581"

	n, err := GetNote(*ses, noteId)
	if err != nil {
		panic(err)
	}

	fmt.Println("note title:", n.Title)
	fmt.Println("note body:", n.Body)
	fmt.Println("note updated_time:", n.Updated_time)
	fmt.Printf("note: %+v\n", n)

	n.Body = n.Body + "\n\nOK!"

	err = PutNote(*ses, noteId, n)
	if err != nil {
		panic(err)
	}
}

func PutNote(ses joplin.Session, id string, note joplin.NoteResponse) (err error) {
	url := fmt.Sprintf("%s/notes/%s?token=%s", ses.Host, id, ses.Token)

	// 1. Prepare the payload
	payload := map[string]any{
		"body": note.Body,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// 2. Create the request
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	// 3. Set headers
	req.Header.Set("Content-Type", "application/json")

	// 4. Execute
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 5. Read response
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\nBody: %s\n", resp.StatusCode, respBody)

	return
}

func GetNote(ses joplin.Session, id string) (note joplin.NoteResponse, err error) {
	req := fmt.Sprintf("%s/notes/%s?token=%s&fields=title,body", ses.Host, id, ses.Token)
	response, err := http.Get(req)
	if err != nil {
		return
	}

	bs, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bs, &note)
	if err != nil {
		return
	}

	return
}
