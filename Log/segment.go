package log

import (
	"fmt"
	"os"
	"path"
)

type segment struct {
	store                  *Store
	index                  *Index
	baseOffset, nextOffset uint64
	config                 Config
}

func NewSegment(dir string, baseOffset uint64, c Config) (*segment, error) {
	s := &segment{
		baseOffset: baseOffset,
		config:     c,
	}
	var err error
	storeFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".store")),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}

	if s.store, err = NewStore(storeFile); err != nil {
		return nil, err
	}
	indexFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".index")),
		os.O_RDWR|os.O_CREATE,
		0644,
	)
	if err != nil {
		return nil, err
	}
	if s.index, err = NewIndex(indexFile, c); err != nil {
		return nil, err
	}
	if off, _, err := s.index.Read(-1); err != nil {
		s.nextOffset = baseOffset
	} else {
		s.nextOffset = baseOffset + uint64(off) + 1
	}

	return s, nil
}

//Append

//Read

// IsMaxed
func (s *segment) IsMaxed() bool {
	return (s.store.size >= s.config.Segment.MaxStoreBytes) || (s.index.size >= s.config.Segment.MaxIndexBytes)
}

//Remove

func (s *segment) Remove() error {
	if err := s.store.Close(); err != nil {
		err = fmt.Errorf("could not close store: %w", err)
		return err
	}
	if err := s.index.Close(); err != nil {
		err = fmt.Errorf("could not close index: %w", err)
		return err
	}
	if err := os.Remove(s.index.file.Name()); err != nil {
		err = fmt.Errorf("could not remove index file: %w", err)
		return err
	}
	if err := os.Remove(s.store.Name()); err != nil {
		err = fmt.Errorf("could not remove store file: %w", err)
		return err
	}
	return nil
}

//Close

func (s *segment) Close() {
	s.store.Close()
	s.index.Close()
}
