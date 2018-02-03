package basicblock

import (
	"testing"
)

func TestgetConseqZeroes(t *testing.T) {
	if getConseqZeroes(byte(0)) != 8 {
		t.Fail()
	}
	if getConseqZeroes(byte(240)) != 0 {
		t.Fail()
	}
	if getConseqZeroes(byte(15)) != 4 {
		t.Fail()
	}
	if getConseqZeroes(byte(1)) != 7 {
		t.Fail()
	}
}

func TestHashesMatchDifficulties(t *testing.T) {
	hash1 := []byte{79, 0}
	hash2 := []byte{0, 1}
	if hashMatchesDifficulty(4, hash1) {
		t.Fail()
	}
	if !hashMatchesDifficulty(0, hash1) {
		t.Fail()
	}
	if hashMatchesDifficulty(16, hash2) {
		t.Fail()
	}
	if !hashMatchesDifficulty(15, hash2) {
		t.Fail()
	}
}

func TestEmptyBlockchain(t *testing.T) {
	blockChain := []BasicBlock{}
	if ValidBasicBlockchain(blockChain) {
		t.Fail()
	}
}

func TestValidBlockchain(t *testing.T) {
	blockChain := []BasicBlock{GenesisBlock}
	for i := 0; i < 5; i++ {
		blockChain = append(blockChain, blockChain[len(blockChain)-1].FindBlock([]byte{}))
	}
	if !ValidBasicBlockchain(blockChain) {
		t.Fail()
	}
}

func TestInvalidExtraBlock(t *testing.T) {
	blockChain := []BasicBlock{GenesisBlock}
	for i := 0; i < 5; i++ {
		blockChain = append(blockChain, blockChain[len(blockChain)-1].FindBlock([]byte{}))
	}
	blockChain = append(blockChain, BasicBlock{})
	if ValidBasicBlockchain(blockChain) {
		t.Fail()
	}
}

func TestInvalidGenesisBlock(t *testing.T) {
	blockChain := []BasicBlock{GenesisBlock}
	for i := 0; i < 5; i++ {
		blockChain = append(blockChain, blockChain[len(blockChain)-1].FindBlock([]byte{}))
	}
	blockChain[0].Data = []byte("DEADBEEF")
	if ValidBasicBlockchain(blockChain) {
		t.Fail()
	}
}

// func TestBlockchainReplaced(t *testing.T) {
// 	blockChainShort := []BasicBlock{GenesisBlock}
// 	blockChainLong := []BasicBlock{GenesisBlock}
//
// 	for i := 0; i < 1; i++ {
// 		blockChainLong = append(blockChainLong, blockChainLong[len(blockChainLong)-1].GenerateNextBasicBlock([]byte{}))
// 	}
//
// 	res := PossiblyReplace(blockChainShort, blockChainLong)
// 	if !deepEqual(res, blockChainShort) {
// 		t.Fail()
// 	}
// }
