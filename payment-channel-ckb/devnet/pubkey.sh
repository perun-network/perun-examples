#!/bin/bash

# Ignore stderr, as it will print some config information.
ckb-cli util key-info --privkey-path $1 2>/dev/null | awk '/pubkey/ { print $2 }' 
