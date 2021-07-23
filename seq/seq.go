package seq

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
)

func CreateRandomSeq(seqType string, seqSize int) (randSeq string) {
	/* creates random nucletide sequence for testing purposes

	input arguments:
	- seqType (sequence type: DNA or RNA)
	- seqSize (size of the sequence in nucleotide bases)

	output:
	- randSeq (a seqType sequence of size seqSize)
	*/

	// initialize bases for DNA or RNA sequence generation
	var bases = []string{"A", "C", "G"}
	switch strings.ToUpper(seqType) {
	case "DNA":
		bases = append(bases, "T")
	case "RNA":
		bases = append(bases, "U")
	default:
		fmt.Printf("Nucleotide sequence of type %v not supported", seqType)
	}

	// add seqSize nucleotide bases to sequence
	for i := 0; i < seqSize; i++ {
		randSeq += bases[rand.Intn(len(bases))]
	}

	return randSeq
}

func CreateRandomLib(libResult, seqType string, libSize, seqSize int) {
	/* creates random sequence library for testing purposes

	input arguments:
	- libResult (name of fasta file to store record ID and sequences)
	- seqType (sequence type: DNA or RNA)
	- libSize (# of sequences to generate for library)
	- seqSize (size of sequence in nucleotide bases)

	output:
	- data written to file specified by libResult
	*/

	fmt.Printf("Building %v library w/ %v %v-base sequences", seqType, libSize, seqSize)
	results, resultOpenErr := os.OpenFile(libResult, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if resultOpenErr != nil {
		log.Fatal(resultOpenErr)
	}

	for i := 0; i < libSize; i++ {
		recordName := fmt.Sprintf(">test-%v-sequence %v\n", seqType, i)
		recordSeq := CreateRandomSeq(seqType, seqSize)
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

func validateFASTA(fileFormat string) bool {
	/* validate FASTA file format

	input argument:
	- fileFormat (file extension as a string)

	output:
	boolean value (true if the extension is valid, false otherwise)
	*/

	var validFormats = []string{".fasta", ".fsa", ".fastq"}
	for _, format := range validFormats {
		if format == fileFormat {
			return true
		}
	}
	return false
}

func BuildMultiFASTA(fsaResult, dataRepo string) {
	/* builds a multi record fasta file from valid files in directory

	input arguments:
	- dataRepo (directory containing multiple FASTA files)
	* All files not in fasta format (i.e. .fasta, .fsa, .fastq) will be ignored
	- fsaResult (name of file to store record ID and sequences)

	output:
	- data writen to file specified by fsaResult
	*/

	fmt.Println("Building mutli record FASTA file...")
	files, err := os.ReadDir(dataRepo)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if validateFASTA((path.Ext(file.Name()))) {
			fsaData, fsaReadErr := os.ReadFile(dataRepo + file.Name())

			if fsaReadErr != nil {
				log.Fatal(fsaReadErr)
			}

			results, resultOpenErr := os.OpenFile(fsaResult, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if resultOpenErr != nil {
				log.Fatal(resultOpenErr)
			}

			if _, writerErr := results.Write([]byte(fsaData)); writerErr != nil {
				results.Close() // ignore error; Write error takes precedence
				log.Fatal(writerErr)
			}

			if closerErr := results.Close(); closerErr != nil {
				log.Fatal(closerErr)
			}
		} else {
			fmt.Printf("Warning: non FASTA formats not supported. %v ignored\n", file.Name())
		}
	}
}

func Compress(sequence string) {
	var basesToBytes = map[string]byte{"A": 00, "C": 01, "T": 11, "U": 11, "G": 10}
	fmt.Println(basesToBytes)
}
