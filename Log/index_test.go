package log

import (
	"os"
	"testing"
)

var (
	config = NewConfig(1024, 1024, 0)
)

func TestIndex(t *testing.T) {

	tempFile, err := os.CreateTemp("", "index_test")

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	index, err := NewIndex(tempFile, *config)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//Read empty index

	_, _, err = index.Read(-1)

	if err == nil {
		t.Fatalf("err: %v", err)
	}

	tables := []struct {
		off uint32
		pos uint64
	}{
		{0, 0},
		{1, 20},
	}

	for _, table := range tables {
		index.Write(table.off, table.pos)

		off, pos, err := index.Read(int64(table.off))

		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if off != table.off {
			t.Fatalf("expected: %v, got: %v", table.off, off)
		}

		if pos != table.pos {
			t.Fatalf("expected: %v, got: %v", table.pos, pos)
		}
	}

	//Read with negative index like in python's lists
	off, pos, _ := index.Read(-1)

	if off != tables[len(tables)-1].off && pos != tables[len(tables)-1].pos {
		t.Fatalf("expected: %v, got: %v", tables[len(tables)-1], off)
	}

	//Read with and offset(index) greater than the last index

	_, _, err = index.Read(int64(len(tables)))

	if err == nil {
		t.Fatalf("err: %v", err)
	}

	err = index.Close()

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//Trying to close the index again, should return nil
	_ = index.Close()

	defer os.Remove(tempFile.Name())

}
