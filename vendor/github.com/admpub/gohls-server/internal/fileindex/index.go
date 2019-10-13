package fileindex

type Entry interface {
	Id() string
	Name() string
	IsDir() bool
	Path() string
	ParentId() string
}

type Index interface {
	Id() string
	Root() string
	Get(id string) (Entry, error)
	WaitForReady() error
	List(parent string) ([]Entry, error)
}
