#!/bin/sh

set -e

# Define ABIGEN and SOLC default values.
ABIGEN=abigen
SOLC=solc

echo 'Please ensure that solc v0.8.15+ and abigen 1.10.18+ are installed.'

if ! $ABIGEN --version
then
    echo "'abigen' not found. Please add to PATH or set ABIGEN='path_to_abigen'."
    exit 1
fi

if ! $SOLC --version
then
    echo "'solc' not found. Please add to PATH or set SOLC='path_to_solc'."
    exit 1
fi

echo "Please ensure that the repository was cloned with submodules: 'git submodule update --init --recursive'."

# Generate golang bindings from solidity contract
# Argument 1: solidity contract file
# Argument 2: golang contract name (used for package and file)
generate_bindings() {
    FILE=$1
    PKG=$2
    CONTRACT=$FILE
    GENDIR=./generated/$PKG

    echo "Generating $PKG bindings..."
    rm -r $GENDIR
    mkdir -p $GENDIR

    # Compile and generate binary runtime
    $SOLC --abi --bin --optimize $FILE.sol -o $GENDIR

    # Generate bindings.
    $ABIGEN --pkg $PKG --abi $GENDIR/$FILE.abi --bin $GENDIR/$FILE.bin --out $GENDIR/$CONTRACT.go
}

generate_bindings "TicTacToeApp" "ticTacToeApp"
# generate_bindings ./perun-eth-contracts/contracts/Adjudicator.sol adjudicator
# generate_bindings ./perun-eth-contracts/contracts/AssetHolderETH.sol assetHolderETH
