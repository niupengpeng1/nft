// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "hardhat/console.sol";

contract MyNFT  is ERC721 ,ERC721URIStorage ,Ownable {
    
    uint256 public totalSupply = 0;

    uint256 public price = 0.01 ether;

    event Mint(address indexed to, uint256 indexed tokenId);

    constructor() ERC721("MyNFT", "MNFT") Ownable(msg.sender) {      
    }

    
    function mint(string memory url) public payable returns (uint256) {
        require(msg.value > 0, "Not enough ETH sent!");
        uint256 newTokenId =  totalSupply++;
        _safeMint(_msgSender(), newTokenId);
        console.log("Minted NFT with tokenId %s for %s", newTokenId, msg.sender);
        _setTokenURI(newTokenId, url);

        emit Mint(_msgSender(), newTokenId);
        return newTokenId;
    }




 /**
     * @dev 重写tokenURI函数
     * @param tokenId Token ID
     * @return 元数据URI
     * @notice 需要重写以解决多重继承的冲突
     */
    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }
    
    /**
     * @dev 检查接口支持
     * @param interfaceId 接口ID
     * @return 是否支持该接口
     * @notice 实现ERC165标准，支持接口查询
     */
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }

    function withdraw() public onlyOwner {
        payable(owner()).transfer(address(this).balance);
    }

    function updatePrice(uint256 _price) external  onlyOwner {
        price = _price;
    }

}