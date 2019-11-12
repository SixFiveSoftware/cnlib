package cnlib

import (
	"strings"

	"github.com/tyler-smith/go-bip39/wordlists"
	"github.com/worldiety/std"
)

func GetWordListStrSlice() *std.StrSlice {
		return &std.StrSlice{
			Slice: wordlists.English,
		}
}

func GetWordListString() string {
	return strings.Join(wordlists.English, ",")
}
