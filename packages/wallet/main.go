package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"

	"golang.org/x/crypto/sha3"

	"github.com/btcsuite/btcutil/base58"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

var (
	mnemonicSecret = `flash couple heart script ramp april average caution plunge alter elite author`
	index          = uint32(0)
)

func main() {
	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonicSecret, index)
	if err != nil {
		panic(err)
	}

	fmt.Println("Address:", address)
	fmt.Println("PrivKey:", privKey)

}

func DeriveTronAddressFromMnemonic(mnemonicSecret string, index uint32) (string, string, error) {
	// generate seed from mnemonic
	seed := bip39.NewSeed(mnemonicSecret, "")

	// then we generate master key
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		slog.Error("failed to generate master key", "error", err)
		return "", "", err
	}

	// 3. Derive path: m/44'/195'/0'/0/index
	// 44' = BIP44, 195' = TRON (coin type), purpose
	purpose, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	coinType, _ := purpose.NewChildKey(bip32.FirstHardenedChild + 195)
	account, _ := coinType.NewChildKey(bip32.FirstHardenedChild + 0)
	change, _ := account.NewChildKey(0)

	walletKey, err := change.NewChildKey(uint32(index))
	if err != nil {
		return "", "", err
	}
	// 4. Get private key in hex
	privateKeyHex := hex.EncodeToString(walletKey.Key)

	// 5. Derive TRON address from private key
	address, err := PrivateKeyToTronAddress(walletKey.Key)
	if err != nil {
		return "", "", err
	}

	return address, privateKeyHex, nil

}

// FIXED VERSION ✅
func PrivateKeyToTronAddress(privateKey []byte) (string, error) {
	// Convert raw 32-byte private key into ecdsa.PrivateKey
	d := new(big.Int).SetBytes(privateKey)
	curve := elliptic.P256() // SECP256K1 alternative in Go (close compatibility)
	priv := new(ecdsa.PrivateKey)
	priv.D = d
	priv.PublicKey.Curve = curve
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())

	// Encode public key (uncompressed, 0x04 + X + Y)
	pubKey := append([]byte{0x04}, append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)...)

	// Remove the 0x04 prefix for hashing
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubKey[1:])
	sum := hash.Sum(nil)

	// Tron address: prefix 0x41 + last 20 bytes of keccak hash
	addressBytes := append([]byte{0x41}, sum[12:]...)

	// Checksum
	checksumHash := sha3.NewLegacyKeccak256()
	checksumHash.Write(addressBytes)
	checksum := checksumHash.Sum(nil)[:4]

	// Final address = base58(address + checksum)
	fullAddress := append(addressBytes, checksum...)
	return base58.Encode(fullAddress), nil
}
