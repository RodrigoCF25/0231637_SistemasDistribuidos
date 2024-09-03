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
	mutex    sync.RWMutex
	buffer   *bufio.Writer
	size     uint64
	isClosed bool
}

func NewStore(f *os.File) (*Store, error) {

	fileInfo, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not get file info: %w", err)
	}

	return &Store{
		mutex:    sync.RWMutex{},
		File:     f,
		size:     uint64(fileInfo.Size()),
		buffer:   bufio.NewWriter(f),
		isClosed: false,
	}, nil
}

func (s *Store) Append(data []byte) (bytesWritten uint64, positionWhereWritten uint64, err error) {
	//So what it is stored is the offset of the data and the data itself
	s.mutex.Lock()
	defer s.mutex.Unlock()

	fmt.Println()
	//Get the initial position where the data will be stored
	positionWhereWritten = s.size

	//Write the length of the data
	err = binary.Write(s.buffer, enc, uint64(len(data)))

	if err != nil {
		err = fmt.Errorf("could not write data length: %w", err)
		return 0, 0, err
	}

	//Write the data
	n, err := s.buffer.Write(data)

	if err != nil {
		err = fmt.Errorf("could not write data: %w", err)
		return 0, 0, err
	}

	bytesWritten = lenWidth + uint64(n)
	s.size += bytesWritten

	//Flush the buffer

	err = s.buffer.Flush()

	if err != nil {
		err = fmt.Errorf("could not flush buffer: %w", err)
		return 0, 0, err
	}

	return bytesWritten, positionWhereWritten, nil

}

func (s *Store) Read(absolute_position uint64) ([]byte, error) {

	data := make([]byte, lenWidth)

	//Read the length of the data
	n, err := s.ReadAt(data, int64(absolute_position))

	if n == 0 {
		return nil, fmt.Errorf("could not read from file: %w", err)
	}

	//Get the length of the data
	size := enc.Uint64(data)

	data = make([]byte, size)

	//Read the data

	n, err = s.ReadAt(data, int64(absolute_position+lenWidth))

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

	return s.File.ReadAt(p, off)
}

func (s *Store) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isClosed {
		return nil
	}

	err := s.File.Close()

	if err != nil {
		err = fmt.Errorf("could not close file: %w", err)
		return err
	}

	s.isClosed = true

	return nil
}
