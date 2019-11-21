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

func TestMetaAddress_Receive_LegacySegwit_Address(t *testing.T) {
	path := NewDerivationPath(49, 0, 0, 0, 0)
	coin := NewBaseCoin(49, 0, 0)
	words := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	wallet := NewHDWalletFromWords(words, coin)
	usableAddress := NewUsableAddress(wallet, path)
	meta := usableAddress.MetaAddress()
	expectedAddr := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	expectedPubkey := "049b3b694b8fc5b5e07fb069c783cac754f5d38c3e08bed1960e31fdb1dda35c2449bdd1f0ae7d37a04991d4f5927efd359c13189437d9eae0faf7d003ffd04c89"

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

func TestMetaAddress_Change_LegacySegwit_Address(t *testing.T) {
	path := NewDerivationPath(49, 0, 0, 1, 0)
	coin := NewBaseCoin(49, 0, 0)
	words := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	wallet := NewHDWalletFromWords(words, coin)
	usableAddress := NewUsableAddress(wallet, path)
	meta := usableAddress.MetaAddress()
	expectedAddr := "34K56kSjgUCUSD8GTtuF7c9Zzwokbs6uZ7"
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
