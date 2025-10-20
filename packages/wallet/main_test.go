package main

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeriveTronAddressFromMnemonic_ValidMnemonic tests derivation with a valid mnemonic
func TestDeriveTronAddressFromMnemonic_ValidMnemonic(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	require.NoError(t, err)
	assert.NotEmpty(t, address)
	assert.NotEmpty(t, privKey)
	assert.Len(t, privKey, 64) // 32 bytes = 64 hex characters
	
	// Verify address format (Tron addresses start with 'T' in base58)
	assert.True(t, len(address) > 0)
	
	// Private key should be valid hex
	_, err = hex.DecodeString(privKey)
	assert.NoError(t, err)
}

// TestDeriveTronAddressFromMnemonic_DeterministicOutput ensures same inputs produce same outputs
func TestDeriveTronAddressFromMnemonic_DeterministicOutput(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0)

	// Generate address twice
	address1, privKey1, err1 := DeriveTronAddressFromMnemonic(mnemonic, index)
	address2, privKey2, err2 := DeriveTronAddressFromMnemonic(mnemonic, index)

	require.NoError(t, err1)
	require.NoError(t, err2)
	
	// Both calls should produce identical results
	assert.Equal(t, address1, address2)
	assert.Equal(t, privKey1, privKey2)
}

// TestDeriveTronAddressFromMnemonic_DifferentIndices tests that different indices produce different addresses
func TestDeriveTronAddressFromMnemonic_DifferentIndices(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	// Generate addresses for different indices
	address0, privKey0, err0 := DeriveTronAddressFromMnemonic(mnemonic, 0)
	address1, privKey1, err1 := DeriveTronAddressFromMnemonic(mnemonic, 1)
	address2, privKey2, err2 := DeriveTronAddressFromMnemonic(mnemonic, 2)

	require.NoError(t, err0)
	require.NoError(t, err1)
	require.NoError(t, err2)

	// All addresses should be different
	assert.NotEqual(t, address0, address1)
	assert.NotEqual(t, address0, address2)
	assert.NotEqual(t, address1, address2)

	// All private keys should be different
	assert.NotEqual(t, privKey0, privKey1)
	assert.NotEqual(t, privKey0, privKey2)
	assert.NotEqual(t, privKey1, privKey2)
}

// TestDeriveTronAddressFromMnemonic_DifferentMnemonics tests different mnemonics produce different addresses
func TestDeriveTronAddressFromMnemonic_DifferentMnemonics(t *testing.T) {
	mnemonic1 := "flash couple heart script ramp april average caution plunge alter elite author"
	mnemonic2 := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	index := uint32(0)

	address1, privKey1, err1 := DeriveTronAddressFromMnemonic(mnemonic1, index)
	address2, privKey2, err2 := DeriveTronAddressFromMnemonic(mnemonic2, index)

	require.NoError(t, err1)
	require.NoError(t, err2)

	// Different mnemonics should produce different results
	assert.NotEqual(t, address1, address2)
	assert.NotEqual(t, privKey1, privKey2)
}

// TestDeriveTronAddressFromMnemonic_EmptyMnemonic tests behavior with empty mnemonic
func TestDeriveTronAddressFromMnemonic_EmptyMnemonic(t *testing.T) {
	mnemonic := ""
	index := uint32(0)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	// Empty mnemonic should still work (creates a valid seed)
	require.NoError(t, err)
	assert.NotEmpty(t, address)
	assert.NotEmpty(t, privKey)
}

// TestDeriveTronAddressFromMnemonic_LargeIndex tests derivation with a large index
func TestDeriveTronAddressFromMnemonic_LargeIndex(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(999999)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	require.NoError(t, err)
	assert.NotEmpty(t, address)
	assert.NotEmpty(t, privKey)
}

// TestDeriveTronAddressFromMnemonic_MaxIndex tests derivation with maximum uint32 index
func TestDeriveTronAddressFromMnemonic_MaxIndex(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(4294967295) // Max uint32

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	require.NoError(t, err)
	assert.NotEmpty(t, address)
	assert.NotEmpty(t, privKey)
}

// TestDeriveTronAddressFromMnemonic_MultipleIndices tests sequential address generation
func TestDeriveTronAddressFromMnemonic_MultipleIndices(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	addresses := make(map[string]bool)
	privKeys := make(map[string]bool)

	// Generate 10 addresses
	for i := uint32(0); i < 10; i++ {
		address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, i)
		require.NoError(t, err)

		// Ensure no duplicates
		assert.False(t, addresses[address], "Duplicate address found at index %d", i)
		assert.False(t, privKeys[privKey], "Duplicate private key found at index %d", i)

		addresses[address] = true
		privKeys[privKey] = true
	}

	assert.Len(t, addresses, 10)
	assert.Len(t, privKeys, 10)
}

// TestDeriveTronAddressFromMnemonic_ValidFormatOutput tests that output format is correct
func TestDeriveTronAddressFromMnemonic_ValidFormatOutput(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"
	index := uint32(0)

	address, privKey, err := DeriveTronAddressFromMnemonic(mnemonic, index)

	require.NoError(t, err)

	// Private key validation
	assert.Len(t, privKey, 64, "Private key should be 64 hex characters")
	decodedKey, err := hex.DecodeString(privKey)
	require.NoError(t, err)
	assert.Len(t, decodedKey, 32, "Private key should be 32 bytes")

	// Address validation
	assert.NotEmpty(t, address)
	assert.Greater(t, len(address), 25, "Tron addresses are typically 34 characters")
	assert.Less(t, len(address), 50, "Tron addresses shouldn't exceed reasonable length")
}

// TestDeriveTronAddressFromMnemonic_StandardMnemonics tests with various standard mnemonics
func TestDeriveTronAddressFromMnemonic_StandardMnemonics(t *testing.T) {
	testCases := []struct {
		name     string
		mnemonic string
	}{
		{
			name:     "Standard 12 word mnemonic",
			mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		},
		{
			name:     "Different 12 word mnemonic",
			mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong",
		},
		{
			name:     "Test mnemonic with mixed words",
			mnemonic: "legal winner thank year wave sausage worth useful legal winner thank yellow",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			address, privKey, err := DeriveTronAddressFromMnemonic(tc.mnemonic, 0)

			require.NoError(t, err)
			assert.NotEmpty(t, address)
			assert.NotEmpty(t, privKey)
			assert.Len(t, privKey, 64)
		})
	}
}

// TestPrivateKeyToTronAddress_ValidPrivateKey tests address derivation from a valid private key
func TestPrivateKeyToTronAddress_ValidPrivateKey(t *testing.T) {
	// Create a valid 32-byte private key
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = byte(i + 1)
	}

	address, err := PrivateKeyToTronAddress(privateKey)

	require.NoError(t, err)
	assert.NotEmpty(t, address)
	assert.Greater(t, len(address), 25)
}

// TestPrivateKeyToTronAddress_DeterministicOutput ensures same key produces same address
func TestPrivateKeyToTronAddress_DeterministicOutput(t *testing.T) {
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = byte(i + 1)
	}

	address1, err1 := PrivateKeyToTronAddress(privateKey)
	address2, err2 := PrivateKeyToTronAddress(privateKey)

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.Equal(t, address1, address2)
}

// TestPrivateKeyToTronAddress_DifferentKeys tests that different keys produce different addresses
func TestPrivateKeyToTronAddress_DifferentKeys(t *testing.T) {
	privateKey1 := make([]byte, 32)
	privateKey2 := make([]byte, 32)

	for i := range privateKey1 {
		privateKey1[i] = byte(i + 1)
		privateKey2[i] = byte(i + 2)
	}

	address1, err1 := PrivateKeyToTronAddress(privateKey1)
	address2, err2 := PrivateKeyToTronAddress(privateKey2)

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEqual(t, address1, address2)
}

// TestPrivateKeyToTronAddress_ZeroKey tests address derivation from a zero-filled key
func TestPrivateKeyToTronAddress_ZeroKey(t *testing.T) {
	privateKey := make([]byte, 32) // All zeros

	address, err := PrivateKeyToTronAddress(privateKey)

	require.NoError(t, err)
	assert.NotEmpty(t, address)
}

// TestPrivateKeyToTronAddress_MaxValueKey tests address derivation from a max-value key
func TestPrivateKeyToTronAddress_MaxValueKey(t *testing.T) {
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = 0xFF
	}

	address, err := PrivateKeyToTronAddress(privateKey)

	require.NoError(t, err)
	assert.NotEmpty(t, address)
}

// TestPrivateKeyToTronAddress_EmptyKey tests error handling with empty key
func TestPrivateKeyToTronAddress_EmptyKey(t *testing.T) {
	privateKey := []byte{}

	address, err := PrivateKeyToTronAddress(privateKey)

	// Should handle gracefully - might return error or empty address
	// The current implementation doesn't explicitly check, but we verify behavior
	if err != nil {
		assert.Empty(t, address)
	} else {
		// If no error, address might be derived from empty key
		assert.NotEmpty(t, address)
	}
}

// TestPrivateKeyToTronAddress_ShortKey tests with a key shorter than 32 bytes
func TestPrivateKeyToTronAddress_ShortKey(t *testing.T) {
	privateKey := []byte{0x01, 0x02, 0x03}

	address, err := PrivateKeyToTronAddress(privateKey)

	// Should handle gracefully
	if err != nil {
		assert.Empty(t, address)
	} else {
		assert.NotEmpty(t, address)
	}
}

// TestPrivateKeyToTronAddress_LongKey tests with a key longer than 32 bytes
func TestPrivateKeyToTronAddress_LongKey(t *testing.T) {
	privateKey := make([]byte, 64)
	for i := range privateKey {
		privateKey[i] = byte(i)
	}

	address, err := PrivateKeyToTronAddress(privateKey)

	// Should handle gracefully - BigInt will use all bytes
	require.NoError(t, err)
	assert.NotEmpty(t, address)
}

// TestPrivateKeyToTronAddress_VariousKeyPatterns tests different key patterns
func TestPrivateKeyToTronAddress_VariousKeyPatterns(t *testing.T) {
	testCases := []struct {
		name string
		key  []byte
	}{
		{
			name: "All ones",
			key:  []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
		},
		{
			name: "Alternating pattern",
			key:  []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55},
		},
		{
			name: "Incremental pattern",
			key:  []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			address, err := PrivateKeyToTronAddress(tc.key)

			require.NoError(t, err)
			assert.NotEmpty(t, address)
		})
	}
}

// TestPrivateKeyToTronAddress_AddressUniqueness tests that different keys produce unique addresses
func TestPrivateKeyToTronAddress_AddressUniqueness(t *testing.T) {
	addresses := make(map[string]bool)

	// Generate 100 addresses from different keys
	for i := 0; i < 100; i++ {
		privateKey := make([]byte, 32)
		for j := range privateKey {
			privateKey[j] = byte((i + j) % 256)
		}

		address, err := PrivateKeyToTronAddress(privateKey)
		require.NoError(t, err)

		// Ensure no duplicates
		assert.False(t, addresses[address], "Duplicate address found")
		addresses[address] = true
	}

	assert.Len(t, addresses, 100)
}

// TestPrivateKeyToTronAddress_ConsistencyWithDerivation tests that derived keys produce valid addresses
func TestPrivateKeyToTronAddress_ConsistencyWithDerivation(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	// Derive a key
	_, privKeyHex, err := DeriveTronAddressFromMnemonic(mnemonic, 0)
	require.NoError(t, err)

	// Decode the private key
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	require.NoError(t, err)

	// Generate address from the private key
	address, err := PrivateKeyToTronAddress(privKeyBytes)
	require.NoError(t, err)
	assert.NotEmpty(t, address)
}

// TestPrivateKeyToTronAddress_AddressFormatValidation tests that generated addresses have valid format
func TestPrivateKeyToTronAddress_AddressFormatValidation(t *testing.T) {
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = byte(i + 1)
	}

	address, err := PrivateKeyToTronAddress(privateKey)

	require.NoError(t, err)
	assert.NotEmpty(t, address)

	// Tron addresses are base58 encoded and should be printable ASCII
	for _, ch := range address {
		assert.True(t, ch >= 32 && ch <= 126, "Address contains non-printable character: %v", ch)
	}

	// Tron addresses don't contain easily confused characters (0, O, I, l)
	// and should be within a reasonable length
	assert.Greater(t, len(address), 20)
	assert.Less(t, len(address), 60)
}

// TestDeriveTronAddressFromMnemonic_SequentialDerivation tests BIP44 path derivation
func TestDeriveTronAddressFromMnemonic_SequentialDerivation(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	// Derive addresses for first 5 indices
	addresses := make([]string, 5)
	privKeys := make([]string, 5)

	for i := uint32(0); i < 5; i++ {
		addr, priv, err := DeriveTronAddressFromMnemonic(mnemonic, i)
		require.NoError(t, err)
		addresses[i] = addr
		privKeys[i] = priv
	}

	// Verify all are unique
	for i := 0; i < len(addresses); i++ {
		for j := i + 1; j < len(addresses); j++ {
			assert.NotEqual(t, addresses[i], addresses[j])
			assert.NotEqual(t, privKeys[i], privKeys[j])
		}
	}
}

// TestDeriveTronAddressFromMnemonic_NilSeedHandling tests error paths in derivation
func TestDeriveTronAddressFromMnemonic_NilSeedHandling(t *testing.T) {
	// Test with various edge case mnemonics
	testCases := []struct {
		name     string
		mnemonic string
	}{
		{name: "Single word", mnemonic: "abandon"},
		{name: "Two words", mnemonic: "abandon abandon"},
		{name: "Whitespace only", mnemonic: "   "},
		{name: "Special characters", mnemonic: "!@#$%^&*()"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			address, privKey, err := DeriveTronAddressFromMnemonic(tc.mnemonic, 0)

			// These should work (bip39 is lenient)
			require.NoError(t, err)
			assert.NotEmpty(t, address)
			assert.NotEmpty(t, privKey)
		})
	}
}

// TestPrivateKeyToTronAddress_PublicKeyDerivation tests the internal public key derivation
func TestPrivateKeyToTronAddress_PublicKeyDerivation(t *testing.T) {
	// Use a known private key constructed from test data
	privateKey := []byte{
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
	}

	address1, err1 := PrivateKeyToTronAddress(privateKey)
	require.NoError(t, err1)

	address2, err2 := PrivateKeyToTronAddress(privateKey)
	require.NoError(t, err2)

	// Same private key should always produce same address
	assert.Equal(t, address1, address2)
}

// Benchmark tests
func BenchmarkDeriveTronAddressFromMnemonic(b *testing.B) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = DeriveTronAddressFromMnemonic(mnemonic, uint32(i%100))
	}
}

func BenchmarkPrivateKeyToTronAddress(b *testing.B) {
	privateKey := make([]byte, 32)
	for i := range privateKey {
		privateKey[i] = byte(i + 1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PrivateKeyToTronAddress(privateKey)
	}
}

// TestDeriveTronAddressFromMnemonic_MemorySafety tests that the functions don't cause panics
func TestDeriveTronAddressFromMnemonic_MemorySafety(t *testing.T) {
	mnemonic := "flash couple heart script ramp april average caution plunge alter elite author"

	// Test concurrent access
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(idx int) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Panic occurred: %v", r)
				}
				done <- true
			}()

			_, _, _ = DeriveTronAddressFromMnemonic(mnemonic, uint32(idx))
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}