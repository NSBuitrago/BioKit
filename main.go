package main

import "github.com/nsbuitrago/biokit/seq"

func main() {
	//seq.BuildMultiFASTA("./testing_data/query/", "test_build.fsa")
	seq.CreateRandomLib("test_lib.fsa", "DNA", 1, 10)
}
