// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract HelloWorld {
    string public message;

    uint[] public numbers;
    
    constructor(string memory _message) {
        message = _message;
    }
    function updateMessage(string memory _message) public {
        message = _message;
    }
    function getMessage() public view returns (string memory) {
        return message;
    }

    function addNumbers(uint[] calldata _nums) external {
        _nums.length;
        for(uint i = 0 ; i < _nums.length; i++){
            numbers.push(_nums[i]);
        }

    }

    function getNumbersSum() external view returns (uint sum) {
        uint length = numbers.length;
        uint total = 0;
        for(uint i = 0 ; i < length; i++){
            total += numbers[i];
        }
        return total;
    }
}