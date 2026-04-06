package bitset

import (
	"database/sql/driver"
	"math"
)

type ReadBitset struct {
	bits []byte
}

func NewReadBitset(n int) *ReadBitset {
	size := int(math.Ceil(float64(n) / 8))
	return &ReadBitset{bits: make([]byte, size)}
}

func FromBytes(data []byte) *ReadBitset {
	cp := make([]byte, len(data))
	copy(cp, data)
	return &ReadBitset{bits: cp}
}

func (r *ReadBitset) Mark(chapterIndex int) {
	bytePos := chapterIndex / 8
	bitPos := uint(chapterIndex % 8)
	for bytePos >= len(r.bits) {
		r.bits = append(r.bits, 0)
	}

	r.bits[bytePos] |= (1 << bitPos)
}

func (r *ReadBitset) Unmark(chapterIndex int) {
	bytePos := chapterIndex / 8
	bitPos := uint(chapterIndex % 8)
	if bytePos >= len(r.bits) {
		return
	}

	r.bits[bytePos] &^= (1 << bitPos)
}

func (r *ReadBitset) IsRead(chapterIndex int) bool {
	bytePos := chapterIndex / 8
	bitPos := uint(chapterIndex % 8)
	if bytePos >= len(r.bits) {
		return false
	}
	return r.bits[bytePos]&(1<<bitPos) != 0
}

func (r *ReadBitset) CountRead() int {
	count := 0
	for _, b := range r.bits {
		for b != 0 {
			b &= b - 1
			count++
		}
	}

	return count
}

func (r *ReadBitset) Resize(newMaxChapter int) {
	newSize := int(math.Ceil(float64(newMaxChapter) / 8))
	newBits := make([]byte, newSize)
	copy(newBits, r.bits)
	if newSize <= len(r.bits) {
		remainder := newMaxChapter % 8
		if remainder != 0 {
			mask := byte((1 << uint(remainder)) - 1)
			newBits[newSize-1] &= mask
		}
	}
	r.bits = newBits
}

func (r *ReadBitset) Cap() int        { return len(r.bits) * 8 }
func (r *ReadBitset) ToBytes() []byte { return r.bits }

// Implementing the Scanner and Valuer interfaces for database serialization
// This allows us to store the ReadBitset as a byte array in the database

func (r *ReadBitset) Scan(value any) error {
	if data, ok := value.([]byte); ok {
		temp := FromBytes(data)
		*r = *temp
	}

	return nil
}

func (r *ReadBitset) Value() (driver.Value, error) {
	return r.ToBytes(), nil
}
