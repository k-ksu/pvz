package json

import (
	"errors"
	"os"
)

func (s *Storage) createFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func (s *Storage) checkFileExistence(filename string) error {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		if errCreateFile := s.createFile(filename); errCreateFile != nil {
			return errCreateFile
		}
		err := os.WriteFile(filename, []byte("{}"), 0666)
		if err != nil {
			return err
		}
	}

	return nil
}
