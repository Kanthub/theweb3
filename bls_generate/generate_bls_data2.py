from py_ecc.bn128 import (
    G1,
    G2,
    add,
    multiply,
    pairing,
    is_on_curve,
    curve_order,
    b,
    b2,
    FQ,
    field_modulus as P
)
from eth_utils import keccak
import hashlib
import os


def print_g1(name, pt):
    x, y = pt
    print(f"{name}:")
    print(f"  X = {x}")
    print(f"  Y = {y}\n")


def print_g2(name, pt):
    x, y = pt
    x0, x1 = x.coeffs
    y0, y1 = y.coeffs
    print(f"{name}:")
    print(f"  X = [{x0}, {x1}]")
    print(f"  Y = [{y0}, {y1}]\n")

def hash_to_G1_solidity_style(msg: bytes):
    x = int.from_bytes(keccak(msg), 'big') % P
    while True:
        beta = (x**3 + 3) % P
        y = pow(beta, (P + 1) // 4, P)

        if (y * y) % P == beta:
            point = (FQ(x), FQ(y))
            assert is_on_curve(point, b)
            return point
        x = (x + 1) % P



# Step 1: Generate private key
sk = int.from_bytes(os.urandom(32), 'big') % curve_order
print(f"Private Key (sk): {sk}\n")

# Step 2: Generate public key
pk_g1 = multiply(G1, sk)
pk_g2 = multiply(G2, sk)

assert is_on_curve(pk_g1,b), "G1 pubkey not on curve"
assert is_on_curve(pk_g2,b2), "G2 pubkey not on curve"

# Step 3: Message hash → scalar → map to G1 (H)
msg = b"test message"
msgHash = keccak(msg)
# msg_scalar = int.from_bytes(hashlib.sha256(msg).digest(), 'big') % curve_order
msg_hash = hash_to_G1_solidity_style(msg)

assert is_on_curve(msg_hash,b), "Message hash not on G1"

# Step 4: Sign with sk
signature = multiply(msg_hash, sk)
assert is_on_curve(signature,b), "Signature not on G1"

# Step 5: Verify signature
assert pairing(G2, signature) == pairing(pk_g2, msg_hash), "Pairing check failed"

# Step 6: Print
print_g1("G1 Public Key", pk_g1)
print_g2("G2 Public Key", pk_g2)
print("Message Hash (in bytes)", msgHash.hex())
print_g1("Message Hash (in G1)", msg_hash)
print_g1("Signature", signature)
