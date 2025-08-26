from eth_account import Account
from eth_utils import keccak, to_hex
from eth_abi import encode

# -------------------------------
# 配置参数（你可以替换为自己的值）
# -------------------------------
name = "CpChain"
chain_id = 1
verifying_contract = "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC"

staker = "0xF39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
cpChainBase = "0x000000000000000000000000000000000000dEaD"
amount = 2 * 10**18
nonce = 0
expiry = 9999999999

# -------------------------------
# 计算 EIP-712 digest hash
# -------------------------------
DOMAIN_TYPEHASH = keccak(text="EIP712Domain(string name,uint256 chainId,address verifyingContract)")
DEPOSIT_TYPEHASH = keccak(text="Deposit(address staker,address cpChainBase,uint256 amount,uint256 nonce,uint256 expiry)")

domain_separator = keccak(encode(
    ["bytes32", "bytes32", "uint256", "address"],
    [DOMAIN_TYPEHASH, keccak(text=name), chain_id, verifying_contract]
))

struct_hash = keccak(encode(
    ["bytes32", "address", "address", "uint256", "uint256", "uint256"],
    [DEPOSIT_TYPEHASH, staker, cpChainBase, amount, nonce, expiry]
))

digest_hash = keccak(b"\x19\x01" + domain_separator + struct_hash)

print("Digest Hash (EIP-712):", to_hex(digest_hash))

# -------------------------------
# 使用私钥签名
# -------------------------------
private_key = "0x4f3edf983ac636a65a842ce7c78d9aa706d3b113b37c9430e6fd8c8d11f8b4e6"  # 示例私钥
signed = Account.signHash(digest_hash, private_key=private_key)

print("Address:", Account.from_key(private_key).address)
print("Signature:", signed.signature.hex())
print("v:", signed.v)
print("r:", hex(signed.r))
print("s:", hex(signed.s))
