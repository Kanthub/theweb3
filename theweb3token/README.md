## The Web3 社区课程之发行自己的第一个 ERC20 代币

## Build

```shell
forge build
```
### Anvil

```shell
anvil
```

### Deploy

```shell
forge script ./script/TheWebThree.s.sol:TreasureManagerScript --rpc-url 127.0.0.1:8545 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 --broadcast
```

### 查看代币发行的总量
```shell
cast call --rpc-url 127.0.0.1:8545 0x5FC8d32690cc91D4c39d9d3abcBD16989F875707 "totalSupply()(uint256)"
```

### 执行 mint 
```shell
cast send --rpc-url 127.0.0.1:8545  --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 0x5FC8d32690cc91D4c39d9d3abcBD16989F875707 "mint(address,uint256)" 0xee2E207D30383430a815390431298EBa3c1C8c2d 100000000
```

### 查看地址里面石有该代币
```shell
cast call --rpc-url 127.0.0.1:8545 0x5FC8d32690cc91D4c39d9d3abcBD16989F875707 "balanceOf(address)(uint256)" 0xee2E207D30383430a815390431298EBa3c1C8c2d
```

### 转账
```shell
cast send --rpc-url 127.0.0.1:8545  --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 0x5FC8d32690cc91D4c39d9d3abcBD16989F875707 "transfer(address,uint256)" 0xee2E207D30383430a815390431298EBa3c1C8c2d 100000000```

