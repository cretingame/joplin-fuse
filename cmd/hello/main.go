package main

import (
	"fmt"
	"joplin-fuse/internal/joplin"
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

	folders, err := joplin.GetItems(*ses, "folders")
	if err != nil {
		panic(err)
	}
	for i, folder := range folders {
		fmt.Println(i, folder)
	}

	notes, err := joplin.GetItems(*ses, "notes")
	if err != nil {
		panic(err)
	}
	for i, item := range notes {
		fmt.Println(i, item)
	}

	note, err := joplin.GetNote(*ses, notes[0].Id)
	if err != nil {
		panic(err)
	}
	fmt.Println(note.Body)

	var nodes []joplin.Node

	tree := joplin.BuildTree(nodes)
	joplin.PrintTree(tree, 0)
}
