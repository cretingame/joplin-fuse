package joplin

import (
	"context"
	"log"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type Node interface {
	Base() NodeBase
	AddChild(n *Node)
}

type NodeBase struct {
	Id        string
	Parent_id string
	Name      string
	Children  []*Node
}

type FolderNode struct {
	Id        string
	Parent_id string
	Name      string
	Children  []*Node
}

func (fn FolderNode) Base() NodeBase {
	return NodeBase{
		Id:        fn.Id,
		Parent_id: fn.Parent_id,
		Name:      fn.Name,
		Children:  fn.Children,
	}
}

func (fn *FolderNode) AddChild(n *Node) {
	fn.Children = append(fn.Children, n)
}

type NoteNode struct {
	*fs.MemRegularFile

	Id        string
	Parent_id string
	Name      string
	Children  []*Node

	Session *Session
}

func (fn NoteNode) Base() NodeBase {
	return NodeBase{
		Id:        fn.Id,
		Parent_id: fn.Parent_id,
		Name:      fn.Name,
		Children:  fn.Children,
	}
}

func (fn *NoteNode) AddChild(n *Node) {
	fn.Children = append(fn.Children, n)
}

var _ = (fs.NodeReader)((*NoteNode)(nil))

func (fn *NoteNode) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	log.Println("get note:", fn.Name)
	noteResponse, err := GetNote(*fn.Session, fn.Id)
	if err != nil {
		return nil, fuse.FOPEN_KEEP_CACHE, syscall.EACCES
	}
	fn.MemRegularFile.Data = []byte(noteResponse.Body)
	return nil, fuse.FOPEN_KEEP_CACHE, fs.OK
}

func (fn *NoteNode) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	// TODO: Get Data from Joplin ?
	// need host and token
	return fn.MemRegularFile.Read(ctx, fh, dest, off)
}

type RessourceNode struct {
	Id        string
	Parent_id string
	Name      string
	Children  []*Node

	File *fs.MemRegularFile
}

func (rn RessourceNode) Base() NodeBase {
	return NodeBase{
		Id:        rn.Id,
		Parent_id: rn.Parent_id,
		Name:      rn.Name,
		Children:  rn.Children,
	}
}

func (rn *RessourceNode) AddChild(n *Node) {
	rn.Children = append(rn.Children, n)
}
