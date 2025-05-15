<h2 align="center">Perun Cross Contract Example</h2>

This example shows how to set up a payment channel on Stellar and Ethereum, which utilizes the [go-perun](https://github.com/perun-network/go-perun) channel library, and also the [Ethereum payment channel backend](https://github.com/perun-network/perun-eth-backend) and the [Stellar payment channel backend](https://github.com/perun-network/perun-stellar-backend).

# Setup

Navigate to the payment-channel-cc directory and run the following commands to install the dependencies.

```sh
  go mod tidy
  go mod vendor
```
Spin up the local Stellar blockchain, serving as a local testnet for demonstration purposes.

```sh
  ./quickstart.sh standalone
```

This will start the Stellar, Horizon and Soroban nodes in the background. This is the platform on which we deploy the Stellar Asset Contract (SAC), and the Perun Payment Channel contract. This allows us to create and utilize L2 channels on Stellar for any customized Stellar asset tokens.

Install [ganache-cli](https://github.com/trufflesuite/ganache-cli) and open a new terminal to run
```sh
KEY_DEPLOYER=0x79ea8f62d97bc0591a4224c1725fca6b00de5b2cea286fe2e0bb35c5e76be46e
KEY_ALICE=0x1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f
KEY_BOB=0xf63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e
BALANCE=100000000000000000000

ganache -h 127.0.0.1 --port 8545 --wallet.accounts $KEY_DEPLOYER,$BALANCE $KEY_ALICE,$BALANCE $KEY_BOB,$BALANCE -b 5 
```
This starts a local ganache node with three prefunded accounts. The first account is used to deploy the contract, and the other two are used as Alice and Bob in the example.

# Using the example

You can start the demo by simply running

```sh
  go run .
```

The stellar accounts for Alice and Bob used in the example are generated randomly and funded at the initialization stage of the demo. 