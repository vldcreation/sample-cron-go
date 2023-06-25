package hashs

import (
	"crypto/sha256"
	"hash"
)

type hashs struct {
	//
	// @param hasher
	// use as hasher to compose hash.Hash interface
	// use sha256 as default
	//
	hasher hash.Hash
	//
	// @param keep
	// keep the hasher state
	// if true, the hasher state will be kept after hash() method is called
	//
	keep bool
	//
	// @param b
	// the result of hash() method
	// if keep is true, the result will be kept here
	// if keep is false, the result will be returned directly
	// the state of hasher will be reset after hash() method is called
	//
	b []byte
}

func NewHashs(hasher hash.Hash, keep bool, source []byte) *hashs {
	if hasher == nil {
		// we use sha256 as default
		hasher = sha256.New()
	}

	return &hashs{
		hasher: hasher,
		keep:   keep,
		b:      source,
	}
}

type hasher interface {
	Hash(data []byte) []byte
}

func (h *hashs) Hash(data []byte) []byte {
	h.hasher.Write(data)
	if h.keep {
		_h := h.hasher.Sum(nil)
		h.b = h.hasher.Sum(h.b)
		h.hasher.Reset()
		return _h
	}
	b := h.hasher.Sum(nil)
	h.hasher.Reset()
	return b
}
