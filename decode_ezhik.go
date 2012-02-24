package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

var n = flag.Int("n", 128, "Number of source blocks in the original message")

func loadBlocks(filenames []string) (blocks [][]byte) {
	for _, filename := range filenames {
		var data []byte
		var err error
		if data, err = ioutil.ReadFile(filename); err != nil {
			log.Fatalf("Could not read input file %s: %v", filename, err)
		}
		blocks = append(blocks, data)
	}
	return
}

func checkBlocks(n int, filenames []string, blocks [][]byte) (blockLen int) {
	if len(blocks) < n {
		log.Fatalf("Too few blocks (%d). Want at least %d, but it's better to have a few more", n, len(blocks))
	}
	blockLen = len(blocks[0])
	for i, block := range blocks {
		if len(block) != blockLen {
			log.Fatalf("Blocks have different length. %d = len(block[\"%s\"]) != len(blocks[\"%s\"]) = %d", len(blocks[0]), filenames[0], filenames[i], len(block))
		}
	}
	return
}

type BitSet []uint64

func NewBitSet(n int) (res BitSet) {
	return make([]uint64, (n+63)/64)
}

func (bs BitSet) Has(i int) bool {
	return (bs[i>>6] >> uint(i&0x3F)) != 0
}

func (bs BitSet) Set(i int, val bool) {
	if val {
		bs[i>>6] |= 1 << uint(i&0x3F)
	} else {
		bs[i>>6] &= ^uint64(1 << uint(i&0x3F))
	}
}

func GetMask(n int, seed int64) (res BitSet) {
	res = NewBitSet(n)
	r := rand.New(rand.NewSource(seed))
	for i := 0; i < n; i++ {
		res.Set(i, r.Intn(2) == 0)
	}
	return
}

func main() {
	flag.Parse()
	filenames := os.Args[1:]
	blocks := loadBlocks(filenames)
	blockLen := checkBlocks(*n, filenames, blocks)
	fmt.Printf("Block len: %d\n", blockLen)
}
