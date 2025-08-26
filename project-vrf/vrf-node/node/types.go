package node

type SignMsgResponse struct {
	Signature       []byte `json:"signature"`
	G2Point         []byte `json:"g2_point"`
	NonSignerPubkey []byte `json:"non_signer_pubkey"`
}
