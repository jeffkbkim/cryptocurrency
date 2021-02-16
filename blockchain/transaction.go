package blockchain

import (
	"fmt"

	"github.com/jeffkbkim/cryptocurrency/helper"
	"github.com/jeffkbkim/cryptocurrency/pki"
)

type Transaction struct {
	From      string
	To        string
	Amount    int
	Signature []byte
}

func (t *Transaction) IsValidSignature() bool {
	if t.isGenesisTxn() {
		return true
	}

	return pki.Verify(t.ToMessage(), t.Signature, t.From)
}

func (t *Transaction) ToMessage() string {
	return helper.Hash(fmt.Sprintf("%s%s%d", t.From, t.To, t.Amount))
}

func (t *Transaction) isGenesisTxn() bool {
	return t.From == ""
}
