package bitcoin

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvd/chaincfg"
	"github.com/bitcoinsv/bsvd/chaincfg/chainhash"
	"github.com/bitcoinsv/bsvd/wire"
	"github.com/bitcoinsv/bsvutil"
)

const (
	// H_BSV is the magic header string required fore Bitcoin Signed Messages
	hBSV string = "Bitcoin Signed Message:\n"
)

// VerifyMessage verifies a string and address against the provided
// signature and assumes Bitcoin Signed Message encoding
//
// Error will occur if verify fails or verification is not successful (no bool)
// Spec: https://docs.moneybutton.com/docs/bsv-message.html
func VerifyMessage(address, sig, data string) error {

	decodedSig, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return err
	}

	// Validate the signature - this just shows that it was valid at all.
	// we will compare it with the key next.
	var buf bytes.Buffer
	if err = wire.WriteVarString(&buf, 0, hBSV); err != nil {
		return err
	}
	if err = wire.WriteVarString(&buf, 0, data); err != nil {
		return err
	}

	// Create the hash
	expectedMessageHash := chainhash.DoubleHashB(buf.Bytes())

	var publicKey *bsvec.PublicKey
	var wasCompressed bool
	if publicKey, wasCompressed, err = bsvec.RecoverCompact(bsvec.S256(), decodedSig, expectedMessageHash); err != nil {
		return err
	}

	// Reconstruct the pubkey hash.
	var serializedPK []byte
	if wasCompressed {
		serializedPK = publicKey.SerializeCompressed()
	} else {
		serializedPK = publicKey.SerializeUncompressed()
	}
	var bsvecAddress *bsvutil.AddressPubKey
	if bsvecAddress, err = bsvutil.NewAddressPubKey(serializedPK, &chaincfg.MainNetParams); err != nil {
		return err
	}

	// Return nil if addresses match.
	if bsvecAddress.EncodeAddress() == address {
		return nil
	}
	return fmt.Errorf("address: %s not found vs %s", address, bsvecAddress.EncodeAddress())
}

// VerifyMessageDER will take a message string, a public key string and a signature string
// (in strict DER format) and verify that the message was signed by the public key.
//
// Copyright (c) 2019 Bitcoin Association
// License: https://github.com/bitcoin-sv/merchantapi-reference/blob/master/LICENSE
//
// Source: https://github.com/bitcoin-sv/merchantapi-reference/blob/master/handler/global.go
func VerifyMessageDER(hash [32]byte, pubKey string, signature string) (verified bool, err error) {

	// Decode the signature string
	var sigBytes []byte
	if sigBytes, err = hex.DecodeString(signature); err != nil {
		return
	}

	// Parse the signature
	var sig *bsvec.Signature
	if sig, err = bsvec.ParseDERSignature(sigBytes, bsvec.S256()); err != nil {
		return
	}

	// Decode the pubKey
	var pubKeyBytes []byte
	if pubKeyBytes, err = hex.DecodeString(pubKey); err != nil {
		return
	}

	// Parse the pubKey
	var rawPubKey *bsvec.PublicKey
	if rawPubKey, err = bsvec.ParsePubKey(pubKeyBytes, bsvec.S256()); err != nil {
		return
	}

	// Verify the signature against the pubKey
	verified = sig.Verify(hash[:], rawPubKey)
	return
}
