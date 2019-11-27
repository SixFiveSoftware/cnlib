package cnlib

import "testing"

func addressHelper() *AddressHelper {
	bc := NewBaseCoin(49, 0, 0)
	ah := NewAddressHelper(bc)
	return ah
}

func TestNewTransactionDataStandard_SingleOutput_SingleInput_SatisfiesAmount(t *testing.T) {
	// given
	paymentAmount := 50000000 // 0.5 BTC
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	utxoAmount := 100000000 // 1.0 BTC
	changePath := NewDerivationPath(84, 0, 0, 1, 0)
	utxoPath := NewDerivationPath(49, 0, 0, 0, 0)
	utxo := NewUTXO("previous txid", 0, utxoAmount, utxoPath, true)
	utxos := []*UTXO{utxo}
	feeRate := 30
	totalBytes, bytesErr := addressHelper().totalBytes(utxos, address, true)
	if bytesErr != nil {
		t.Errorf("Expected total bytes, got error: %v", bytesErr)
	}

	expectedFeeAmount := feeRate * totalBytes // 4,980
	expectedChangeAmount := (utxoAmount - paymentAmount - expectedFeeAmount)
	expectedNumberOfUTXOs := 1
	expectedLocktime := 500000

	// when
	data := NewTransactionDataStandard(
		address, addressHelper().Basecoin, paymentAmount, feeRate, changePath, 500000, MustBeRBF,
	)
	data.TransactionData.AddUTXO(utxo)
	success, err := data.Generate()

	// then
	if !success {
		t.Errorf("Failed to generate transaction. Error: %v", err)
	}

	if data.TransactionData.Amount != paymentAmount {
		t.Errorf("Expected amount to be %v, got %v", paymentAmount, data.TransactionData.Amount)
	}
	if data.TransactionData.FeeAmount != expectedFeeAmount {
		t.Errorf("Expected fee amount to be %v, got %v", expectedFeeAmount, data.TransactionData.FeeAmount)
	}
	if expectedChangeAmount != 49995020 {
		t.Errorf("Expected change amount to be %v, got %v", 49995020, data.TransactionData.ChangeAmount)
	}
	if data.TransactionData.ChangeAmount != expectedChangeAmount {
		t.Errorf("Expected change amount to be %v, got %v", expectedChangeAmount, data.TransactionData.ChangeAmount)
	}
	if data.TransactionData.utxoCount() != expectedNumberOfUTXOs {
		t.Errorf("Expected number of UTXOs to be %v, got %v", expectedNumberOfUTXOs, data.TransactionData.utxoCount())
	}
	if data.TransactionData.Locktime != expectedLocktime {
		t.Errorf("Expected locktime to be %v, got %v", expectedLocktime, data.TransactionData.Locktime)
	}
	if !data.TransactionData.shouldAddChangeToTransaction() {
		t.Errorf("Expected to add change to transaction.")
	}
}

func TestTransactionDataStandard_SingleOutput_DoubleInput_WithChange(t *testing.T) {
	// given
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	ah := addressHelperTestHelpers() // helper uses 84 purpose to ensure correct input size is calculated
	paymentAmount := 50000000        // 0.5 BTC
	utxoAmount := 30000000           // 0.3 BTC
	changePath := NewDerivationPath(84, 0, 0, 1, 0)
	utxoPath := NewDerivationPath(49, 0, 0, 0, 0)
	utxo1 := NewUTXO("previous txid", 0, utxoAmount, utxoPath, true)
	utxo2 := NewUTXO("previous txid", 1, utxoAmount, utxoPath, true)
	utxos := []*UTXO{utxo1, utxo2}
	feeRate := 30
	totalBytes, tbErr := ah.totalBytes(utxos, address, true)
	if tbErr != nil {
		t.Errorf("Expected to get total bytes from helper, got error: %v", tbErr)
	}
	expectedFeeAmount := feeRate * totalBytes // 7,680
	amountFromUTXOs := 0
	for _, utxo := range utxos {
		amountFromUTXOs += utxo.Amount
	}
	expectedChangeAmount := amountFromUTXOs - paymentAmount - expectedFeeAmount
	expectedNumberOfUTXOs := 2
	expectedLocktime := 500000
	expectedRBFOption := RBFOption(AllowedToBeRBF)

	// when
	data := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, AllowedToBeRBF)
	for _, utxo := range utxos {
		data.TransactionData.AddUTXO(utxo)
	}
	success, dataErr := data.Generate()

	// then
	if !success {
		t.Errorf("Expected to build tx data, got error: %v", dataErr)
	}
	if data.TransactionData.Amount != paymentAmount {
		t.Errorf("Expected amount to be %v, got %v", paymentAmount, data.TransactionData.Amount)
	}
	if data.TransactionData.FeeAmount != expectedFeeAmount {
		t.Errorf("Expected fee amount to be %v, got %v", expectedFeeAmount, data.TransactionData.FeeAmount)
	}
	if data.TransactionData.ChangeAmount != expectedChangeAmount {
		t.Errorf("Expected change amount to be %v, got %v", expectedChangeAmount, data.TransactionData.ChangeAmount)
	}
	if data.TransactionData.utxoCount() != expectedNumberOfUTXOs {
		t.Errorf("Expected number of UTXOs to be %v, got %v", expectedNumberOfUTXOs, data.TransactionData.utxoCount())
	}
	if data.TransactionData.Locktime != expectedLocktime {
		t.Errorf("Expected locktime to be %v, got %v", expectedLocktime, data.TransactionData.Locktime)
	}
	if !data.TransactionData.shouldAddChangeToTransaction() {
		t.Errorf("Expected to add change to transaction.")
	}
	if data.TransactionData.RBFOption != expectedRBFOption {
		t.Errorf("Expected RBFOption to be %v, got %v", expectedRBFOption, data.TransactionData.RBFOption)
	}
}

func TestNewTransactionDataStandard_SingleInput_SingleOutput_NoChange(t *testing.T) {
	// given
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	ah := addressHelperTestHelpers() // helper uses 84 purpose to ensure correct input size is calculated
	paymentAmount := 50000000        // 0.5 BTC
	utxoAmount := 50004020           // 0.50004020 BTC
	changePath := NewDerivationPath(84, 0, 0, 1, 0)
	utxoPath := NewDerivationPath(49, 0, 0, 0, 0)
	utxo := NewUTXO("previous txid", 0, utxoAmount, utxoPath, true)
	utxos := []*UTXO{utxo}
	feeRate := 30
	totalBytes, tbErr := ah.totalBytes(utxos, address, false)
	if tbErr != nil {
		t.Errorf("Expected to get total bytes from helper, got error: %v", tbErr)
	}
	expectedFeeAmount := feeRate * totalBytes // 4,020
	amountFromUTXOs := 0
	for _, utxo := range utxos {
		amountFromUTXOs += utxo.Amount
	}
	expectedChangeAmount := amountFromUTXOs - paymentAmount - expectedFeeAmount
	expectedNumberOfUTXOs := 1
	expectedLocktime := 500000
	expectedRBFOption := RBFOption(AllowedToBeRBF)

	// when
	data := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, AllowedToBeRBF)
	for _, utxo := range utxos {
		data.TransactionData.AddUTXO(utxo)
	}
	success, dataErr := data.Generate()

	// then
	if !success {
		t.Errorf("Expected to build tx data, got error: %v", dataErr)
	}
	if data.TransactionData.Amount != paymentAmount {
		t.Errorf("Expected amount to be %v, got %v", paymentAmount, data.TransactionData.Amount)
	}
	if data.TransactionData.FeeAmount != expectedFeeAmount {
		t.Errorf("Expected fee amount to be %v, got %v", expectedFeeAmount, data.TransactionData.FeeAmount)
	}
	if data.TransactionData.ChangeAmount != expectedChangeAmount {
		t.Errorf("Expected change amount to be %v, got %v", expectedChangeAmount, data.TransactionData.ChangeAmount)
	}
	if data.TransactionData.utxoCount() != expectedNumberOfUTXOs {
		t.Errorf("Expected number of UTXOs to be %v, got %v", expectedNumberOfUTXOs, data.TransactionData.utxoCount())
	}
	if data.TransactionData.Locktime != expectedLocktime {
		t.Errorf("Expected locktime to be %v, got %v", expectedLocktime, data.TransactionData.Locktime)
	}
	if data.TransactionData.shouldAddChangeToTransaction() {
		t.Errorf("Expected to not add change to transaction.")
	}
	if data.TransactionData.RBFOption != expectedRBFOption {
		t.Errorf("Expected RBFOption to be %v, got %v", expectedRBFOption, data.TransactionData.RBFOption)
	}
}

func TestNewTransactionStandard_SingleOutput_DoubleInput_NoChange(t *testing.T) {
	// given
	ah := addressHelperTestHelpers()
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	paymentAmount := 50000000 // 0.50000000 BTC
	utxoAmount1 := 20001750   // 0.20001750 BTC
	utxoAmount2 := 30005000   // 0.30005000 BTC
	changePath := NewDerivationPath(84, 0, 0, 1, 0)
	utxoPath := NewDerivationPath(49, 0, 0, 0, 0)
	utxo1 := NewUTXO("previous txid", 0, utxoAmount1, utxoPath, true)
	utxo2 := NewUTXO("previous txid", 1, utxoAmount2, utxoPath, true)
	utxos := []*UTXO{utxo1, utxo2}
	feeRate := 30
	totalBytes, tbErr := ah.totalBytes(utxos, address, false)
	if tbErr != nil {
		t.Errorf("Expected to get total bytes from helper, got error: %v", tbErr)
	}
	expectedFeeAmount := feeRate * totalBytes // 6, 750
	amountFromUTXOs := 0
	for _, utxo := range utxos {
		amountFromUTXOs += utxo.Amount
	}
	expectedChangeAmount := 0
	expectedNumberOfUTXOs := len(utxos)
	expectedLocktime := 500000
	expectedRBFOption := RBFOption(AllowedToBeRBF)

	// when
	data := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, AllowedToBeRBF)
	for _, utxo := range utxos {
		data.TransactionData.AddUTXO(utxo)
	}
	success, err := data.Generate()

	// then
	if !success {
		t.Errorf("Expected to generate transaction data, got error: %v", err)
	}
	if data.TransactionData.Amount != paymentAmount {
		t.Errorf("Expected amount to be %v, got %v", paymentAmount, data.TransactionData.Amount)
	}
	if data.TransactionData.FeeAmount != expectedFeeAmount {
		t.Errorf("Expected fee amount to be %v, got %v", expectedFeeAmount, data.TransactionData.FeeAmount)
	}
	if data.TransactionData.ChangeAmount != expectedChangeAmount {
		t.Errorf("Expected change amount to be %v, got %v", expectedChangeAmount, data.TransactionData.ChangeAmount)
	}
	if data.TransactionData.utxoCount() != expectedNumberOfUTXOs {
		t.Errorf("Expected number of UTXOs to be %v, got %v", expectedNumberOfUTXOs, data.TransactionData.utxoCount())
	}
	if data.TransactionData.Locktime != expectedLocktime {
		t.Errorf("Expected locktime to be %v, got %v", expectedLocktime, data.TransactionData.Locktime)
	}
	if data.TransactionData.shouldAddChangeToTransaction() {
		t.Errorf("Expected to not add change to transaction.")
	}
	if data.TransactionData.RBFOption != expectedRBFOption {
		t.Errorf("Expected RBFOption to be %v, got %v", expectedRBFOption, data.TransactionData.RBFOption)
	}
}

func TestNewTransactionStandard_SingleOutput_DoubleInput_InsufficientFunds(t *testing.T) {
	// given
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	ah := addressHelperTestHelpers()
	paymentAmount := 50000000 // 0.50000000 BTC
	utxoAmount1 := 20000000   // 0.20000000 BTC
	utxoAmount2 := 10000000   // 0.10000000 BTC
	changePath := NewDerivationPath(84, 0, 0, 1, 0)
	utxoPath := NewDerivationPath(49, 0, 0, 0, 0)
	utxo1 := NewUTXO("previous txid", 0, utxoAmount1, utxoPath, true)
	utxo2 := NewUTXO("previous txid", 1, utxoAmount2, utxoPath, true)
	utxos := []*UTXO{utxo1, utxo2}
	feeRate := 30

	// when
	data := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, AllowedToBeRBF)
	for _, utxo := range utxos {
		data.TransactionData.AddUTXO(utxo)
	}
	success, _ := data.Generate()

	// then
	if success {
		t.Error("Should have failed to create transaction data with insufficient funds.")
	}
}

func TestNewTransactionDataStandard_SingleBIP84Output_SingleBIP49Input(t *testing.T) {
	// given
	ah := addressHelperTestHelpers()
	address := "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4"
	paymentAmount := 50000000 // 0.5 BTC
	utxoAmount1 := 30000000   // 0.3 BTC
	utxoAmount2 := 30000000   // 0.3 BTC
	changePath := NewDerivationPath(84, 0, 0, 1, 0)
	utxoPath := NewDerivationPath(49, 0, 0, 0, 0)
	utxo1 := NewUTXO("previous txid", 0, utxoAmount1, utxoPath, true)
	utxo2 := NewUTXO("previous txid", 1, utxoAmount2, utxoPath, true)
	utxos := []*UTXO{utxo1, utxo2}
	feeRate := 30
	totalBytes, tbErr := ah.totalBytes(utxos, address, true)
	if tbErr != nil {
		t.Errorf("Expected to get total bytes from helper, got error: %v", tbErr)
	}
	expectedFeeAmount := feeRate * totalBytes // 7,680
	expectedChangeAmount := (utxoAmount1 + utxoAmount2) - paymentAmount - expectedFeeAmount
	expectedNumberOfUTXOs := len(utxos)
	expectedLocktime := 500000
	expectedRBFOption := RBFOption(AllowedToBeRBF)

	// when
	data := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, AllowedToBeRBF)
	for _, utxo := range utxos {
		data.TransactionData.AddUTXO(utxo)
	}
	success, err := data.Generate()

	// then
	if !success {
		t.Errorf("Expected to generate transaction data, got error: %v", err)
	}
	if totalBytes != 255 {
		t.Errorf("Expected total bytes to be 255, got %v", totalBytes)
	}
	if data.TransactionData.Amount != paymentAmount {
		t.Errorf("Expected amount to be %v, got %v", paymentAmount, data.TransactionData.Amount)
	}
	if data.TransactionData.FeeAmount != expectedFeeAmount {
		t.Errorf("Expected fee amount to be %v, got %v", expectedFeeAmount, data.TransactionData.FeeAmount)
	}
	if data.TransactionData.FeeAmount != 7650 {
		t.Errorf("Expected fee amount to be 7650, got %v", data.TransactionData.FeeAmount)
	}
	if data.TransactionData.ChangeAmount != expectedChangeAmount {
		t.Errorf("Expected change amount to be %v, got %v", expectedChangeAmount, data.TransactionData.ChangeAmount)
	}
	if data.TransactionData.utxoCount() != expectedNumberOfUTXOs {
		t.Errorf("Expected number of UTXOs to be %v, got %v", expectedNumberOfUTXOs, data.TransactionData.utxoCount())
	}
	if data.TransactionData.Locktime != expectedLocktime {
		t.Errorf("Expected locktime to be %v, got %v", expectedLocktime, data.TransactionData.Locktime)
	}
	if !data.TransactionData.shouldAddChangeToTransaction() {
		t.Errorf("Expected to add change to transaction.")
	}
	if data.TransactionData.RBFOption != expectedRBFOption {
		t.Errorf("Expected RBFOption to be %v, got %v", expectedRBFOption, data.TransactionData.RBFOption)
	}
}

func TestNewTransactionDataStandard_CostOfChangeIsBeneficial(t *testing.T) {
	// given
	ah := addressHelperTestHelpers()
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	path1 := NewDerivationPath(49, 0, 0, 1, 3)
	utxo1 := NewUTXO("909ac6e0a31c68fe345cc72d568bbab75afb5229b648753c486518f11c0d0009", 1, 100000, path1, true)
	path2 := NewDerivationPath(49, 0, 0, 0, 2)
	utxo2 := NewUTXO("419a7a7d27e0c4341ca868d0b9744ae7babb18fd691e39be608b556961c00ade", 0, 100000, path2, true)
	utxos := []*UTXO{utxo1, utxo2}
	changePath := NewDerivationPath(49, 0, 0, 1, 5)
	feeRate := 10
	totalBytes, tbErr := ah.totalBytes(utxos, address, false)
	if tbErr != nil {
		t.Errorf("Expected to get total bytes for tx, got error: %v", tbErr)
	}
	dustyChange := 1100
	expectedFeeAmount := feeRate*totalBytes + dustyChange
	paymentAmount := utxo1.Amount + utxo2.Amount - expectedFeeAmount
	expectedLocktime := 500000
	expectedRBFOption := RBFOption(AllowedToBeRBF)

	// when
	data := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, AllowedToBeRBF)
	for _, utxo := range utxos {
		data.TransactionData.AddUTXO(utxo)
	}
	success1, err1 := data.Generate()

	// then
	if !success1 {
		t.Errorf("Expected to generate transaction, got error: %v", err1)
	}
	if data.TransactionData.Amount != paymentAmount {
		t.Errorf("Expected amount to be %v, got %v", paymentAmount, data.TransactionData.Amount)
	}
	if data.TransactionData.FeeAmount != expectedFeeAmount {
		t.Errorf("Expected fee amount to be %v, got %v", expectedFeeAmount, data.TransactionData.FeeAmount)
	}
	if data.TransactionData.ChangeAmount != 0 {
		t.Errorf("Expected change amount to be %v, got %v", 0, data.TransactionData.ChangeAmount)
	}
	if data.TransactionData.utxoCount() != len(utxos) {
		t.Errorf("Expected number of UTXOs to be %v, got %v", len(utxos), data.TransactionData.utxoCount())
	}
	if data.TransactionData.Locktime != expectedLocktime {
		t.Errorf("Expected locktime to be %v, got %v", expectedLocktime, data.TransactionData.Locktime)
	}
	if data.TransactionData.shouldAddChangeToTransaction() {
		t.Errorf("Expected to not add change to transaction.")
	}
	if data.TransactionData.RBFOption != expectedRBFOption {
		t.Errorf("Expected RBFOption to be %v, got %v", expectedRBFOption, data.TransactionData.RBFOption)
	}

	// when again
	paymentAmount = 194000
	expectedFeeAmount = 2560
	expectedChange := 3440
	goodData := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, AllowedToBeRBF)
	for _, utxo := range utxos {
		goodData.TransactionData.AddUTXO(utxo)
	}
	success2, err2 := goodData.Generate()

	// then again
	if !success2 {
		t.Errorf("Expected transaction to be generated, got error: %v", err2)
	}
	if goodData.TransactionData.Amount != paymentAmount {
		t.Errorf("Expected amount to be %v, got %v", paymentAmount, goodData.TransactionData.Amount)
	}
	if goodData.TransactionData.FeeAmount != expectedFeeAmount {
		t.Errorf("Expected fee amount to be %v, got %v", expectedFeeAmount, goodData.TransactionData.FeeAmount)
	}
	if goodData.TransactionData.ChangeAmount != expectedChange {
		t.Errorf("Expected change amount to be %v, got %v", expectedChange, goodData.TransactionData.ChangeAmount)
	}
	if goodData.TransactionData.utxoCount() != len(utxos) {
		t.Errorf("Expected number of UTXOs to be %v, got %v", len(utxos), goodData.TransactionData.utxoCount())
	}
	if goodData.TransactionData.Locktime != expectedLocktime {
		t.Errorf("Expected locktime to be %v, got %v", expectedLocktime, goodData.TransactionData.Locktime)
	}
	if !goodData.TransactionData.shouldAddChangeToTransaction() {
		t.Errorf("Expected to add change to transaction.")
	}
	if goodData.TransactionData.RBFOption != expectedRBFOption {
		t.Errorf("Expected RBFOption to be %v, got %v", expectedRBFOption, goodData.TransactionData.RBFOption)
	}
}
