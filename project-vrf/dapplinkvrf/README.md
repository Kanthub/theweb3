## The Web3 VRF 小型项目实战

学习目的：
- 链上如何使用随机数
- 理解 BLS 签名验证签名流程，

## Build
```
$ forge build
```

## Deploy
```shell
forge script ./script/DappLinkVRFDepoly.s.sol:DappLinkVRFDepolyScript --rpc-url https://ethereum-holesky-rpc.publicnode.com --private-key $PRIVATE_KEY --broadcast
```

## Call Function
```
cast send --rpc-url https://ethereum-holesky-rpc.publicnode.com --private-key $PRIVATE_KEY 0xF752D67b60fC154Ef3BAD976AA80Ff2FaD367946 "requestRandomWords(uint256,uint256)" 10000 3
```

## Address
```
  dapplink blsApkRegistry contract deployed at: 0x4d33d437ACA1c028B92963F09AeBe1c7F91e8922
  dapplink base contract deployed at: 0x1F933548fD9A8841eCb48F087871Ed805B8492b6
  DappLink Proxy Factory contract deployed at: 0x981C51c1678a8E61e24958988d24b9391870b0d8
  DappLink Proxy contract deployed at: 0xF752D67b60fC154Ef3BAD976AA80Ff2FaD367946
```
