package seq

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
)

// CreateRandomSeq returns a random uncompressed DNA or RNA sequence
func CreateRandomSeq(seqType string, seqLength int) string {
	var bases = []string{"A", "C", "G"}
	switch strings.ToUpper(seqType) {
	case "DNA":
		bases = append(bases, "T")
	case "RNA":
		bases = append(bases, "U")
	default:
		log.Fatalf("Nucleotide sequence of type %v not supported", seqType)
	}

	var randSeq string
	for i := 0; i < seqLength; i++ {
		randSeq += bases[rand.Intn(len(bases))]
	}

	return randSeq
}

// CreateRandomLib creates a DNA or RNA library and writes sequences to a fasta file
func CreateRandomLib(libResult, seqType string, libSize, seqLength int) {
	fmt.Printf("Creating %v library w/ %v %v-base sequences", seqType, libSize, seqLength)
	out, err := os.OpenFile(libResult, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < libSize; i++ {
		recordName := fmt.Sprintf(">test-%v-sequence %v\n", seqType, i)
		recordSeq := CreateRandomSeq(seqType, seqLength)
		fullRecord := recordName + recordSeq

		if _, err := out.Write([]byte(fullRecord)); err != nil {
			out.Close() // ignore error; Write error takes precedence
			log.Fatal(err)
		}
	}

	if err := out.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v library written in %v", seqType, libResult)
}

// CompressSeq performs byte-packed compression on a DNA or RNA sequence and returns its byte representation.
func CompressSeq(sequence string) *bytes.Buffer {
	var basesToBytes = map[string]byte{"A": 00, "C": 01, "T": 11, "U": 11, "G": 10, "N": 111}
	compressedSeq := bytes.NewBuffer(make([]byte, 0, len(sequence)))

	// for i in range(0, len(sequence), 4)
	// 4bases := sequence[i] && sequence[i+1] && sequence[i+2] && sequence[i+3]
	for _, base := range sequence {
		compressedSeq.WriteByte(basesToBytes[string(base)])
	}

	return compressedSeq
}

// merge multiple FASTA files
func MergeFASTA(fsaOut, fsaDir string) {
	files, err := os.ReadDir(fsaDir)

	if err != nil {
		log.Fatalf("could not read directory: %v", err)
	}

	out, err := os.OpenFile(fsaOut, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("could not open %v: %v", fsaOut, err)
	}

	for _, file := range files {
		if validFormat := ValidateFASTA(path.Ext(file.Name())); validFormat {
			fsaRecord, err := os.Open(fsaDir + file.Name())
			if err != nil {
				fmt.Printf("could not open %v for reading:, %v\n", file.Name(), err)
			}

			scanner := bufio.NewScanner(fsaRecord)
			scanner.Split(ScanFASTA)

			for scanner.Scan() {
				if _, err := out.Write([]byte(scanner.Text())); err != nil {
					out.Close()
					log.Fatal(err)
				}
			}

			if err := fsaRecord.Close(); err != nil {
				fmt.Printf("warning: error closing %v: %v", fsaRecord, err)
			}

		} else {
			fmt.Printf("Warning: non fasta files not supported. %v ignored\n", file.Name())
		}
	}

	if err := out.Close(); err != nil {
		log.Fatalf("could not close %v: %v", out, err)
	}
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

// ScanFASTA is a split function for a Scanner, same as ScanLines, but returns each line
// of text, NOT stripped of any trailing end-of-line marker.
func ScanFASTA(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0 : i+1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
