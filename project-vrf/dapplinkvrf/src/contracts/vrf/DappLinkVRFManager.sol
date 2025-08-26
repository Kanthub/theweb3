// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import "@openzeppelin-upgrades/contracts/proxy/utils/Initializable.sol";
import "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";

import "./DappLinkVRFStorage.sol";
import "../../interface/IDappLinkVRFManager.sol";

contract DappLinkVRFManager is
    Initializable,
    OwnableUpgradeable,
    DappLinkVRFStorage,
    IDappLinkVRFManager
{
    event RequestSent(uint256 requestId, uint256 _numWords, address current);

    event FillRandomWords(uint256 requestId, uint256[] randomWords);

    modifier onlyDappLink() {
        require(msg.sender == dappLinkAddress, "DappLinkVRF.onlyDappLink");
        _;
    }

    constructor() {
        _disableInitializers();
    }

    function initialize(
        address initialOwner,
        address _dappLinkAddress,
        address _blsRegistry
    ) public initializer {
        __Ownable_init(initialOwner);
        blsRegistry = IBLSApkRegistry(_blsRegistry);
        dappLinkAddress = _dappLinkAddress;
    }

    function requestRandomWords(
        uint256 _requestId,
        uint256 _numWords
    ) external onlyOwner {
        requestMapping[_requestId] = RequestStatus({
            randomWords: new uint256[](0),
            fulfilled: false
        });
        requestIds.push(_requestId);
        lastRequestId = _requestId;
        emit RequestSent(_requestId, _numWords, address(this));
    }

    function fulfillRandomWords(
        uint256 _requestId,
        uint256[] memory _randomWords,
        bytes32 msgHash,
        uint256 referenceBlockNumber,
        IBLSApkRegistry.VrfNoSignerAndSignature memory params
    ) external onlyDappLink {
        blsRegistry.checkSignatures(msgHash, referenceBlockNumber, params);
        requestMapping[_requestId] = RequestStatus({
            fulfilled: true,
            randomWords: _randomWords
        });
        emit FillRandomWords(_requestId, _randomWords);
    }

    function getRequestStatus(
        uint256 _requestId
    ) external view returns (bool fulfilled, uint256[] memory randomWords) {
        return (
            requestMapping[_requestId].fulfilled,
            requestMapping[_requestId].randomWords
        );
    }

    function setDappLink(address _dappLinkAddress) external onlyOwner {
        dappLinkAddress = _dappLinkAddress;
    }
}
