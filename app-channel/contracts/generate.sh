#!/bin/sh

set -e

if [ -z "$ABIGEN" ]
then
    ABIGEN=abigen
fi

if [ -z "$SOLC" ]
then
    SOLC=solc
fi

# Generate golang bindings from solidity contract
# Argument 1: solidity contract file
# Argument 2: golang contract name (used for package and file)
generate_bindings() {
    CONTRACT_SOL_FILE=$1
    CONTRACT_GO_NAME=$2
    PKG=$CONTRACT_GO_NAME
    GENDIR=./generated/$PKG
    mkdir -p $GENDIR
    $ABIGEN --pkg $PKG --sol $CONTRACT_SOL_FILE --out $GENDIR/$CONTRACT_GO_NAME.go --solc $SOLC
}

generate_bindings ./TicTacToeApp.sol ticTacToeApp
# generate_bindings ./perun-eth-contracts/contracts/Adjudicator.sol adjudicator
# generate_bindings ./perun-eth-contracts/contracts/AssetHolderETH.sol assetHolderETH
