package cnlib

import "testing"

func addressHelperTestHelpers() *AddressHelper {
	bc := NewBaseCoin(84, 0, 0)
	ah := NewAddressHelper(bc)
	return ah
}

func TestBase58CheckEncoding_ValidAddress_ReturnsTrue(t *testing.T) {
	addresses := []string{
		"12vRFewBpbdiS5HXDDLEfVFtJnpA2x8NV8",
		"16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM",
		"3EH9Wj6KWaZBaYXhVCa8ZrwpHJYtk44bGX",
		"3Cd4xEu2VvM352BVgd9cb1Ct5vxz318tVT",
	}

	for _, addr := range addresses {
		valid := addressHelperTestHelpers().AddressIsBase58CheckEncoded(addr)
		if !valid {
			t.Errorf("Expected %v to be base58Check encoded", addr)
		}
	}

}

func TestBase58CheckEncoding_InvalidAddresses_ReturnFalse(t *testing.T) {
	addresses := []string{
		"12vRFewBpbdiS5HXDDLEfVFtJnpA2",
		"12vRFewBpbdiS5HXDDLEfVFt",
		"diS5HXDDLEfVFtJnpA2x8NV8",
		"212vRFewBpbdiS5HXDDLEfVFtJnpA2x8NV8",
		"42vRFewBpbdiS5HXDDLEfVFtJnpA2x8NV8",
		"3EH9Wj6KWaZBaYXhVCa8ZrwpHJYtk",
		"j6KWaZBaYXhVCa8ZrwpHJYtk44bGX",
		"ZBaYXhVCa8ZrwpHJYtk44bGX",
		"23EH9Wj6KWaZBaYXhVCa8ZrwpHJYtk44bGX",
		"3Cd4xEu2VvM352BVgd9cb1Ct5vxz3",
		"3Cd4xEu2VvM352BVgd9cb1Ct",
		"Eu2VvM352BVgd9cb1Ct5vxz318tVT",
		"M352BVgd9cb1Ct5vxz318tVT",
		"23Cd4xEu2VvM352BVgd9cb1Ct5vxz318tVT",
		"4Cd4xEu2VvM352BVgd9cb1Ct5vxz318tVT",
		"0xF26C29D25a1E1696c5CC54DE4bf2AEc906EB4F79",
		"qr45rul6luexjgg5h8p26c0cs6rrhwzrkg6e0hdvrf",
		"Jenny86753098675309IgotIt",
		"31415926535ILikePi89793238462643",
		"foo",
		"",
	}

	for _, addr := range addresses {
		valid := addressHelperTestHelpers().AddressIsBase58CheckEncoded(addr)
		if valid {
			t.Errorf("Expected %v to not be base58Check encoded", addr)
		}
	}
}

func TestSegwitAddress_ValidAddresses_ReturnTrue(t *testing.T) {
	addresses := []string{
		"bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu",                     // demo wallet first receive address
		"bc1q8c6fshw2dlwun7ekn9qwf37cu2rn755upcp6el",                     // demo wallet first change address
		"bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",                     // p2wpkh sipa demo
		"bc1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3qccfmv3", // p2wsh sipa demo
	}

	for _, addr := range addresses {
		valid := addressHelperTestHelpers().AddressIsValidSegwitAddress(addr)
		if !valid {
			t.Errorf("Expected %v to be valid segwit address.", addr)
		}
	}
}

func TestSegwitAddress_InvalidAddresses_ReturnFalse(t *testing.T) {
	addresses := []string{
		"BC1QW508D6QEJXTDG4Y5R3ZARVAYR0C5XW7KV8F3T4", // p2wsh sipa demo invalid p2wpkh, YR transposed
		"3Cd4xEu2VvM352BVgd9cb1Ct5vxz318tVT",
	}

	for _, addr := range addresses {
		valid := addressHelperTestHelpers().AddressIsValidSegwitAddress(addr)
		if valid {
			t.Errorf("Expected %v to be invalid segwit address.", addr)
		}
	}
}

func TestSegwitAddressHRP(t *testing.T) {
	bcAddr := "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu"
	rtAddr := "bcrt1q6rz28mcfaxtmd6v789l9rrlrusdprr9pz3cppk"
	legacyAddr := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	helper := addressHelperTestHelpers()

	bcHrp, bcErr := helper.HRPFromAddress(bcAddr)
	if bcErr != nil || bcHrp != "bc" {
		t.Errorf("Expected hrp of bc, got %v. Error: %v", bcHrp, bcErr)
	}
	rtHrp, rtErr := helper.HRPFromAddress(rtAddr)
	if rtErr != nil || rtHrp != "bcrt" {
		t.Errorf("Expected hrp of bcrt, got %v. Error: %v", rtHrp, rtErr)
	}
	laHrp, laErr := helper.HRPFromAddress(legacyAddr)
	if laErr == nil || laHrp != "" {
		t.Errorf("Expected error, got %v. Error: %v", laHrp, laErr)
	}
}
