# Simple Client (Archived ⚠️)

This example demonstrates how a simple payment channel client is created with [go-perun](https://github.com/hyperledger-labs/go-perun).

A corresponding tutorial can be found at https://tutorial.perun.network/.

## Dependency
This example requires [ganache-cli](https://github.com/trufflesuite/ganache-cli) `v6.12.2` and [go-perun](https://github.com/hyperledger-labs/go-perun) `v0.6.0` . Newer versions of the dependencies might not be compatible.

## Execution

Running the example requires [golang](https://golang.org) and [ganache-cli](https://github.com/trufflesuite/ganache-cli).
In one terminal window, start ganache-cli with:
```sh
ganache-cli -b 5 -a 2 -m "pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic"
```
In a second terminal window, navigate to the example directory and enter:
```sh
go run .
```