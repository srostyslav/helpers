package helpers

import (
	"io/ioutil"
	"os"
)

type File struct {
	Name string
}

func (f *File) Write(data []byte) error {
	f.Create()

	if file, err := os.OpenFile(f.Name, os.O_RDWR, 0644); err != nil {
		return err
	} else {
		defer file.Close()

		if _, err = file.WriteString(string(data)); err != nil {
			return err
		}
		return file.Sync()
	}
}

func (f *File) Remove() {
	if _, err := os.Stat(f.Name); err == nil {
		os.Remove(f.Name)
	}
}

func (f *File) Create() {
	f.Remove()
	file, _ := os.Create(f.Name)
	file.Close()
}

func (f *File) CreateAndGet() (*os.File, error) {
	f.Remove()
	return os.Create(f.Name)
}

func (f *File) Content() (string, error) {

	if file, err := os.Open(f.Name); err != nil {

		ErrorLogger.Println(err)
		return "", err

	} else {
		defer file.Close()

		if b, err := ioutil.ReadAll(file); err != nil {
			ErrorLogger.Println(err)
			return "", err
		} else {
			return string(b), nil
		}
	}
}
