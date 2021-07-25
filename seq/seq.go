// TODO:
// [x] refactor comments above functions
// [ ] write detailed comments for each function

package seq

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
)

// CreateRandomSeq returns a random DNA or RNA sequence
func CreateRandomSeq(seqType string, seqLength int) string {
	var bases = []string{"A", "C", "G"}
	switch strings.ToUpper(seqType) {
	case "DNA":
		bases = append(bases, "T")
	case "RNA":
		bases = append(bases, "U")
	default:
		fmt.Printf("Nucleotide sequence of type %v not supported", seqType)
	}

	var randSeq string
	for i := 0; i < seqLength; i++ {
		randSeq += bases[rand.Intn(len(bases))]
	}

	return randSeq
}

// CreateRandomLib creates a DNA or RNA library and writes sequences to a fasta file
func CreateRandomLib(libResult, seqType string, libSize, seqLength int) {
	fmt.Printf("Building %v library w/ %v %v-base sequences", seqType, libSize, seqLength)
	results, resultOpenErr := os.OpenFile(libResult, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if resultOpenErr != nil {
		log.Fatal(resultOpenErr)
	}

	for i := 0; i < libSize; i++ {
		recordName := fmt.Sprintf(">test-%v-sequence %v\n", seqType, i)
		recordSeq := CreateRandomSeq(seqType, seqLength)
		fullRecord := recordName + recordSeq

		if _, writerErr := results.Write([]byte(fullRecord)); writerErr != nil {
			results.Close() // ignore error; Write error takes precedence
			log.Fatal(writerErr)
		}

		if closerErr := results.Close(); closerErr != nil {
			log.Fatal(closerErr)
		}
	}
	fmt.Println("Job Complete")
}

// ValidateFASTA returns true if passed extension is a fasta-like extenstion, returns false otherwise.
func ValidateFASTA(fileFormat string) bool {
	var validFormats = []string{".fasta", ".fsa", ".fastq"}
	for _, format := range validFormats {
		if format == fileFormat {
			return true
		}
	}
	return false
}

// BuildMultiFASTA writes a fasta file containing all records found in a given directory
func BuildMultiFASTA(fsaResult, dataRepo string) {
	fmt.Println("Building mutli record FASTA file...")
	files, err := os.ReadDir(dataRepo)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if ValidateFASTA((path.Ext(file.Name()))) {
			fsaData, err := os.ReadFile(dataRepo + file.Name())

			if err != nil {
				log.Fatal(err)
			}

			results, err := os.OpenFile(fsaResult, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				log.Fatal(err)
			}

			if _, err := results.Write([]byte(fsaData)); err != nil {
				results.Close() // ignore error; Write error takes precedence
				log.Fatal(err)
			}

			if err := results.Close(); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("Warning: non FASTA formats not supported. %v ignored\n", file.Name())
		}
	}
}

// CompressSeq performs byte-packed compression on a DNA or RNA sequence and returns its byte representation.
func CompressSeq(sequence string) *bytes.Buffer {
	var basesToBytes = map[string]byte{"A": 00, "C": 01, "T": 11, "U": 11, "G": 10}
	compressedSeq := bytes.NewBuffer(make([]byte, 0, len(sequence)))

	for _, base := range sequence {
		compressedSeq.WriteByte(basesToBytes[string(base)])
	}

	return compressedSeq
}
