import { Interface } from "@ethersproject/abi";
import { ethers } from "ethers";
import { signTransaction } from "../wallet";
import BigNumber from "bignumber.js";
import * as dotenv from "dotenv";
dotenv.config();
import { CrossChainMessenger, MessageStatus } from "@eth-optimism/sdk";

// 配置私钥
const PRIVATE_KEY = process.env.PRIVATE_KEY;
if (!PRIVATE_KEY) {
  throw new Error("PRIVATE_KEY not set in .env");
}

// 配置sepolia URL, provider 和 wallet
const sepoliaRpcUrl = process.env.SEPOLIA_RPC_URL;
if (!sepoliaRpcUrl) {
  throw new Error("sepoliaRpcUrl is missing");
}
const sepoliaProvider = new ethers.providers.JsonRpcProvider(sepoliaRpcUrl);
const sepoliaWallet = new ethers.Wallet(PRIVATE_KEY, sepoliaProvider);
const sepoliaWallet1 = new ethers.Wallet(PRIVATE_KEY);

// 配置op-sepolia URL, provider 和 wallet
const opSepoliaRpcUrl = process.env.OP_SEPOLIA_RPC_URL;
if (!opSepoliaRpcUrl) {
  throw new Error("opSepoliaRpcUrl is missing");
}
const opProvider = new ethers.providers.JsonRpcProvider(opSepoliaRpcUrl);
const opWallet = new ethers.Wallet(PRIVATE_KEY, opProvider);

const walletAddress = opWallet.address;
console.log("wallet address", walletAddress);

describe("op-sepolia", () => {
  test("op-sepolia transfer", async () => {
    const nonce = await opProvider.getTransactionCount(walletAddress);
    console.log("OP Current nonce:", nonce);

    const rawHex = signTransaction({
      privateKey: PRIVATE_KEY?.slice(2),
      nonce: nonce,
      from: walletAddress,
      to: "0x79731E63A00fc987a507DccD6dF4612d7febf31B",
      gasLimit: 91000,
      maxFeePerGas: 327993150328,
      maxPriorityFeePerGas: 32799315032,
      gasPrice: 0,
      amount: "0.002",
      decimal: 18,
      chainId: 11155420,
      tokenAddress: "",
      callData: "",
    });

    console.log("rawHex: ", rawHex);

    const tx = await opProvider.sendTransaction(rawHex);
    console.log("Transaction hash:", tx.hash);

    const receipt = await tx.wait(1);
    console.log("Transaction confirmed in block:", receipt.blockNumber);
    console.log(
      "Transaction status:",
      receipt.status === 1 ? "Success" : "Failed" // 务必检查 receipt.status
    );
  });

  test("op-sepolia erc20 token transfer", async () => {
    const nonce = await opProvider.getTransactionCount(walletAddress);
    console.log("OP Current nonce:", nonce);

    const rawHex = signTransaction({
      privateKey: PRIVATE_KEY?.slice(2),
      nonce: nonce,
      from: walletAddress,
      to: "0x79731E63A00fc987a507DccD6dF4612d7febf31B",
      gasLimit: 91000,
      maxFeePerGas: 327993150328,
      maxPriorityFeePerGas: 32799315032,
      gasPrice: 0,
      amount: "0.002",
      decimal: 18,
      chainId: 11155420,
      tokenAddress: "0xMyTokenAddress",
      callData: "",
    });

    console.log("rawHex: ", rawHex);

    const tx = await opProvider.sendTransaction(rawHex);
    console.log("Transaction hash:", tx.hash);

    const receipt = await tx.wait(1);
    console.log("Transaction confirmed in block:", receipt.blockNumber);
    console.log(
      "Transaction status:",
      receipt.status === 1 ? "Success" : "Failed" // 务必检查 receipt.status
    );
  });

  test("op-sepolia L1 deposit", async () => {
    // OP L1官方桥的合约
    const bridgeContract = "0xFBb0621E0B23b5478B630BD55a5f21f67730B0F1";
    const nonce = await sepoliaProvider.getTransactionCount(walletAddress);
    console.log("address sepolia Current nonce:", nonce);

    // 注意这里要设置
    const l1GasLimit = 2000000;
    const l2GasLimit = 2000000;
    const data = "0x";

    // L1 Standard Bridge ABI (只要 depositETH)
    const bridgeAbi = ["function depositETH(uint32, bytes calldata)"];
    const iface = new Interface(bridgeAbi);
    const callData = iface.encodeFunctionData("depositETH", [l2GasLimit, "0x"]);
    console.log("input data", callData);

    const amount = "0.05"; // ETH
    const decimal = 18;
    const valueInWei = new BigNumber(amount).times("1e18");
    const from = walletAddress;

    const feeData = await sepoliaProvider.getFeeData();
    const estimate = await sepoliaProvider.estimateGas({
      from,
      to: bridgeContract,
      value: valueInWei.toFixed(0),
      data: callData,
    });

    const rawHex = signTransaction({
      privateKey: PRIVATE_KEY,
      nonce: nonce,
      from: walletAddress,
      to: bridgeContract,
      gasLimit: estimate.toNumber(),
      maxPriorityFeePerGas: parseInt(feeData.maxPriorityFeePerGas!.toString()),
      maxFeePerGas: parseInt(feeData.maxFeePerGas!.toString()),
      amount: "0.05",

      gasPrice: 0,

      decimal: 18,
      chainId: 11155111,
      tokenAddress: "0x00",
      callData: callData,
    });

    console.log("rawHex: ", rawHex);

    const tx = await sepoliaProvider.sendTransaction(rawHex);
    console.log("Transaction hash:", tx.hash);

    const receipt = await tx.wait(1);
    console.log("Transaction confirmed in block:", receipt.blockNumber);
    console.log(
      "Transaction status:",
      receipt.status === 1 ? "Success" : "Failed" // 务必检查 receipt.status
    );
  });

  // test op bridge deposit
  test("sign eth tx provider", async () => {
    // 你的钱包私钥（测试用）
    const privateKey =
      "cadf7450f8a7f15b5a3c9eb094b7de3fcef85465ee3de3d3a3bb687f31c79289";
    const wallet = new ethers.Wallet(privateKey);
    const from = wallet.address;

    const provider = new ethers.providers.JsonRpcProvider(
      "https://ethereum-sepolia-rpc.publicnode.com"
    );

    // Optimism L1 Bridge Sepolia 地址（官方地址）
    const bridgeAddress = "0xFBb0621E0B23b5478B630BD55a5f21f67730B0F1";

    // 构建 calldata
    const bridgeInterface = new ethers.utils.Interface([
      "function depositETH(uint32 _gas, bytes _data)",
    ]);
    const calldata = bridgeInterface.encodeFunctionData("depositETH", [
      2000000, // l2 gas
      "0x", // empty data
    ]);

    const amount = "0.05"; // ETH
    const decimal = 18;
    const valueInWei = new BigNumber(amount).times("1e18");

    const nonce = await provider.getTransactionCount(from, "latest");
    const feeData = await provider.getFeeData();
    const estimate = await provider.estimateGas({
      from,
      to: bridgeAddress,
      value: valueInWei.toFixed(0),
      data: calldata,
    });

    const signed = signTransaction({
      privateKey,
      nonce,
      from,
      to: bridgeAddress,
      gasLimit: estimate.toNumber(),
      amount,
      gasPrice: 0,
      decimal,
      chainId: 11155111,
      tokenAddress: "0x00",
      callData: calldata,
      maxPriorityFeePerGas: parseInt(feeData.maxPriorityFeePerGas!.toString()),
      maxFeePerGas: parseInt(feeData.maxFeePerGas!.toString()),
    });

    console.log("📤 Broadcasting tx...");
    const tx = await provider.sendTransaction(signed);
    console.log("📨 Sent! Tx Hash:", tx.hash);

    const receipt = await tx.wait();
    console.log("✅ Transaction confirmed");
    console.log(receipt.status === 1 ? "✅ 成功" : "❌ 执行失败");
    console.log(receipt);
  });

  test("deposit eth to another address", async () => {
    // 备用RPC
    const provider = new ethers.providers.JsonRpcProvider(
      "https://ethereum-sepolia-rpc.publicnode.com"
    );

    // Optimism L1 Bridge Sepolia 地址（官方地址）
    const bridgeAddress = "0xFBb0621E0B23b5478B630BD55a5f21f67730B0F1";

    // 构建 calldata
    const bridgeInterface = new ethers.utils.Interface([
      "function depositETHTo(address _address, uint32 _gas, bytes _data)",
    ]);
    const calldata = bridgeInterface.encodeFunctionData("depositETHTo", [
      "0xAe4Eb06A79F922Fc491C347baAF40F0562333895",
      2000000, // l2 gas
      "0x", // empty data
    ]);

    const amount = "0.05"; // ETH
    const decimal = 18;
    const valueInWei = new BigNumber(amount).times("1e18");

    // 获取nonce
    const nonce = await provider.getTransactionCount(walletAddress, "latest");

    // 预估手续费和gas
    const feeData = await provider.getFeeData();
    const estimate = await provider.estimateGas({
      from: walletAddress,
      to: bridgeAddress,
      value: valueInWei.toFixed(0),
      data: calldata,
    });

    const signed = signTransaction({
      privateKey: PRIVATE_KEY,
      nonce,
      from: walletAddress,
      to: bridgeAddress,
      gasLimit: estimate.toNumber(),
      amount,
      gasPrice: 0,
      decimal,
      chainId: 11155111,
      tokenAddress: "0x00",
      callData: calldata,
      maxPriorityFeePerGas: parseInt(feeData.maxPriorityFeePerGas!.toString()),
      maxFeePerGas: parseInt(feeData.maxFeePerGas!.toString()),
    });

    console.log("📤 Broadcasting tx...");
    const tx = await provider.sendTransaction(signed);
    console.log("📨 Sent! Tx Hash:", tx.hash);

    const receipt = await tx.wait();
    console.log("✅ Transaction confirmed");
    console.log(receipt.status === 1 ? "✅ 成功" : "❌ 执行失败");
    console.log(receipt);
  });

  test("op-sepolia get L2 balance", async () => {
    const l2Balance = await opProvider.getBalance(walletAddress);
    console.log("L2 余额:", ethers.utils.formatEther(l2Balance), "ETH");
  });

  test("op sepolia l2 approve", async () => {
    // 备用RPC
    const provider = new ethers.providers.JsonRpcProvider(
      "https://ethereum-sepolia-rpc.publicnode.com"
    );

    // Optimism Sepolia ETH 地址
    const TokenAddress = "0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000";

    const amount = "1"; // ETH
    const decimal = 18;
    const valueInWei = new BigNumber(amount).times("1e18");

    // 构建 calldata
    const bridgeInterface = new ethers.utils.Interface([
      "function approve(address spender, uint256 amount)",
    ]);
    const calldata = bridgeInterface.encodeFunctionData("approve", [
      "0x4200000000000000000000000000000000000010",
      ethers.utils.parseEther(amount), // amount
    ]);

    // 获取nonce
    const nonce = await opProvider.getTransactionCount(walletAddress, "latest");
    console.log("nonce:", nonce);
    const l2Balance = await opProvider.getBalance(walletAddress);
    console.log("L2 余额:", ethers.utils.formatEther(l2Balance), "ETH");

    // 预估手续费和gas
    const feeData = await provider.getFeeData();
    const estimate = await provider.estimateGas({
      from: walletAddress,
      to: TokenAddress,
      value: "0", // bigNumber转为字符串
      data: calldata,
    });
    console.log("estimate.toNumber(),", estimate.toNumber());
    console.log(
      "maxPriorityFeePerGas",
      parseInt(feeData.maxPriorityFeePerGas!.toString())
    );
    console.log("maxFeePerGas", parseInt(feeData.maxFeePerGas!.toString()));

    const signed = signTransaction({
      privateKey: PRIVATE_KEY,
      nonce,
      from: walletAddress,
      to: TokenAddress,
      gasLimit: estimate.toNumber(), // 22374
      amount: "0",
      gasPrice: 0,
      decimal,
      chainId: 11155420,
      tokenAddress: "0x00",
      callData: calldata,
      maxPriorityFeePerGas: parseInt(feeData.maxPriorityFeePerGas!.toString()),
      maxFeePerGas: parseInt(feeData.maxFeePerGas!.toString()),
    });

    console.log("📤 Broadcasting tx...");
    const tx = await opProvider.sendTransaction(signed);
    console.log("📨 Sent! Tx Hash:", tx.hash);

    const receipt = await tx.wait();
    console.log("✅ Transaction confirmed");
    console.log(receipt.status === 1 ? "✅ 成功" : "❌ 执行失败");
    console.log(receipt);
  });

  test("op sepolia l2 withdraw", async () => {
    // 备用RPC
    const provider = new ethers.providers.JsonRpcProvider(
      "https://ethereum-sepolia-rpc.publicnode.com"
    );

    // Optimism L2 Bridge Sepolia 地址（官方地址）
    const bridgeAddress = "0x4200000000000000000000000000000000000010";

    const l2TokenAddress = "0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000";

    const amount = "0.01"; // ETH
    const decimal = 18;
    const valueInWei = new BigNumber(amount).times("1e18");

    // 构建 calldata
    const bridgeInterface = new ethers.utils.Interface([
      "function withdrawTo(address _token, address _toAddress, uint256 _amount, uint32 _gas, bytes _data)",
    ]);
    const calldata = bridgeInterface.encodeFunctionData("withdrawTo", [
      l2TokenAddress,
      "0xAe4Eb06A79F922Fc491C347baAF40F0562333895",
      ethers.utils.parseEther(amount), // amount
      2000000, // l2 gas
      "0x", // empty data
    ]);

    // 获取nonce
    const nonce = await opProvider.getTransactionCount(walletAddress, "latest");
    console.log("nonce:", nonce);

    // 预估手续费和gas
    const feeData = await provider.getFeeData();
    const estimate = await provider.estimateGas({
      from: walletAddress,
      to: bridgeAddress,
      value: valueInWei.toFixed(0), // bigNumber转为字符串
      data: calldata,
    });

    console.log("estimate.toNumber(),", estimate.toNumber());
    console.log(
      "maxPriorityFeePerGas",
      parseInt(feeData.maxPriorityFeePerGas!.toString())
    );
    console.log("maxFeePerGas", parseInt(feeData.maxFeePerGas!.toString()));

    const signed = signTransaction({
      privateKey: PRIVATE_KEY,
      nonce,
      from: walletAddress,
      to: bridgeAddress,
      gasLimit: estimate.toNumber() * 5,
      amount,
      gasPrice: 0,
      decimal,
      chainId: 11155420,
      tokenAddress: "0x00",
      callData: calldata,
      maxPriorityFeePerGas:
        parseInt(feeData.maxPriorityFeePerGas!.toString()) * 3,
      maxFeePerGas: parseInt(feeData.maxFeePerGas!.toString()) * 3,
    });

    console.log("📤 Broadcasting tx...");
    const tx = await opProvider.sendTransaction(signed);
    console.log("📨 Sent! Tx Hash:", tx.hash);

    const receipt = await tx.wait();
    console.log("✅ Transaction confirmed");
    console.log(receipt.status === 1 ? "✅ 成功" : "❌ 执行失败");
    console.log(receipt);
  });

  test(
    "op-sepolia L2 withdraw",
    async () => {
      const bridgeAddress = "0x4200000000000000000000000000000000000010";

      const bridgeAbi = [
        "function withdraw(address _l2Token,uint256 _amount,uint32 _minGasLimit,bytes _extraData) external payable returns (uint256)",
      ];

      const bridgeContract = new ethers.utils.Interface(bridgeAbi);

      // ETH token地址 (Bedrock)
      const l2TokenAddress = "0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000";

      const amountStr = "0.01";
      const amount = ethers.utils.parseEther(amountStr);

      const l1GasLimit = 1000000;
      const l2GasLimit = 1000000;
      const data = "0x";

      const callData = bridgeContract.encodeFunctionData("withdraw", [
        l2TokenAddress,
        amount,
        l1GasLimit,
        data,
      ]);
      console.log("calldata", callData);

      const nonce = await opProvider.getTransactionCount(walletAddress);
      console.log("Current nonce:", nonce);

      console.log(`开始跨链转账 ${amountStr} ETH...`);

      const rawHex = signTransaction({
        privateKey: PRIVATE_KEY,
        nonce: nonce,
        from: walletAddress,
        to: bridgeAddress,
        gasLimit: l2GasLimit,
        maxFeePerGas: 327993150328,
        maxPriorityFeePerGas: 32799315032,
        gasPrice: 0,
        amount: amountStr,
        decimal: 18,
        chainId: 11155420,
        tokenAddress: "0x00",
        callData: callData,
      });

      const tx = await opProvider.sendTransaction(rawHex);

      console.log(`L2提现交易hash: ${tx.hash}`);
      await tx.wait();
      console.log(`✅ L2提现交易已确认`);

      console.log(`下一步：等待挑战期后在L1 finalize`);
    },
    1000 * 600
  );

  test(
    "withdraw status",
    async () => {
      const txHash =
        "0xc506d8be44b361d704857e0161007a4686adae2d703246441d073a3be7267f56";
      // const txHash = "0xc506d8be44b361d704857e0161007a4686adae2d703246441d073a3be7267f56";
      // const txHash = "0x9f6c0453153b660a8ce90968228ea4718e9d7d1269f273e0cd20b44ba6137ecc";
      // 创建 CrossChainMessenger
      const messenger = new CrossChainMessenger({
        l1ChainId: 11155111, // Sepolia
        l2ChainId: 11155420, // Optimism Sepolia
        l1SignerOrProvider: sepoliaWallet,
        l2SignerOrProvider: opWallet,
      });
      // 等待挑战期结束
      console.log(`等待消息状态变为 READY_FOR_RELAY...`);
      let status = await messenger.getMessageStatus(txHash);
      console.log(`当前状态: ${MessageStatus[status]}`);
      while (status !== MessageStatus.READY_FOR_RELAY) {
        console.log(`当前状态: ${MessageStatus[status]}，每60s检查一次...`);
        await new Promise((r) => setTimeout(r, 60000));
        status = await messenger.getMessageStatus(txHash);
      }
      console.log(`✅ 挑战期结束，准备在L1 finalize`);

      // 在L1 finalize withdrawal
      const finalizeTx = await messenger.finalizeMessage(txHash);
      console.log(`L1 finalize交易hash: ${finalizeTx.hash}`);
      await finalizeTx.wait();
      console.log(`✅ 完成提现！资金现在在L1可用`);
    },
    1000 * 60 * 30
  );

  test("send prove", async () => {
    // const txHash = "0xad82fffa08fb85bb15f3d4fb5501451f9e0d527a0313d16ae7bb2e4b106ba9e4";
    const txHash =
      "0xc506d8be44b361d704857e0161007a4686adae2d703246441d073a3be7267f56";
    // 创建 CrossChainMessenger
    const messenger = new CrossChainMessenger({
      l1ChainId: 11155111, // Sepolia
      l2ChainId: 11155420, // Optimism Sepolia
      l1SignerOrProvider: sepoliaWallet,
      l2SignerOrProvider: opWallet,
    });
    const response = await messenger.proveMessage(txHash);
    console.log(`responseHash:${response.hash}`);
    await response.wait();
    console.log(`✅ 完成跨链消息的证明`);
  });

  test("nodejs", async () => {
    try {
      setTimeout(function () {
        throw new Error("异步出错");
      }, 0);
    } catch (error) {
      console.log("error");
    }
    console.log("out try catch");
  });
});
