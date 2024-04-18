# payment-channel on polkadot

This demo connects to our [Pallet] that runs on a [Polkadot Node] by using our [Polkadot Backend].

## Example Walkthrough

In another terminal, clone and start a local [Polkadot Node]: 

**WARNING**: only working with a newer version of Polkadot Node from branch `update_substrate`

To build the Polkadot Node, you need [Rust Developer Environment].

```sh
git clone git@github.com:perun-network/perun-polkadot-node.git

cd perun-polkadot-node/node

cargo run --release -- --dev --tmp
```

Wait around 10 seconds for the Node to start-up.

In a second terminal, run the demo
```sh
cd payment-channel-dot/
go run .
```

<!-- Links -->
[Pallet]: https://github.com/perun-network/perun-polkadot-pallet/
[Polkadot Node]: https://github.com/perun-network/perun-polkadot-node
[Polkadot Backend]: https://github.com/perun-network/perun-polkadot-backend
[Rust Developer Environment]: https://github.com/perun-network/perun-polkadot-node/blob/update_substrate/node/doc/rust-setup.md