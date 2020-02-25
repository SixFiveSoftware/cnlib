package cnlib

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountExtendedKeyPrefix_m_44_0(t *testing.T) {
	bc := NewBaseCoin(44, 0, 0)
	key, err := bc.defaultExtendedPubkeyType()
	assert.Nil(t, err)
	assert.Equal(t, "xpub", key)
}

func TestAccountExtendedKeyPrefix_m_49_0(t *testing.T) {
	bc := NewBaseCoin(49, 0, 0)
	key, err := bc.defaultExtendedPubkeyType()
	assert.Nil(t, err)
	assert.Equal(t, "ypub", key)
}

func TestAccountExtendedKeyPrefix_m_84_0(t *testing.T) {
	bc := NewBaseCoin(84, 0, 0)
	key, err := bc.defaultExtendedPubkeyType()
	assert.Nil(t, err)
	assert.Equal(t, "zpub", key)
}

func TestAccountExtendedKeyPrefix_m_44_1(t *testing.T) {
	bc := NewBaseCoin(44, 1, 0)
	key, err := bc.defaultExtendedPubkeyType()
	assert.Nil(t, err)
	assert.Equal(t, "tpub", key)
}

func TestAccountExtendedKeyPrefix_m_49_1(t *testing.T) {
	bc := NewBaseCoin(49, 1, 0)
	key, err := bc.defaultExtendedPubkeyType()
	assert.Nil(t, err)
	assert.Equal(t, "upub", key)
}

func TestAccountExtendedKeyPrefix_m_84_1(t *testing.T) {
	bc := NewBaseCoin(84, 1, 0)
	key, err := bc.defaultExtendedPubkeyType()
	assert.Nil(t, err)
	assert.Equal(t, "vpub", key)
}

func TestAccountExtendedKeyPrefix_m_45_0(t *testing.T) {
	bc := NewBaseCoin(45, 0, 0)
	key, err := bc.defaultExtendedPubkeyType()
	assert.NotNil(t, err)
	assert.EqualError(t, errors.New("invalid basecoin purpose value"), err.Error())
	assert.Equal(t, "", key)
}

func TestAccountExtendedKeyPrefix_m_44_2(t *testing.T) {
	bc := NewBaseCoin(44, 2, 0)
	key, err := bc.defaultExtendedPubkeyType()
	assert.NotNil(t, err)
	assert.EqualError(t, errors.New("invalid basecoin coin value"), err.Error())
	assert.Equal(t, "", key)
}

func TestNewBaseCoin_FromAcctPubKey_M_44_0_0(t *testing.T) {
	keyStr := "xpub6BosfCnifzxcFwrSzQiqu2DBVTshkCXacvNsWGYJVVhhawA7d4R5WSWGFNbi8Aw6ZRc1brxMyWMzG3DSSSSoekkudhUd9yLb6qx39T9nMdj"
	bc, err := NewBaseCoinFromAccountPubKey(keyStr)
	assert.Nil(t, err)
	assert.Equal(t, 44, bc.Purpose)
	assert.Equal(t, 0, bc.Coin)
	assert.Equal(t, 0, bc.Account)
}

func TestNewBaseCoin_FromAcctPubKey_M_44_0_1(t *testing.T) {
	keyStr := "xpub6BosfCnifzxcJJ1wYuntGJfF2zPJkDeG9ELNHcKNjezuea4tumswN9sH1psMdSVqCMoJC21Bv8usSeqSP4Sp1tLzW7aY59fGn9GCYzx5UTo"
	bc, err := NewBaseCoinFromAccountPubKey(keyStr)
	assert.Nil(t, err)
	assert.Equal(t, 44, bc.Purpose)
	assert.Equal(t, 0, bc.Coin)
	assert.Equal(t, 1, bc.Account)
}

func TestNewBaseCoin_FromAcctPubKey_M_49_0_0(t *testing.T) {
	keyStr := "ypub6Ww3ibxVfGzLrAH1PNcjyAWenMTbbAosGNB6VvmSEgytSER9azLDWCxoJwW7Ke7icmizBMXrzBx9979FfaHxHcrArf3zbeJJJUZPf663zsP"
	bc, err := NewBaseCoinFromAccountPubKey(keyStr)
	assert.Nil(t, err)
	assert.Equal(t, 49, bc.Purpose)
	assert.Equal(t, 0, bc.Coin)
	assert.Equal(t, 0, bc.Account)
}

func TestNewBaseCoin_FromAcctPubKey_M_84_0_0(t *testing.T) {
	keyStr := "zpub6rFR7y4Q2AijBEqTUquhVz398htDFrtymD9xYYfG1m4wAcvPhXNfE3EfH1r1ADqtfSdVCToUG868RvUUkgDKf31mGDtKsAYz2oz2AGutZYs"
	bc, err := NewBaseCoinFromAccountPubKey(keyStr)
	assert.Nil(t, err)
	assert.Equal(t, 84, bc.Purpose)
	assert.Equal(t, 0, bc.Coin)
	assert.Equal(t, 0, bc.Account)
}
