package cnlib

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	utxo := NewUTXO("previous txid", 0, utxoAmount, utxoPath, nil, true)
	utxos := []*UTXO{utxo}
	feeRate := 30
	totalBytes, err := addressHelper().totalBytes(utxos, address, true)
	assert.Nil(t, err)

	expectedFeeAmount := feeRate * totalBytes // 4,980
	expectedChangeAmount := (utxoAmount - paymentAmount - expectedFeeAmount)
	expectedNumberOfUTXOs := 1
	expectedLocktime := 500000

	// when
	rbf := NewRBFOption(MustBeRBF)
	data := NewTransactionDataStandard(
		address, addressHelper().Basecoin, paymentAmount, feeRate, changePath, 500000, rbf,
	)
	data.AddUTXO(utxo)
	success, err := data.Generate()

	// then
	assert.Nil(t, err)
	assert.True(t, success)
	assert.Equal(t, paymentAmount, data.TransactionData.Amount)
	assert.Equal(t, expectedFeeAmount, data.TransactionData.FeeAmount)
	assert.Equal(t, 49995020, expectedChangeAmount)
	assert.Equal(t, expectedChangeAmount, data.TransactionData.ChangeAmount)
	assert.Equal(t, expectedNumberOfUTXOs, data.TransactionData.UtxoCount())
	assert.Equal(t, expectedLocktime, data.TransactionData.Locktime)
	assert.True(t, data.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, rbf.Value, data.TransactionData.RBFOption.Value)
}

func TestTransactionDataStandard_SingleOutput_DoubleInput_WithChange(t *testing.T) {
	// given
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	ah := addressHelperTestHelpers() // helper uses 84 purpose to ensure correct input size is calculated
	paymentAmount := 50000000        // 0.5 BTC
	utxoAmount := 30000000           // 0.3 BTC
	changePath := NewDerivationPath(84, 0, 0, 1, 0)
	utxoPath := NewDerivationPath(49, 0, 0, 0, 0)
	utxo1 := NewUTXO("previous txid", 0, utxoAmount, utxoPath, nil, true)
	utxo2 := NewUTXO("previous txid", 1, utxoAmount, utxoPath, nil, true)
	utxos := []*UTXO{utxo1, utxo2}
	feeRate := 30
	totalBytes, err := ah.totalBytes(utxos, address, true)
	assert.Nil(t, err)

	expectedFeeAmount := feeRate * totalBytes // 7,680
	amountFromUTXOs := 0
	for _, utxo := range utxos {
		amountFromUTXOs += utxo.Amount
	}
	expectedChangeAmount := amountFromUTXOs - paymentAmount - expectedFeeAmount
	expectedNumberOfUTXOs := 2
	expectedLocktime := 500000
	expectedRBFOption := NewRBFOption(AllowedToBeRBF)

	// when
	data := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, expectedRBFOption)
	for _, utxo := range utxos {
		data.AddUTXO(utxo)
	}
	success, err := data.Generate()

	// then
	assert.Nil(t, err)
	assert.True(t, success)
	assert.Equal(t, paymentAmount, data.TransactionData.Amount)
	assert.Equal(t, expectedFeeAmount, data.TransactionData.FeeAmount)
	assert.Equal(t, expectedChangeAmount, data.TransactionData.ChangeAmount)
	assert.Equal(t, expectedNumberOfUTXOs, data.TransactionData.UtxoCount())
	assert.Equal(t, expectedLocktime, data.TransactionData.Locktime)
	assert.True(t, data.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, expectedRBFOption.Value, data.TransactionData.RBFOption.Value)
}

func TestNewTransactionDataStandard_SingleInput_SingleOutput_NoChange(t *testing.T) {
	// given
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	ah := addressHelperTestHelpers() // helper uses 84 purpose to ensure correct input size is calculated
	paymentAmount := 50000000        // 0.5 BTC
	utxoAmount := 50004020           // 0.50004020 BTC
	changePath := NewDerivationPath(84, 0, 0, 1, 0)
	utxoPath := NewDerivationPath(49, 0, 0, 0, 0)
	utxo := NewUTXO("previous txid", 0, utxoAmount, utxoPath, nil, true)
	utxos := []*UTXO{utxo}
	feeRate := 30
	totalBytes, err := ah.totalBytes(utxos, address, false)
	assert.Nil(t, err)

	expectedFeeAmount := feeRate * totalBytes // 4,020
	amountFromUTXOs := 0
	for _, utxo := range utxos {
		amountFromUTXOs += utxo.Amount
	}
	expectedChangeAmount := amountFromUTXOs - paymentAmount - expectedFeeAmount
	expectedNumberOfUTXOs := 1
	expectedLocktime := 500000
	expectedRBFOption := NewRBFOption(AllowedToBeRBF)

	// when
	data := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, expectedRBFOption)
	for _, utxo := range utxos {
		data.AddUTXO(utxo)
	}
	success, err := data.Generate()

	// then
	assert.Nil(t, err)
	assert.True(t, success)
	assert.Equal(t, paymentAmount, data.TransactionData.Amount)
	assert.Equal(t, expectedFeeAmount, data.TransactionData.FeeAmount)
	assert.Equal(t, expectedChangeAmount, data.TransactionData.ChangeAmount)
	assert.Equal(t, expectedNumberOfUTXOs, data.TransactionData.UtxoCount())
	assert.Equal(t, expectedLocktime, data.TransactionData.Locktime)
	assert.False(t, data.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, expectedRBFOption.Value, data.TransactionData.RBFOption.Value)
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
	utxo1 := NewUTXO("previous txid", 0, utxoAmount1, utxoPath, nil, true)
	utxo2 := NewUTXO("previous txid", 1, utxoAmount2, utxoPath, nil, true)
	utxos := []*UTXO{utxo1, utxo2}
	feeRate := 30
	totalBytes, err := ah.totalBytes(utxos, address, false)
	assert.Nil(t, err)

	expectedFeeAmount := feeRate * totalBytes // 6, 750
	amountFromUTXOs := 0
	for _, utxo := range utxos {
		amountFromUTXOs += utxo.Amount
	}
	expectedChangeAmount := 0
	expectedNumberOfUTXOs := len(utxos)
	expectedLocktime := 500000
	expectedRBFOption := NewRBFOption(AllowedToBeRBF)

	// when
	data := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, expectedRBFOption)
	for _, utxo := range utxos {
		data.AddUTXO(utxo)
	}
	success, err := data.Generate()

	// then
	assert.Nil(t, err)
	assert.True(t, success)
	assert.Equal(t, paymentAmount, data.TransactionData.Amount)
	assert.Equal(t, expectedFeeAmount, data.TransactionData.FeeAmount)
	assert.Equal(t, expectedChangeAmount, data.TransactionData.ChangeAmount)
	assert.Equal(t, expectedNumberOfUTXOs, data.TransactionData.UtxoCount())
	assert.Equal(t, expectedLocktime, data.TransactionData.Locktime)
	assert.False(t, data.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, expectedRBFOption.Value, data.TransactionData.RBFOption.Value)
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
	utxo1 := NewUTXO("previous txid", 0, utxoAmount1, utxoPath, nil, true)
	utxo2 := NewUTXO("previous txid", 1, utxoAmount2, utxoPath, nil, true)
	utxos := []*UTXO{utxo1, utxo2}
	feeRate := 30
	rbf := NewRBFOption(AllowedToBeRBF)

	// when
	data := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, rbf)
	for _, utxo := range utxos {
		data.AddUTXO(utxo)
	}
	success, err := data.Generate()

	// then
	assert.EqualError(t, errors.New("insufficient funds"), err.Error())
	assert.False(t, success)
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
	utxo1 := NewUTXO("previous txid", 0, utxoAmount1, utxoPath, nil, true)
	utxo2 := NewUTXO("previous txid", 1, utxoAmount2, utxoPath, nil, true)
	utxos := []*UTXO{utxo1, utxo2}
	feeRate := 30
	totalBytes, err := ah.totalBytes(utxos, address, true)
	assert.Nil(t, err)

	expectedFeeAmount := feeRate * totalBytes // 7,680
	expectedChangeAmount := (utxoAmount1 + utxoAmount2) - paymentAmount - expectedFeeAmount
	expectedNumberOfUTXOs := len(utxos)
	expectedLocktime := 500000
	expectedRBFOption := NewRBFOption(AllowedToBeRBF)

	// when
	data := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, expectedRBFOption)
	for _, utxo := range utxos {
		data.AddUTXO(utxo)
	}
	success, err := data.Generate()

	// then
	assert.Nil(t, err)
	assert.True(t, success)
	assert.Equal(t, 255, totalBytes)
	assert.Equal(t, paymentAmount, data.TransactionData.Amount)
	assert.Equal(t, expectedFeeAmount, data.TransactionData.FeeAmount)
	assert.Equal(t, 7650, data.TransactionData.FeeAmount)
	assert.Equal(t, expectedChangeAmount, data.TransactionData.ChangeAmount)
	assert.Equal(t, expectedNumberOfUTXOs, data.TransactionData.UtxoCount())
	assert.Equal(t, expectedLocktime, data.TransactionData.Locktime)
	assert.True(t, data.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, expectedRBFOption.Value, data.TransactionData.RBFOption.Value)
}

func TestNewTransactionDataStandard_CostOfChangeIsBeneficial(t *testing.T) {
	// given
	ah := addressHelperTestHelpers()
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	path1 := NewDerivationPath(49, 0, 0, 1, 3)
	utxo1 := NewUTXO("909ac6e0a31c68fe345cc72d568bbab75afb5229b648753c486518f11c0d0009", 1, 100000, path1, nil, true)
	path2 := NewDerivationPath(49, 0, 0, 0, 2)
	utxo2 := NewUTXO("419a7a7d27e0c4341ca868d0b9744ae7babb18fd691e39be608b556961c00ade", 0, 100000, path2, nil, true)
	utxos := []*UTXO{utxo1, utxo2}
	changePath := NewDerivationPath(49, 0, 0, 1, 5)
	feeRate := 10
	totalBytes, err := ah.totalBytes(utxos, address, false)
	assert.Nil(t, err)

	dustyChange := 1100
	expectedFeeAmount := feeRate*totalBytes + dustyChange
	paymentAmount := utxo1.Amount + utxo2.Amount - expectedFeeAmount
	expectedLocktime := 500000
	expectedRBFOption := NewRBFOption(AllowedToBeRBF)

	// when
	data := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, expectedRBFOption)
	for _, utxo := range utxos {
		data.AddUTXO(utxo)
	}
	success1, err := data.Generate()

	// then
	assert.Nil(t, err)
	assert.True(t, success1)
	assert.Equal(t, paymentAmount, data.TransactionData.Amount)
	assert.Equal(t, expectedFeeAmount, data.TransactionData.FeeAmount)
	assert.Equal(t, 0, data.TransactionData.ChangeAmount)
	assert.Equal(t, len(utxos), data.TransactionData.UtxoCount())
	assert.Equal(t, expectedLocktime, data.TransactionData.Locktime)
	assert.False(t, data.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, expectedRBFOption.Value, data.TransactionData.RBFOption.Value)

	// when again
	paymentAmount = 194000
	expectedFeeAmount = 2560
	expectedChange := 3440
	goodData := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 500000, expectedRBFOption)
	for _, utxo := range utxos {
		goodData.AddUTXO(utxo)
	}
	success2, err := goodData.Generate()

	// then again
	assert.Nil(t, err)
	assert.True(t, success2)
	assert.Equal(t, paymentAmount, goodData.TransactionData.Amount)
	assert.Equal(t, expectedFeeAmount, goodData.TransactionData.FeeAmount)
	assert.Equal(t, expectedChange, goodData.TransactionData.ChangeAmount)
	assert.Equal(t, len(utxos), goodData.TransactionData.UtxoCount())
	assert.Equal(t, expectedLocktime, goodData.TransactionData.Locktime)
	assert.True(t, goodData.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, expectedRBFOption.Value, goodData.TransactionData.RBFOption.Value)
}

func TestNewTransactionDataFlatFee_WithChange(t *testing.T) {
	// given
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	ah := addressHelper()
	path1 := NewDerivationPath(49, 0, 0, 1, 3)
	path2 := NewDerivationPath(49, 0, 0, 0, 2)
	path3 := NewDerivationPath(49, 0, 0, 0, 8)
	utxo1 := NewUTXO("909ac6e0a31c68fe345cc72d568bbab75afb5229b648753c486518f11c0d0009", 1, 2221, path1, nil, true)
	utxo2 := NewUTXO("419a7a7d27e0c4341ca868d0b9744ae7babb18fd691e39be608b556961c00ade", 0, 15935, path2, nil, true)
	utxo3 := NewUTXO("3013fcd9ea8fd65a69709f07fed2c1fd765d57664486debcb72ef47f2ea415f6", 0, 15526, path3, nil, true)
	utxos := []*UTXO{utxo1, utxo2, utxo3}
	changePath := NewDerivationPath(49, 0, 0, 1, 5)
	paymentAmount := 20000
	flatFeeAmount := 10000
	expectedChange := 3682
	expectedRBFOption := NewRBFOption(MustBeRBF)

	// when
	data := NewTransactionDataFlatFee(address, ah.Basecoin, paymentAmount, flatFeeAmount, changePath, 500000)
	for _, utxo := range utxos {
		data.AddUTXO(utxo)
	}
	success, err := data.Generate()

	// then
	assert.Nil(t, err)
	assert.True(t, success)
	assert.Equal(t, address, data.TransactionData.PaymentAddress)
	assert.Equal(t, paymentAmount, data.TransactionData.Amount)
	assert.Equal(t, flatFeeAmount, data.TransactionData.FeeAmount)
	assert.Equal(t, expectedChange, data.TransactionData.ChangeAmount)
	assert.Equal(t, len(utxos), data.TransactionData.UtxoCount())
	assert.Equal(t, 500000, data.TransactionData.Locktime)
	assert.True(t, data.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, expectedRBFOption.Value, data.TransactionData.RBFOption.Value)
}

func TestNewTransactionDataFlatFee_DustyTransaction_NoChange(t *testing.T) {
	// given
	ah := addressHelper()
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	path1 := NewDerivationPath(49, 0, 0, 1, 3)
	path2 := NewDerivationPath(49, 0, 0, 0, 2)
	utxo1 := NewUTXO("909ac6e0a31c68fe345cc72d568bbab75afb5229b648753c486518f11c0d0009", 1, 20000, path1, nil, true)
	utxo2 := NewUTXO("419a7a7d27e0c4341ca868d0b9744ae7babb18fd691e39be608b556961c00ade", 0, 10100, path2, nil, true)
	utxos := []*UTXO{utxo1, utxo2}
	changePath := NewDerivationPath(49, 0, 0, 1, 5)
	paymentAmount := 20000
	expectedFeeAmount := 10000
	expectedChange := 0
	expectedRBFOption := NewRBFOption(MustBeRBF)

	// when
	data := NewTransactionDataFlatFee(address, ah.Basecoin, paymentAmount, expectedFeeAmount, changePath, 500000)
	for _, utxo := range utxos {
		data.AddUTXO(utxo)
	}
	success, err := data.Generate()

	// then
	assert.Nil(t, err)
	assert.True(t, success)
	assert.Equal(t, address, data.TransactionData.PaymentAddress)
	assert.Equal(t, paymentAmount, data.TransactionData.Amount)
	assert.Equal(t, expectedFeeAmount, data.TransactionData.FeeAmount)
	assert.Equal(t, len(utxos), data.TransactionData.UtxoCount())
	assert.Equal(t, expectedChange, data.TransactionData.ChangeAmount)
	assert.False(t, data.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, expectedRBFOption.Value, data.TransactionData.RBFOption.Value)
}

func TestNewTransactionDataSendMax_UsesAllUTXOs_AmountIsTotalMinusFee(t *testing.T) {
	// given
	ah := addressHelper()
	feeRate := 5
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	path1 := NewDerivationPath(49, 0, 0, 1, 3)
	path2 := NewDerivationPath(49, 0, 0, 0, 2)
	utxo1 := NewUTXO("909ac6e0a31c68fe345cc72d568bbab75afb5229b648753c486518f11c0d0009", 1, 20000, path1, nil, true)
	utxo2 := NewUTXO("419a7a7d27e0c4341ca868d0b9744ae7babb18fd691e39be608b556961c00ade", 0, 10000, path2, nil, true)
	utxos := []*UTXO{utxo1, utxo2}
	inputAmount := utxo1.Amount + utxo2.Amount
	totalBytes, err := ah.totalBytes(utxos, address, false)
	assert.Nil(t, err)

	expectedFeeAmount := feeRate * totalBytes // 1,125
	expectedAmount := inputAmount - expectedFeeAmount
	expectedRBFOption := NewRBFOption(MustNotBeRBF)

	// when
	data := NewTransactionDataSendingMax(address, ah.Basecoin, feeRate, 500000)
	for _, utxo := range utxos {
		data.AddUTXO(utxo)
	}
	success, err := data.Generate()

	// then
	assert.Nil(t, err)
	assert.True(t, success)
	assert.Equal(t, address, data.TransactionData.PaymentAddress)
	assert.Equal(t, expectedAmount, data.TransactionData.Amount)
	assert.Equal(t, expectedFeeAmount, data.TransactionData.FeeAmount)
	assert.Equal(t, 0, data.TransactionData.ChangeAmount)
	assert.False(t, data.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, expectedRBFOption.Value, data.TransactionData.RBFOption.Value)
}

func TestNewTransactionDataSendMax_JustEnoughFunds(t *testing.T) {
	// given
	ah := addressHelper()
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	feeRate := 5
	path1 := NewDerivationPath(49, 0, 0, 1, 3)
	utxo := NewUTXO("909ac6e0a31c68fe345cc72d568bbab75afb5229b648753c486518f11c0d0009", 1, 670, path1, nil, true)
	utxos := []*UTXO{utxo}
	inputAmount := utxo.Amount
	totalBytes, err := ah.totalBytes(utxos, address, false)
	assert.Nil(t, err)

	expectedFeeAmount := feeRate * totalBytes         // 670
	expectedAmount := inputAmount - expectedFeeAmount // 0
	expectedRBFOption := NewRBFOption(MustNotBeRBF)

	// when
	data := NewTransactionDataSendingMax(address, ah.Basecoin, feeRate, 500000)
	for _, u := range utxos {
		data.AddUTXO(u)
	}
	success, err := data.Generate()

	// then
	assert.Nil(t, err)
	assert.True(t, success)
	assert.Equal(t, address, data.TransactionData.PaymentAddress)
	assert.Equal(t, expectedAmount, data.TransactionData.Amount)
	assert.Equal(t, 0, data.TransactionData.Amount)
	assert.Equal(t, expectedFeeAmount, data.TransactionData.FeeAmount)
	assert.Equal(t, 670, data.TransactionData.FeeAmount)
	assert.Equal(t, 0, data.TransactionData.ChangeAmount)
	assert.False(t, data.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, expectedRBFOption.Value, data.TransactionData.RBFOption.Value)
}

func TestNewTransactionDataSendMax_InsufficientFunds(t *testing.T) {
	// given
	ah := addressHelper()
	address := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	feeRate := 5
	path1 := NewDerivationPath(39, 0, 0, 1, 3)
	utxo := NewUTXO("909ac6e0a31c68fe345cc72d568bbab75afb5229b648753c486518f11c0d0009", 1, 100, path1, nil, true)
	utxos := []*UTXO{utxo}

	// when
	data := NewTransactionDataSendingMax(address, ah.Basecoin, feeRate, 500000)
	for _, u := range utxos {
		data.AddUTXO(u)
	}
	success, err := data.Generate()

	// then
	assert.EqualError(t, errors.New("insufficient funds"), err.Error())
	assert.False(t, success)
}

func TestNewTransactionDataSendMax_ToNativeSegwit(t *testing.T) {
	// given
	ah := addressHelper()
	address := "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4"
	feeRate := 5
	path1 := NewDerivationPath(49, 0, 0, 1, 3)
	path2 := NewDerivationPath(49, 0, 0, 0, 2)
	utxo1 := NewUTXO("909ac6e0a31c68fe345cc72d568bbab75afb5229b648753c486518f11c0d0009", 1, 20000, path1, nil, true)
	utxo2 := NewUTXO("419a7a7d27e0c4341ca868d0b9744ae7babb18fd691e39be608b556961c00ade", 0, 10000, path2, nil, true)
	utxos := []*UTXO{utxo1, utxo2}
	inputAmount := utxo1.Amount + utxo2.Amount
	totalBytes, err := ah.totalBytes(utxos, address, false) // 224
	assert.Nil(t, err)

	expectedFeeAmount := feeRate * totalBytes // 1,120
	expectedAmount := inputAmount - expectedFeeAmount
	expectedRBFOption := NewRBFOption(MustNotBeRBF)

	// when
	data := NewTransactionDataSendingMax(address, ah.Basecoin, feeRate, 500000)
	for _, utxo := range utxos {
		data.AddUTXO(utxo)
	}
	success, err := data.Generate()

	// then
	assert.Nil(t, err)
	assert.True(t, success)
	assert.Equal(t, 224, totalBytes)
	assert.Equal(t, address, data.TransactionData.PaymentAddress)
	assert.Equal(t, expectedAmount, data.TransactionData.Amount)
	assert.Equal(t, expectedFeeAmount, data.TransactionData.FeeAmount)
	assert.Equal(t, 0, data.TransactionData.ChangeAmount)
	assert.False(t, data.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, expectedRBFOption.Value, data.TransactionData.RBFOption.Value)
}
