package utils

import (
	"bytes"
	"log"
	"os"
)

func WriteResultToFile(resultBuffer bytes.Buffer, file string) {
	f, err := os.Create("../resources/tests/results/" + file)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.Write(resultBuffer.Bytes())

	if err2 != nil {
		log.Fatal(err2)
	}
}
