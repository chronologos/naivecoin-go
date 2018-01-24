package basicblock

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

// Hashable32 is wrapper over our implementations TODO
type Hashable32 interface {
	calculateHash() [32]byte
	getHash() [32]byte
}

// BasicBlock - Implementation of a block of cryptocurrency!
type BasicBlock struct {
	Index        uint32
	Hash         [32]byte
	PreviousHash [32]byte
	Timestamp    time.Time
	Data         []byte
}

func (bb *BasicBlock) String() string {
	return fmt.Sprintf("(i: %d, h: %x, ph: %x, ts: %s, d: %x)", bb.Index, bb.Hash, bb.PreviousHash, bb.Timestamp.Format(time.RFC3339), bb.Data)
}

func (bb *BasicBlock) deepEqual(bb2 *BasicBlock) bool {
	debug("gb comp %t\n", len(bb.Data) == len(bb2.Data))
	if bb.Index == bb2.Index &&
		bb.Hash == bb2.Hash &&
		bb.PreviousHash == bb2.PreviousHash &&
		bb.Timestamp.Equal(bb2.Timestamp) && //works across timezones
		len(bb.Data) == len(bb2.Data) {
		for i, b := range bb.Data {
			if b != bb2.Data[i] {
				debug("gb comp data byte %d different\n", i)
				return false
			}
		}
		return true
	}
	return false
}

// GenesisBlock is the very first block, duh! Package globals are usually bad!
var GenesisBlock BasicBlock

func init() {
	here, err := time.LoadLocation("Singapore")
	if err != nil {
		debug("error!") // TODO
		return
	}
	GenesisBlock = BasicBlock{
		Index:     1,
		Timestamp: time.Date(1, time.January, 1, 1, 1, 1, 1, here),
		// previousHash takes on weird default value of "01000000"...
		Data: []byte("this is the genesis block"),
	}
}

// calculateHash for BasicBlock uses all fields of block except Hash, which is being calculated, to build the hash.
func (bb *BasicBlock) calculateHash() [32]byte {
	/*
		h := sha256.New()
		fmt.Print(len(bb.PreviousHash) + len(bb.Data))
		in := make([]byte, len(bb.PreviousHash)+len(bb.Data))
		in = append(in, bb.PreviousHash...)
		in = append(in, bb.Data...)
		fmt.Printf("in: %x\n previoushash: %x\n data: %x\n", in, bb.PreviousHash, bb.Data)
		h.Write(in)
		fmt.Printf("hash1: %x\n", h.Sum(nil))
	*/
	h := sha256.New()
	timeStr := bb.Timestamp.Format(time.RFC3339)
	var in2 bytes.Buffer
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, bb.Index)
	in2.Write(bs)
	var prevHash []byte
	prevHash = bb.PreviousHash[:]
	in2.Write(prevHash)
	debug("previoushash: %x\n", in2.Bytes())
	in2.WriteString(timeStr)
	in2.Write(bb.Data)
	debug("in: %x\n", in2.Bytes())
	h.Write(in2.Bytes())
	res := h.Sum(nil)
	var ret [32]byte
	for i, d := range res {
		ret[i] = d
	}
	return ret
}

func (bb *BasicBlock) getHash() [32]byte {
	return bb.Hash
}

// ValidBasicBlock makes sure that the current BasicBlock has the correct Hash and PreviousHash.
func (bb *BasicBlock) ValidBasicBlock(prev *BasicBlock) bool {
	computedHash := bb.calculateHash()
	return bb.PreviousHash == prev.Hash && computedHash == bb.Hash
}

// ValidBasicBlockchain makes sure that the entire blockChain is valid
func ValidBasicBlockchain(bc []BasicBlock) bool {
	if len(bc) < 1 {
		debug("ValidBasicBlockchain: Length of blockchain is 0.\n")
		return false
	}
	if !bc[0].deepEqual(&GenesisBlock) {
		debug("ValidBasicBlockchain: Wrong genesis block.\n")
		return false
	}
	for i, blk := range bc {
		if i == 0 { // genesis block is already verified.
			continue
		} else {
			if !blk.ValidBasicBlock(&bc[i-1]) {
				debug("ValidBasicBlockchain: Block %d was invalid.\n", i)
				return false
			}
		}
	}
	return true
}

// GenerateNextBasicBlock computes the next basic block in the blockChain
func (bb *BasicBlock) GenerateNextBasicBlock(d []byte) (next BasicBlock) {
	next = BasicBlock{
		Index:        bb.Index + 1,
		PreviousHash: bb.getHash(),
		Data:         d,
		Timestamp:    time.Now()}
	next.Hash = next.calculateHash()
	return next
}

// PossiblyReplace accepts a "contender blockchain", if the contender is valid AND longer than the blockchain we currently have, we replace it. Assumption: orig is valid.
func PossiblyReplace(orig []BasicBlock, next []BasicBlock) []BasicBlock {
	if !ValidBasicBlockchain(next) || !(len(next) > len(orig)) {
		return orig
	}
	return next
}
