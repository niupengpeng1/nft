// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.24;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";


import "hardhat/console.sol";

contract MultiSigWallet is Initializable {
    
    event Deposit( address indexed sender, uint256 amount);

    event SubmitTransaction(
        uint256 indexed txIndex,
        address indexed to,
        uint256 value,
        bytes data
    );

    event ConfirmTransaction(address indexed owner,uint256 indexed txIndex);
    event RevokerTransaction(address indexed owner,uint256 indexed txIndex);
    event ExecuteTransaction(uint256 indexed txIndex);
    event OwnerAdded(address indexed owner);
    event OwerRemoved(address indexed owner);
    event ThresholdChanged(uint256 indexed newThreshold );



    struct Transaction {
        address to;
        uint256 value;
        bytes data;
        bool executed;
        uint256 numConfirmations;
    }

    //状态变量
    address[] public owners;
    mapping(address => bool) public isOwner;
    uint256 public numConfirmationsRequired;

    Transaction[] public transactions;
    mapping(uint256 => mapping(address => bool)) public isConfirmed;


    modifier onlyOwner() {
        require(isOwner[msg.sender], "not owner");
        _;
    }
    modifier txExists(uint256 _txIndex) {
        require(_txIndex < transactions.length, "tx does not exist");
        _;
    }   
    modifier notConfirmed(uint256 _txIndex) {
        require(!isConfirmed[_txIndex][msg.sender], "tx already confirmed");
        _;
    }   
    modifier notExecuted(uint256 _txIndex) {
        require(!transactions[_txIndex].executed, "tx already executed");
        _;
    }   
    modifier notRevoked(uint256 _txIndex) {
        require(isConfirmed[_txIndex][msg.sender], "tx already revoked");
        _;
    }   


  /*   constructor(address[] memory _owners,uint256 _numConfirmationsRequired){
        require(_owners.length > 0, "owners required");
        require(
            _numConfirmationsRequired > 0 &&
                _numConfirmationsRequired <= _owners.length,
            "invalid number of required confirmations"
        );

        for (uint256 i = 0; i < _owners.length; i++) {
            address owner = _owners[i];
            owners.push(owner);
            require(owner != address(0), "invalid owner");
        }
    } */

    function initialize(address[] calldata _owners,uint256 _num) public initializer {
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
    }
    
    //新增用户
    function addOwner(address newOwner) public onlyOwner {
        require(newOwner != address(0),"Invalid address");
        require(!isOwner[newOwner],"Already an owner");

        isOwner[newOwner] = true;

        emit OwnerAdded(newOwner);
    }

    //delete owner
    function removerOwner(address owner) public onlyOwner {
        require(isOwner[owner],"Not an owner");

        isOwner[owner] = false;
        uint256 length = owners.length;
        for(uint256 i=0 ; i <length; i++){
            if(owners[i]==owner){
                owners[i] = owners[length-1];
                owners.pop();
                break;
            }
        }
    
        emit OwerRemoved(owner);
    }

    //change confrimation threshold
    function changeThreshold(uint256 _newThreshold  ) public onlyOwner {
        require(_newThreshold > 0 && _newThreshold <= owners.length,"Invalid threshold");
        numConfirmationsRequired = _newThreshold;
        emit ThresholdChanged(_newThreshold);
    }

    //提交提案
    function submitTransaction(
        address to,
        uint256 value,
        bytes memory data
    ) public onlyOwner { 
        uint256 txIndex = transactions.length;

        transactions.push(
            Transaction(
                {
                    to:to,
                    value:value,
                    data:data,
                    executed:false,
                    numConfirmations:0
                }
            )
        );

        emit SubmitTransaction(txIndex, to, value, data);
    }

    //获取填信息
    function getTransaction(uint256 _txIndex) public view txExists(_txIndex) returns(Transaction memory){

        return(transactions[_txIndex]);
    }

    //获取交易总数
    function getTransactionCount() public view returns(uint256){
        return transactions.length;
    }

    //确认提案
    function firmTransaction(uint256 _txIndex) public onlyOwner txExists(_txIndex) notConfirmed(_txIndex) notExecuted(_txIndex) { 
        Transaction storage transaction = transactions[_txIndex];
        transaction.numConfirmations +=1;
        isConfirmed[_txIndex][msg.sender] = true;

        emit ConfirmTransaction(msg.sender,_txIndex);
    }
    
    //取消确认状态
    function revokeConfirmation( uint256 _txIndex) public onlyOwner txExists(_txIndex) notRevoked(_txIndex) notExecuted(_txIndex) { 
        Transaction storage transaction = transactions[_txIndex];
        transaction.numConfirmations -=1;
        isConfirmed[_txIndex][msg.sender] = false;

        emit RevokerTransaction(msg.sender,_txIndex);
    }

    //执行交易
    function executeTransaction(uint256 _txIndex) public onlyOwner txExists(_txIndex) notExecuted(_txIndex) { 
        Transaction storage transaction = transactions[_txIndex];
        require(
            transaction.numConfirmations >= numConfirmationsRequired,
            "cannot execute tx"
        );
        transaction.executed = true;
        (bool success, ) = transaction.to.call{value: transaction.value}(
            transaction.data
        );
        require(success, "tx failed");

        emit ExecuteTransaction(_txIndex);
    }

    function getOwners() public view returns(address[] memory){
        return owners;
    }
    function getThreashold() public view returns(uint256){ 
        return numConfirmationsRequired;
    }
    function getOwnersCount() public view returns(uint256){ 
        return owners.length;
    }


    receive() external payable {
        if(msg.value > 0){
            emit Deposit(msg.sender,msg.value);
        }
    }

    fallback() external payable {
        if(msg.value > 0){
            emit Deposit(msg.sender,msg.value);
        }
    }
}