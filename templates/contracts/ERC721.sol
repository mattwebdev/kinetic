// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

contract {{.ContractName}} is ERC721, ERC721URIStorage, ERC721Enumerable, Ownable {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIds;

    {{#if .HasMaxSupply}}
    uint256 public immutable maxSupply;
    {{/if}}

    {{#if .HasBaseURI}}
    string private _baseTokenURI;
    {{/if}}

    {{#if .IsPausable}}
    bool public paused;
    event Paused(address account);
    event Unpaused(address account);
    {{/if}}

    constructor(
        string memory name,
        string memory symbol{{#if .HasMaxSupply}},
        uint256 _maxSupply{{/if}}{{#if .HasBaseURI}},
        string memory baseURI{{/if}}
    ) ERC721(name, symbol) Ownable(msg.sender) {
        {{#if .HasMaxSupply}}
        require(_maxSupply > 0, "Max supply must be greater than 0");
        maxSupply = _maxSupply;
        {{/if}}
        {{#if .HasBaseURI}}
        _baseTokenURI = baseURI;
        {{/if}}
    }

    {{#if .IsMintable}}
    function mint(address to{{#if .HasCustomURI}}, string memory tokenURI{{/if}}) public {{#if .OnlyOwnerCanMint}}onlyOwner{{/if}} returns (uint256) {
        {{#if .HasMaxSupply}}
        require(_tokenIds.current() < maxSupply, "Max supply reached");
        {{/if}}
        {{#if .IsPausable}}
        require(!paused, "Minting is paused");
        {{/if}}

        _tokenIds.increment();
        uint256 newTokenId = _tokenIds.current();
        _safeMint(to, newTokenId);
        
        {{#if .HasCustomURI}}
        _setTokenURI(newTokenId, tokenURI);
        {{/if}}

        return newTokenId;
    }
    {{/if}}

    {{#if .IsBurnable}}
    function burn(uint256 tokenId) public {
        require(_isApprovedOrOwner(msg.sender, tokenId), "Caller is not owner nor approved");
        _burn(tokenId);
    }
    {{/if}}

    {{#if .HasBaseURI}}
    function _baseURI() internal view virtual override returns (string memory) {
        return _baseTokenURI;
    }

    function setBaseURI(string memory newBaseURI) public onlyOwner {
        _baseTokenURI = newBaseURI;
    }
    {{/if}}

    {{#if .IsPausable}}
    modifier whenNotPaused() {
        require(!paused, "Contract is paused");
        _;
    }

    modifier whenPaused() {
        require(paused, "Contract is not paused");
        _;
    }

    function pause() public onlyOwner whenNotPaused {
        paused = true;
        emit Paused(msg.sender);
    }

    function unpause() public onlyOwner whenPaused {
        paused = false;
        emit Unpaused(msg.sender);
    }
    {{/if}}

    // Override required functions
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 tokenId,
        uint256 batchSize
    ) internal {{#if .IsPausable}}whenNotPaused{{/if}} override(ERC721, ERC721Enumerable) {
        super._beforeTokenTransfer(from, to, tokenId, batchSize);
    }

    function _burn(uint256 tokenId) internal override(ERC721, ERC721URIStorage) {
        super._burn(tokenId);
    }

    function tokenURI(uint256 tokenId) public view override(ERC721, ERC721URIStorage) returns (string memory) {
        return super.tokenURI(tokenId);
    }

    function supportsInterface(bytes4 interfaceId) public view override(ERC721, ERC721Enumerable) returns (bool) {
        return super.supportsInterface(interfaceId);
    }
} 