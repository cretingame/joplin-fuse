package main

import (
	"fmt"
	"joplin-fuse/internal/joplin"
	"log"
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

	id := "73ca02a59f2741bcb98a77516b85b9a5"

	bs, err := joplin.GetRessourceFile(*ses, id)
	if err != nil {
		panic(err)
	}
	log.Println(bs)

	resp, err := joplin.GetRessource(*ses, id)
	if err != nil {
		panic(err)
	}
	log.Println(resp)
}
