package main

import (
	"fmt"
	"joplin-fuse/internal/joplin"
	"time"
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

	n, err := joplin.GetNote(*ses, noteId)
	if err != nil {
		panic(err)
	}

	fmt.Println("note title:", n.Title)
	fmt.Println("note body:", n.Body)
	fmt.Println("note updated_time:", n.Updated_time)

	t := time.Unix(int64(n.Updated_time/1000), int64((n.Updated_time%1000)*1000_000))
	fmt.Println("note updated date", t)
}
