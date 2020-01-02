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
	err = data.Generate()

	// then
	assert.Nil(t, err)
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
	err = data.Generate()

	// then
	assert.Nil(t, err)
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
	err = data.Generate()

	// then
	assert.Nil(t, err)
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
	err = data.Generate()

	// then
	assert.Nil(t, err)
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
	err := data.Generate()

	// then
	assert.EqualError(t, errors.New("insufficient funds"), err.Error())
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
	err = data.Generate()

	// then
	assert.Nil(t, err)
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
	err = data.Generate()

	// then
	assert.Nil(t, err)
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
	err = goodData.Generate()

	// then again
	assert.Nil(t, err)
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
	err := data.Generate()

	// then
	assert.Nil(t, err)
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
	err := data.Generate()

	// then
	assert.Nil(t, err)
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
	err = data.Generate()

	// then
	assert.Nil(t, err)
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
	err = data.Generate()

	// then
	assert.Nil(t, err)
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
	err := data.Generate()

	// then
	assert.EqualError(t, errors.New("insufficient funds"), err.Error())
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
	err = data.Generate()

	// then
	assert.Nil(t, err)
	assert.Equal(t, 224, totalBytes)
	assert.Equal(t, address, data.TransactionData.PaymentAddress)
	assert.Equal(t, expectedAmount, data.TransactionData.Amount)
	assert.Equal(t, expectedFeeAmount, data.TransactionData.FeeAmount)
	assert.Equal(t, 0, data.TransactionData.ChangeAmount)
	assert.False(t, data.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, expectedRBFOption.Value, data.TransactionData.RBFOption.Value)
}

func TestNewTransactionDataStandard_TwoSegwitInputs_TwoSegwitOutputs(t *testing.T) {
	ah := addressHelperTestHelpers()
	address := "bc1q2myn4sqfwcjdgn8xqpeuq77277gj5ngmda5uk8"
	feeRate := 1
	path1 := NewDerivationPath(84, 0, 0, 0, 15)
	path2 := NewDerivationPath(84, 0, 0, 1, 19)
	changePath := NewDerivationPath(84, 0, 0, 1, 20)
	utxo1 := NewUTXO("ca470899cad4aa48487e5cabb6abd387b0ff7a4ef380d3544a6a738f3c101e37", 0, 13770, path1, nil, true)
	utxo2 := NewUTXO("16ce8aaf23d15f3440e4369600a3004e47ca0940d4756eb45a655c538dcaaa4a", 1, 197171, path2, nil, true)
	utxos := []*UTXO{utxo1, utxo2}
	inputAmount := utxo1.Amount + utxo2.Amount
	paymentAmount := 200000
	totalBytes, err := ah.totalBytes(utxos, address, true)
	assert.Nil(t, err)
	assert.Equal(t, 209, totalBytes)

	expectedFeeAmount := feeRate * totalBytes // 209
	expectedChangeAmount := 10732             //196791
	expectedAmount := inputAmount - expectedFeeAmount - expectedChangeAmount
	expectedRBFOption := NewRBFOption(AllowedToBeRBF)

	// when
	data := NewTransactionDataStandard(address, ah.Basecoin, paymentAmount, feeRate, changePath, 610518, expectedRBFOption)
	data.AddUTXO(utxo1)
	data.AddUTXO(utxo2)
	err = data.Generate()

	// then
	assert.Nil(t, err)
	assert.Equal(t, address, data.TransactionData.PaymentAddress)
	assert.Equal(t, expectedAmount, data.TransactionData.Amount)
	assert.Equal(t, expectedFeeAmount, data.TransactionData.FeeAmount)
	assert.Equal(t, expectedChangeAmount, data.TransactionData.ChangeAmount)
	assert.True(t, data.TransactionData.shouldAddChangeToTransaction())
	assert.Equal(t, expectedRBFOption.Value, data.TransactionData.RBFOption.Value)

	wallet := NewHDWalletFromWords(w, ah.Basecoin)
	metadata, err := wallet.BuildTransactionMetadata(data.TransactionData)
	assert.Nil(t, err)
	expectedTxid := "4683df1447daec29bfab1514803304b722f4890cbdbaaec0f9cdfd7bc74681ca"
	expectedEncodedTx := "01000000000102371e103c8f736a4a54d380f34e7affb087d3abb6ab5c7e4848aad4ca990847ca0000000000ffffffff4aaaca8d535c655ab46e75d44009ca474e00a3009636e440345fd123af8ace160100000000ffffffff02400d03000000000016001456c93ac0097624d44ce60073c07bcaf7912a4d1bec290000000000001600145b8585924dc44505ed40d8a127e792fa4e68cbfd02483045022100d05e99f619084e76edcd04595af4e0a31bb05efa9d9cab831578d63e8a388442022044e43fb1b4df85e97fe2cfe7d9fb7bc922af6e0516fbfa0d51ca5686db01b5a9012102b05e67ab098575526f23a7c4f3b69449125604c34a9b34909def7432a792fbf60248304502210088213160aa8b43fdee2fbcc8da497fdca8e4adc5f9028b01cf59f019af502c3c02202bbe894e35391befc91ae4fefb5afa63e01fb02ab326670e56864ea20facd3dc012103020d7c261fb5c6103a8f8f4c73b3fbed228c981869e68b6e9c6f6973b0550659d6500900"
	assert.Equal(t, expectedTxid, metadata.Txid)
	assert.Equal(t, expectedEncodedTx, metadata.EncodedTx)
}
