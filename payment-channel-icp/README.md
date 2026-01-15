# payment-channel on internet computer

This demo connects to our [Canister](https://github.com/perun-network/perun-icp-canister) that runs on a by using our [ICP Backend](https://github.com/perun-network/perun-icp-backend) to enables a Perun Payment Channel.

## Example Walkthrough


1. Install [dfx](https://internetcomputer.org/docs/current/references/cli-reference/dfx-parent), the DFINITY command-line execution environment:
```sh
sh -ci "$(curl -fsSL https://internetcomputer.org/install.sh)"
```

2. Run the deployment script to start local devnet and deploy Perun Canister.
```sh
cd payment-channel-icp
./startdeploy.sh
```
3. Start the example.
```sh
go run .
```

4. Stop the the local blockchain with.
```sh
./stopdfx.sh
``` 