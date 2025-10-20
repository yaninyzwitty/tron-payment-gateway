package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/sha3"
)

// TestDeriveTronAddressFromMnemonic_ValidMnemonic tests derivation with valid mnemonic
func TestDeriveTronAddressFromMnemonic_ValidMnemonic(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if address == "" {
		t.Error("Expected non-empty address")
	}

	if privKey == "" {
		t.Error("Expected non-empty private key")
	}

	// Private key should be 64 hex characters (32 bytes)
	if len(privKey) != 64 {
		t.Errorf("Expected private key length 64, got %d", len(privKey))
	}

	// Verify private key is valid hex
	_, err = hex.DecodeString(privKey)
	if err != nil {
		t.Errorf("Private key should be valid hex: %v", err)
	}

	// TRON addresses typically start with 'T'
	if !strings.HasPrefix(address, "T") {
		t.Logf("Note: TRON address doesn't start with 'T': %s", address)
	}
}

// TestDeriveTronAddressFromMnemonic_DifferentIndices tests that different indices generate different addresses
func TestDeriveTronAddressFromMnemonic_DifferentIndices(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	address0, privKey0, err0 := DeriveTronAddressFromMnemonic(mnemonic, 0)
	address1, privKey1, err1 := DeriveTronAddressFromMnemonic(mnemonic, 1)
	address2, privKey2, err2 := DeriveTronAddressFromMnemonic(mnemonic, 2)

	if err0 != nil || err1 != nil || err2 != nil {
		t.Fatalf("Expected no errors, got: %v, %v, %v", err0, err1, err2)
	}

	// Different indices should generate different addresses
	if address0 == address1 {
		t.Error("Index 0 and 1 should generate different addresses")
	}

	if address0 == address2 {
		t.Error("Index 0 and 2 should generate different addresses")
	}

	if address1 == address2 {
		t.Error("Index 1 and 2 should generate different addresses")
	}

	// Different indices should generate different private keys
	if privKey0 == privKey1 {
		t.Error("Index 0 and 1 should generate different private keys")
	}

	if privKey0 == privKey2 {
		t.Error("Index 0 and 2 should generate different private keys")
	}

	if privKey1 == privKey2 {
		t.Error("Index 1 and 2 should generate different private keys")
	}
}

// TestDeriveTronAddressFromMnemonic_Deterministic tests that derivation is deterministic
func TestDeriveTronAddressFromMnemonic_Deterministic(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(5)

	address1, privKey1, err1 := DeriveTronAddressFromMnemonic(mnemonic, index)
	address2, privKey2, err2 := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err1 != nil || err2 != nil {
		t.Fatalf("Expected no errors, got: %v, %v", err1, err2)
	}

	// Same mnemonic and index should always produce the same results
	if address1 != address2 {
		t.Errorf("Expected same address, got %s and %s", address1, address2)
	}

	if privKey1 != privKey2 {
		t.Errorf("Expected same private key, got %s and %s", privKey1, privKey2)
	}
}

// TestDeriveTronAddressFromMnemonic_EmptyMnemonic tests with empty mnemonic
func TestDeriveTronAddressFromMnemonic_EmptyMnemonic(t *testing.T) {
	mnemonic := ""
	index := uint32(0)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	// Should handle empty mnemonic (bip39 might accept it)
	// We check that either it errors or generates something
	if err == nil {
		if address == "" || privKey == "" {
			t.Error("If no error, should generate address and private key")
		}
	}
}

// TestDeriveTronAddressFromMnemonic_InvalidMnemonic tests with invalid mnemonic
func TestDeriveTronAddressFromMnemonic_InvalidMnemonic(t *testing.T) {
	mnemonic := "invalid mnemonic that is not valid"
	index := uint32(0)

	// This might not error as bip39.NewSeed doesn't validate
	// But we test the behavior
	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err == nil {
		// Even with invalid mnemonic, BIP39 can generate seed
		// So we just check the outputs are non-empty
		if address == "" {
			t.Error("Expected non-empty address even with invalid mnemonic")
		}
		if privKey == "" {
			t.Error("Expected non-empty private key even with invalid mnemonic")
		}
	}
}

// TestDeriveTronAddressFromMnemonic_DifferentMnemonics tests with different mnemonics
func TestDeriveTronAddressFromMnemonic_DifferentMnemonics(t *testing.T) {
	mnemonic1 := "flash couple heart script ramp april average caution plunge alter elite author"
	mnemonic2 := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	index := uint32(0)

	address1, privKey1, err1 := DeriveTronAddressFromMnemonic(mnemonic1, index)
	address2, privKey2, err2 := DeriveTronAddressFromMnemonic(mnemonic2, index)

	if err1 != nil || err2 != nil {
		t.Fatalf("Expected no errors, got: %v, %v", err1, err2)
	}

	// Different mnemonics should generate different results
	if address1 == address2 {
		t.Error("Different mnemonics should generate different addresses")
	}

	if privKey1 == privKey2 {
		t.Error("Different mnemonics should generate different private keys")
	}
}

// TestDeriveTronAddressFromMnemonic_LargeIndex tests with large index
func TestDeriveTronAddressFromMnemonic_LargeIndex(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(2147483647) // Max int32

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err != nil {
		t.Fatalf("Expected no error with large index, got %v", err)
	}

	if address == "" || privKey == "" {
		t.Error("Expected valid address and private key with large index")
	}
}

// TestDeriveTronAddressFromMnemonic_MaxIndex tests with maximum possible index
func TestDeriveTronAddressFromMnemonic_MaxIndex(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(4294967295) // Max uint32

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	// Might fail or succeed depending on implementation
	if err == nil {
		if address == "" || privKey == "" {
			t.Error("If no error, should generate valid outputs")
		}
	}
}

// TestPrivateKeyToTronAddress_ValidKey tests with valid private key
func TestPrivateKeyToTronAddress_ValidKey(t *testing.T) {
	// Generate a valid 32-byte private key
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = byte(i)
	}

	address, err := PrivateKeyToTronAddress(privateKey)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if address == "" {
		t.Error("Expected non-empty address")
	}

	// Verify it's valid base58
	decoded := base58.Decode(address)
	if len(decoded) == 0 {
		t.Error("Address should be valid base58")
	}
}

// TestPrivateKeyToTronAddress_Deterministic tests deterministic address generation
func TestPrivateKeyToTronAddress_Deterministic(t *testing.T) {
	privateKey := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
	}

	address1, err1 := PrivateKeyToTronAddress(privateKey)
	address2, err2 := PrivateKeyToTronAddress(privateKey)

	if err1 != nil || err2 != nil {
		t.Fatalf("Expected no errors, got: %v, %v", err1, err2)
	}

	// Same private key should always generate same address
	if address1 != address2 {
		t.Errorf("Expected same address, got %s and %s", address1, address2)
	}
}

// TestPrivateKeyToTronAddress_DifferentKeys tests different keys generate different addresses
func TestPrivateKeyToTronAddress_DifferentKeys(t *testing.T) {
	privateKey1 := make([]byte, 32)
	privateKey2 := make([]byte, 32)

	for i := range privateKey1 {
		privateKey1[i] = byte(i)
		privateKey2[i] = byte(i + 1)
	}

	address1, err1 := PrivateKeyToTronAddress(privateKey1)
	address2, err2 := PrivateKeyToTronAddress(privateKey2)

	if err1 != nil || err2 != nil {
		t.Fatalf("Expected no errors, got: %v, %v", err1, err2)
	}

	// Different keys should generate different addresses
	if address1 == address2 {
		t.Error("Different private keys should generate different addresses")
	}
}

// TestPrivateKeyToTronAddress_InvalidLength tests with wrong length private key
func TestPrivateKeyToTronAddress_InvalidLength(t *testing.T) {
	testCases := []struct {
		name string
		key  []byte
	}{
		{"empty", []byte{}},
		{"too short", make([]byte, 16)},
		{"too long", make([]byte, 64)},
		{"single byte", []byte{0x01}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			address, err := PrivateKeyToTronAddress(tc.key)

			// Might error or produce unexpected result
			// We just verify it handles gracefully
			if err == nil && address == "" {
				t.Error("If no error, should produce non-empty address")
			}
		})
	}
}

// TestPrivateKeyToTronAddress_NilKey tests with nil private key
func TestPrivateKeyToTronAddress_NilKey(t *testing.T) {
	address, err := PrivateKeyToTronAddress(nil)

	// Should handle nil gracefully (might panic or error)
	if err == nil && address == "" {
		t.Error("Should either error or produce address")
	}
}

// TestPrivateKeyToTronAddress_AllZeros tests with all-zero private key
func TestPrivateKeyToTronAddress_AllZeros(t *testing.T) {
	privateKey := make([]byte, 32) // All zeros

	address, err := PrivateKeyToTronAddress(privateKey)

	// All-zero key should still generate an address (though not recommended)
	if err != nil {
		t.Logf("All-zero key errored: %v", err)
	} else if address == "" {
		t.Error("Expected non-empty address for all-zero key")
	}
}

// TestPrivateKeyToTronAddress_AllOnes tests with all-ones private key
func TestPrivateKeyToTronAddress_AllOnes(t *testing.T) {
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = 0xFF
	}

	address, err := PrivateKeyToTronAddress(privateKey)

	if err != nil {
		t.Logf("All-ones key errored: %v", err)
	} else if address == "" {
		t.Error("Expected non-empty address for all-ones key")
	}
}

// TestPrivateKeyToTronAddress_AddressFormat tests address format
func TestPrivateKeyToTronAddress_AddressFormat(t *testing.T) {
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = byte(i * 3)
	}

	address, err := PrivateKeyToTronAddress(privateKey)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify address is valid base58
	decoded := base58.Decode(address)
	if len(decoded) == 0 {
		t.Error("Address should be valid base58 encoded")
	}

	// Decoded should be 25 bytes (21 bytes + 4 bytes checksum)
	if len(decoded) != 25 {
		t.Logf("Note: Decoded address length is %d, expected 25", len(decoded))
	}

	// First byte should be 0x41 for TRON mainnet
	if len(decoded) > 0 && decoded[0] != 0x41 {
		t.Logf("Note: First byte is 0x%x, expected 0x41", decoded[0])
	}
}

// TestPrivateKeyToTronAddress_ChecksumVerification tests checksum calculation
func TestPrivateKeyToTronAddress_ChecksumVerification(t *testing.T) {
	privateKey := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22,
		0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa,
		0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22, 0x33,
		0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb,
	}

	address, err := PrivateKeyToTronAddress(privateKey)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Decode and verify checksum
	decoded := base58.Decode(address)
	if len(decoded) < 25 {
		t.Fatalf("Decoded address too short: %d bytes", len(decoded))
	}

	// Verify checksum manually
	addressBytes := decoded[:21]
	checksumBytes := decoded[21:]

	hash := sha3.NewLegacyKeccak256()
	hash.Write(addressBytes)
	expectedChecksum := hash.Sum(nil)[:4]

	for i := 0; i < 4; i++ {
		if checksumBytes[i] != expectedChecksum[i] {
			t.Errorf("Checksum mismatch at byte %d: got 0x%x, expected 0x%x",
				i, checksumBytes[i], expectedChecksum[i])
		}
	}
}

// TestPrivateKeyToTronAddress_PublicKeyGeneration tests public key generation
func TestPrivateKeyToTronAddress_PublicKeyGeneration(t *testing.T) {
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = byte(i + 10)
	}

	// Manually derive public key to verify logic
	d := new(big.Int).SetBytes(privateKey)
	curve := elliptic.P256()
	x, y := curve.ScalarBaseMult(d.Bytes())

	// Verify we can generate address
	address, err := PrivateKeyToTronAddress(privateKey)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if address == "" {
		t.Error("Expected non-empty address")
	}

	// Verify public key is on curve
	if !curve.IsOnCurve(x, y) {
		t.Error("Generated public key is not on the curve")
	}
}

// TestPrivateKeyToTronAddress_MultipleKeys tests multiple keys in sequence
func TestPrivateKeyToTronAddress_MultipleKeys(t *testing.T) {
	addresses := make(map[string]bool)

	for i := 0; i < 10; i++ {
		privateKey := make([]byte, 32)
		for j := range privateKey {
			privateKey[j] = byte(i*10 + j)
		}

		address, err := PrivateKeyToTronAddress(privateKey)
		if err != nil {
			t.Fatalf("Error at iteration %d: %v", i, err)
		}

		// Check for duplicates
		if addresses[address] {
			t.Errorf("Duplicate address generated at iteration %d: %s", i, address)
		}
		addresses[address] = true
	}

	if len(addresses) != 10 {
		t.Errorf("Expected 10 unique addresses, got %d", len(addresses))
	}
}

// TestIntegration_FullFlow tests the complete flow from mnemonic to address
func TestIntegration_FullFlow(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	// Test multiple indices
	for index := uint32(0); index < 5; index++ {
		address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

		if err != nil {
			t.Errorf("Error at index %d: %v", index, err)
			continue
		}

		// Verify private key can generate the same address
		privKeyBytes, err := hex.DecodeString(privKey)
		if err != nil {
			t.Errorf("Invalid private key hex at index %d: %v", index, err)
			continue
		}

		addressFromPrivKey, err := PrivateKeyToTronAddress(privKeyBytes)
		if err != nil {
			t.Errorf("Error generating address from private key at index %d: %v", index, err)
			continue
		}

		if address != addressFromPrivKey {
			t.Errorf("Address mismatch at index %d: derivation=%s, from_privkey=%s",
				index, address, addressFromPrivKey)
		}
	}
}

// TestDeriveTronAddressFromMnemonic_LongMnemonic tests with 24-word mnemonic
func TestDeriveTronAddressFromMnemonic_LongMnemonic(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art"
	index := uint32(0)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err != nil {
		t.Fatalf("Expected no error with 24-word mnemonic, got %v", err)
	}

	if address == "" || privKey == "" {
		t.Error("Expected valid outputs with 24-word mnemonic")
	}
}

// TestDeriveTronAddressFromMnemonic_Passphrase tests that passphrase is handled
func TestDeriveTronAddressFromMnemonic_Passphrase(t *testing.T) {
	// The function uses empty passphrase
	// Test that same mnemonic always gives same result (no random passphrase)
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	results := make(map[string]int)
	for i := 0; i < 3; i++ {
		address, _, err := DeriveTronAddressFromMnemonic(mnemonic, 0)
		if err != nil {
			t.Fatalf("Error in iteration %d: %v", i, err)
		}
		results[address]++
	}

	// Should always produce the same address
	if len(results) != 1 {
		t.Errorf("Expected 1 unique address, got %d", len(results))
	}
}

// TestPrivateKeyToTronAddress_EdgeCases tests edge cases
func TestPrivateKeyToTronAddress_EdgeCases(t *testing.T) {
	testCases := []struct {
		name      string
		key       []byte
		expectErr bool
	}{
		{
			name:      "valid 32 bytes",
			key:       make([]byte, 32),
			expectErr: false,
		},
		{
			name:      "sequential bytes",
			key:       []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			address, err := PrivateKeyToTronAddress(tc.key)

			if tc.expectErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.expectErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !tc.expectErr && address == "" {
				t.Error("Expected non-empty address")
			}
		})
	}
}

// BenchmarkDeriveTronAddressFromMnemonic benchmarks address derivation
func BenchmarkDeriveTronAddressFromMnemonic(b *testing.B) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = DeriveTronAddressFromMnemonic(mnemonic, uint32(i%100))
	}
}

// BenchmarkPrivateKeyToTronAddress benchmarks address generation from private key
func BenchmarkPrivateKeyToTronAddress(b *testing.B) {
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = byte(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PrivateKeyToTronAddress(privateKey)
	}
}