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
//
// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.7.0;
pragma experimental ABIEncoderV2;

import "./CollateralAssetHolder.sol";

/**
 * @notice CollateralAssetHolderETH is implementation of CollateralAssetHolder
 * for the asset type Ether.
 */
contract CollateralAssetHolderETH is CollateralAssetHolder {

    /**
    * @notice Must be initialized with references to an adjudicator and a
    * collateral app.
    */
    constructor(address _adjudicator, address _app, uint256 _collateralWithdrawalDelay) CollateralAssetHolder(_adjudicator, _app, _collateralWithdrawalDelay) {} // solhint-disable-line no-empty-blocks

    // depositCheck checks that the amount is attached to the transaction.
    function depositCheck(bytes32, uint256 amount) internal override view {
        require(msg.value == amount, "wrong amount of ETH for deposit");
    }

    // withdrawEnact performs the authorized transaction.
    function withdrawEnact(WithdrawalAuth calldata authorization, bytes calldata) internal override {
        authorization.receiver.transfer(authorization.amount);
    }

    // withdrawCollateralEnact performs the collateral withdrawal transaction.
    function withdrawCollateralEnact(address payable receiver, uint256 amount) internal override {
        receiver.transfer(amount);
    }
}
