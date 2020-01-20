package cnlib

import "errors"

import "github.com/btcsuite/btcd/wire"

/// Type Definitions

// Following constants are used for RBFOption.
const (
	MustBeRBF      int = 0
	MustNotBeRBF   int = 1
	AllowedToBeRBF int = 2
)

// PlaceholderDestination is a constant which can be used to indicate a destination is not yet selected, but tx size needs to be estimated.
const PlaceholderDestination = "---placeholder---"

const dustThreshold = 1000

// RBFOption is a struct wrapping an int for RBF preferred value. Value should be `MustBeRBF` (0), `MustNotBeRBF` (1), or `AllowedToBeRBF` (2).
type RBFOption struct {
	Value int
}

// NewRBFOption returns a pointer to RBFOption.
func NewRBFOption(value int) *RBFOption {
	return &RBFOption{Value: value}
}

// TransactionData is the main object containing all info necessary to build a bitcoin transaction.
// Will retain references to all pointers, no need to carry on externally.
type TransactionData struct {
	PaymentAddress string
	availableUtxos []*UTXO
	requiredUtxos  []*UTXO
	basecoin       *BaseCoin
	Amount         int
	FeeAmount      int
	feeRate        int
	ChangeAmount   int
	ChangePath     *DerivationPath
	Locktime       int
	RBFOption      *RBFOption
}

// TransactionDataStandard adopts the Transaction interface, customizing the generation of the transaction.
type TransactionDataStandard struct {
	TransactionData *TransactionData
}

// TransactionDataFlatFee adopts the Transaction interface, customizing the generation of the transaction.
type TransactionDataFlatFee struct {
	TransactionData *TransactionData
}

// TransactionDataSendMax adopts the Transaction interface, customizing the generation of the transaction.
type TransactionDataSendMax struct {
	TransactionData *TransactionData
}

/// Constructors

/*
NewTransactionDataStandard Create transaction data object using a fee rate, calculating fee via number of inputs and outputs.

Once created, add all available utxos one at a time using `addUTXO` function, as gomobile does not support custom arrays/slices. This method will select the ones needed.

@param paymentAddress The address to which you want to send currency.
@param coin The coin representing the current user's wallet.
@param amount The amount which you would like to send to the receipient.
@param feeRate The fee rate to be multiplied by the estimated transaction size.
@param changePath The derivative path for receiving change, if any. Retains reference.
@param blockHeight The current block height, used to calculate the locktime (blockHeight + 1).
@param rbfOption A ref to a RBFOption object passed to the transaction builder to determind replaceability. Retains reference.
@return Returns an instantiated object if fully able to satisfy amount+fee with UTXOs, or nil if insufficient funds.
*/
func NewTransactionDataStandard(
	paymentAddress string,
	basecoin *BaseCoin,
	amount int,
	feeRate int,
	changePath *DerivationPath,
	blockHeight int,
	rbfOption *RBFOption,
) *TransactionDataStandard {
	td := TransactionData{
		PaymentAddress: paymentAddress,
		availableUtxos: []*UTXO{},
		requiredUtxos:  []*UTXO{},
		basecoin:       basecoin,
		Amount:         amount,
		FeeAmount:      0,
		feeRate:        feeRate,
		ChangeAmount:   0,
		ChangePath:     changePath,
		Locktime:       blockHeight,
		RBFOption:      rbfOption,
	}
	tsd := TransactionDataStandard{TransactionData: &td}

	return &tsd
}

/*
NewTransactionDataFlatFee Create transaction data object with a flat fee, versus calculated via number of inputs/outputs.

Once created, add all available utxos one at a time using `addUTXO` function, as gomobile does not support custom arrays/slices. This method will select the ones needed.

Default RBFOption is MustBeRBF.

@param paymentAddress The address to which you want to send currency.
@param coin The coin representing the current user's wallet.
@param amount The amount which you would like to send to the receipient.
@param flatFee The flat-fee to pay, NOT a rate. This fee, added to amount, will equal the total deducted from the wallet.
@param changePath The derivative path for receiving change, if any. Retains reference.
@param blockHeight The current block height, used to calculate the locktime (blockHeight + 1).
@return Returns an instantiated object if fully able to satisfy amount+fee with UTXOs, or nil if insufficient funds.
*/
func NewTransactionDataFlatFee(
	paymentAddress string,
	basecoin *BaseCoin,
	amount int,
	flatFee int,
	changePath *DerivationPath,
	blockHeight int,
) *TransactionDataFlatFee {
	rbf := NewRBFOption(MustBeRBF)
	td := TransactionData{
		PaymentAddress: paymentAddress,
		availableUtxos: []*UTXO{},
		requiredUtxos:  []*UTXO{},
		basecoin:       basecoin,
		Amount:         amount,
		FeeAmount:      flatFee,
		feeRate:        0,
		ChangeAmount:   0,
		ChangePath:     changePath,
		Locktime:       blockHeight,
		RBFOption:      rbf,
	}
	tdff := TransactionDataFlatFee{TransactionData: &td}
	return &tdff
}

/*
NewTransactionDataSendingMax Send max amount to a given address, minus the calculated fee based on size of transaction times feeRate.

Once created, add all available utxos one at a time using `addUTXO` function, as gomobile does not support custom arrays/slices. This method will select the ones needed.

Default RBFOption is MustNotBeRBF.

@param paymentAddress The address to which you want to send currency.
@param coin The coin representing the current user's wallet.
@param feeRate The fee rate to be multiplied by the estimated transaction size.
@param blockHeight The current block height, used to calculate the locktime (blockHeight + 1).
@return Returns an instantiated object if fully able to satisfy amount+fee with UTXOs, or nil if insufficient funds. This would only be
nil if the funding amount is less than the fee.
*/
func NewTransactionDataSendingMax(
	paymentAddress string,
	basecoin *BaseCoin,
	feeRate int,
	blockHeight int,
) *TransactionDataSendMax {
	rbf := NewRBFOption(MustNotBeRBF)
	td := TransactionData{
		PaymentAddress: paymentAddress,
		availableUtxos: []*UTXO{},
		requiredUtxos:  []*UTXO{},
		basecoin:       basecoin,
		Amount:         0,
		FeeAmount:      0,
		feeRate:        feeRate,
		ChangeAmount:   0,
		ChangePath:     nil,
		Locktime:       blockHeight,
		RBFOption:      rbf,
	}
	tdsm := TransactionDataSendMax{TransactionData: &td}
	return &tdsm
}

/// Receiver Functions

// AddUTXO Adds a utxo to the private array.
func (td *TransactionData) AddUTXO(utxo *UTXO) {
	td.availableUtxos = append(td.availableUtxos, utxo)
}

// RequiredUTXOAtIndex returns a utxo that has been selected to be included in the outgoing transaction, or error if out of bounds.
func (td *TransactionData) RequiredUTXOAtIndex(index int) (*UTXO, error) {
	if index < 0 {
		return nil, errors.New("index must be greater than 0")
	}

	if index > len(td.requiredUtxos)-1 {
		return nil, errors.New("index must be within range of utxos")
	}

	return td.requiredUtxos[index], nil
}

// AddUTXO Adds a utxo to the private array.
func (t *TransactionDataStandard) AddUTXO(utxo *UTXO) {
	t.TransactionData.AddUTXO(utxo)
}

// AddUTXO Adds a utxo to the private array.
func (t *TransactionDataFlatFee) AddUTXO(utxo *UTXO) {
	t.TransactionData.AddUTXO(utxo)
}

// AddUTXO Adds a utxo to the private array.
func (t *TransactionDataSendMax) AddUTXO(utxo *UTXO) {
	t.TransactionData.AddUTXO(utxo)
}

// Generate is called after all available utxo's have been added, to configure the transaction data. Builds a standard transaction with a fee rate.
func (t *TransactionDataStandard) Generate() error {

	err := t.TransactionData.validate()
	if err != nil {
		t.TransactionData = nil
		return err
	}

	totalFromUTXOs := 0
	totalSendingValue := 0
	currentFee := 0
	tempUTXOs := make([]*UTXO, 0)

	for i := 0; i < len(t.TransactionData.availableUtxos); i++ {
		utxo := t.TransactionData.availableUtxos[i]
		bytes, err := t.TransactionData.basecoin.bytesPerInput(utxo)
		if err != nil {
			t.TransactionData = nil
			return err
		}
		feePerInput := t.TransactionData.feeRate * bytes
		totalSendingValue = t.TransactionData.Amount + currentFee

		if totalSendingValue > totalFromUTXOs {
			tempUTXOs = append(tempUTXOs, utxo)
			totalFromUTXOs += utxo.Amount
			totalBytes, err := t.TransactionData.basecoin.totalBytes(tempUTXOs, t.TransactionData.PaymentAddress, false)
			if err != nil {
				return err
			}
			currentFee = t.TransactionData.feeRate * totalBytes
			totalSendingValue = t.TransactionData.Amount + currentFee

			changeValue := totalFromUTXOs - totalSendingValue

			if (totalFromUTXOs < totalSendingValue) || (changeValue < 0) {
				continue
			}

			if (changeValue > 0) && (changeValue < (feePerInput + dustThreshold)) {
				// it is not beneficial to add change, would just dust self with change
				currentFee += changeValue
				break
			} else if changeValue > 0 {
				estBytes, err := t.TransactionData.basecoin.totalBytes(tempUTXOs, t.TransactionData.PaymentAddress, true)
				if err != nil {
					return err
				}
				totalBytes = estBytes
				currentFee = t.TransactionData.feeRate * totalBytes
				changeValue = totalFromUTXOs - t.TransactionData.Amount - currentFee
				t.TransactionData.ChangeAmount = changeValue
				break
			} else if changeValue < 0 {
				currentFee += changeValue
				changeValue = 0
				t.TransactionData.ChangeAmount = changeValue
			}
		} else {
			break
		}
	}

	t.TransactionData.FeeAmount = currentFee
	t.TransactionData.requiredUtxos = tempUTXOs

	if totalFromUTXOs < totalSendingValue {
		return errors.New("insufficient funds")
	}

	return nil
}

// Generate is called after all available utxo's have been added, to configure the transaction data. Builds a standard transaction with a flat fee.
func (t *TransactionDataFlatFee) Generate() error {

	err := t.TransactionData.validate()
	if err != nil {
		t.TransactionData = nil
		return err
	}

	totalFromUTXOs := 0
	tempUTXOs := make([]*UTXO, 0)

	for i := 0; i < len(t.TransactionData.availableUtxos); i++ {
		utxo := t.TransactionData.availableUtxos[i]
		tempUTXOs = append(tempUTXOs, utxo)
		totalFromUTXOs += utxo.Amount

		possibleChange := totalFromUTXOs - t.TransactionData.Amount - t.TransactionData.FeeAmount
		tempChangeAmount := Max(0, possibleChange)
		t.TransactionData.ChangeAmount = tempChangeAmount

		if totalFromUTXOs >= t.TransactionData.Amount && tempChangeAmount > 0 {
			if tempChangeAmount < dustThreshold {
				t.TransactionData.ChangeAmount = 0
			}
		}
	}

	if totalFromUTXOs < (t.TransactionData.FeeAmount + t.TransactionData.Amount) {
		return errors.New("insufficient funds")
	}

	t.TransactionData.requiredUtxos = tempUTXOs

	return nil
}

// Generate is called after all available utxo's have been added, to configure the transaction data. Builds a transaction sending max with a fee rate.
func (t *TransactionDataSendMax) Generate() error {
	tempUTXOs := t.TransactionData.availableUtxos
	totalFromUTXOs := 0
	for _, utxo := range t.TransactionData.availableUtxos {
		totalFromUTXOs += utxo.Amount
	}

	totalBytes, err := t.TransactionData.basecoin.totalBytes(tempUTXOs, t.TransactionData.PaymentAddress, false)
	if err != nil {
		return err
	}

	feeAmount := t.TransactionData.feeRate * totalBytes
	amountForValidation := totalFromUTXOs - feeAmount
	if amountForValidation < 0 {
		return errors.New("insufficient funds")
	}
	t.TransactionData.Amount = amountForValidation
	t.TransactionData.FeeAmount = feeAmount
	t.TransactionData.requiredUtxos = tempUTXOs

	err = t.TransactionData.validate()
	if err != nil {
		t.TransactionData = nil
		return err
	}

	return nil
}

// UtxoCount returns count of UTXOs required to satisfy the transaction, not all UTXOs passed in before calling `Generate`.
func (td *TransactionData) UtxoCount() int {
	return len(td.requiredUtxos)
}

/// Unexported Functions

func (td *TransactionData) shouldAddChangeToTransaction() bool {
	return td.ChangeAmount > 0
}

func (td *TransactionData) getSuggestedSequence() uint32 {
	if td.RBFOption.Value == MustBeRBF {
		return wire.MaxTxInSequenceNum - 2
	}
	if td.RBFOption.Value == MustNotBeRBF {
		return wire.MaxTxInSequenceNum
	}
	if td.RBFOption.Value == AllowedToBeRBF {
		includesUnconfirmedUTXOs := false
		for _, utxo := range td.requiredUtxos {
			includesUnconfirmedUTXOs = includesUnconfirmedUTXOs || !utxo.IsConfirmed
		}
		if includesUnconfirmedUTXOs {
			return wire.MaxTxInSequenceNum - 2
		}
		return wire.MaxTxInSequenceNum
	}
	return wire.MaxTxInSequenceNum
}

func (td *TransactionData) validate() error {
	if td.Amount < 546 {
		return errors.New("transaction too small")
	}
	return nil
}
