package models

import (
	"bytes"
	"encoding/gob"
)

// File структра для хранения файлов.
// Name - имя файла.
// Data - содержимое файла.
// IsCompressed - признак сжатия данных.
type File struct {
	Name         []byte
	Data         []byte
	IsCompressed bool
}

func NewFile(name, data []byte, isCompressed bool) *File {
	return &File{name, data, isCompressed}
}

func (f *File) Decode(d []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(d))
	err := dec.Decode(f)
	return err
}

func (f *File) Type() MType {
	return FILE
}
