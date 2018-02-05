package basicblock

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"time"
)

// BlockGenerationInterval in seconds, defines how often a block should be found. (in Bitcoin this value is 10 minutes)
const BlockGenerationInterval int = 10

// DifficultyAdjustmentInterval in blocks, defines how often the difficulty should adjust to the increasing or decreasing network hashrate. (in Bitcoin this value is 2016 blocks)
const DifficultyAdjustmentInterval int = 2

var Difficulty int32 = 2

// GenesisBlock is the very first block, duh! Package globals are usually bad!
var GenesisBlock BasicBlock

func init() {
	here, err := time.LoadLocation("UTC")
	if err != nil {
		debug("error!") // TODO
		return
	}
	GenesisBlock = BasicBlock{
		Index:     1,
		Timestamp: time.Date(1, time.January, 1, 1, 1, 1, 1, here),
		// previousHash takes on weird default value of "01000000"...
		Data:       []byte("this is the genesis block"),
		Difficulty: 2,
	}
}

// BasicBlock - Implementation of a block of cryptocurrency!
type BasicBlock struct {
	Index        int32
	Hash         [32]byte
	PreviousHash [32]byte
	Timestamp    time.Time
	Data         []byte
	Difficulty   int32
	Nonce        []byte
}

// BlockChain basic implementation
type BlockChain []BasicBlock

func (bb *BasicBlock) String() string {
	return fmt.Sprintf("(Index: %d, Hash: %x, PreviousHash: %x, Timestamp: %s, Data: %x, Difficulty: %d, Nonce %x)", bb.Index, bb.Hash, bb.PreviousHash, bb.Timestamp.Format(time.RFC3339), bb.Data, bb.Difficulty, bb.Nonce)
}

func (bc BlockChain) String() string {
	var s string
	for _, blk := range bc {
		s += blk.String() + "\n"
	}
	return s
}

func (bb *BasicBlock) deepEqual(bb2 *BasicBlock) bool {
	debug("deepEqual returned %t\n", len(bb.Data) == len(bb2.Data))
	if bb.Index == bb2.Index &&
		bb.Hash == bb2.Hash &&
		bb.PreviousHash == bb2.PreviousHash &&
		bb.Timestamp.Equal(bb2.Timestamp) && // works across timezones
		len(bb.Data) == len(bb2.Data) {
		for i, b := range bb.Data {
			if b != bb2.Data[i] {
				debug("deepEqual data byte %d different\n", i)
				return false
			}
		}
		return true
	}
	return false
}

func deepEqual(bc1, bc2 []BasicBlock) bool {
	if len(bc1) != len(bc2) {
		return false
	}
	for i, b := range bc1 {
		if !b.deepEqual(&bc2[i]) {
			return false
		}
	}
	return true
}

func (bb *BasicBlock) calculateHash() [32]byte {
	h := sha256.New()
	var hashInput bytes.Buffer
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, bb.Index)
	if err != nil {
		log.Fatalf("Int to binary conversion failed with error %v", err)
	}
	hashInput.Write(buf.Bytes())

	var prevHash []byte
	prevHash = bb.PreviousHash[:]
	hashInput.Write(prevHash)

	debug("previoushash: %x\n", hashInput.Bytes())
	timeStr := bb.Timestamp.Format(time.RFC3339)
	hashInput.WriteString(timeStr)

	hashInput.Write(bb.Data)
	debug("in: %x\n", hashInput.Bytes())

	hashInput.Write(bb.Nonce)

	h.Write(hashInput.Bytes())

	res := h.Sum(nil)
	var ret [32]byte
	for i, d := range res {
		ret[i] = d
	}
	return ret
}

// IsValid makes sure that the current BasicBlock has the correct Hash and PreviousHash.
func (bb *BasicBlock) IsValid(prev *BasicBlock) bool {
	computedHash := bb.calculateHash()
	return bb.PreviousHash == prev.Hash && computedHash == bb.Hash && hashMatchesDifficulty(bb.Difficulty, bb.Hash[:]) && bb.isValidTimestamp(prev)
}

// IsValid makes sure that the entire blockChain is valid
func (bc BlockChain) IsValid() bool {
	if len(bc) < 1 {
		debug("IsValidBasicBlockchain: Length of blockchain is 0.\n")
		return false
	}
	if !bc[0].deepEqual(&GenesisBlock) {
		debug("IsValidBasicBlockchain: Wrong genesis block.\n")
		return false
	}
	for i, blk := range bc {
		if i == 0 { // genesis block is already verified.
			continue
		} else {
			if !blk.IsValid(&bc[i-1]) {
				debug("IsValidBasicBlockchain: Block %d was invalid.\n", i)
				return false
			}
		}
	}
	return true
}

// isValidTimestamp is used to mitigate attacks in which a false timestamp is introduced in order to manipulate the difficulty. A block is valid, if the timestamp is at most 1 min in the future from the time we perceive. A block in the chain is valid, if the timestamp is at most 1 min in the past of the previous block.
func (bb *BasicBlock) isValidTimestamp(prev *BasicBlock) bool {
	return bb.Timestamp.After(prev.Timestamp.Add(-60*time.Second)) && bb.Timestamp.Before(time.Now().Add(60*time.Second))
}

// PossiblyReplace accepts a "contender blockchain", if the contender is valid AND has a larger cumulative difficulty than the blockchain we currently have, we replace it. Assumption: orig is valid.
func PossiblyReplace(orig BlockChain, next BlockChain) []BasicBlock {
	if !next.IsValid() {
		return orig
	}
	var cumOrig int32
	var cumNext int32
	for _, x := range orig {
		cumOrig += int32(math.Pow(2, float64(x.Difficulty)))
	}
	for _, x := range next {
		cumNext += int32(math.Pow(2, float64(x.Difficulty)))
	}
	if cumOrig > cumNext {
		return orig
	}
	return next
}

func getConseqZeroes(hash byte) int32 {
	b := byte(hash)
	if b&255 == 0 {
		return 8
	}
	if b&254 == 0 {
		return 7
	}
	if b&252 == 0 {
		return 6
	}
	if b&248 == 0 {
		return 5
	}
	if b&240 == 0 {
		return 4
	}
	if b&224 == 0 {
		return 3
	}
	if b&192 == 0 {
		return 2
	}
	if b&128 == 0 {
		return 1
	}
	return 0
}

// HashMatchesDifficulty makes sure that hash has at least (difficulty # of) leading zeroes.
func hashMatchesDifficulty(difficulty int32, hash []byte) bool {
	zeroes := int32(0)
	for _, x := range hash {
		debug("x is %08b\n", x)
		cz := getConseqZeroes(x)
		zeroes += cz
		if cz != 8 {
			break
		}
	}
	debug("hmd: hash had = %d, difficulty : %d\n\n", zeroes, difficulty)
	return zeroes >= difficulty
}

// FindBlock finds the next block with the expected difficulty.
func (bb *BasicBlock) FindBlock(data []byte) BasicBlock {
	nonceInt := int32(0) // TODO this is a problem! we may not always be able to find a solution with a limited number of bits
	result := &BasicBlock{
		Index:        bb.Index + 1,
		PreviousHash: bb.Hash,
		Timestamp:    time.Now(),
		Difficulty:   Difficulty,
		Nonce:        []byte{0},
		Data:         data,
	}
	for {
		var buf bytes.Buffer
		err := binary.Write(&buf, binary.LittleEndian, nonceInt)
		if err != nil {
			log.Fatal("Int to binary conversion failed")
		}

		result.Nonce = buf.Bytes()
		debug("n: %x\n", result.Nonce)

		hash := result.calculateHash()
		debug("h: %08b\n", hash)

		if hashMatchesDifficulty(result.Difficulty, hash[:]) {
		}
		nonceInt++
	}
}

// GetDifficulty calculates if the current difficulty needs to be adjusted.
func GetDifficulty(bc BlockChain) (int32, error) {
	if len(bc) == 0 {
		return 0, nil // TODO implement errors
	}
	latestBlock := bc[len(bc)-1]
	if latestBlock.Index%int32(DifficultyAdjustmentInterval) == 0 && latestBlock.Index != 0 {
		return getAdjustedDifficulty(latestBlock, bc)
	} else {
		return latestBlock.Difficulty, nil
	}
}

func getAdjustedDifficulty(latestBlock BasicBlock, bc BlockChain) (int32, error) {
	prevAdjustmentBlock := bc[len(bc)-DifficultyAdjustmentInterval]
	timeExpected := BlockGenerationInterval * DifficultyAdjustmentInterval
	timeTaken := latestBlock.Timestamp.Second() - prevAdjustmentBlock.Timestamp.Second()
	if timeTaken < timeExpected/2 {
		log.Println("difficulty up.")
		return prevAdjustmentBlock.Difficulty + 1, nil
	} else if timeTaken > timeExpected*2 {
		log.Println("difficulty down.")
		return prevAdjustmentBlock.Difficulty - 1, nil
	} else {
		log.Println("difficulty unchanged.")
		return prevAdjustmentBlock.Difficulty, nil
	}
}
