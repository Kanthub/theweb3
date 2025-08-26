const bitcoinjs = require('bitcoinjs-lib');
const bip39 = require('bip39');
const ecc = require('tiny-secp256k1');
const { BIP32Factory } = require('bip32');
const bip32 = BIP32Factory(ecc);

// 初始化 ecc 库
bitcoinjs.initEccLib(ecc);





const mnemonic = "around dumb spend sample oil crane plug embrace outdoor panel rhythm salon";
const seed = bip39.mnemonicToSeedSync(mnemonic, "")
const param = {
    seedHex: seed.toString("hex"),
    receiveOrChange: "0",
    addressIndex: "0",
    network: "mainnet",
}
const account = createSchnorrAddress(param)
console.log(account.address);
console.log(account.privateKey);
console.log(account.publicKey);