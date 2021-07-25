package main

import (
	"fmt"

	"github.com/nsbuitrago/biokit/seq"
)

func main() {
	sequence := "ACTG"
	compressedSeq := seq.CompressSeq(sequence)
	fmt.Println(compressedSeq)
}
