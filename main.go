package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"google.golang.org/protobuf/proto"

	api "github.com/RodrigoCF25/0231637_SistemasDistribuidos/api/v1"
)

func main() {

	/*
		myLog, err := log.NewLog("Archivos", *log.NewConfig(1024, 1023, 16))

		if err != nil {
			return
		}

		defer myLog.Close()

		record := api.Record{
			Value: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor. Cras elementum ultrices diam. Maecenas ligula massa, varius a, semper congue, euismod non, mi."),
		}

		off, err := myLog.Append(&record)

		if err != nil {
			return
		}

		record2, err := myLog.Read(off)

		if err != nil {
			return
		}

		fmt.Println(record2)

		record = api.Record{
			Value: []byte("StolasGoetia"),
		}

		off, err = myLog.Append(&record)

		if err != nil {

			fmt.Println(err)

			return
		}

		record2, err = myLog.Read(off)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(record2)

		record2, err = myLog.Read(100)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(record2)
	*/

	// Crea un nuevo mensaje Record
	record := &api.Record{
		Value: []byte("hello world"),
	}

	fmt.Println(record.Value)

	// Serializa el mensaje
	data, err := proto.Marshal(record)
	fmt.Println(data)
	if err != nil {
		log.Fatalf("Failed to marshal: %v", err)
	}

	// Simula la lectura desde un Reader
	lenWidth := 8
	simulatedReaderData := append([]byte{0, 0, 0, 0, 0, 0, 0, 0}, data...) // Prepend lenWidth zeros
	reader := bytes.NewReader(simulatedReaderData)

	// Lee los datos
	b, err := io.ReadAll(reader)
	fmt.Println(b)
	if err != nil {
		log.Fatalf("Failed to read: %v", err)
	}

	b = b[lenWidth:] // Skip the width
	fmt.Println(b)
	read := &api.Record{}

	err = proto.Unmarshal(b, read)
	if err != nil {
		log.Fatalf("Failed to unmarshal: %v", err)
	}

	fmt.Printf("Deserialized Record: %+v\n", read)
}
