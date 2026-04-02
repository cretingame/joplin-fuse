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

var _ = (fs.NodeOpener)((*NoteNode)(nil))
var _ = (fs.NodeReader)((*NoteNode)(nil))
var _ = (fs.NodeWriter)((*NoteNode)(nil))

func (fn *NoteNode) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	log.Println("get note:", fn.Name)
	noteResponse, err := GetNote(*fn.Session, fn.Id)
	if err != nil {
		return nil, fuse.FOPEN_KEEP_CACHE, syscall.EIO
	}
	fn.MemRegularFile.Data = []byte(noteResponse.Body)
	return fn.MemRegularFile.Open(ctx, flags)
}

func (fn *NoteNode) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	// TODO: Check integrity
	return fn.MemRegularFile.Read(ctx, fh, dest, off)
}

func (fn *NoteNode) Write(ctx context.Context, fh fs.FileHandle, data []byte, off int64) (uint32, syscall.Errno) {
	written, errno := fn.MemRegularFile.Write(ctx, fh, data, off)

	err := PutNoteBody(*fn.Session, fn.Id, string(fn.MemRegularFile.Data))
	if err != nil {
		return 0, syscall.EIO
	}

	return written, errno
}

type RessourceNode struct {
	*fs.MemRegularFile

	Id        string
	Parent_id string
	Name      string
	Children  []*Node

	Session *Session
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

var _ = (fs.NodeOpener)((*RessourceNode)(nil))
var _ = (fs.NodeReader)((*RessourceNode)(nil))

func (rn *RessourceNode) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	log.Println("get ressource:", rn.Name)
	ressourceBytes, err := GetRessourceFile(*rn.Session, rn.Id)
	if err != nil {
		return fh, fuseFlags, syscall.EIO
	}
	rn.Data = ressourceBytes
	return rn.MemRegularFile.Open(ctx, flags)
}

func (rn *RessourceNode) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	return rn.MemRegularFile.Read(ctx, fh, dest, off)
}
