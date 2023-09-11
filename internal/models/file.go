package models

import (
	"fmt"
	"os"
)

// File структра для хранения файлов.
type File struct {
	Name string
	Data []byte
}

// NewFile returns File.
func NewFile(name string, data []byte) *File {
	return &File{name, data}
}

/*
func (f *File) Decode(d []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(d))
	err := dec.Decode(f)
	return err
}

func (f *File) Type() MType {
	return FILE
}
*/
// String prints file information.
func (f *File) String() string {
	return fmt.Sprintf("%s : %d", f.Name, len(f.Data))
}

// Show shows file information.
func (f *File) Show() string {
	return f.String()
}

// IsReadyForStorage check File for ready to store.
func (f *File) IsReadyForStorage() (bool, string) {
	const (
		fileNameIsEmpty = "file name is empty"
		fileIsNotExist  = "file is not exist"
	)
	if f.Name == "" {
		return false, fileNameIsEmpty
	}
	data, err := os.ReadFile(f.Name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Sprintf("file %s %s", f.Name, fileIsNotExist)
		}
		//TODO
		return false, err.Error()
	}
	f.Data = data
	return true, ""
}
