package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

var (
	seed      = flag.Int64("seed", 0, "Seed for encoded block")
	seedCount = flag.Int("seedCount", 1, "Number of blocks with consequentive seeds to generate")
	n         = flag.Int("n", 128, "Number of source blocks")
	output    = flag.String("output", "output", "Prefix for output files. The actual filenames will be like output.seed.eshik")
)

func encodeAndWriteToFile(output string, data []byte, n int, seed int64) {
	blockLen := (len(data) + n - 1) / n

	block := make([]byte, blockLen)
	r := rand.New(rand.NewSource(seed))
	for i := 0; i < n; i++ {
		if r.Intn(2) == 0 {
			var cur []byte
			if i < n-1 {
				cur = data[i*blockLen : (i+1)*blockLen]
			} else {
				cur = make([]byte, blockLen)
				copy(cur, data[i*blockLen:])
			}
			for i, v := range cur {
				block[i] ^= v
			}
		}
	}

	var f *os.File
	var err error
	filename := fmt.Sprintf("%s.%d.ezhik", output, seed)
	if f, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err != nil {
		log.Fatalf("Could not create the output file %s: %v", filename, err)
	}
	defer f.Close()
	if _, err = f.Write(block); err != nil {
		log.Fatalf("Unable to write to file: %v", err)
	}
}

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
