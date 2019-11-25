package cnlib

import "testing"

const (
	words = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
)

func TestMetaAddress_Receive_Segwit_Address(t *testing.T) {
	path := NewDerivationPath(84, 0, 0, 0, 0)
	coin := NewBaseCoin(84, 0, 0)
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

func TestMetaAddress_RegTestAddresses(t *testing.T) {
	bc := NewBaseCoin(84, 1, 0)

	rpath0 := NewDerivationPath(84, 1, 0, 0, 0)
	rpath1 := NewDerivationPath(84, 1, 0, 0, 1)
	rpath2 := NewDerivationPath(84, 1, 0, 0, 2)
	cpath0 := NewDerivationPath(84, 1, 0, 1, 0)
	cpath1 := NewDerivationPath(84, 1, 0, 1, 1)
	cpath2 := NewDerivationPath(84, 1, 0, 1, 2)

	rexp0 := "bcrt1q6rz28mcfaxtmd6v789l9rrlrusdprr9pz3cppk"
	rexp1 := "bcrt1qd7spv5q28348xl4myc8zmh983w5jx32cs707jh"
	rexp2 := "bcrt1qxdyjf6h5d6qxap4n2dap97q4j5ps6ua8jkxz0z"
	cexp0 := "bcrt1q9u62588spffmq4dzjxsr5l297znf3z6jkgnhsw"
	cexp1 := "bcrt1qkwgskuzmmwwvqajnyr7yp9hgvh5y45kg984qvy"
	cexp2 := "bcrt1q2vma00td2g9llw8hwa8ny3r774rtt7ae3q2e44"

	wallet := NewHDWalletFromWords(w, bc)

	rua0 := NewUsableAddress(wallet, rpath0).MetaAddress().Address
	rua1 := NewUsableAddress(wallet, rpath1).MetaAddress().Address
	rua2 := NewUsableAddress(wallet, rpath2).MetaAddress().Address
	cua0 := NewUsableAddress(wallet, cpath0).MetaAddress().Address
	cua1 := NewUsableAddress(wallet, cpath1).MetaAddress().Address
	cua2 := NewUsableAddress(wallet, cpath2).MetaAddress().Address

	if rua0 != rexp0 {
		t.Errorf("Expected address %v, got %v", rexp0, rua0)
	}
	if rua1 != rexp1 {
		t.Errorf("Expected address %v, got %v", rexp1, rua1)
	}
	if rua2 != rexp2 {
		t.Errorf("Expected address %v, got %v", rexp2, rua2)
	}
	if cua0 != cexp0 {
		t.Errorf("Expected address %v, got %v", cexp0, cua0)
	}
	if cua1 != cexp1 {
		t.Errorf("Expected address %v, got %v", cexp1, cua1)
	}
	if cua2 != cexp2 {
		t.Errorf("Expected address %v, got %v", cexp2, cua2)
	}
}

func TestMetaAddress_Receive_LegacySegwit_Address(t *testing.T) {
	path := NewDerivationPath(49, 0, 0, 0, 0)
	coin := NewBaseCoin(49, 0, 0)
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
