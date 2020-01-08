# Cnlib

Library used by Coin Ninja for use in DropBit app for iOS and Android. This library uses gomobile to produce frameworks consumable by iOS and Android, respectively. Dependencies used can be seen in each file's import statements.

## Getting Started

1. Clone the project
2. Ensure Go is installed (golang.org)
3. Use [gomobile](https://godoc.org/golang.org/x/mobile/cmd/gomobile) command relative to your needs

### Prerequisites

What things you need to install the software and how to install them
* Go >= 1.13.4
* macOS >= 10.14 Mojave (if using macOS)

## Running the tests

Run `go test`.

## Usage

Here are general guidelines for using the library.  

### Wallet
The wallet will require a BaseCoin object, which dictates what type of wallet the client is concerned with (i.e. SegWit/BIP84, Script Nested Segwit/BIP49, mainnet, regtest, etc). From the BaseCoin, the wallet will generate the necessary child keys and associated bitcoin addresses.  

### Transactions
To create a transaction, use one of three available constructors:  

* `NewTransactionDataStandard`  
* `NewTransactionDataFlatFee`  
* `NewTransactionDataSendingMax`  

Each type provides functions to add UTXOs (`AddUTXO(utxo *UTXO)`) as well as subsequently generating the needed transaction information to be broadcasted (`Generate()`). Once initialized, iterate over **all** available UTXOs in your wallet and add them all to the transaction data object:  

(pseudocode example)

```
for utxo in utxos {  
    data.AddUTXO(utxo)  
}  

err := data.Generate()
```

Once generated, the selected UTXOs needed to satisfy the amount + fee + change will be in an array called `requiredUtxos`. A client needing to get the required UTXO count selected for use in the transaction can call `data.utxoCount()`.  

A client is expected to broadcast the transaction on their own, so a function on the HDWallet type called `BuildTransactionMetadata` should be called with the transaction data's embedded `TransactionData` object, which will return the encoded transaction, associated txid, and any change information needed, if any.  

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for our contribution policy.  

Please read [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md) for details on our code of conduct.  

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags).

## Authors

* **BJ Miller** - *Initial work, client implementation considerations* - [Coin Ninja](https://coinninja.com)
* **Zach Brown** - *Encryption/Decryption, CI/gomobile implementation, and consultation* - [Coin Ninja](https://coinninja.com)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
