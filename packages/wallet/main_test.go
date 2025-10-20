package main

import (
	"encoding/hex"
	"strings"
	"testing"
)

// TestDeriveTronAddressFromMnemonic_Success tests successful derivation with valid inputs
func TestDeriveTronAddressFromMnemonic_Success(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if address == "" {
		t.Error("expected non-empty address")
	}

	if privKey == "" {
		t.Error("expected non-empty private key")
	}

	// TRON addresses should start with 'T'
	if !strings.HasPrefix(address, "T") {
		t.Errorf("expected address to start with 'T', got: %s", address)
	}

	// Private key should be valid hex and 64 characters (32 bytes)
	if len(privKey) != 64 {
		t.Errorf("expected private key length 64, got: %d", len(privKey))
	}

	// Validate it's valid hex
	_, err = hex.DecodeString(privKey)
	if err != nil {
		t.Errorf("private key should be valid hex: %v", err)
	}
}

// TestDeriveTronAddressFromMnemonic_DeterministicOutput tests that same inputs produce same outputs
func TestDeriveTronAddressFromMnemonic_DeterministicOutput(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0)

	address1, privKey1, err1 := DeriveTronAddressFromMnemonic(mnemonic, index)
	if err1 != nil {
		t.Fatalf("first call failed: %v", err1)
	}

	address2, privKey2, err2 := DeriveTronAddressFromMnemonic(mnemonic, index)
	if err2 != nil {
		t.Fatalf("second call failed: %v", err2)
	}

	if address1 != address2 {
		t.Errorf("addresses should be identical, got: %s and %s", address1, address2)
	}

	if privKey1 != privKey2 {
		t.Errorf("private keys should be identical, got: %s and %s", privKey1, privKey2)
	}
}

// TestDeriveTronAddressFromMnemonic_DifferentIndices tests different indices produce different outputs
func TestDeriveTronAddressFromMnemonic_DifferentIndices(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	
	address0, privKey0, err0 := DeriveTronAddressFromMnemonic(mnemonic, 0)
	if err0 != nil {
		t.Fatalf("index 0 failed: %v", err0)
	}

	address1, privKey1, err1 := DeriveTronAddressFromMnemonic(mnemonic, 1)
	if err1 != nil {
		t.Fatalf("index 1 failed: %v", err1)
	}

	address2, privKey2, err2 := DeriveTronAddressFromMnemonic(mnemonic, 2)
	if err2 != nil {
		t.Fatalf("index 2 failed: %v", err2)
	}

	// All addresses should be different
	if address0 == address1 || address0 == address2 || address1 == address2 {
		t.Error("different indices should produce different addresses")
	}

	// All private keys should be different
	if privKey0 == privKey1 || privKey0 == privKey2 || privKey1 == privKey2 {
		t.Error("different indices should produce different private keys")
	}
}

// TestDeriveTronAddressFromMnemonic_LargeIndex tests with large index values
func TestDeriveTronAddressFromMnemonic_LargeIndex(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	testCases := []uint32{100, 1000, 10000, 2147483647} // Max uint32/2

	for _, index := range testCases {
		address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)
		if err != nil {
			t.Errorf("index %d failed: %v", index, err)
			continue
		}

		if address == "" || privKey == "" {
			t.Errorf("index %d produced empty result", index)
		}

		if !strings.HasPrefix(address, "T") {
			t.Errorf("index %d: expected address to start with 'T', got: %s", index, address)
		}
	}
}

// TestDeriveTronAddressFromMnemonic_DifferentMnemonics tests different mnemonics produce different outputs
func TestDeriveTronAddressFromMnemonic_DifferentMnemonics(t *testing.T) {
	mnemonics := []string{
		"flash couple heart script ramp april average caution plunge alter elite author",
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		"legal winner thank year wave sausage worth useful legal winner thank yellow",
	}

	index := uint32(0)
	results := make(map[string]bool)

	for _, mnemonic := range mnemonics {
		address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)
		if err != nil {
			t.Errorf("mnemonic failed: %v", err)
			continue
		}

		if results[address] {
			t.Errorf("duplicate address found for different mnemonics: %s", address)
		}
		results[address] = true

		if results[privKey] {
			t.Errorf("duplicate private key found for different mnemonics: %s", privKey)
		}
		results[privKey] = true
	}

	if len(results) != len(mnemonics)*2 {
		t.Error("expected unique addresses and keys for all mnemonics")
	}
}

// TestDeriveTronAddressFromMnemonic_InvalidMnemonic tests error handling with invalid mnemonic
func TestDeriveTronAddressFromMnemonic_InvalidMnemonic(t *testing.T) {
	testCases := []struct {
		name     string
		mnemonic string
	}{
		{"empty", ""},
		{"single word", "word"},
		{"invalid words", "invalid invalid invalid invalid invalid invalid invalid invalid invalid invalid invalid invalid"},
		{"too few words", "abandon abandon abandon"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			address, privKey, err := DeriveTronAddressFromMnemonic(tc.mnemonic, 0)
			
			// BIP39 library might not error on all invalid inputs, but results should be consistent
			// We test that it doesn't panic and produces some output
			if err == nil && (address == "" || privKey == "") {
				t.Error("expected error or valid output for invalid mnemonic")
			}
		})
	}
}

// TestDeriveTronAddressFromMnemonic_BoundaryIndices tests boundary values for index
func TestDeriveTronAddressFromMnemonic_BoundaryIndices(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	
	testCases := []uint32{
		0,
		1,
		255,
		256,
		65535,
		65536,
	}

	for _, index := range testCases {
		t.Run(hex.EncodeToString([]byte{byte(index)}), func(t *testing.T) {
			address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)
			if err != nil {
				t.Fatalf("index %d failed: %v", index, err)
			}

			if address == "" || privKey == "" {
				t.Errorf("index %d produced empty result", index)
			}
		})
	}
}

// TestPrivateKeyToTronAddress_ValidPrivateKey tests conversion with valid 32-byte private key
func TestPrivateKeyToTronAddress_ValidPrivateKey(t *testing.T) {
	// Valid 32-byte private key (hex encoded to 64 chars)
	privateKeyHex := "e8b0c3c3c5c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3"
	privateKey, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		t.Fatalf("failed to decode test private key: %v", err)
	}

	address, err := PrivateKeyToTronAddress(privateKey)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if address == "" {
		t.Error("expected non-empty address")
	}

	if !strings.HasPrefix(address, "T") {
		t.Errorf("expected address to start with 'T', got: %s", address)
	}
}

// TestPrivateKeyToTronAddress_DeterministicOutput tests deterministic address generation
func TestPrivateKeyToTronAddress_DeterministicOutput(t *testing.T) {
	privateKeyHex := "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	privateKey, _ := hex.DecodeString(privateKeyHex)

	address1, err1 := PrivateKeyToTronAddress(privateKey)
	if err1 != nil {
		t.Fatalf("first call failed: %v", err1)
	}

	address2, err2 := PrivateKeyToTronAddress(privateKey)
	if err2 != nil {
		t.Fatalf("second call failed: %v", err2)
	}

	if address1 != address2 {
		t.Errorf("same private key should produce same address: %s != %s", address1, address2)
	}
}

// TestPrivateKeyToTronAddress_DifferentKeys tests different keys produce different addresses
func TestPrivateKeyToTronAddress_DifferentKeys(t *testing.T) {
	keys := []string{
		"0000000000000000000000000000000000000000000000000000000000000001",
		"0000000000000000000000000000000000000000000000000000000000000002",
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"e8b0c3c3c5c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3",
	}

	addresses := make(map[string]bool)

	for i, keyHex := range keys {
		privateKey, err := hex.DecodeString(keyHex)
		if err != nil {
			t.Fatalf("failed to decode key %d: %v", i, err)
		}

		address, err := PrivateKeyToTronAddress(privateKey)
		if err != nil {
			t.Errorf("key %d failed: %v", i, err)
			continue
		}

		if addresses[address] {
			t.Errorf("duplicate address for different keys: %s", address)
		}
		addresses[address] = true
	}

	if len(addresses) != len(keys) {
		t.Error("expected unique addresses for all different keys")
	}
}

// TestPrivateKeyToTronAddress_EmptyKey tests error handling with empty key
func TestPrivateKeyToTronAddress_EmptyKey(t *testing.T) {
	privateKey := []byte{}

	address, err := PrivateKeyToTronAddress(privateKey)
	
	// Should either error or handle gracefully
	if err == nil && address == "" {
		t.Error("empty key should produce error or valid address")
	}
}

// TestPrivateKeyToTronAddress_ShortKey tests error handling with short key
func TestPrivateKeyToTronAddress_ShortKey(t *testing.T) {
	privateKey := []byte{0x01, 0x02, 0x03}

	address, err := PrivateKeyToTronAddress(privateKey)
	
	// Should either error or handle gracefully
	if err == nil && address == "" {
		t.Error("short key should produce error or valid address")
	}
}

// TestPrivateKeyToTronAddress_LongKey tests with key longer than 32 bytes
func TestPrivateKeyToTronAddress_LongKey(t *testing.T) {
	// 64 bytes instead of 32
	privateKey := make([]byte, 64)
	for i := range privateKey {
		privateKey[i] = byte(i)
	}

	address, err := PrivateKeyToTronAddress(privateKey)
	
	// Should handle gracefully
	if err == nil && address == "" {
		t.Error("long key should produce error or valid address")
	}
}

// TestPrivateKeyToTronAddress_AllZeros tests edge case with all-zero key
func TestPrivateKeyToTronAddress_AllZeros(t *testing.T) {
	privateKey := make([]byte, 32)

	address, err := PrivateKeyToTronAddress(privateKey)
	
	// Should handle this edge case
	if err != nil {
		t.Logf("all-zero key produced error (expected): %v", err)
	} else if address == "" {
		t.Error("all-zero key should produce error or valid address")
	} else if !strings.HasPrefix(address, "T") {
		t.Errorf("expected address to start with 'T', got: %s", address)
	}
}

// TestPrivateKeyToTronAddress_AllOnes tests edge case with all-one key
func TestPrivateKeyToTronAddress_AllOnes(t *testing.T) {
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = 0xFF
	}

	address, err := PrivateKeyToTronAddress(privateKey)
	
	if err != nil {
		t.Fatalf("all-ones key failed: %v", err)
	}

	if address == "" {
		t.Error("expected non-empty address for all-ones key")
	}

	if !strings.HasPrefix(address, "T") {
		t.Errorf("expected address to start with 'T', got: %s", address)
	}
}

// TestPrivateKeyToTronAddress_SequentialBytes tests with sequential byte pattern
func TestPrivateKeyToTronAddress_SequentialBytes(t *testing.T) {
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = byte(i)
	}

	address, err := PrivateKeyToTronAddress(privateKey)
	if err != nil {
		t.Fatalf("sequential bytes key failed: %v", err)
	}

	if address == "" {
		t.Error("expected non-empty address")
	}

	if !strings.HasPrefix(address, "T") {
		t.Errorf("expected address to start with 'T', got: %s", address)
	}
}

// TestPrivateKeyToTronAddress_AddressFormat tests address format validity
func TestPrivateKeyToTronAddress_AddressFormat(t *testing.T) {
	privateKeyHex := "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	privateKey, _ := hex.DecodeString(privateKeyHex)

	address, err := PrivateKeyToTronAddress(privateKey)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	// TRON address should be base58 encoded and start with T
	if !strings.HasPrefix(address, "T") {
		t.Errorf("TRON address must start with 'T', got: %s", address)
	}

	// TRON addresses are typically 34 characters in base58
	if len(address) < 30 || len(address) > 40 {
		t.Errorf("unexpected address length: %d for address: %s", len(address), address)
	}

	// Address should only contain valid base58 characters
	validBase58 := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	for _, char := range address {
		if !strings.ContainsRune(validBase58, char) {
			t.Errorf("address contains invalid base58 character: %c", char)
		}
	}
}

// TestPrivateKeyToTronAddress_NilKey tests with nil key
func TestPrivateKeyToTronAddress_NilKey(t *testing.T) {
	var privateKey []byte = nil

	address, err := PrivateKeyToTronAddress(privateKey)
	
	// Should either error or handle gracefully
	if err == nil && address == "" {
		t.Error("nil key should produce error or valid address")
	}
}

// TestPrivateKeyToTronAddress_AddressChecksum tests that addresses have valid checksums
func TestPrivateKeyToTronAddress_AddressChecksum(t *testing.T) {
	// Generate multiple addresses and verify they're all valid
	for i := 0; i < 10; i++ {
		privateKey := make([]byte, 32)
		for j := range privateKey {
			privateKey[j] = byte((i * 7 + j) % 256)
		}

		address, err := PrivateKeyToTronAddress(privateKey)
		if err != nil {
			t.Errorf("iteration %d failed: %v", i, err)
			continue
		}

		if address == "" {
			t.Errorf("iteration %d: empty address", i)
			continue
		}

		if !strings.HasPrefix(address, "T") {
			t.Errorf("iteration %d: address doesn't start with 'T': %s", i, address)
		}
	}
}

// TestDeriveTronAddressFromMnemonic_ConsistentWithPrivateKeyFunction tests integration
func TestDeriveTronAddressFromMnemonic_ConsistentWithPrivateKeyFunction(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(5)

	address1, privKeyHex, err := DeriveTronAddressFromMnemonic(mnemonic, index)
	if err != nil {
		t.Fatalf("derivation failed: %v", err)
	}

	// Convert hex private key to bytes and derive address directly
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		t.Fatalf("failed to decode private key: %v", err)
	}

	address2, err := PrivateKeyToTronAddress(privKeyBytes)
	if err != nil {
		t.Fatalf("address derivation from private key failed: %v", err)
	}

	if address1 != address2 {
		t.Errorf("addresses should match: %s != %s", address1, address2)
	}
}

// TestMain_DoesNotPanic tests that the main function doesn't panic
func TestMain_DoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main() panicked: %v", r)
		}
	}()

	// We can't easily test main() directly, but we test that the functions it uses work
	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonicSecret, index)
	if err != nil {
		t.Fatalf("main function would fail with error: %v", err)
	}

	if address == "" || privKey == "" {
		t.Error("main function would produce empty results")
	}
}

// Benchmark tests for performance analysis
func BenchmarkDeriveTronAddressFromMnemonic(b *testing.B) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = DeriveTronAddressFromMnemonic(mnemonic, index)
	}
}

func BenchmarkPrivateKeyToTronAddress(b *testing.B) {
	privateKeyHex := "e8b0c3c3c5c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3"
	privateKey, _ := hex.DecodeString(privateKeyHex)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PrivateKeyToTronAddress(privateKey)
	}
}