package log

import (
	"os"
	"testing"

	api "github.com/RodrigoCF25/0231637_SistemasDistribuidos/api/v1"
)

var (
	write = []byte("Stolas Goetia")
	width = lenWidth + uint64(len(write))
)

//Append, Read and Close

func TestAppendReadClose(t *testing.T) {

	tempFile, err := os.CreateTemp("", "store_test")

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	store, err := NewStore(tempFile)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	tables := []struct {
		record *api.Record
	}{
		{&api.Record{Value: write}},
		{&api.Record{Value: write}},
	}

	expectedPosition := 0
	for _, table := range tables {

		_, positionWhereWritten, err := store.Append(table.record.Value)

		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if positionWhereWritten != uint64(expectedPosition) {
			t.Fatalf("expected: %v, got: %v", expectedPosition, positionWhereWritten)
		}

		record, err := store.Read(uint64(expectedPosition))

		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if string(record) != string(table.record.Value) {
			t.Fatalf("expected: %v, got: %v", string(table.record.Value), string(record))
		}

		expectedPosition += int(width)

	}

	err = store.Close()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//Trying to close the store again
	err = store.Close()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	defer os.Remove(tempFile.Name())
}
