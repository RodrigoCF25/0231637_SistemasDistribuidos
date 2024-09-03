package main

import (
	"fmt"

	Log "github.com/RodrigoCF25/0231637_SistemasDistribuidos/Log"

	api "github.com/RodrigoCF25/0231637_SistemasDistribuidos/api/v1"
)

func main() {

	config := Log.NewConfig(1024, 1024, 16)

	fmt.Println(config)

	segment, err := Log.NewSegment("Archivos", 16, *config)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer segment.Close()

	fmt.Println(segment)

	record := api.Record{
		Value: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor. Cras elementum ultrices diam. Maecenas ligula massa, varius a, semper congue, euismod non, mi."),
	}

	off, err := segment.Append(&record)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(off)

	record2, err := segment.Read(uint32(off))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(record2)

	/*
		record = api.Record{
			Value: []byte("A"),
		}

		off, err = segment.Append(&record)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(off)

		record2, err = segment.Read(uint32(off))

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(record2)

		fmt.Println(segment.Read(10))
	*/

}
