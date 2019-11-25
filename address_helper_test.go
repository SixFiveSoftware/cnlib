package cnlib

import "testing"

func TestBase58CheckEncoding_ValidAddress_ReturnsTrue(t *testing.T) {
	addresses := []string{
		"12vRFewBpbdiS5HXDDLEfVFtJnpA2x8NV8",
		"16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM",
		"3EH9Wj6KWaZBaYXhVCa8ZrwpHJYtk44bGX",
		"3Cd4xEu2VvM352BVgd9cb1Ct5vxz318tVT",
	}

	for _, addr := range addresses {
		valid := AddressIsBase58CheckEncoded(addr)
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
		valid := AddressIsBase58CheckEncoded(addr)
		if valid {
			t.Errorf("Expected %v to not be base58Check encoded", addr)
		}
	}
}
