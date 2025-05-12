# payment-channel-icp

## Requirement
Install [dfx](https://internetcomputer.org/docs/building-apps/getting-started/install) with maximum version v0.25.1.

```sh
sh -ci "$(curl -fsSL https://internetcomputer.org/install.sh)"
```

Choose customize installation and specify the version **`0.25.1`**.

## Example Walkthrough

Start a local development [Polkadot Node]: 

```sh
./startdeploy.sh
```

Wait around 10 seconds for the Node to start-up.

Run the demo
```sh
cd payment-channel-dot/
go run .
```

Stop the devnet
```sh
./stopdfx.sh
```