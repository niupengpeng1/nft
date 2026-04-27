// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "hardhat/console.sol";
import "./MultiSigWallet.sol";

contract MyVersion_v1 is MultiSigWallet {
    uint256 public version;
    uint256 public x;
    

    function initialize() public initializer {
        version = 1;
    }

    function reinitialize(uint64 _version,address[] calldata _owners,uint256 _num) public reinitializer(_version) {
        console.log("reinitialize");
        version = _version;
        require(_owners.length > 0, "owners required");
        require(
            _num > 0 &&
                _num <= _owners.length,
            "invalid number of required confirmations"
        );

        for (uint256 i = 0; i < _owners.length; i++) {
            address owner = _owners[i];
            owners.push(owner);
            isOwner[owner] = true;
            require(owner != address(0), "invalid owner");
        }
        numConfirmationsRequired = _num;
    }
}
