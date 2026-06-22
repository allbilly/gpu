package l1v

import (
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/util/ca"
)

type bankActionType int

const (
	bankActionReadHit bankActionType = iota
	bankActionWrite
	bankActionWriteFetched
)

type transaction struct {
	id string

	read         *mem.ReadReq
	readToBottom *mem.ReadReq

	write         *mem.WriteReq
	writeToBottom *mem.WriteReq

	preCoalesceTransactions []*transaction

	bankAction            bankActionType
	block                 *cache.Block
	data                  []byte
	writeFetchedDirtyMask []bool

	fetchAndWrite bool
	done          bool
}

func (t *transaction) Address() uint64 {
	if t.read != nil {
		return t.read.Address
	}
	return t.write.Address
}

func (t *transaction) PID() ca.PID {
	if t.read != nil {
		return t.read.PID
	}
	return t.write.PID
}
