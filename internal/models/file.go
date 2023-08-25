package models

import (
	"bytes"
	"encoding/gob"
)

// File структра для хранения файлов.
// Name - имя файла.
// Data - содержимое файла.
type File struct {
	Name []byte
	Data []byte
}

func NewFile(name, data []byte) *File {
	return &File{name, data}
}

func (f *File) Decode(d []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(d))
	err := dec.Decode(f)
	return err
}

func (f *File) Type() MType {
	return FILE
}
