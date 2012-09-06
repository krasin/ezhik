package ezhik

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
)

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
	return bs.a[i>>6]&(1<<uint(i&0x3F)) != 0
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
	res = NewBitSet(2 * n)
	r := rand.New(rand.NewSource(seed))
	for i := 0; i < n; i++ {
		res.Set(i, r.Intn(2) == 1)
	}
	return
}

type LinearSystem struct {
	n     int
	lines []BitSet
	y     [][]byte
	pos   int
	ready bool
}

func (ls *LinearSystem) Eliminate(dst, src int) {
	if !ls.lines[dst].Has(src) {
		return
	}
	ls.lines[dst].XorWith(ls.lines[src])
}

func (ls *LinearSystem) EliminateSrcRange(dst, src, count int) {
	for i := 0; i < count; i++ {
		ls.Eliminate(dst, src+i)
	}
}

func (ls *LinearSystem) EliminateDstRange(dst, count, src int) {
	for i := 0; i < count; i++ {
		ls.Eliminate(dst+i, src)
	}
}

func (ls *LinearSystem) FindOne(base, count, index int) int {
	for i := 0; i < count; i++ {
		if ls.lines[base+i].Has(index) {
			return base + i
		}
	}
	return -1
}

func (ls *LinearSystem) Promote(index int) {
	fmt.Fprintf(os.Stderr, "Promote(index=%d, pos=%d)\n", index, ls.pos)
	ls.lines[ls.pos], ls.lines[index] = ls.lines[index], ls.lines[ls.pos]
	ls.y[ls.pos], ls.y[index] = ls.y[index], ls.y[ls.pos]
	ls.lines[ls.pos].Set(ls.n+ls.pos, true)
	ls.pos++
	ls.EliminateDstRange(ls.pos, len(ls.lines)-ls.pos, ls.pos-1)
}

func FormatSlice(line BitSet, from, to int) string {
	buf := new(bytes.Buffer)
	for i := from; i < to; i++ {
		if i > 0 {
			fmt.Fprintf(buf, " ")
		}
		val := 0
		if line.Has(i) {
			val = 1
		}
		fmt.Fprintf(buf, "%d", val)
	}
	return buf.String()
}

func (ls *LinearSystem) Add(line BitSet, y []byte) bool {
	fmt.Fprintf(os.Stderr, "Add(line[0:10]: %s)\n", FormatSlice(line, 0, 10))
	if ls.pos >= ls.n {
		return true
	}
	ls.lines = append(ls.lines, line)
	ls.y = append(ls.y, y)
	index := len(ls.lines) - 1
	ls.EliminateSrcRange(index, 0, ls.pos)
	fmt.Fprintf(os.Stderr, "Partially eliminated line[0:10]: %s\n", FormatSlice(ls.lines[index], 0, 10))
	if !ls.lines[index].Has(ls.pos) {
		fmt.Fprintf(os.Stderr, "Add does not lead to Promote. index=%d, ls.pos=%d, line[0:10]: %s\n", index, ls.pos, FormatSlice(ls.lines[index], 0, 10))
		return false
	}
	ls.Promote(index)
	for ls.pos < ls.n {
		i := ls.FindOne(ls.pos, len(ls.lines)-ls.pos, ls.pos)
		if i == -1 {
			break
		}
		ls.Promote(i)
	}
	return ls.pos == ls.n
}

func (ls *LinearSystem) Determined() bool {
	return ls.pos == ls.n
}

func (ls *LinearSystem) PrintMatrix(title string) {
	fmt.Fprintf(os.Stderr, "PrintMatrix, %s\n", title)
	for i := 0; i < ls.n; i++ {
		fmt.Fprintf(os.Stderr, "%s\n", FormatSlice(ls.lines[i], 0, 2*ls.n))
	}
	fmt.Fprintf(os.Stderr, "\n")
}

func (ls *LinearSystem) Backtrack() {
	if !ls.Determined() {
		panic("Backtrack: linear system is not determined")
	}
	ls.PrintMatrix("Before backtrack")
	for i := ls.n - 1; i > 0; i-- {
		ls.EliminateDstRange(0, i, i)
	}
	ls.PrintMatrix("After backtrack")
	ls.ready = true
}

func (ls *LinearSystem) Pos() int {
	return ls.pos
}

func (ls *LinearSystem) Solve() (x [][]byte) {
	if !ls.ready {
		panic("Solve: !ls.ready")
	}
	blockLen := len(ls.y[0])
	x = make([][]byte, ls.n)
	for i := 0; i < ls.n; i++ {
		x[i] = make([]byte, blockLen)
		for j := 0; j < ls.n; j++ {
			if ls.lines[i].Has(ls.n + j) {
				XorBytes(x[i], ls.y[j])
			}
		}
	}
	return
}

func XorBytes(dst, src []byte) {
	for i, v := range src {
		dst[i] ^= v
	}
}

func Decode(n int, seeds []int64, blocks [][]byte) (data []byte, err error) {
	ls := &LinearSystem{n: n}
	for i, block := range blocks {
		seed := seeds[i]
		line := GetMask(n, seed)
		fmt.Fprintf(os.Stderr, "Decode, i = %d\n", i)
		if ls.Add(line, block) {
			// The linear system is determined now
			break
		}
	}
	if !ls.Determined() {
		return nil, fmt.Errorf("Decode failed: the linear system is not determined. May be more blocks are needed. Pos=%d", ls.Pos())
	}
	ls.Backtrack()
	xs := ls.Solve()
	buf := new(bytes.Buffer)
	for _, x := range xs {
		buf.Write([]byte(x))
	}
	return buf.Bytes(), nil
}
