// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;
import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

contract MultiSigWalletProxy is TransparentUpgradeableProxy{

    constructor(
        address _logic,
        address initialOwner,
        bytes memory   _data
    )  TransparentUpgradeableProxy( _logic,  initialOwner,  _data){

    }

}