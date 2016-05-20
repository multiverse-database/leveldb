package bloom

import (
	"bytes"
	"io"

	"github.com/kezhuw/leveldb/internal/filter"
	"github.com/kezhuw/leveldb/internal/hash"
	"github.com/kezhuw/leveldb/internal/util"
)

const (
	filterName = "leveldb.BuiltinBloomFilter2"
)

// bloomFilter is a bloom filter implements methods used by leveldb storage.
type bloomFilter struct {
	// Number of hashing bits for a key in filter data.
	k int
	// Used to calculate the number of bytes of filter data to store a bunch of keys.
	bitsPerKey int
}

type bloomGenerator struct {
	k          int
	bitsPerKey int
	hashs      []uint32
}

var _ filter.Filter = (*bloomFilter)(nil)
var _ filter.Generator = (*bloomGenerator)(nil)

func (g *bloomGenerator) alignedSize(n int) (uint32, int) {
	bits := g.bitsPerKey * n
	switch {
	case bits < 64:
		return 64, 8
	default:
		bits = (bits + 7) & (^7)
		return uint32(bits), bits / 8
	}
}

func (g *bloomGenerator) Name() string {
	return filterName
}

func (g *bloomGenerator) Add(key []byte) {
	g.hashs = append(g.hashs, hash.Hash(key))
}

func (g *bloomGenerator) Reset() {
	g.hashs = g.hashs[:0]
}

func (g *bloomGenerator) Empty() bool {
	return len(g.hashs) == 0
}

func (g *bloomGenerator) Append(buf *bytes.Buffer) {
	n := len(g.hashs)
	if n == 0 {
		return
	}

	bits, numBytes := g.alignedSize(n)
	buf.Grow(numBytes + 1)

	l := buf.Len()
	buf.ReadFrom(io.LimitReader(util.ZeroByteReader, int64(numBytes)))
	k := g.k
	data := buf.Bytes()[l : l+numBytes]
	for i := 0; i < n; i++ {
		h := g.hashs[i]
		delta := (h >> 17) | (h << 15)
		for j := 0; j < k; j++ {
			pos := h % bits
			data[pos/8] |= 1 << (pos % 8)
			h += delta
		}
	}
	buf.WriteByte(byte(k))
	g.Reset()
}

// Name returns the name of this filter.
func (f *bloomFilter) Name() string {
	return filterName
}

// Append appends filter data generated by hashing keys to buf.
func (f *bloomFilter) Append(buf *bytes.Buffer, keys [][]byte) {
	var g bloomGenerator
	g.k = f.k
	g.bitsPerKey = f.bitsPerKey
	g.hashs = make([]uint32, 0, len(keys))
	for _, k := range keys {
		g.Add(k)
	}
	g.Append(buf)
}

// Contains reports whether key is hashed in provided filter data.
// False positive is allowed.
func (f *bloomFilter) Contains(data, key []byte) bool {
	n := len(data)
	if n <= 1 {
		return false
	}
	n -= 1

	k := int(data[n])
	if k > 30 {
		// Reserved for potentially new encodings for short bloom filters.
		// Consider it a match.
		return true
	}

	data = data[:n]
	bits := uint32(n * 8)

	h := hash.Hash(key)
	delta := (h >> 17) | (h << 15)
	for j := 0; j < k; j++ {
		pos := h % bits
		if (data[pos/8] & (1 << (pos % 8))) == 0 {
			return false
		}
		h += delta
	}
	return true
}

func (f *bloomFilter) NewGenerator() filter.Generator {
	return &bloomGenerator{k: f.k, bitsPerKey: f.bitsPerKey}
}

// NewFilter creates a bloom filter.
func NewFilter(bitsPerKey int) filter.Filter {
	if bitsPerKey < 0 {
		bitsPerKey = 0
	}
	// Round down intentionally to reduce probing cost a little bit.
	k := bitsPerKey * 69 / 100 // 0.69 =~ ln(2)
	switch {
	case k < 1:
		k = 1
	case k > 30:
		k = 30
	}
	return &bloomFilter{k: k, bitsPerKey: bitsPerKey}
}
