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
)
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
msg_scalar = int.from_bytes(hashlib.sha256(msg).digest(), 'big') % curve_order
msg_hash = multiply(G1, msg_scalar)

assert is_on_curve(msg_hash,b), "Message hash not on G1"

# Step 4: Sign with sk
signature = multiply(msg_hash, sk)
assert is_on_curve(signature,b), "Signature not on G1"

# Step 5: Verify signature
assert pairing(G2, signature) == pairing(pk_g2, msg_hash), "Pairing check failed"

# Step 6: Print
print_g1("G1 Public Key", pk_g1)
print_g2("G2 Public Key", pk_g2)
print_g1("Message Hash (in G1)", msg_hash)
print_g1("Signature", signature)
