package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/krasin/ezhik"
)

var n = flag.Int("n", 128, "Number of source blocks in the original message")

var seedRe = regexp.MustCompile(`\.([0-9]+)\.ezhik`)

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

func GetSeed(filename string) (val int64) {
	res := seedRe.FindStringSubmatch(filename)
	if len(res) < 2 {
		log.Fatalf("Unable to parse filename: %s. Filenames should have the form <anything>.<numerical seed>.ezhik", filename)
	}
	str := res[1]
	var err error
	if val, err = strconv.ParseInt(str, 10, 64); err != nil {
		log.Fatalf("Unable to parse filename: %s. Extracted seed string: %s, strconv error: %v", filename, str, err)
	}
	return
}

/*func TestLinearSystem() {
	ls := &ezhik.LinearSystem{n: 3}
	ls.Add(&bitSet{a: []uint64{5}, n: 6}, []byte{1})
	ls.Add(&bitSet{a: []uint64{7}, n: 6}, []byte{0})
	ls.Add(&bitSet{a: []uint64{1}, n: 6}, []byte{0})
	fmt.Fprintf(os.Stderr, "TesLinearSystem: ls.Determined: %v\n", ls.Determined())
	ls.Backtrack()
	data := ls.Solve()
	fmt.Fprintf(os.Stderr, "data: %v\n", data)
	log.Fatalf("Test Stub")
}*/

func checkBlocks(n int, filenames []string, blocks [][]byte) (blockLen int) {
	//	if len(blocks) < n {
	//		log.Fatalf("Too few blocks (%d). Want at least %d, but it's better to have a few more", len(blocks), n)
	//	}
	blockLen = len(blocks[0])
	for i, block := range blocks {
		if len(block) != blockLen {
			log.Fatalf("Blocks have different length. %d = len(block[\"%s\"]) != len(blocks[\"%s\"]) = %d", len(blocks[0]), filenames[0], filenames[i], len(block))
		}
	}
	return
}

func GetSeeds(filenames []string) []int64 {
	res := make([]int64, len(filenames))
	for i, filename := range filenames {
		res[i] = GetSeed(filename)
	}
	return res
}

func main() {
	flag.Parse()

	//	TestLinearSystem()

	filenames := os.Args[1:]
	blocks := loadBlocks(filenames)
	/* blockLen := */ checkBlocks(*n, filenames, blocks)
	seeds := GetSeeds(filenames)
	//	fmt.Printf("Block len: %d\n", blockLen)
	data, err := ezhik.Decode(*n, seeds, blocks)
	if err != nil {
		log.Fatalf("Decode: %v", err)
	}
	os.Stdout.Write(data)
}
