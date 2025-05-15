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

import "./perun-eth-contracts/contracts/Channel.sol";
import "./App.sol";
import "./CollateralAssetHolder.sol";

/**
 * @notice CollateralApp is a channel app that enables collateralized channels.
 * A collateralized channel is a partially funded channel that is backed by
 * collateral. The channel balances are tracked as signed integers in the
 * appData field of the channel state. A positive asset balance corresponds to a
 * credit and a negative balance corresponds to a debit.
 */
contract CollateralApp is App {

    // The adjudicator that is allowed to call this contract.
    address public adjudicator;

    /**
     * @notice The onlyAdjudicator modifier ensures that the decorated function
     * can only be called from the defined app contract.
     */
    modifier onlyAdjudicator {
        require(msg.sender == adjudicator, "can only be called by the defined adjudicator"); // solhint-disable-line reason-string
        _;
    }

    /**
     * @notice Sets the adjudicator contract that is allowed to call this
     * contract.
     */
    constructor(address _adjudicator) {
        adjudicator = _adjudicator;
    }

    /**
     * @notice onConclude is called when an app channel is concluded. It settles
     * the channel balances of each asset against the asset holdings.
     */
    function onConclude(
        Channel.Params calldata params,
        Channel.State calldata state)
    external override onlyAdjudicator
    {
        int256[][] memory bals = abi.decode(state.appData, (int256[][]));
        address[] memory assets = state.outcome.assets;
        for (uint256 i = 0; i < assets.length; i++) {
            CollateralAssetHolder a = CollateralAssetHolder(assets[i]);
            a.settleChannel(state.channelID, bals[i], params.participants);
        }
    }

    /**
     * @notice ValidTransition checks if there was a valid transition between two states.
     * @param params The parameters of the channel.
     * @param from The current state.
     * @param to The potenrial next state.
     * @param actorIdx Index of the actor who signed this transition.
     */
    function validTransition(
        Channel.Params calldata params,
        Channel.State calldata from,
        Channel.State calldata to,
        uint256 actorIdx)
    external pure override
    {
        // We do not allow on-chain progression of collateralized channels.
        require(false, "on-chain progression prohibited");
    }
}
