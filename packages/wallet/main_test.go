package main

import (
	"testing"

	"github.com/btcsuite/btcutil/base58"
)

func TestDeriveTronAddressFromMnemonic_Deterministic(t *testing.T) {
	mn := "flash couple heart script ramp april average caution plunge alter elite author"
	addr1, pk1, err := DeriveTronAddressFromMnemonic(mn, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	addr2, pk2, err := DeriveTronAddressFromMnemonic(mn, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr1 == "" || pk1 == "" {
		t.Fatalf("expected non-empty address/private key")
	}
	if addr1 != addr2 || pk1 != pk2 {
		t.Fatalf("expected deterministic output for same mnemonic/index")
	}
	addr3, pk3, err := DeriveTronAddressFromMnemonic(mn, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr3 == addr1 {
		t.Errorf("expected different address for different index")
	}
	if pk3 == pk1 {
		t.Errorf("expected different private key for different index")
	}
	if len(pk1) != 64 {
		t.Errorf("expected hex private key length 64, got %d", len(pk1))
	}
}

func TestPrivateKeyToTronAddress_Format(t *testing.T) {
	// 32-byte private key
	priv := make([]byte, 32)
	for i := range priv {
		priv[i] = byte(i + 1)
	}
	addr, err := PrivateKeyToTronAddress(priv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr == "" {
		t.Fatal("expected non-empty address")
	}
	decoded := base58.Decode(addr)
	if len(decoded) != 25 {
		t.Errorf("expected 25-byte payload (0x41 + 20 bytes + 4 checksum), got %d", len(decoded))
	}
	if len(decoded) > 0 && decoded[0] != 0x41 {
		t.Errorf("expected Tron prefix 0x41, got 0x%02x", decoded[0])
	}
}

func TestDeriveTronAddressFromMnemonic_EmptyMnemonic(t *testing.T) {
	addr, pk, err := DeriveTronAddressFromMnemonic("", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr == "" || pk == "" {
		t.Fatalf("expected non-empty address and private key even with empty mnemonic")
	}
}