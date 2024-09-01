package log

import (
	"fmt"
	"os"
	"path"

	api "github.com/RodrigoCF25/0231637_SistemasDistribuidos/api/v1"
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

// Append
func (s *segment) Append(record *api.Record) (off uint64, err error) {
	off = s.nextOffset

	if s.IsMaxed() {
		err = fmt.Errorf("segment's store/index is maxed")
		return 0, err
	}

	var pos uint64
	if _, pos, err = s.store.Append(record.Value); err != nil {
		return 0, err
	}

	if err = s.index.Write(uint32(s.nextOffset-uint64(s.baseOffset)), pos); err != nil {
		return 0, err
	}

	s.nextOffset++
	return off, nil
}

//Read

func (s *segment) Read(off uint32) (*api.Record, error) {
	_, pos, err := s.index.Read(int64(off - uint32(s.baseOffset)))

	if err != nil {
		return nil, err
	}

	data, err := s.store.Read(pos)

	if err != nil {
		return nil, err
	}

	return &api.Record{
		Value:  data,
		Offset: uint64(off),
	}, nil

}

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

func (s *segment) Close() error {

	if err := s.store.Close(); err != nil {
		err = fmt.Errorf("could not close store: %w", err)
		return err
	}
	if err := s.index.Close(); err != nil {
		err = fmt.Errorf("could not close index: %w", err)
		return err
	}
	return nil
}
