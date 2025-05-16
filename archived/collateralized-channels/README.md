# Collateralized Channels (Archived ⚠️)

## Dependencies
This example requires the [ganache-cli](https://github.com/trufflesuite/ganache-cli) `v6.12.2` and [go-perun](https://github.com/hyperledger-labs/go-perun) `v0.6.0` . Newer versions of the dependencies might not be compatible.

## Run the example

Install [ganache-cli](https://github.com/trufflesuite/ganache-cli) `v6.12.2` and run
```
go test -v perun.network/perun-collateralized-channels
```

The test starts a ganache-cli in the background. If the test does not shut down in order, make sure to kill the corresponding `node` process manually.
