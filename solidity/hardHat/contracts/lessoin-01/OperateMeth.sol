// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract OperateMeth {
    uint public num;

    mapping(address => uint) public balances;

    mapping (bytes32=>bool) public isUsed;

    function add(uint _num) public {
        require(_num > 0," must >0 ");
        
        num += _num;
    }

    function sub(uint _num) public {
        num -= _num;
    }

    function mul(uint _num) public {
        num *= _num;
    }

    function div(uint _num) public {
        require(_num != 0, "Cannot divide by zero");
        num /= _num;
    }
}