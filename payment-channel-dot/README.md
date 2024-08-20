# payment-channel on polkadot

This demo connects to our [Pallet] that runs on a [Polkadot Node] by using our [Polkadot Backend].

## Example Walkthrough

In a first terminal, start a development [Polkadot Node]:
```sh
docker run --rm -it -p 9944:9944 ghcr.io/perun-network/polkadot-test-node
```

In a second terminal, run the demo
```sh
cd payment-channel-dot/
go run .
```

<!-- Links -->
[Pallet]: https://github.com/perun-network/perun-polkadot-pallet/
[Polkadot Node]: https://github.com/perun-network/perun-polkadot-node
[Polkadot Backend]: https://github.com/perun-network/perun-polkadot-backend