package main

import (
	"fmt"
	"log"

	sf "github.com/kawabatas/go-id-generator"
)

func main() {
	id, err := sf.NewSnowflakeID(
		sf.WithRandomEnabled(),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("int64: %d\n bits: %064b\n", id, id)
}
