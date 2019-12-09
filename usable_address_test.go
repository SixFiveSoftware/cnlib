package cnlib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetaAddress_Receive_Segwit_Address(t *testing.T) {
	path := NewDerivationPath(84, 0, 0, 0, 0)
	coin := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, coin)

	usableAddress, err := NewUsableAddressWithDerivationPath(wallet, path)
	assert.Nil(t, err)

	meta, err := usableAddress.MetaAddress()
	assert.Nil(t, err)

	expectedAddr := "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu"
	expectedPubkey := "0430d54fd0dd420a6e5f8d3624f5f3482cae350f79d5f0753bf5beef9c2d91af3c04717159ce0828a7f686c2c7510b7aa7d4c685ebc2051642ccbebc7099e2f679"

	assert.Equal(t, expectedAddr, meta.Address)
	assert.Equal(t, path, meta.DerivationPath)
	assert.Equal(t, expectedPubkey, meta.UncompressedPublicKey)
}

func TestMetaAddress_Change_Segwit_Address(t *testing.T) {
	path := NewDerivationPath(84, 0, 0, 1, 0)
	coin := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, coin)
	usableAddress, err := NewUsableAddressWithDerivationPath(wallet, path)
	assert.Nil(t, err)

	meta, err := usableAddress.MetaAddress()
	assert.Nil(t, err)

	expectedAddr := "bc1q8c6fshw2dlwun7ekn9qwf37cu2rn755upcp6el"
	expectedPubkey := ""

	assert.Equal(t, expectedAddr, meta.Address)
	assert.Equal(t, path, meta.DerivationPath)
	assert.Equal(t, expectedPubkey, meta.UncompressedPublicKey)
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

	rua0, err := NewUsableAddressWithDerivationPath(wallet, rpath0)
	assert.Nil(t, err)
	rua1, err := NewUsableAddressWithDerivationPath(wallet, rpath1)
	assert.Nil(t, err)
	rua2, err := NewUsableAddressWithDerivationPath(wallet, rpath2)
	assert.Nil(t, err)
	cua0, err := NewUsableAddressWithDerivationPath(wallet, cpath0)
	assert.Nil(t, err)
	cua1, err := NewUsableAddressWithDerivationPath(wallet, cpath1)
	assert.Nil(t, err)
	cua2, err := NewUsableAddressWithDerivationPath(wallet, cpath2)
	assert.Nil(t, err)

	rua0meta, err := rua0.MetaAddress()
	assert.Nil(t, err)
	rua1meta, err := rua1.MetaAddress()
	assert.Nil(t, err)
	rua2meta, err := rua2.MetaAddress()
	assert.Nil(t, err)
	cua0meta, err := cua0.MetaAddress()
	assert.Nil(t, err)
	cua1meta, err := cua1.MetaAddress()
	assert.Nil(t, err)
	cua2meta, err := cua2.MetaAddress()
	assert.Nil(t, err)

	assert.Equal(t, rexp0, rua0meta.Address)
	assert.Equal(t, rexp1, rua1meta.Address)
	assert.Equal(t, rexp2, rua2meta.Address)
	assert.Equal(t, cexp0, cua0meta.Address)
	assert.Equal(t, cexp1, cua1meta.Address)
	assert.Equal(t, cexp2, cua2meta.Address)
}

func TestMetaAddress_Receive_LegacySegwit_Address(t *testing.T) {
	path := NewDerivationPath(49, 0, 0, 0, 0)
	coin := NewBaseCoin(49, 0, 0)
	wallet := NewHDWalletFromWords(w, coin)

	usableAddress, err := NewUsableAddressWithDerivationPath(wallet, path)
	assert.Nil(t, err)

	meta, err := usableAddress.MetaAddress()
	assert.Nil(t, err)

	expectedAddr := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	expectedPubkey := "049b3b694b8fc5b5e07fb069c783cac754f5d38c3e08bed1960e31fdb1dda35c2449bdd1f0ae7d37a04991d4f5927efd359c13189437d9eae0faf7d003ffd04c89"

	assert.Equal(t, expectedAddr, meta.Address)
	assert.Equal(t, path, meta.DerivationPath)
	assert.Equal(t, expectedPubkey, meta.UncompressedPublicKey)
}

func TestMetaAddress_Change_LegacySegwit_Address(t *testing.T) {
	path := NewDerivationPath(49, 0, 0, 1, 0)
	coin := NewBaseCoin(49, 0, 0)
	wallet := NewHDWalletFromWords(w, coin)

	usableAddress, err := NewUsableAddressWithDerivationPath(wallet, path)
	assert.Nil(t, err)

	meta, err := usableAddress.MetaAddress()
	assert.Nil(t, err)

	expectedAddr := "34K56kSjgUCUSD8GTtuF7c9Zzwokbs6uZ7"
	expectedPubkey := ""

	assert.Equal(t, expectedAddr, meta.Address)
	assert.Equal(t, path, meta.DerivationPath)
	assert.Equal(t, expectedPubkey, meta.UncompressedPublicKey)
}
