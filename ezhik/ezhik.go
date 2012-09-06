package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/krasin/ezhik"
)

var (
	seed      = flag.Int64("seed", 0, "Seed for encoded block")
	seedCount = flag.Int("seedCount", 1, "Number of blocks with consequentive seeds to generate")
	n         = flag.Int("n", 128, "Number of source blocks")
	output    = flag.String("output", "output", "Prefix for output files. The actual filenames will be like output.seed.eshik")
)

func encodeAndWriteToFile(output string, data []byte, n int, seed int64) {
	block := ezhik.Encode(data, n, seed)
	filename := fmt.Sprintf("%s.%d.ezhik", output, seed)
	if err := ioutil.WriteFile(filename, block, 0644); err != nil {
		log.Fatalf("iotil.WriteFile: %v", err)
	}
}

// This is the comment
func main() {
	flag.Parse()
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Could not read from stdin: %v", err)
	}
	for i := 0; i < *seedCount; i++ {
		encodeAndWriteToFile(*output, data, *n, *seed+int64(i))
	}
}
