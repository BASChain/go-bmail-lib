pragma solidity >=0.5.0;
contract MailName{
    struct BMail{
        string Address;
        string CName;
    }
    mapping(bytes32=>BMail) public DNS;
    
    function setDns(bytes32 emailHash, 
                string calldata bmailAddress, 
                string calldata cName) external{
        DNS[emailHash] = BMail(bmailAddress, cName);
    }
}