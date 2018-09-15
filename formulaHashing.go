package api

import (
	"crypto/sha512"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/cbor"
	"github.com/polydawn/refmt/misc"
)

// FormulaSetupHash is an opaque string derived from a cryptographic hash of
// the deterministic serialization of a Formula.  Which is a fancy way of
// saying it's a fantastic primary key for memoizing computations.
type FormulaSetupHash string

/*
	Returns a hash covering parts of the formula such that the hash may be
	expected to converge for formulae that describe identical setups.

	Specifically, this hash includes the inputs, actions, and output slot specs;
	it excludes any actual output ware hashes, and excludes any fields which
	are incidental to correctly reproducing the task, such as warehouse URLs.

	The returned string is the base58 encoding of a SHA-384 hash, though
	there is no reason you should treat it as anything but opaque.
	The returned string may be relied upon to be all alphanumeric characters.
*/
func (frm Formula) SetupHash() FormulaSetupHash {
	msg, err := refmt.MarshalAtlased(
		cbor.EncodeOptions{},
		frm,
		Atlas_Formula,
	)
	if err != nil {
		panic(err)
	}
	hasher := sha512.New384()
	hasher.Write(msg)
	return FormulaSetupHash(misc.Base58Encode(hasher.Sum(nil)))
}
