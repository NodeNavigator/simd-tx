# Simd Transaction Example (with example of ethscep256k1 algo)

`Simd` is the full node client or daemon built by [SimApp](https://github.com/cosmos/cosmos-sdk/blob/v0.50.9/simapp/README.md), which is an application built using the Cosmos SDK for testing and educational purposes. This repository contains sample code to demonstrate how to wrap message(s) in a transaction.

## Version

| Type | Version |
|----------------|---------|
| Cosmos SDK  | v0.50.9     |
| CometBFT    | v0.38.12     |

## Configuration

This project necessitates a configuration file named `config.toml`. To create one, duplicate the `example.toml` and customize the field values according to your preferences. You can locate the configuration source code in the `/config/config.go file`.

## Usage

```bash
go run main.go
```
