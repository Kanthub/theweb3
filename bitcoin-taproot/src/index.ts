import BIP32Factory from "bip32";
import * as ecc from "tiny-secp256k1";
import * as bitcoin from "bitcoinjs-lib";

import {
  toXOnly,
  tapTreeToList,
  tapTreeFromList,
} from "bitcoinjs-lib/src/psbt/bip371";

bitcoin.initEccLib(ecc);
const bip32 = BIP32Factory(ecc);

export async function createSchnorrAddress(params: {
  seedHex: any;
  receiveOrChange: any;
  addressIndex: any;
  network: any;
}) {
  const { seedHex, receiveOrChange, addressIndex, network } = params;
  const root = bip32.fromSeed(
    Buffer.from(seedHex, "hex"),
    bitcoin.networks.bitcoin
  );
  let path = "m/86'/0'/0'/0/" + addressIndex + "";
  if (receiveOrChange === "1") {
    path = "m/86'/0'/0'/1/" + addressIndex + "";
  }

  const childKey = root.derivePath(path);
  const privateKey = childKey.privateKey;
  if (!privateKey) throw new Error("No private key found");

  const publicKey = childKey.publicKey;

  // 生成 P2TR 地址
  const { address } = bitcoin.payments.p2tr({
    internalPubkey: childKey.publicKey!.slice(1, 33),
    network: bitcoin.networks.bitcoin,
  });

  return {
    privateKey: Buffer.from(privateKey).toString("hex"),
    publicKey: Buffer.from(childKey.publicKey).toString("hex"),
    address,
  };
}

export async function signBtcTaprootTransaction(params: {
  privateKey: any;
  signObj: any;
  network?: string;
}) {
  const { signObj, privateKey } = params;
  const psbt = new bitcoin.Psbt({ network: bitcoin.networks.bitcoin });

  const inputs = signObj.inputs.map(
    (input: {
      txid: any;
      index: any;
      amount: any;
      output: any;
      publicKey: Buffer;
    }) => {
      return {
        hash: input.txid,
        index: input.index,
        witnessUtxo: { value: input.amount, script: input.output! },
        tapInternalKey: toXOnly(input.publicKey),
      };
    }
  );
  psbt.addInputs(inputs);

  /////////////    改动1
  const dummyChainCode = Buffer.alloc(32, 0); // 全 0 填充的 32 字节 buffer
  const sendInternalKey = bip32.fromPrivateKey(privateKey, dummyChainCode);
  //   const sendInternalKey = bip32.fromPrivateKey(privateKey, Buffer.from("0"));

  //////////////// 改动2
  const output = signObj.output.map(
    (output: { value: any; sendAddress: any; sendPubKey: any }) => {
      return {
        value: output.value,
        address: output.sendAddress!,
        // tapInternalKey: output.sendPubKey,
      };
    }
  );

  psbt.addOutputs(output);

  const tweakedSigner = sendInternalKey.tweak(
    bitcoin.crypto.taggedHash("TapTweak", toXOnly(sendInternalKey.publicKey))
  );

  await psbt.signInputAsync(0, tweakedSigner);
  psbt.finalizeAllInputs();
  const tx = psbt.extractTransaction();
  return tx.toBuffer().toString("hex");
}
