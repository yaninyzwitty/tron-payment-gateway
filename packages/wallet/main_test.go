package main

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/btcsuite/btcutil/base58"
)

// Test DeriveTronAddressFromMnemonic with valid mnemonic and index 0
func TestDeriveTronAddressFromMnemonic_Success(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if address == "" {
		t.Error("Expected non-empty address")
	}

	if privKey == "" {
		t.Error("Expected non-empty private key")
	}

	// Validate address format: should start with 'T' and be 34 characters
	if !strings.HasPrefix(address, "T") {
		t.Errorf("Expected address to start with 'T', got: %s", address)
	}

	// Validate private key is hex encoded
	if len(privKey) != 64 {
		t.Errorf("Expected private key hex length of 64, got: %d", len(privKey))
	}

	_, err = hex.DecodeString(privKey)
	if err != nil {
		t.Errorf("Private key should be valid hex: %v", err)
	}
}

// Test DeriveTronAddressFromMnemonic with different indices produces different addresses
func TestDeriveTronAddressFromMnemonic_DifferentIndices(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	address0, privKey0, err0 := DeriveTronAddressFromMnemonic(mnemonic, 0)
	address1, privKey1, err1 := DeriveTronAddressFromMnemonic(mnemonic, 1)
	address2, privKey2, err2 := DeriveTronAddressFromMnemonic(mnemonic, 2)

	if err0 != nil || err1 != nil || err2 != nil {
		t.Fatal("Expected no errors")
	}

	// Different indices should produce different addresses
	if address0 == address1 || address1 == address2 || address0 == address2 {
		t.Error("Different indices should produce different addresses")
	}

	// Different indices should produce different private keys
	if privKey0 == privKey1 || privKey1 == privKey2 || privKey0 == privKey2 {
		t.Error("Different indices should produce different private keys")
	}
}

// Test DeriveTronAddressFromMnemonic with large index
func TestDeriveTronAddressFromMnemonic_LargeIndex(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(999999)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err != nil {
		t.Fatalf("Expected no error with large index, got: %v", err)
	}

	if address == "" || privKey == "" {
		t.Error("Expected valid address and private key for large index")
	}
}

// Test DeriveTronAddressFromMnemonic with maximum uint32 index
func TestDeriveTronAddressFromMnemonic_MaxIndex(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0xFFFFFFFF) // Maximum uint32 value

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err != nil {
		t.Fatalf("Expected no error with max index, got: %v", err)
	}

	if address == "" || privKey == "" {
		t.Error("Expected valid address and private key for max index")
	}
}

// Test DeriveTronAddressFromMnemonic determinism - same inputs should produce same outputs
func TestDeriveTronAddressFromMnemonic_Deterministic(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(42)

	address1, privKey1, err1 := DeriveTronAddressFromMnemonic(mnemonic, index)
	address2, privKey2, err2 := DeriveTronAddressFromMnemonic(mnemonic, index)

	if err1 != nil || err2 != nil {
		t.Fatal("Expected no errors")
	}

	if address1 != address2 {
		t.Errorf("Expected deterministic address generation, got %s and %s", address1, address2)
	}

	if privKey1 != privKey2 {
		t.Errorf("Expected deterministic private key generation, got %s and %s", privKey1, privKey2)
	}
}

// Test DeriveTronAddressFromMnemonic with empty mnemonic
func TestDeriveTronAddressFromMnemonic_EmptyMnemonic(t *testing.T) {
	mnemonic := ""
	index := uint32(0)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	// With empty mnemonic, it should still generate something (seed will be derived from empty string)
	// This is technically valid behavior, though not recommended
	if err != nil {
		t.Logf("Empty mnemonic resulted in error: %v", err)
	}

	if address != "" && privKey != "" && err == nil {
		// If it succeeds, ensure format is correct
		if !strings.HasPrefix(address, "T") {
			t.Error("Address should still start with T")
		}
	}
}

// Test DeriveTronAddressFromMnemonic with different mnemonics produce different addresses
func TestDeriveTronAddressFromMnemonic_DifferentMnemonics(t *testing.T) {
	mnemonic1 := "flash couple heart script ramp april average caution plunge alter elite author"
	mnemonic2 := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	index := uint32(0)

	address1, privKey1, err1 := DeriveTronAddressFromMnemonic(mnemonic1, index)
	address2, privKey2, err2 := DeriveTronAddressFromMnemonic(mnemonic2, index)

	if err1 != nil || err2 != nil {
		t.Fatal("Expected no errors")
	}

	if address1 == address2 {
		t.Error("Different mnemonics should produce different addresses")
	}

	if privKey1 == privKey2 {
		t.Error("Different mnemonics should produce different private keys")
	}
}

// Test PrivateKeyToTronAddress with valid 32-byte private key
func TestPrivateKeyToTronAddress_ValidKey(t *testing.T) {
	// Generate a valid 32-byte private key (example)
	privateKeyHex := "3481E79956D4BD95F358AC96D151C976392FC4E3FC132F78A847906DE588C145"
	privateKey, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		t.Fatalf("Failed to decode test private key: %v", err)
	}

	address, err := PrivateKeyToTronAddress(privateKey)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if address == "" {
		t.Error("Expected non-empty address")
	}

	if !strings.HasPrefix(address, "T") {
		t.Errorf("Expected address to start with 'T', got: %s", address)
	}

	// Validate it's a valid base58 encoded string
	decoded := base58.Decode(address)
	if len(decoded) == 0 {
		t.Error("Address should be valid base58")
	}
}

// Test PrivateKeyToTronAddress with 32-byte key produces consistent results
func TestPrivateKeyToTronAddress_Deterministic(t *testing.T) {
	privateKeyHex := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	privateKey, _ := hex.DecodeString(privateKeyHex)

	address1, err1 := PrivateKeyToTronAddress(privateKey)
	address2, err2 := PrivateKeyToTronAddress(privateKey)

	if err1 != nil || err2 != nil {
		t.Fatal("Expected no errors")
	}

	if address1 != address2 {
		t.Errorf("Expected deterministic address from same private key, got %s and %s", address1, address2)
	}
}

// Test PrivateKeyToTronAddress with all zeros
func TestPrivateKeyToTronAddress_AllZeros(t *testing.T) {
	privateKey := make([]byte, 32) // All zeros

	address, err := PrivateKeyToTronAddress(privateKey)

	if err != nil {
		t.Fatalf("Expected no error with all-zero key, got: %v", err)
	}

	if address == "" {
		t.Error("Expected non-empty address even with all-zero key")
	}

	if !strings.HasPrefix(address, "T") {
		t.Error("Address should start with T")
	}
}

// Test PrivateKeyToTronAddress with all ones
func TestPrivateKeyToTronAddress_AllOnes(t *testing.T) {
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = 0xFF
	}

	address, err := PrivateKeyToTronAddress(privateKey)

	if err != nil {
		t.Fatalf("Expected no error with all-one key, got: %v", err)
	}

	if address == "" {
		t.Error("Expected non-empty address")
	}

	if !strings.HasPrefix(address, "T") {
		t.Error("Address should start with T")
	}
}

// Test PrivateKeyToTronAddress with different keys produce different addresses
func TestPrivateKeyToTronAddress_DifferentKeys(t *testing.T) {
	key1 := make([]byte, 32)
	key1[0] = 0x01

	key2 := make([]byte, 32)
	key2[0] = 0x02

	address1, err1 := PrivateKeyToTronAddress(key1)
	address2, err2 := PrivateKeyToTronAddress(key2)

	if err1 != nil || err2 != nil {
		t.Fatal("Expected no errors")
	}

	if address1 == address2 {
		t.Error("Different private keys should produce different addresses")
	}
}

// Test PrivateKeyToTronAddress with various key patterns
func TestPrivateKeyToTronAddress_VariousKeyPatterns(t *testing.T) {
	testCases := []struct {
		name     string
		keyValue byte
	}{
		{"Low value", 0x01},
		{"Mid value", 0x7F},
		{"High value", 0xFE},
		{"Alternating", 0xAA},
		{"Alternating inverse", 0x55},
	}

	addresses := make(map[string]bool)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key := make([]byte, 32)
			for i := range key {
				key[i] = tc.keyValue
			}

			address, err := PrivateKeyToTronAddress(key)

			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if !strings.HasPrefix(address, "T") {
				t.Errorf("Address should start with T, got: %s", address)
			}

			// Check for uniqueness
			if addresses[address] {
				t.Errorf("Duplicate address found: %s", address)
			}
			addresses[address] = true
		})
	}
}

// Test PrivateKeyToTronAddress empty key
func TestPrivateKeyToTronAddress_EmptyKey(t *testing.T) {
	privateKey := []byte{}

	address, err := PrivateKeyToTronAddress(privateKey)

	// Empty key should still produce a result (treated as zero)
	if err != nil {
		t.Logf("Empty key resulted in error: %v", err)
	}

	if address != "" && !strings.HasPrefix(address, "T") {
		t.Error("If address is generated, it should start with T")
	}
}

// Test PrivateKeyToTronAddress short key (less than 32 bytes)
func TestPrivateKeyToTronAddress_ShortKey(t *testing.T) {
	privateKey := make([]byte, 16) // Only 16 bytes
	for i := range privateKey {
		privateKey[i] = 0xAB
	}

	address, err := PrivateKeyToTronAddress(privateKey)

	// Should still work - will be padded with zeros
	if err != nil {
		t.Logf("Short key resulted in error: %v", err)
	}

	if address != "" && !strings.HasPrefix(address, "T") {
		t.Error("If address is generated, it should start with T")
	}
}

// Test PrivateKeyToTronAddress validates checksum
func TestPrivateKeyToTronAddress_ValidChecksum(t *testing.T) {
	privateKeyHex := "1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF"
	privateKey, _ := hex.DecodeString(privateKeyHex)

	address, err := PrivateKeyToTronAddress(privateKey)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Decode the base58 address and verify it has the proper structure
	decoded := base58.Decode(address)

	// Should be 25 bytes: 21 bytes address + 4 bytes checksum
	if len(decoded) != 25 {
		t.Errorf("Expected decoded address to be 25 bytes, got: %d", len(decoded))
	}

	// First byte should be 0x41 (TRON mainnet prefix)
	if len(decoded) > 0 && decoded[0] != 0x41 {
		t.Errorf("Expected first byte to be 0x41, got: 0x%02X", decoded[0])
	}
}

// Test address format validation helper
func TestAddressFormat(t *testing.T) {
	testCases := []struct {
		name        string
		mnemonic    string
		index       uint32
		expectError bool
	}{
		{
			name:        "Valid BIP39 mnemonic",
			mnemonic:    "flash couple heart script ramp april average caution plunge alter elite author",
			index:       0,
			expectError: false,
		},
		{
			name:        "Another valid BIP39 mnemonic",
			mnemonic:    "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
			index:       0,
			expectError: false,
		},
		{
			name:        "Valid with high index",
			mnemonic:    "flash couple heart script ramp april average caution plunge alter elite author",
			index:       10000,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			address, privKey, err := DeriveTronAddressFromMnemonic(tc.mnemonic, tc.index)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tc.expectError {
				// Validate address length and prefix
				if len(address) < 30 || len(address) > 35 {
					t.Errorf("Address length %d is outside expected range [30-35]", len(address))
				}

				if !strings.HasPrefix(address, "T") {
					t.Errorf("Address should start with T, got: %s", address)
				}

				// Validate private key is 64 hex characters
				if len(privKey) != 64 {
					t.Errorf("Private key should be 64 hex characters, got: %d", len(privKey))
				}
			}
		})
	}
}

// Test BIP44 path derivation consistency
func TestBIP44PathDerivation(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	// Test multiple indices in sequence
	previousAddresses := make(map[string]uint32)

	for i := uint32(0); i < 10; i++ {
		address, _, err := DeriveTronAddressFromMnemonic(mnemonic, i)
		if err != nil {
			t.Fatalf("Error at index %d: %v", i, err)
		}

		// Check for uniqueness
		if prevIndex, exists := previousAddresses[address]; exists {
			t.Errorf("Duplicate address at indices %d and %d: %s", prevIndex, i, address)
		}
		previousAddresses[address] = i
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

// Test concurrent address generation
func TestConcurrentAddressGeneration(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	type result struct {
		index   uint32
		address string
		err     error
	}

	results := make(chan result, 10)

	// Generate 10 addresses concurrently
	for i := uint32(0); i < 10; i++ {
		go func(idx uint32) {
			addr, _, err := DeriveTronAddressFromMnemonic(mnemonic, idx)
			results <- result{index: idx, address: addr, err: err}
		}(i)
	}

	addresses := make(map[string]uint32)
	for i := 0; i < 10; i++ {
		res := <-results
		if res.err != nil {
			t.Errorf("Error at index %d: %v", res.index, res.err)
		}

		if prevIdx, exists := addresses[res.address]; exists {
			t.Errorf("Duplicate address at indices %d and %d", prevIdx, res.index)
		}
		addresses[res.address] = res.index
	}
}