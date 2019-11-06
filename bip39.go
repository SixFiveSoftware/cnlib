package golib

import (
	"github.com/tyler-smith/go-bip39/wordlists"
)

func GetWordList() []string {
	return wordlists.English
}
