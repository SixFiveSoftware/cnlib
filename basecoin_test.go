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
