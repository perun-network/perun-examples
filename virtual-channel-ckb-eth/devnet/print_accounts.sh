#!/bin/bash

echo "Extract SUDT owner lock-hash into own file"
awk '/lock_hash:/ { print $2 }' ./accounts/genesis-1.txt > ./accounts/sudt-owner-lock-hash1.txt
awk '/lock_hash:/ { print $2 }' ./accounts/genesis-2.txt > ./accounts/sudt-owner-lock-hash2.txt
