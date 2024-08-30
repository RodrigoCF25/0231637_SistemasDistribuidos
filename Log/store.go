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

func NewStore(f *os.File) (*Store, error) {

	return &Store{
		mutex:  sync.RWMutex{},
		File:   f,
		size:   0,
		buffer: bufio.NewWriter(f),
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

func (s *Store) Read(absolute_position uint64) ([]byte, error) {

	data := make([]byte, lenWidth)

	n, err := s.ReadAt(data, int64(absolute_position))

	if n == 0 {
		return nil, fmt.Errorf("could not read from file: %w", err)
	}

	return data, nil

}

func (s *Store) ReadAt(p []byte, off int64) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.buffer.Flush(); err != nil {
		return 0, err
	}

	fmt.Println(("Reading at position: "), off)

	return s.File.ReadAt(p, off)
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
