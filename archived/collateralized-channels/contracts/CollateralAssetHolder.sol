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

import "./perun-eth-contracts/vendor/openzeppelin-contracts/contracts/math/SafeMath.sol";
import "./AssetHolder.sol";
import "./App.sol";

/**
 * @notice CollateralAssetHolder is an abstract assetholder for channels that
 * are backed by collateral.
 */
abstract contract CollateralAssetHolder is AssetHolder {
    using SafeMath for uint256;

    address public app;

    event CollateralOverdrawn(address indexed peer, bytes32 channelID, uint256 amount);
    event CollateralWithdrawn(address peer, uint256 amount);

    /**
     * @notice The onlyApp modifier ensures that the decorated function can only
     * be called from the defined app contract.
     */
    modifier onlyApp {
        require(msg.sender == app, "can only be called by the defined app"); // solhint-disable-line reason-string
        _;
    }

    /**
     * @notice Sets the adjudicator and app contracts that are able to this
     * contract.
     */
    constructor(address _adjudicator, address _app, uint256 _collateralWithdrawalDelay) AssetHolder(_adjudicator) {
        app = _app;
        collateralWithdrawalDelay = _collateralWithdrawalDelay;
    }

    /**
     * @notice settleChannel settles channel balances against channel funding
     * and peer collateral.
     */
    function settleChannel(bytes32 channelID, int256[] memory bals, address[] memory parts) external onlyApp {
        // Settle debit against funding and collateral.
        for (uint256 i = 0; i < bals.length; i++) {
            address peer = parts[i];
            int256 bal = bals[i];
            if (bal < 0) {
                uint256 debit = uint256(bal * -1); // The debit is the absolute value of a negative balance.
                uint256 leftoverDebit = settleChannelDebit(channelID, peer, debit);
                // If we could not settle the debit, return without touching the asset holdings.
                if (leftoverDebit > 0) {
                    return;
                }
            }
        }

        // Settle credit against holdings.
        for (uint256 i = 0; i < bals.length; i++) {
            address peer = parts[i];
            int256 bal = bals[i];
            if (bal > 0) {
                bytes32 fundingID = calcFundingID(channelID, peer);
                holdings[fundingID] = holdings[fundingID].add(uint256(bal));
            }
        }
    }
    
    /**
     * @notice settleChannelDebit settles the peer's channel debit against its
     * collateral. Returns leftover debit amount.
     */
    function settleChannelDebit(bytes32 channelID, address peer, uint256 debit) internal returns (uint256) {
        bytes32 collateralID = calcCollateralID(peer);
        bytes32 fundingID = calcFundingID(channelID, peer);
        uint256 funding = holdings[fundingID];
        uint256 collateral = holdings[collateralID];

        // If channel funding is insufficient, we try to take the missing amount
        // from the peer's collateral.
        if (funding < debit) {
            uint256 missingAmount = debit.sub(funding);

            // If the collateral is not sufficient, we cannot settle. We notify
            // and return.
            if (missingAmount > collateral) {
                emit CollateralOverdrawn(peer, channelID, missingAmount);
                return missingAmount;
            }
            
            // The collateral is sufficient. We move the asset from the
            // collateral to the funding.
            holdings[collateralID] = collateral.sub(missingAmount);
            holdings[fundingID] = funding.add(missingAmount);
            emit CollateralWithdrawn(peer, missingAmount);
        }

        // There is sufficient funding available. We settle against the funding.
        holdings[fundingID] = holdings[fundingID].sub(debit);
        return 0;
    }

    /**
     * @notice calcCollateralID calculates the collateralID for a peer.
     */
    function calcCollateralID(address participant) internal pure returns (bytes32) {
        return keccak256(abi.encode(participant));
    }

    struct CollateralWithdrawal {
        uint256 amount;
        uint256 registered;
    }

    // collateralWithdrawals holds information about ongoing collateral
    // withdrawals.
    mapping(address => CollateralWithdrawal) public collateralWithdrawals;

    // collateralWithdrawalDelay is the minimum duration between the
    // registration of a collateral withdrawal and its execution.
    uint256 public collateralWithdrawalDelay;

    /**
     * @notice registerCollateralWithdrawal registers the withdrawal of
     * collateral.
     */
    function registerCollateralWithdrawal(uint256 amount) public {
        address client = msg.sender;
        bytes32 id = calcCollateralID(client);
        require(holdings[id] >= amount, "insufficient collateral");
        collateralWithdrawals[client] = CollateralWithdrawal({
            amount: amount,
            registered: block.timestamp
        });
    }

    /**
     * @notice performCollateralWithdrawal executes the registered collateral
     * withdrawal if the withdrawal delay has passed.
     */
    function performCollateralWithdrawal() external {
        address payable client = msg.sender;
        bytes32 id = calcCollateralID(client);
        CollateralWithdrawal memory withdrawal = collateralWithdrawals[client];
        require(withdrawal.registered.add(collateralWithdrawalDelay) >= block.timestamp, "delay not passed");
        withdrawCollateralEnact(client, withdrawal.amount);
        holdings[id] = holdings[id].sub(withdrawal.amount);
        emit CollateralWithdrawn(client, withdrawal.amount);
    }

    function withdrawCollateralEnact(address payable receiver, uint256 amount) internal virtual
    {} // solhint-disable no-empty-blocks
}
