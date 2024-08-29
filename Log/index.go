package log

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/edsrzf/mmap-go"
	gommap "github.com/edsrzf/mmap-go"
)

var (
	offWidth uint64 = 4
	posWidth uint64 = 8
	entWidth uint64 = offWidth + posWidth
)

type Index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

var MaxSize uint64 = 100

func NewIndex(filePath string) (*Index, error) {

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		err = fmt.Errorf("could not open file: %w", err)
		return nil, err
	}

	fileInfo, err := file.Stat()

	if err != nil {
		err = fmt.Errorf("could not get file info: %w", err)
		return nil, err
	}

	size := uint64(fileInfo.Size())

	file.Truncate(int64(MaxSize))

	mmap, err := gommap.Map(file, mmap.RDWR, 0)

	if err != nil {
		err = fmt.Errorf("could not mmap file: %w", err)
		if mmap != nil {
			mmap.Unmap()
		}
		return nil, err
	}

	return &Index{
		file: file,
		mmap: mmap,
		size: size,
	}, nil
}

func (i *Index) Close() error {

	err := i.mmap.Unmap()

	if err != nil {
		err = fmt.Errorf("could not unmap file: %w", err)
		return err
	}
	i.file.Truncate(int64(i.size))

	err = i.file.Close()
	if err != nil {
		err = fmt.Errorf("could not close file: %w", err)
		return err
	}
	return nil
}

func (i *Index) Write(offset uint32, position uint64) error {

	if i.size+entWidth > MaxSize {
		return fmt.Errorf("index size exceeded")
	}

	offBytes := make([]byte, offWidth)
	posBytes := make([]byte, posWidth)

	binary.BigEndian.PutUint32(offBytes, offset)
	binary.BigEndian.PutUint64(posBytes, position)

	copy(i.mmap[i.size:], offBytes)
	copy(i.mmap[i.size+offWidth:], posBytes)

	i.size += entWidth

	return nil
}

func (i *Index) Read(index int64) (offset uint32, position uint64, err error) {
	if i.size == 0 {
		return offset, position, fmt.Errorf("index is empty")
	}

	if index < 0 {
		index = int64(i.size/entWidth) + index
	}

	if index >= int64(i.size/entWidth) || index < 0 {
		return offset, position, fmt.Errorf("index out of bounds")
	}

	position = uint64(index) * entWidth

	offset = enc.Uint32(i.mmap[position : position+offWidth])
	position = enc.Uint64(i.mmap[position+offWidth : position+entWidth])

	return offset, position, nil

}
