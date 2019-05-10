package ini

import (
	"bytes"
	"io/ioutil"
)

//==========================
//added by admpub
//==========================

// Clone copy
func (s *Section) Clone(newName string) *Section {
	clone, _ := s.f.NewSection(newName)
	for _, key := range s.KeyStrings() {
		clone.Key(key).SetValue(s.Key(key).String())
	}
	return clone
}

// LoadContent load content
func LoadContent(content string, others ...interface{}) (*File, error) {
	return Load(ioutil.NopCloser(bytes.NewBufferString(content)), others...)
}

// SaveToIndent writes content to file system with given value indention.
func (f *File) String() string {
	// Note: Because we are truncating with os.Create,
	// 	so it's safer to save to a temporary file location and rename afte done.
	buf, err := f.writeToBuffer("")
	if err != nil {
		return err.Error()
	}
	return buf.String()
}
