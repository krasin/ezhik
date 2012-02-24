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
		log.Fatalf("Too few blocks (%d). Want at least %d, but it's better to have a few more", len(blocks), n)
	}
	blockLen = len(blocks[0])
	for i, block := range blocks {
		if len(block) != blockLen {
			log.Fatalf("Blocks have different length. %d = len(block[\"%s\"]) != len(blocks[\"%s\"]) = %d", len(blocks[0]), filenames[0], filenames[i], len(block))
		}
	}
	return
}

type bitSet struct {
	a []uint64
	n int
}

type BitSet interface {
	Has(i int) bool
	Set(i int, val bool)
	XorWith(set BitSet)
	Len() int
}

func NewBitSet(n int) (res BitSet) {
	return &bitSet{a: make([]uint64, (n+63)/64), n: n}
}

func (bs *bitSet) Len() int {
	return bs.n
}

func (bs *bitSet) Has(i int) bool {
	return (bs.a[i>>6] >> uint(i&0x3F)) != 0
}

func (bs *bitSet) Set(i int, val bool) {
	if val {
		bs.a[i>>6] |= 1 << uint(i&0x3F)
	} else {
		bs.a[i>>6] &= ^uint64(1 << uint(i&0x3F))
	}
}

func (bs *bitSet) XorWith(set BitSet) {
	if bs.Len() != set.Len() {
		panic("XorWith: different lengths")
	}
	// TODO: speed up if BitSet is *bitSet
	for i := 0; i < bs.n; i++ {
		bs.Set(i, bs.Has(i) != set.Has(i))
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

type LinearSystem struct {
	Lines   []BitSet
	Pos     int
	MaxRank int
}

func (ls *LinearSystem) Eliminate(dst, src int) {
	if !ls.Lines[dst].Has(src) {
		return
	}
	ls.Lines[dst].XorWith(ls.Lines[src])
}

func (ls *LinearSystem) Add(line BitSet) {
	ls.Lines = append(ls.Lines, line)
	index := len(ls.Lines) - 1
	for i := 0; i < ls.Pos; i++ {
		ls.Eliminate(index, i)
	}

	if ls.Pos >= ls.MaxRank-1 || !ls.Lines[index].Has(ls.Pos) {
		return
	}
	ls.Lines[ls.Pos], ls.Lines[index] = ls.Lines[index], ls.Lines[ls.Pos]
	for i := ls.Pos + 1; i < len(ls.Lines); i++ {
		ls.Eliminate(i, ls.Pos)
	}
	ls.Pos++
}

func (ls *LinearSystem) BackPropagate() {
	if ls.Pos != ls.MaxRank {
		panic("BackPropagate: ls.Pos != ls.MaxRank")
	}
	for i := ls.Pos - 1; i > 0; i-- {
		for j := i - 1; j >= 0; j-- {
			ls.Eliminate(j, i)
		}
	}
}

func (ls *LinearSystem) Solve(y []BitSet) (x []BitSet) {
	if ls.Pos != ls.MaxRank {
		panic("BackPropagate: ls.Pos != ls.MaxRank")
	}
	// TODO: complete this method
}

func main() {
	flag.Parse()
	filenames := os.Args[1:]
	blocks := loadBlocks(filenames)
	blockLen := checkBlocks(*n, filenames, blocks)
	fmt.Printf("Block len: %d\n", blockLen)
}
