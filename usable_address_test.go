package cnlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaAddress_Receive_Segwit_Address(t *testing.T) {
	path := NewDerivationPath(BaseCoinBip84MainNet, 0, 0)
	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)

	usableAddress, err := newUsableAddressWithDerivationPath(wallet, path)
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
	path := NewDerivationPath(BaseCoinBip84MainNet, 1, 0)
	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	usableAddress, err := newUsableAddressWithDerivationPath(wallet, path)
	assert.Nil(t, err)

	meta, err := usableAddress.MetaAddress()
	assert.Nil(t, err)

	expectedAddr := "bc1q8c6fshw2dlwun7ekn9qwf37cu2rn755upcp6el"
	expectedPubkey := ""

	assert.Equal(t, expectedAddr, meta.Address)
	assert.Equal(t, path, meta.DerivationPath)
	assert.Equal(t, expectedPubkey, meta.UncompressedPublicKey)
}

func TestMetaAddress_TestNetAddresses(t *testing.T) {
	rpath0 := NewDerivationPath(BaseCoinBip84TestNet, 0, 0)
	rpath1 := NewDerivationPath(BaseCoinBip84TestNet, 0, 1)
	rpath2 := NewDerivationPath(BaseCoinBip84TestNet, 0, 2)
	cpath0 := NewDerivationPath(BaseCoinBip84TestNet, 1, 0)
	cpath1 := NewDerivationPath(BaseCoinBip84TestNet, 1, 1)
	cpath2 := NewDerivationPath(BaseCoinBip84TestNet, 1, 2)

	rexp0 := "tb1q6rz28mcfaxtmd6v789l9rrlrusdprr9pqcpvkl"
	rexp1 := "tb1qd7spv5q28348xl4myc8zmh983w5jx32cjhkn97"
	rexp2 := "tb1qxdyjf6h5d6qxap4n2dap97q4j5ps6ua8sll0ct"
	cexp0 := "tb1q9u62588spffmq4dzjxsr5l297znf3z6j5p2688"
	cexp1 := "tb1qkwgskuzmmwwvqajnyr7yp9hgvh5y45kg8wvdmd"
	cexp2 := "tb1q2vma00td2g9llw8hwa8ny3r774rtt7aenfn5zu"

	wallet := NewHDWalletFromWords(w, BaseCoinBip84TestNet)

	rua0, err := newUsableAddressWithDerivationPath(wallet, rpath0)
	assert.Nil(t, err)
	rua1, err := newUsableAddressWithDerivationPath(wallet, rpath1)
	assert.Nil(t, err)
	rua2, err := newUsableAddressWithDerivationPath(wallet, rpath2)
	assert.Nil(t, err)
	cua0, err := newUsableAddressWithDerivationPath(wallet, cpath0)
	assert.Nil(t, err)
	cua1, err := newUsableAddressWithDerivationPath(wallet, cpath1)
	assert.Nil(t, err)
	cua2, err := newUsableAddressWithDerivationPath(wallet, cpath2)
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
	path := NewDerivationPath(BaseCoinBip49MainNet, 0, 0)
	wallet := NewHDWalletFromWords(w, BaseCoinBip49MainNet)

	usableAddress, err := newUsableAddressWithDerivationPath(wallet, path)
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
	path := NewDerivationPath(BaseCoinBip49MainNet, 1, 0)
	wallet := NewHDWalletFromWords(w, BaseCoinBip49MainNet)

	usableAddress, err := newUsableAddressWithDerivationPath(wallet, path)
	assert.Nil(t, err)

	meta, err := usableAddress.MetaAddress()
	assert.Nil(t, err)

	expectedAddr := "34K56kSjgUCUSD8GTtuF7c9Zzwokbs6uZ7"
	expectedPubkey := ""

	assert.Equal(t, expectedAddr, meta.Address)
	assert.Equal(t, path, meta.DerivationPath)
	assert.Equal(t, expectedPubkey, meta.UncompressedPublicKey)
}
