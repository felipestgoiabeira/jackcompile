package utils

import (
	"bytes"
	"log"
	"os"
)

func writeResultToFile(resultBuffer bytes.Buffer) {
	f, err := os.Create("../resources/tests/results/SquareResult.xml")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.Write(resultBuffer.Bytes())

	if err2 != nil {
		log.Fatal(err2)
	}
}
