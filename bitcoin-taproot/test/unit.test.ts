import * as bip39 from "bip39";
import { createSchnorrAddress, signBtcTaprootTransaction } from "../src";
// import test, { describe } from "node:test";
import { describe, it, expect, test } from "@jest/globals";

describe("Bitcoin Taproot test", () => {
  test("create address", async () => {
    const mnemonic =
      "sort what document outdoor plastic little country witness output beauty upon pudding";
    const params_1 = {
      mnemonic: mnemonic,
      password: "",
    };
    const seedHex = bip39.mnemonicToSeedSync(mnemonic);
    const params = {
      seedHex: seedHex,
      receiveOrChange: 0,
      addressIndex: 0,
      network: "mainnet",
    };
    const account = await createSchnorrAddress(params);
    console.log(account);
  });

  test("sign transaction", async () => {
    const data = {
      outputs: [
        {
          amount: 3000,
          address:
            "bc1pp4802jvzr9vswddm4msstqdfgsjg5f9yvfh458q4e9xw0qcacrzs49sw25",
        },
        {
          amount: 1000,
          address:
            "bc1p5p6ptfzjfm4dy6vey8zcqk747cnqa35cwggza6fd6qw7g0mucqnq5l6jnc",
        },
      ],
      inputs: [
        {
          address:
            "bc1pp4802jvzr9vswddm4msstqdfgsjg5f9yvfh458q4e9xw0qcacrzs49sw25",
          txid: "b00771c6acc9d84e503edb1cab32325dee4d261762e84d23fb11fab26143ff18",
          vout: 1,
          amount: 5000,
        },
      ],
    };

    var ss1 = await signBtcTaprootTransaction({
      privateKey: "L5QKkjddguEJQjnTkfm9Fbm6QnLxQvi7RPo9CNj4xNo1xPfq99jp",
      signObj: data,
      network: "mainnet",
    });

    console.log(ss1);
  });
});

import * as bitcoin from "bitcoinjs-lib";
import * as bip32 from "bip32";
import { ECPairFactory } from "ecpair";
import * as ecc from "tiny-secp256k1";
import { toXOnly } from "bitcoinjs-lib/src/psbt/bip371";

const ECPair = ECPairFactory(ecc);

describe("Bitcoin Taproot Offline Sign Test", () => {
  test("sign transaction", async () => {
    const network = bitcoin.networks.bitcoin;

    // ✅ 使用一个已经生成过的 Taproot 私钥（WIF 或 raw）
    const keyPair = ECPair.fromWIF(
      "L5QKkjddguEJQjnTkfm9Fbm6QnLxQvi7RPo9CNj4xNo1xPfq99jp",
      network
    );
    const privateKey = keyPair.privateKey!;
    const publicKey = keyPair.publicKey; // compressed (33 bytes)

    // ✅ 模拟离线输入 UTXO
    const mockInput = {
      txid: "f9d8999d08b2a8e9f7a1743486f26813c94c2c1cfbeaa8db52e95a7f4c8c6be3",
      index: 0,
      amount: 10000, // sats
      output: bitcoin.payments.p2tr({
        internalPubkey: publicKey.slice(1, 33),
        network,
      }).output!,
      publicKey: publicKey, // required by signBtcTaprootTransaction
    };

    // ✅ 构造输出信息
    const mockOutput = {
      value: 9000, // 10000 - 1000 fee
      sendAddress:
        "bc1p5p6ptfzjfm4dy6vey8zcqk747cnqa35cwggza6fd6qw7g0mucqnq5l6jnc",
      //   sendPubKey: toXOnly(publicKey), // optional, 可选用于 tapOutputKey（你代码中用了）
    };

    const txHex = await signBtcTaprootTransaction({
      privateKey,
      signObj: {
        inputs: [mockInput],
        output: [mockOutput],
      },
    });

    console.log("Raw signed tx:", txHex);

    // ✅ 校验 hex 有效性
    expect(txHex).toMatch(/^02/); // 比特币交易 hex 开头一般为 02
  });
});
