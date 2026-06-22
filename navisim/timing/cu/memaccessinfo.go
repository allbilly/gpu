package cu

import (
	"gitlab.com/akita/mem"
	"gitlab.com/akita/navisim/rdnainsts"
	"gitlab.com/akita/navisim/timing/wavefront"
)

type vectorMemAccessLaneInfo struct {
	laneID                int
	reg                   *rdnainsts.Reg
	regCount              int
	addrOffsetInCacheLine uint64
}

// VectorMemAccessInfo defines access info
type VectorMemAccessInfo struct {
	ID            string
	Read          *mem.ReadReq
	Write         *mem.WriteReq
	GL0Invalidate *mem.GL0InvalidateReq
	Wavefront     *wavefront.Wavefront
	Inst          *wavefront.Inst
	laneInfo      []vectorMemAccessLaneInfo
}

// TaskID returns the ID of the VectorMemAccess transaction
func (i VectorMemAccessInfo) TaskID() string {
	return i.ID
}

// InstFetchReqInfo defines request info
type InstFetchReqInfo struct {
	Req       *mem.ReadReq
	Wavefront *wavefront.Wavefront
	Address   uint64
}

// ScalarMemAccessInfo defines request info
type ScalarMemAccessInfo struct {
	Req       *mem.ReadReq
	Wavefront *wavefront.Wavefront
	DstSGPR   *rdnainsts.Reg
	Inst      *wavefront.Inst
}
