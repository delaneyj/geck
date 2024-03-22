package out

import (
	"github.com/RoaringBitmap/roaring/roaring64"
)

const ID_BIT_SIZE = 31
const TargetBitMask = 1<<ID_BIT_SIZE - 1
const SourceBitMask = TargetBitMask << 32
const PairMask = 1 << 31

type ID uint64

const UserDefined = 1000

func (id ID) U64() uint64 {
	return uint64(id)
}

func (id ID) IsPair() bool {
	// check if the second highest bit is set
	return id&PairMask != 0
}

func (id ID) SplitPair() (source, target ID, wasPair bool) {
	if !id.IsPair() {
		return 0, id.Target(), false
	}
	return id.Source(), id.Target(), true
}

func (id ID) Index() uint32 {
	// lower 32 bits
	return uint32(id) & TargetBitMask
}

func (id ID) Generation() uint16 {
	// upper 16 bits
	return uint16(id >> 32)
}

func (id ID) UpdateGeneration() ID {
	return id + 1<<32
}

func NewEntity(index uint32, generation uint16) ID {
	return ID(generation)<<32 | ID(index)
}

func NewPair(source, target ID) ID {
	source &= TargetBitMask
	target &= TargetBitMask
	return source<<32 | target | PairMask
}

func (id ID) Source() ID {
	return id & SourceBitMask >> 32
}

func (id ID) Target() ID {
	return id & TargetBitMask
}

type IDSet struct {
	bits *roaring64.Bitmap
}

func NewIDSetFromUint64s(ids ...uint64) *IDSet {
	bits := roaring64.NewBitmap()
	for _, id := range ids {
		bits.Add(id)
	}
	return &IDSet{bits}
}

func NewIDSet(ids ...ID) *IDSet {
	bits := roaring64.NewBitmap()
	for _, id := range ids {
		bits.Add(uint64(id))
	}
	return &IDSet{bits}
}

func NewIDSetFromBase64(b64 string) (*IDSet, error) {
	bits := roaring64.New()
	if _, err := bits.FromBase64(b64); err != nil {
		return nil, err
	}
	set := &IDSet{bits: bits}
	return set, nil
}

func (set *IDSet) String() string {
	return set.bits.String()
}

func (set *IDSet) ToUint64s() []uint64 {
	return set.bits.ToArray()
}

func (set *IDSet) ToBase64() string {
	b64, err := set.bits.ToBase64()
	if err != nil {
		panic(err)
	}
	return b64
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for the bitmap
func (set *IDSet) MarshalBinary() []byte {
	b, err := set.bits.ToBytes()
	if err != nil {
		panic(err)
	}
	return b
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for the bitmap
func (set *IDSet) UnmarshalBinary(data []byte) error {
	return set.bits.UnmarshalBinary(data)
}

// Range iterates over the set and calls the specified function for each id
func (set *IDSet) Range(fn func(id ID)) {
	iter := set.bits.Iterator()
	for iter.HasNext() {
		fn(ID(iter.Next()))
	}
}

// ConditionalRange iterates over the set and calls the specified function for each id
func (set *IDSet) ConditionalRange(fn func(id ID, i int) bool) {
	iter := set.bits.Iterator()
	i := 0
	for iter.HasNext() {
		id := ID(iter.Next())
		keepGoing := fn(id, i)
		if !keepGoing {
			break
		}
		i++
	}
}

// Contains returns true if the set contains the specified id
func (set *IDSet) Contains(id ID) bool {
	return set.bits.Contains(uint64(id))
}

// Clone returns a deep copy of the Bitmap
func (set *IDSet) Clone() *IDSet {
	return &IDSet{set.bits.Clone()}
}

// Add adds the specified id to the set
func (set *IDSet) Add(ids ...ID) *IDSet {
	for _, id := range ids {
		set.bits.Add(uint64(id))
	}
	return set
}

// AddRange adds the integers in [rangeStart, rangeEnd) to the bitmap.
func (set *IDSet) AddRange(from, to ID) *IDSet {
	set.bits.AddRange(uint64(from), uint64(to))
	return set
}

// Remove removes the specified id from the set
func (set *IDSet) Remove(ids ...ID) *IDSet {
	for _, id := range ids {
		set.bits.Remove(uint64(id))
	}
	return set
}

// RemoveRange removes the integers in [rangeStart, rangeEnd) from the bitmap.
func (set *IDSet) RemoveRange(from, to ID) *IDSet {
	set.bits.RemoveRange(uint64(from), uint64(to))
	return set
}

// Pop removes the last integer from the bitmap and returns it
func (set *IDSet) Pop() ID {
	max := set.bits.Maximum()
	set.bits.Remove(max)
	return ID(max)
}

// IsEmpty returns true if the bitmap is empty (contains no integers)
func (set *IDSet) IsEmpty() bool {
	return set.bits.IsEmpty()
}

// Equals returns true if the two bitmaps contain the same integers
func (set *IDSet) Equals(other *IDSet) bool {
	return set.bits.Equals(other.bits)
}

// Cardinality returns the number of integers contained in the bitmap
func (set *IDSet) Cardinality() int {
	return int(set.bits.GetCardinality())
}

// And computes the intersection between two bitmaps and stores the result in the current bitmap
func (set *IDSet) And(other *IDSet) *IDSet {
	set.bits.And(other.bits)
	return set
}

// Or computes the union between two bitmaps and stores the result in the current bitmap
func (set *IDSet) Or(other *IDSet) *IDSet {
	set.bits.Or(other.bits)
	return set
}

// Xor computes the symmetric difference between two bitmaps and stores the result in the current bitmap
func (set *IDSet) Xor(other *IDSet) *IDSet {
	set.bits.Xor(other.bits)
	return set
}

// AndNot computes the difference between two bitmaps and stores the result in the current bitmap
func (set *IDSet) Not(other *IDSet) *IDSet {
	set.bits.AndNot(other.bits)
	return set
}

// AndCardinality returns the cardinality of the intersection between two bitmaps, bitmaps are not modified
func (set *IDSet) AndCardinality(other *IDSet) int {
	return int(set.bits.AndCardinality(other.bits))
}

// OrCardinality returns the cardinality of the union between two bitmaps, bitmaps are not modified
func (set *IDSet) OrCardinality(other *IDSet) int {
	return int(set.bits.OrCardinality(other.bits))
}

// Clear resets the Bitmap to be logically empty, but may retain some memory allocations that may speed up future operations
func (set *IDSet) Clear() *IDSet {
	set.bits.Clear()
	return set
}

// Minimum returns the smallest integer in the bitmap
func (set *IDSet) Minimum() ID {
	return ID(set.bits.Minimum())
}

// Maximum returns the largest integer in the bitmap
func (set *IDSet) Maximum() ID {
	return ID(set.bits.Maximum())
}
