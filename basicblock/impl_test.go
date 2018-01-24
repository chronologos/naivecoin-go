package basicblock

import (
	"testing"
)

func TestEmptyBlockchain(t *testing.T) {
	blockChain := []BasicBlock{}
	if ValidBasicBlockchain(blockChain) {
		t.Fail()
	}
}

func TestValidBlockchain(t *testing.T) {
	blockChain := []BasicBlock{GenesisBlock}
	for i := 0; i < 5; i++ {
		blockChain = append(blockChain, blockChain[len(blockChain)-1].GenerateNextBasicBlock([]byte{}))
	}
	if !ValidBasicBlockchain(blockChain) {
		t.Fail()
	}
}

func TestInvalidExtraBlock(t *testing.T) {
	blockChain := []BasicBlock{GenesisBlock}
	for i := 0; i < 5; i++ {
		blockChain = append(blockChain, blockChain[len(blockChain)-1].GenerateNextBasicBlock([]byte{}))
	}
	blockChain = append(blockChain, BasicBlock{})
	if ValidBasicBlockchain(blockChain) {
		t.Fail()
	}
}

func TestInvalidGenesisBlock(t *testing.T) {
	blockChain := []BasicBlock{GenesisBlock}
	for i := 0; i < 5; i++ {
		blockChain = append(blockChain, blockChain[len(blockChain)-1].GenerateNextBasicBlock([]byte{}))
	}
	blockChain[0].Data = []byte("DEADBEEF")
	if ValidBasicBlockchain(blockChain) {
		t.Fail()
	}
}

func TestInvalidGenesisBlock2(t *testing.T) {
	blockChain := []BasicBlock{GenesisBlock}
	for i := 0; i < 5; i++ {
		blockChain = append(blockChain, blockChain[len(blockChain)-1].GenerateNextBasicBlock([]byte{}))
	}
	blockChain[0].Data = []byte("this is not genesis block")
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
