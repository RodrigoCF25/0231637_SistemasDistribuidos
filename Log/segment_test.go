package log

import (
	"os"
	"testing"

	api "github.com/RodrigoCF25/0231637_SistemasDistribuidos/api/v1"
)

func TestSegment(t *testing.T) {

	write := []byte("Stolas")
	recordToWrite := api.Record{
		Value: write,
	}

	config := NewConfig(uint64(8*(len(write)+lenWidth)), 24, 16)

	//fmt.Println(recordToWrite.Value)
	//fmt.Println(string(recordToWrite.Value))
	tempDir, err := os.MkdirTemp("", "SegmentTestDir")

	defer os.RemoveAll(tempDir)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	segment, err := NewSegment(tempDir, 16, *config)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	defer segment.Close()

	//Write
	_, err = segment.Append(&recordToWrite)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//Read
	recordRead, err := segment.Read(16)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if string(recordRead.Value) != string(recordToWrite.Value) {
		t.Fatalf("expected: %v, got: %v", string(recordToWrite.Value), string(recordRead.Value))
	}

	//Close segment
	err = segment.Close()

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//Open segment so we can test the nextOffset actualization

	segment, err = NewSegment(tempDir, 16, *config)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//Write
	_, err = segment.Append(&recordToWrite)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//Try to write more than the segment can handle, so Maxed error should be returned

	_, err = segment.Append(&recordToWrite)

	if err == nil {
		t.Fatalf("err: %v", err)
	}

	//Close segment
	err = segment.Close()

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//Close segment again, should return nil

	err = segment.Close()

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//Remove segment
	err = segment.Remove()

	if err != nil {
		t.Fatalf("err: %v", err)
	}

}
