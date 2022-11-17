// Copyright 2021 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.7.0;
pragma experimental ABIEncoderV2;

import "./perun-eth-contracts/contracts/App.sol";

/**
 * @notice TicTacToeApp is a channel app for playing tic tac toe.
 * The data is encoded as follows:
 * - data[0]: The index of the next actor.
 * - data[i], i in [1,10]: The value of field i. 0 means no tick, 1 means tick by player 1, 2 means tick by player 2.
 */
contract TicTacToeApp is App {
    uint8 constant actorDataIndex = 0;
    uint8 constant actorDataLength = 1;
    uint8 constant gridDataIndex = actorDataIndex + actorDataLength;
    uint8 constant gridDataLength = 9;
    uint8 constant appDataLength = gridDataIndex + gridDataLength; // Actor index + grid.
    uint8 constant numParts = 2;
    uint8 constant notSet = 0;
    uint8 constant firstPlayer = 1;
    uint8 constant secondPlayer = 2;

    /**
     * @notice ValidTransition checks if there was a valid transition between two states.
     * @param params The parameters of the channel.
     * @param from The current state.
     * @param to The potential next state.
     * @param signerIdx Index of the participant who signed this transition.
     */
    function validTransition(
        Channel.Params calldata params,
        Channel.State calldata from,
        Channel.State calldata to,
        uint256 signerIdx)
    external pure override
    {
        require(params.participants.length == numParts, "number of participants");

        uint8 actorIndex = uint8(from.appData[actorDataIndex]);
        require(to.appData.length == appDataLength, "data length");
        require(actorIndex == signerIdx, "actor not signer");
        require((actorIndex + 1) % numParts == uint8(to.appData[actorDataIndex]), "next actor");

        // Test valid action.
        bool changed = false;
        for (uint i = gridDataIndex; i < gridDataIndex + gridDataLength; i++) {
            require(uint8(to.appData[i]) <= 2, "grid value");
            if (to.appData[i] != from.appData[i]) {
                require(uint8(from.appData[i]) == notSet, "overwrite");
                require(!changed, "two actions");
                changed = true;
            }
        }

        // Test final state.
        (bool isFinal, bool hasWinner, uint8 winner) = checkFinal(to.appData);
        require(to.isFinal == isFinal, "final flag");
        Array.requireEqualAddressArray(to.outcome.assets, from.outcome.assets);
        Channel.requireEqualSubAllocArray(to.outcome.locked, from.outcome.locked);
        uint256[][] memory expectedBalances = from.outcome.balances;
        if (hasWinner) {
            uint8 loser = 1 - winner;
            expectedBalances = new uint256[][](expectedBalances.length);
            for (uint i = 0; i < expectedBalances.length; i++) {
                expectedBalances[i] = new uint256[](numParts);
                expectedBalances[i][winner] = from.outcome.balances[i][0] + from.outcome.balances[i][1];
                expectedBalances[i][loser] = 0;
            }
        }
        requireEqualUint256ArrayArray(to.outcome.balances, expectedBalances);
    }

    /// @dev Asserts that a and b are equal.
    function requireEqualAddressArray(
        address[] memory a,
        address[] memory b
    )
    internal
    pure
    {
        require(a.length == b.length, "address[]: unequal length");
        for (uint i = 0; i < a.length; i++) {
            require(a[i] == b[i], "address[]: unequal item");
        }
    }

    function checkFinal(bytes memory d) internal pure returns (bool isFinal, bool hasWinner, uint8 winner) {
        // 0 1 2
        // 3 4 5
        // 6 7 8

        // Check winner.
        uint8[3][8] memory winningRows = [
        [0, 1, 2], [3, 4, 5], [6, 7, 8], // horizontal
        [0, 3, 6], [1, 4, 7], [2, 5, 8], // vertical
        [0, 4, 8], [2, 4, 6]             // diagonal
        ];
        for (uint i = 0; i < winningRows.length; i++) {
            (bool ok, uint8 v) = sameValue(d, winningRows[i]);
            if (ok) {
                if (v == firstPlayer) {
                    return (true, true, 0);
                } else if (v == secondPlayer) {
                    return (true, true, 1);
                }
            }
        }

        // Check all set.
        for (uint i = 0; i < d.length; i++) {
            if (uint8(d[i]) != notSet) {
                return (false, false, 0);
            }
        }
        return (true, false, 0);
    }

    function sameValue(bytes memory d, uint8[3] memory gridIndices) internal pure returns (bool ok, uint8 v) {
        bytes1 first = d[gridDataIndex + gridIndices[0]];
        for (uint i = 1; i < gridIndices.length; i++) {
            if (d[gridDataIndex + gridIndices[i]] != first) {
                return (false, 0);
            }
        }
        return (true, uint8(first));
    }

    function requireEqualUint256ArrayArray(
        uint256[][] memory a,
        uint256[][] memory b
    )
    internal pure
    {
        require(a.length == b.length, "uint256[][]: unequal length");
        for (uint i = 0; i < a.length; i++) {
            Array.requireEqualUint256Array(a[i], b[i]);
        }
    }
}