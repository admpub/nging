package fileindex

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func NewMemIndex(root string, id string, filter Filter) (Index, error) {
	rootPath := filepath.Clean(root)
	fi, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%v is not a directory", root)
	}
	idx := &memIndex{id, rootPath, nil}
	go func() {
		for {
			go idx.update()
			<-time.After(300 * time.Second)
		}
	}()
	return idx, nil
}

/* Index Entry */

func entryId(path string, fi os.FileInfo) string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%v\n", path)))
	h.Write([]byte(fmt.Sprintf("%v\n", fi.Size())))
	h.Write([]byte(fmt.Sprintf("%v\n", fi.ModTime())))
	h.Write([]byte(fmt.Sprintf("%v\n", fi.IsDir())))
	return fmt.Sprintf("%x", h.Sum(nil))
}

type memEntry struct {
	index    *memIndex
	id       string
	path     string
	parentId string
	isDir    bool
}

func newMemEntry(index *memIndex, path string, fi os.FileInfo, parentId string) Entry {
	id := entryId(path, fi)
	return &memEntry{index, id, path, parentId, fi.IsDir()}
}

func (e *memEntry) Id() string {
	return e.id
}

func (e *memEntry) Name() string {
	return filepath.Base(e.path)
}

func (e *memEntry) IsDir() bool {
	return e.isDir
}

func (e *memEntry) ParentId() string {
	return e.parentId
}

func (e *memEntry) Path() string {
	return path.Join(e.index.root, e.path)
}

/* Index Entry */

type memIndex struct {
	id   string
	root string
	data *memIndexData
}

func (i *memIndex) Id() string {
	return i.id
}

func (i *memIndex) Root() string {
	return i.id
}

func (i *memIndex) Get(id string) (Entry, error) {
	return i.data.entries[id], nil
}

func (i *memIndex) WaitForReady() error {
	for {
		if i.data != nil {
			return nil
		}
		<-time.After(1 * time.Second)
	}
}

func (i *memIndex) List(parent string) ([]Entry, error) {
	return i.data.children[parent], nil
}

func (i *memIndex) updateDir(d *memIndexData, path string, parentId string) error {
	dir, err := readDirInfo(path)
	if err != nil {
		return err
	}

	dirEntry, err := i.entryFromInfo(dir.info, dir.path, parentId)
	if err != nil {
		return err
	}

	d.add(dirEntry.(*memEntry))

	err = i.updateChildren(d, dir, dirEntry.Id())

	return nil
}

func (i *memIndex) updateChildren(d *memIndexData, dir *dirInfo, parentId string) error {
	for _, fi := range dir.children {
		if fi.IsDir() {

			err := i.updateDir(d, filepath.Join(dir.path, fi.Name()), parentId)
			if err != nil {
				return err
			}

		} else {

			fileEntry, err := i.entryFromInfo(fi, filepath.Join(dir.path, fi.Name()), parentId)
			if err != nil {
				return err
			}

			d.add(fileEntry.(*memEntry))
		}
	}
	return nil
}

func (i *memIndex) entryFromInfo(fi os.FileInfo, path string, parentId string) (Entry, error) {
	rp, err := filepath.Rel(i.root, path)
	if err != nil {
		return nil, fmt.Errorf("Could not determine relative path of %v in %v", path, i)
	}
	e := newMemEntry(i, rp, fi, parentId)
	return e, nil
}

func (i *memIndex) update() {
	log.Infof("Starting index scan for %v", i)
	d := newMemIndexData()
	dir, err := readDirInfo(i.root)
	if err != nil {
		log.Errorf("Error during index scan for %v: %v", i, err)
		return
	}
	err = i.updateChildren(d, dir, "")
	if err != nil {
		log.Errorf("Error during index scan for %v: %v", i, err)
		return
	}
	i.data = d
	log.Infof("Finished index scan for %v. Found %v entries", i, len(i.data.entries))
}

func (i *memIndex) String() string {
	return fmt.Sprintf("filepath.MemIndex(%v)", i.root)
}

/* Index Data */

type memIndexData struct {
	entries  map[string]Entry
	children map[string][]Entry
}

func newMemIndexData() *memIndexData {
	return &memIndexData{make(map[string]Entry), make(map[string][]Entry)}
}

func (d *memIndexData) add(e *memEntry) {
	log.Infof("Adding index entry %v", e.Path())
	d.entries[e.Id()] = e
	d.children[e.ParentId()] = append(d.children[e.ParentId()], e)
}
