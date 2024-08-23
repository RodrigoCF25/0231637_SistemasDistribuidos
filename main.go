package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

const (
	port      = ":8080"
	writePath = "/write"
	readPath  = "/read"
)

type Log struct {
	mutex   sync.RWMutex
	records []Record
}

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"-"`
}

// Constructor for Log
func GetNewLog() *Log {
	return &Log{
		mutex:   sync.RWMutex{},
		records: make([]Record, 0, 10),
	}
}

// Method to append a record to the log
func (l *Log) Append(record Record) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	record.Offset = uint64(len(l.records))
	fmt.Println("Record Appended: ", record)
	l.records = append(l.records, record)
}

// Method to read a record from the log
func (l *Log) Read(offset uint64) (Record, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	if offset >= uint64(len(l.records)) {
		return Record{}, fmt.Errorf("offset out of range")
	}

	return l.records[offset], nil
}

// Stringer interface, basically the __repr__ of python
func (l *Log) String() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return fmt.Sprintf("Log: %v", l.records)
}

func (r Record) String() string {
	return fmt.Sprintf("{value: %s, offset: %d}", string(r.Value), r.Offset)
}

// Function to write to the log when a POST request is made to /write
func WriteToLog(l *Log, w http.ResponseWriter, r *http.Request) error {

	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed")
	}

	var record Record

	err := json.NewDecoder(r.Body).Decode(&record)

	// Check if the body is empty or if record is not empty
	if err != nil || len(record.Value) == 0 {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return fmt.Errorf("error reading body")
	}

	l.Append(record)

	return nil
}

// Function to read from the log when a GET request is made to /read
func ReadFromLog(l *Log, w http.ResponseWriter, r *http.Request) error {

	// Check if the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed")
	}

	var Offset struct {
		Offset uint64 `json:"offset"`
	}

	// Check if the offset is invalid
	if err := json.NewDecoder(r.Body).Decode(&Offset); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return fmt.Errorf("invalid offset")
	}

	// Read the record from the log
	record, err := l.Read(Offset.Offset)

	// Check if the record is not found
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, "Record not found", http.StatusNotFound)
		return fmt.Errorf("record not found")
	}

	// Encode the record to JSON and write it to the response
	err = json.NewEncoder(w).Encode(record)

	// Check if the encoding failed
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Error encoding record", http.StatusInternalServerError)
		return fmt.Errorf("error encoding record")
	}

	fmt.Println(record)

	return nil

}

func main() {

	myLog := GetNewLog()

	http.HandleFunc(writePath, func(w http.ResponseWriter, r *http.Request) {
		WriteToLog(myLog, w, r)
	})

	http.HandleFunc(readPath, func(w http.ResponseWriter, r *http.Request) {
		ReadFromLog(myLog, w, r)
	})

	log.Fatal(http.ListenAndServe(port, nil))

}
