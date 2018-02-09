package basicblock

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
)

// Coinbase is a transaction that contains only an output, but no inputs. This means that a coinbase transaction adds new coins to circulation. The coinbase transaction is always the first transaction in the block and it is included by the miner of the block. The coinbase reward acts as an incentive for the miners: if you find the block, you are able to collect 50 coins.

const CoinbaseAmount = 50

// TxOut consists of an address and an amount of coins. The address is an ECDSA public-key. This means that the user having the private-key of the referenced public-key (=address) will be able to access the coins.
type TxOut struct {
	address ecdsa.PublicKey
	amount  int32
}

// TxIn provides the information "where" the coins are coming from. Each TxIn refers to an earlier output, from which the coins are 'unlocked', with the signature. These unlocked coins are now 'available' for the TxOuts. The signature gives proof that only the user, that has the private-key of the referred public-key ( =address) could have created the transaction.
type TxIn struct {
	txOutID    [32]byte
	txOutIndex int32
	r          *big.Int
	s          *big.Int
}

// Transaction consists of two components: inputs and outputs. Outputs specify where the coins are sent and inputs give a proof that the coins that are actually sent exists in the first place and are owned by the sender.
type Transaction struct {
	// id is hash of txIns and txOuts
	id     [32]byte
	txIns  []TxIn
	txOuts []TxOut
}

type UnspentTxOut struct {
	txOutId    [32]byte // Transaction id
	txOutIndex int32    // index of txOut in Transaction.txOuts
	address    ecdsa.PublicKey
	amount     int32
}

type TxErrorClass int

const (
	Generic TxErrorClass = iota
	TxNotFound
	SigningError
)

type TxError struct {
	msg  string
	kind TxErrorClass
}

func (txerror TxError) Error() string {
	return txerror.msg
}

func findUnspentTxOut(txOutId [32]byte, txOutIndex int32, aUnspentTxOuts []UnspentTxOut) (UnspentTxOut, error) {
	for _, aUnspentTxOut := range aUnspentTxOuts {
		if aUnspentTxOut.txOutId == txOutId && aUnspentTxOut.txOutIndex == txOutIndex {
			return aUnspentTxOut, nil
		}
	}
	return UnspentTxOut{}, TxError{"Tx not found", TxNotFound}
}

func checkBinaryWrite(err error) {
	if err != nil {
		log.Fatalf("binary.Write failed with error %v", err)
	}
}

func checkGobEncode(err error) {
	if err != nil {
		log.Fatalf("GobEncode failed with error %v", err)
	}
}

func (tx Transaction) getID() [32]byte {
	h := sha256.New()
	var hashInput bytes.Buffer
	for _, txIn := range tx.txIns {
		hashInput.Write(txIn.txOutID[:])
		var idx bytes.Buffer
		err := binary.Write(&idx, binary.LittleEndian, txIn.txOutIndex)
		checkBinaryWrite(err)
		hashInput.Write(idx.Bytes())
	}
	for _, txOut := range tx.txOuts {
		var b bytes.Buffer
		xMarshalled, err := txOut.address.X.GobEncode()
		checkGobEncode(err)
		err = binary.Write(&b, binary.LittleEndian, xMarshalled)
		checkBinaryWrite(err)
		yMarshalled, err := txOut.address.Y.GobEncode()
		checkGobEncode(err)
		err = binary.Write(&b, binary.LittleEndian, yMarshalled)
		checkBinaryWrite(err)

		err = binary.Write(&b, binary.LittleEndian, txOut.amount)
		checkBinaryWrite(err)
		hashInput.Write(b.Bytes())
	}
	_, err := h.Write(hashInput.Bytes())
	if err != nil {
		log.Fatalln("sha256 failed")
	}
	var res [32]byte
	hashOutput := h.Sum(nil)
	if len(hashOutput) != 32 {
		log.Fatalln("SHA256 hash output is larger than 32 bytes")
	}
	for i := 0; i < 32; i++ {
		res[i] = hashOutput[i]
	}
	return res
}

func (tx Transaction) signTxIn(txInIndex int32, privateKey ecdsa.PrivateKey, aUnspentTxOuts []UnspentTxOut) (*big.Int, *big.Int, error) {
	txIn := tx.txIns[txInIndex]
	dataToSign := tx.id
	referencedUnspentTxOut, err := findUnspentTxOut(txIn.txOutID, txIn.txOutIndex, aUnspentTxOuts)
	if err != nil {
		return nil, nil, err
	}
	referencedAddress := referencedUnspentTxOut.address
	if privateKey.Public() != referencedAddress {
		return nil, nil, TxError{"trying to sign an input with private key that does not match the address that is referenced in txIn", SigningError}
	}
	r, s, err := ecdsa.Sign(nil, &privateKey, dataToSign[:])
	if err != nil {
		return nil, nil, err
	}
	return r, s, nil
}

func updateUnspentTxOuts(txs []Transaction, aUnspentTxOuts []UnspentTxOut) []UnspentTxOut {
	var newUnspentTxOuts []UnspentTxOut
	for _, tx := range txs {
		for idx, txOut := range tx.txOuts {
			newUnspentTxOuts = append(newUnspentTxOuts, UnspentTxOut{tx.id, int32(idx), txOut.address, txOut.amount})
		}
	}
	var consumedTxOuts []UnspentTxOut
	for _, tx := range txs {
		for _, txIn := range tx.txIns {
			consumedTxOuts = append(consumedTxOuts, UnspentTxOut{txIn.txOutID, txIn.txOutIndex, ecdsa.PublicKey{}, 0})
		}
	}
	var resultingUnspentTxOuts []UnspentTxOut
	for _, utxo := range aUnspentTxOuts {
		utxo, err := findUnspentTxOut(utxo.txOutId, utxo.txOutIndex, consumedTxOuts)
		txerr, ok := err.(TxError)
		if ok && txerr.kind == TxNotFound {
			resultingUnspentTxOuts = append(resultingUnspentTxOuts, utxo)
		}
	}
	resultingUnspentTxOuts = append(resultingUnspentTxOuts, newUnspentTxOuts...)
	return resultingUnspentTxOuts
}

// blockHeight is the number of blocks in the chain between it and the genesis block. (So the genesis block has height 0.)
func validateCoinbaseTx(tx Transaction, blockHeight int32) bool {
	if tx.getID() != tx.id || tx.txIns[0].txOutIndex != blockHeight || len(tx.txIns) != 1 || len(tx.txOuts) != 1 || tx.txOuts[0].amount != CoinbaseAmount {
		fmt.Printf("validateCoinbaseTx failed \n id not equal = %t, txOutIndex not equal blockHeight = %t, length txIns not equal 1 = %t, length txOuts not equal 1 = %t, amount not equal CoinbaseAmount = %t \n", tx.getID() != tx.id, tx.txIns[0].txOutIndex != blockHeight, len(tx.txIns) != 1, len(tx.txOuts) != 1, tx.txOuts[0].amount != CoinbaseAmount)
		return false
	}
	return true
}
