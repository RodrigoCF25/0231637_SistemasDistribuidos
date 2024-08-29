package log

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8
)

type Store struct {
	*os.File
	mutex  sync.RWMutex
	buffer *bufio.Writer
	size   uint64
}

func NewStore(filePath string) (*Store, error) {

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		err = fmt.Errorf("could not open file: %w", err)
		return nil, err
	}

	return &Store{
		mutex:  sync.RWMutex{},
		File:   file,
		size:   0,
		buffer: bufio.NewWriter(file),
	}, nil
}

func (s *Store) Append(data []byte) (bytesWritten uint64, actualPosition uint64, err error) {
	//So what it is stored is the offset of the data and the data itself
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//Get the initial position where the data will be stored
	actualPosition = s.size

	//Write in the buffer the data
	var n int
	n, err = s.buffer.Write(data)

	if err != nil {
		err = fmt.Errorf("could not write to buffer: %w", err)
		return 0, 0, err
	}

	bytesWritten = uint64(n)
	s.size += bytesWritten

	//Flush the buffer

	err = s.buffer.Flush()

	if err != nil {
		err = fmt.Errorf("could not flush buffer: %w", err)
		return 0, 0, err
	}

	actualPosition = s.size

	return bytesWritten, actualPosition, nil

}

func (s *Store) Read(element int64) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

}

func (s *Store) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	err := s.File.Close()

	if err != nil {
		err = fmt.Errorf("could not close file: %w", err)
		return err
	}

	return nil
}
