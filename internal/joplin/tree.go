package joplin

import (
	"context"
	"log"
	"regexp"
	"strings"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type JoplinRoot struct {
	fs.Inode
	items []Node
}

var _ = (fs.NodeGetattrer)((*JoplinRoot)(nil))
var _ = (fs.NodeOnAdder)((*JoplinRoot)(nil))

func NewRoot(host string, tokenLocation string) (JoplinRoot, error) {
	var items []Node

	ses, err := Authenticate(host, tokenLocation)
	if err != nil {
		return JoplinRoot{}, err
	}

	folders, err := GetItems(*ses, "folders")
	if err != nil {
		return JoplinRoot{}, err
	}
	for i := range folders {
		name := folders[i].Title
		name = removeSpecialCharacters(name)
		name = sanitizeFilename(name)
		folderNode := FolderNode{
			Id:        folders[i].Id,
			Parent_id: folders[i].Parent_id,
			Name:      name,
		}
		items = append(items, &folderNode)
	}

	notes, err := GetItems(*ses, "notes")
	if err != nil {
		return JoplinRoot{}, err
	}
	for i := range notes {
		name := notes[i].Title
		name = removeSpecialCharacters(name)
		name = sanitizeFilename(name)
		name = name + ".md"

		noteNode := NoteNode{
			Id:        notes[i].Id,
			Parent_id: notes[i].Parent_id,
			Name:      name,
			MemRegularFile: &fs.MemRegularFile{
				Data: []byte{},
				Attr: fuse.Attr{
					Mode:  0644,
					Owner: *fuse.CurrentOwner(),
				},
			},
			Session: ses,
		}

		items = append(items, &noteNode)
	}

	resourceFolderNode := FolderNode{
		Id:        ":",
		Parent_id: "",
		Name:      ":",
	}
	items = append(items, &resourceFolderNode)

	resources, err := GetItems(*ses, "resources")
	if err != nil {
		return JoplinRoot{}, err
	}
	for i := range resources {
		ressourceNode := RessourceNode{
			Id:        resources[i].Id,
			Parent_id: resourceFolderNode.Id,
			Name:      resources[i].Id,
			MemRegularFile: &fs.MemRegularFile{
				Data: []byte{},
				Attr: fuse.Attr{
					Mode:  0444,
					Owner: *fuse.CurrentOwner(),
				},
			},
			Session: ses,
		}

		items = append(items, &ressourceNode)
	}

	return JoplinRoot{
		items: items,
	}, nil
}

func (r *JoplinRoot) OnAdd(ctx context.Context) {
	tree := BuildTree(r.items)

	addNode(ctx, &r.Inode, tree, 0)

	log.Println("OnAdd finished")
}

func (r *JoplinRoot) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = 0755
	return 0
}

func addNode(ctx context.Context, parentInode *fs.Inode, items []*Node, level int) {
	for i := range items {
		child := items[i]

		switch v := (*child).(type) {
		case *FolderNode:
			childInode := parentInode.NewPersistentInode(
				ctx, &fs.Inode{}, fs.StableAttr{Mode: syscall.S_IFDIR})

			parentInode.AddChild(v.Name, childInode, false)

			// NOTE: in progress
			link := ".."
			for i := 0; i < level; i++ {
				link = link + "/.."
			}
			link = link + "/:"
			l := &fs.MemSymlink{
				Data: []byte(link),
			}
			symInode := childInode.NewPersistentInode(ctx, l, fs.StableAttr{Mode: syscall.S_IFLNK})
			childInode.AddChild(":", symInode, false)
			addNode(ctx, childInode, v.Children, level+1)
		case *NoteNode:
			childInode := parentInode.NewPersistentInode(
				ctx, v, v.StableAttr())

			parentInode.AddChild(v.Name, childInode, false)
			addNode(ctx, childInode, v.Children, level+1)
		case *RessourceNode:
			childInode := parentInode.NewPersistentInode(
				ctx, v, v.StableAttr())

			parentInode.AddChild(v.Name, childInode, false)
			addNode(ctx, childInode, v.Children, level+1)
		default:
			panic("not handled")
		}
	}
}

// sanitizeFilename removes or replaces characters not allowed in filenames
func sanitizeFilename(name string) string {
	// Replace forbidden characters with underscore
	// Forbidden: \ / : * ? " < > | (Windows)
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	name = invalidChars.ReplaceAllString(name, "_")

	// Trim spaces and dots (Windows does not allow filenames ending with them)
	name = strings.Trim(name, " .")

	// Avoid reserved filenames on Windows (case-insensitive)
	reserved := map[string]bool{
		"CON": true, "PRN": true, "AUX": true, "NUL": true,
		"COM1": true, "COM2": true, "COM3": true, "COM4": true,
		"COM5": true, "COM6": true, "COM7": true, "COM8": true, "COM9": true,
		"LPT1": true, "LPT2": true, "LPT3": true, "LPT4": true,
		"LPT5": true, "LPT6": true, "LPT7": true, "LPT8": true, "LPT9": true,
	}
	upper := strings.ToUpper(name)
	if reserved[upper] {
		name = "_" + name
	}

	// Limit length (255 bytes is a common limit)
	if len(name) > 255 {
		name = name[:255]
	}

	if name == "" {
		return "unnamed"
	}

	return name
}

func removeSpecialCharacters(input string) string {
	var sb strings.Builder
	for _, r := range input {
		if r < 0xFFFF {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
