package main

import (
	"encoding/hex"
	"strings"
	"testing"
)

// Test DeriveTronAddressFromMnemonic function

func TestDeriveTronAddressFromMnemonic_ValidMnemonic(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if address == "" {
		t.Error("Expected non-empty address")
	}

	if privKey == "" {
		t.Error("Expected non-empty private key")
	}

	// TRON addresses start with 'T'
	if !strings.HasPrefix(address, "T") {
		t.Errorf("Expected address to start with 'T', got %s", address)
	}

	// Private key should be 64 hex characters (32 bytes)
	if len(privKey) != 64 {
		t.Errorf("Expected private key length 64, got %d", len(privKey))
	}

	// Verify private key is valid hex
	_, err = hex.DecodeString(privKey)
	if err != nil {
		t.Errorf("Private key is not valid hex: %v", err)
	}
}

func TestDeriveTronAddressFromMnemonic_DifferentIndices(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	addresses := make(map[string]bool)
	privateKeys := make(map[string]bool)

	// Test indices 0-9
	for i := uint32(0); i < 10; i++ {
		address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, i)

		if err != nil {
			t.Errorf("Expected no error for index %d, got %v", i, err)
		}

		// Each index should produce unique address and private key
		if addresses[address] {
			t.Errorf("Duplicate address generated for index %d: %s", i, address)
		}
		addresses[address] = true

		if privateKeys[privKey] {
			t.Errorf("Duplicate private key generated for index %d", i)
		}
		privateKeys[privKey] = true

		// Verify format
		if !strings.HasPrefix(address, "T") {
			t.Errorf("Address at index %d doesn't start with 'T': %s", i, address)
		}

		if len(privKey) != 64 {
			t.Errorf("Private key at index %d has incorrect length: %d", i, len(privKey))
		}
	}
}

func TestDeriveTronAddressFromMnemonic_ConsistentResults(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(5)

	// Call multiple times with same inputs
	address1, privKey1, err1 := DeriveTronAddressFromMnemonic(mnemonic, index)
	address2, privKey2, err2 := DeriveTronAddressFromMnemonic(mnemonic, index)
	address3, privKey3, err3 := DeriveTronAddressFromMnemonic(mnemonic, index)

	// Should produce identical results
	if err1 != nil || err2 != nil || err3 != nil {
		t.Error("Expected no errors")
	}

	if address1 != address2 || address1 != address3 {
		t.Errorf("Addresses not consistent: %s, %s, %s", address1, address2, address3)
	}

	if privKey1 != privKey2 || privKey1 != privKey3 {
		t.Error("Private keys not consistent")
	}
}

func TestDeriveTronAddressFromMnemonic_EmptyMnemonic(t *testing.T) {
	mnemonic := ""
	index := uint32(0)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	// Should still work (will derive from empty seed)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if address == "" {
		t.Error("Expected non-empty address even with empty mnemonic")
	}

	if privKey == "" {
		t.Error("Expected non-empty private key even with empty mnemonic")
	}
}

func TestDeriveTronAddressFromMnemonic_DifferentMnemonics(t *testing.T) {
	mnemonics := []string{
		"flash couple heart script ramp april average caution plunge alter elite author",
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		"legal winner thank year wave sausage worth useful legal winner thank yellow",
	}

	index := uint32(0)
	addresses := make(map[string]bool)
	privateKeys := make(map[string]bool)

	for _, mnemonic := range mnemonics {
		address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

		if err != nil {
			t.Errorf("Expected no error for mnemonic, got %v", err)
		}

		// Each mnemonic should produce unique address and key
		if addresses[address] {
			t.Errorf("Duplicate address for different mnemonic: %s", address)
		}
		addresses[address] = true

		if privateKeys[privKey] {
			t.Error("Duplicate private key for different mnemonic")
		}
		privateKeys[privKey] = true
	}
}

func TestDeriveTronAddressFromMnemonic_LargeIndex(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	
	testCases := []uint32{
		100,
		1000,
		10000,
		2147483647, // Max uint32/2
	}

	for _, index := range testCases {
		address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

		if err != nil {
			t.Errorf("Expected no error for index %d, got %v", index, err)
		}

		if address == "" || privKey == "" {
			t.Errorf("Empty result for index %d", index)
		}

		if !strings.HasPrefix(address, "T") {
			t.Errorf("Invalid address format for index %d", index)
		}
	}
}

func TestDeriveTronAddressFromMnemonic_WithPassphrase(t *testing.T) {
	// Current implementation uses empty passphrase
	// This test documents expected behavior
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Results should be consistent
	address2, privKey2, _ := DeriveTronAddressFromMnemonic(mnemonic, index)

	if address != address2 || privKey != privKey2 {
		t.Error("Results not consistent with empty passphrase")
	}
}

// Test PrivateKeyToTronAddress function

func TestPrivateKeyToTronAddress_ValidPrivateKey(t *testing.T) {
	// Generate a valid 32-byte private key
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = byte(i)
	}

	address, err := PrivateKeyToTronAddress(privateKey)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if address == "" {
		t.Error("Expected non-empty address")
	}

	if !strings.HasPrefix(address, "T") {
		t.Errorf("Expected address to start with 'T', got %s", address)
	}

	// TRON addresses in base58 are typically 34 characters
	if len(address) < 30 || len(address) > 40 {
		t.Errorf("Unexpected address length: %d", len(address))
	}
}

func TestPrivateKeyToTronAddress_DifferentPrivateKeys(t *testing.T) {
	addresses := make(map[string]bool)

	// Test with different private keys
	for i := 0; i < 10; i++ {
		privateKey := make([]byte, 32)
		for j := range privateKey {
			privateKey[j] = byte((i * 10) + j)
		}

		address, err := PrivateKeyToTronAddress(privateKey)

		if err != nil {
			t.Errorf("Expected no error for iteration %d, got %v", i, err)
		}

		if addresses[address] {
			t.Errorf("Duplicate address generated: %s", address)
		}
		addresses[address] = true

		if !strings.HasPrefix(address, "T") {
			t.Errorf("Invalid address format at iteration %d", i)
		}
	}
}

func TestPrivateKeyToTronAddress_ConsistentResults(t *testing.T) {
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = byte(i * 7)
	}

	// Call multiple times with same input
	address1, err1 := PrivateKeyToTronAddress(privateKey)
	address2, err2 := PrivateKeyToTronAddress(privateKey)
	address3, err3 := PrivateKeyToTronAddress(privateKey)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Error("Expected no errors")
	}

	if address1 != address2 || address1 != address3 {
		t.Errorf("Addresses not consistent: %s, %s, %s", address1, address2, address3)
	}
}

func TestPrivateKeyToTronAddress_ZeroPrivateKey(t *testing.T) {
	privateKey := make([]byte, 32)

	address, err := PrivateKeyToTronAddress(privateKey)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if address == "" {
		t.Error("Expected non-empty address")
	}

	if !strings.HasPrefix(address, "T") {
		t.Errorf("Expected address to start with 'T', got %s", address)
	}
}

func TestPrivateKeyToTronAddress_MaxPrivateKey(t *testing.T) {
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = 0xFF
	}

	address, err := PrivateKeyToTronAddress(privateKey)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if address == "" {
		t.Error("Expected non-empty address")
	}

	if !strings.HasPrefix(address, "T") {
		t.Errorf("Expected address to start with 'T', got %s", address)
	}
}

func TestPrivateKeyToTronAddress_IncorrectLength(t *testing.T) {
	testCases := []int{0, 16, 31, 33, 64}

	for _, length := range testCases {
		privateKey := make([]byte, length)

		address, err := PrivateKeyToTronAddress(privateKey)

		// Implementation may or may not handle incorrect lengths
		// Document the actual behavior
		if length == 32 {
			if err != nil {
				t.Errorf("Expected no error for length 32, got %v", err)
			}
		} else {
			// For other lengths, behavior depends on implementation
			// Just verify it doesn't panic
			_ = address
			_ = err
		}
	}
}

func TestPrivateKeyToTronAddress_SequentialKeys(t *testing.T) {
	addresses := make(map[string]bool)

	// Test with sequential private keys
	for i := 0; i < 5; i++ {
		privateKey := make([]byte, 32)
		privateKey[31] = byte(i) // Only vary the last byte

		address, err := PrivateKeyToTronAddress(privateKey)

		if err != nil {
			t.Errorf("Expected no error for sequence %d, got %v", i, err)
		}

		if addresses[address] {
			t.Errorf("Duplicate address in sequence: %s", address)
		}
		addresses[address] = true
	}
}

// Integration tests

func TestIntegration_MnemonicToAddress(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0)

	// Derive address using the high-level function
	address, privKeyHex, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err != nil {
		t.Fatalf("Failed to derive from mnemonic: %v", err)
	}

	// Decode the private key
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		t.Fatalf("Failed to decode private key: %v", err)
	}

	// Derive address from the private key directly
	address2, err := PrivateKeyToTronAddress(privKeyBytes)

	if err != nil {
		t.Fatalf("Failed to derive from private key: %v", err)
	}

	// Both methods should produce the same address
	if address != address2 {
		t.Errorf("Addresses don't match: %s vs %s", address, address2)
	}
}

func TestIntegration_MultipleDerivations(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	results := make(map[uint32]struct {
		address string
		privKey string
	})

	// Derive addresses for indices 0-4
	for i := uint32(0); i < 5; i++ {
		address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, i)

		if err != nil {
			t.Fatalf("Failed to derive for index %d: %v", i, err)
		}

		results[i] = struct {
			address string
			privKey string
		}{address, privKey}
	}

	// Verify all addresses are unique
	addresses := make(map[string]uint32)
	for idx, result := range results {
		if existingIdx, exists := addresses[result.address]; exists {
			t.Errorf("Duplicate address between index %d and %d: %s", idx, existingIdx, result.address)
		}
		addresses[result.address] = idx
	}

	// Verify all private keys are unique
	privKeys := make(map[string]uint32)
	for idx, result := range results {
		if existingIdx, exists := privKeys[result.privKey]; exists {
			t.Errorf("Duplicate private key between index %d and %d", idx, existingIdx)
		}
		privKeys[result.privKey] = idx
	}
}

func TestIntegration_AddressValidation(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0)

	address, _, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err != nil {
		t.Fatalf("Failed to derive address: %v", err)
	}

	// Basic validation checks
	if len(address) < 30 {
		t.Errorf("Address too short: %s", address)
	}

	if !strings.HasPrefix(address, "T") {
		t.Errorf("Invalid address prefix: %s", address)
	}

	// TRON addresses should be base58 encoded
	// Check that it only contains valid base58 characters
	validBase58 := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	for _, char := range address {
		if !strings.ContainsRune(validBase58, char) {
			t.Errorf("Invalid character in address: %c", char)
		}
	}
}

// Benchmark tests

func BenchmarkDeriveTronAddressFromMnemonic(b *testing.B) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = DeriveTronAddressFromMnemonic(mnemonic, index)
	}
}

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

func BenchmarkDeriveTronAddressFromMnemonic_DifferentIndices(b *testing.B) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := uint32(i % 1000)
		_, _, _ = DeriveTronAddressFromMnemonic(mnemonic, index)
	}
}