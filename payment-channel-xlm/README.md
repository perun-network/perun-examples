<h2 align="center">Perun Stellar Example</h2>

This example shows how to set up a payment channel on Stellar, which utilizes the [go-perun](https://github.com/perun-network/go-perun) channel library, and also the [Stellar payment channel backend](https://github.com/perun-network/perun-stellar-backend).

# Setup

Spin up the local Stellar blockchain, serving as a local testnet for demonstration purposes.

```
  $ ./quickstart.sh standalone
```

This will start the Stellar, Horizon and Soroban nodes in the background. This is the platform on which we deploy the Stellar Asset Contract (SAC), and the Perun Payment Channel contract. This allows us to create and utilize L2 channels on Stellar for any customized Stellar asset tokens.

# Using the example

You can start the demo by simply running

```
  $ go run main.go
```

The accounts for Alice and Bob used in the example are generated randomly and funded at the initialization stage of the demo. 