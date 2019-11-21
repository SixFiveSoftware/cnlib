package cnlib

import "testing"

func TestMetaAddress_Receive_Segwit_Address(t *testing.T) {
	path := NewDerivationPath(84, 0, 0, 0, 0)
	coin := NewBaseCoin(84, 0, 0)
	words := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	wallet := NewHDWalletFromWords(words, coin)
	usableAddress := NewUsableAddress(wallet, path)
	meta := usableAddress.MetaAddress()
	expectedAddr := "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu"
	expectedPubkey := "0430d54fd0dd420a6e5f8d3624f5f3482cae350f79d5f0753bf5beef9c2d91af3c04717159ce0828a7f686c2c7510b7aa7d4c685ebc2051642ccbebc7099e2f679"

	if meta.Address != expectedAddr {
		t.Errorf("Expected address %v, got %v", expectedAddr, meta.Address)
	}

	if meta.DerivationPath != path {
		t.Errorf("Expected path %v, got %v", path, meta.DerivationPath)
	}

	if meta.UncompressedPublicKey != expectedPubkey {
		t.Errorf("Expected pubkey %v, got %v", expectedPubkey, meta.UncompressedPublicKey)
	}
}

func TestMetaAddress_Change_Segwit_Address(t *testing.T) {
	path := NewDerivationPath(84, 0, 0, 1, 0)
	coin := NewBaseCoin(84, 0, 0)
	words := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	wallet := NewHDWalletFromWords(words, coin)
	usableAddress := NewUsableAddress(wallet, path)
	meta := usableAddress.MetaAddress()
	expectedAddr := "bc1q8c6fshw2dlwun7ekn9qwf37cu2rn755upcp6el"
	expectedPubkey := ""

	if meta.Address != expectedAddr {
		t.Errorf("Expected address %v, got %v", expectedAddr, meta.Address)
	}

	if meta.DerivationPath != path {
		t.Errorf("Expected path %v, got %v", path, meta.DerivationPath)
	}

	if meta.UncompressedPublicKey != expectedPubkey {
		t.Errorf("Expected pubkey empty string, got %v", meta.UncompressedPublicKey)
	}
}
