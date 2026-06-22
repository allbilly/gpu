package rdnainsts

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDisassembler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RDNA Disassembler")
}

var _ = Describe("Disassembler", func() {
	var (
		disassembler *Disassembler
	)

	BeforeEach(func() {
		disassembler = NewDisassembler()
	})

	It("should disassembler BE803C6A", func() {
		buf := []byte{0x6A, 0x3C, 0x80, 0xBE}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).To(Equal("s_and_saveexec_b32 s0, vcc_lo"))
	})

	It("should disassembler 50060201", func() {
		buf := []byte{0x01, 0x02, 0x06, 0x50}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_add_co_ci_u32_e32 v3, vcc_lo, s1, v1, vcc_lo"))
	})

	It("should disassemble DC308000 037D0005", func() {
		buf := []byte{0x00, 0x80, 0x30, 0xDC, 0x05, 0x00, 0x7D, 0x03}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).To(Equal("global_load_dword v3, v[5:6], off"))
	})

	It("should disassemble D7010000 0002029E", func() {
		buf := []byte{0x00, 0x00, 0x01, 0xD7, 0x9E, 0x02, 0x02, 0x00}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).To(Equal("v_ashrrev_i64 v[0:1], 30, v[1:2]"))
	})

	It("should decode F4000082 FA000004", func() {
		buf := []byte{0x82, 0x00, 0x00, 0xF4, 0x04, 0x00, 0x00, 0xFA}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("s_load_dword s2, s[4:5], 0x4"))
	})

	It("should decode 7E040280", func() {
		buf := []byte{0x80, 0x02, 0x04, 0x7E}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_mov_b32_e32 v2, 0"))
	})

	It("should decode D5280001 00010005", func() {
		buf := []byte{0x01, 0x00, 0x28, 0xD5, 0x05, 0x00, 0x01, 0x00}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_add_co_ci_u32_e64 v1, s0, s5, 0, s0"))
	})

	It("should decode 38060205", func() {
		buf := []byte{0x05, 0x02, 0x06, 0x38}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_or_b32_e32 v3, s5, v1"))
	})

	It("should decode D4040000 00020B0F", func() {
		buf := []byte{0x00, 0x00, 0x04, 0xD4, 0x0F, 0x0B, 0x02, 0x00}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cmp_gt_f32_e64 s0, v15, v5"))
	})

	It("should decode D4860001 0002090B", func() {
		buf := []byte{0x01, 0x00, 0x86, 0xD4, 0x0B, 0x09, 0x02, 0x00}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cmp_ge_i32_e64 s1, v11, v4"))
	})

	It("should decode 7DCA0480", func() {
		buf := []byte{0x80, 0x04, 0xCA, 0x7D}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cmp_ne_u64_e32 vcc_lo, 0, v[2:3]"))
	})

	It("should decode 8904007E", func() {
		buf := []byte{0x7E, 0x00, 0x04, 0x89}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("s_xor_b32 s4, exec_lo, s0"))
	})

	It("should decode 7E040C10", func() {
		buf := []byte{0x10, 0x0C, 0x04, 0x7E}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cvt_f32_u32_e32 v2, s16"))
	})

	It("should decode 3E0406FF 4F800000", func() {
		buf := []byte{0xFF, 0x06, 0x04, 0x3E, 0x00, 0x00, 0x80, 0x4F}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_mac_f32_e32 v2, 0x4f800000, v3"))

	})

	It("should decode 7E045502", func() {
		buf := []byte{0x02, 0x55, 0x04, 0x7E}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_rcp_f32_e32 v2, v2"))
	})

	It("should decode 100404FF 5F7FFFFC", func() {
		buf := []byte{0xFF, 0x04, 0x04, 0x10, 0xFC, 0xFF, 0x7F, 0x5F}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_mul_f32_e32 v2, 0x5f7ffffc, v2"))
	})

	It("should decode 7E064303", func() {
		buf := []byte{0x03, 0x43, 0x06, 0x7E}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_trunc_f32_e32 v3, v3"))
	})

	It("should decode 7E060F03", func() {
		buf := []byte{0x03, 0x0F, 0x06, 0x7E}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cvt_u32_f32_e32 v3, v3"))
	})

	It("should decode D7106A06 00020F00", func() {
		buf := []byte{0x06, 0x6A, 0x10, 0xD7, 0x00, 0x0F, 0x02, 0x00}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_sub_co_u32_e64 v6, vcc_lo, v0, v7"))
	})

	It("should decode 4C0A0901", func() {
		buf := []byte{0x01, 0x09, 0x0A, 0x4C}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_sub_nc_u32_e32 v5, v1, v4"))
	})

	It("should decode D52A0005 01AA0A05", func() {
		buf := []byte{0x05, 0x00, 0x2A, 0xD5, 0x05, 0x0A, 0xAA, 0x01}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_subrev_co_ci_u32_e64 v5, s0, s5, v5, vcc_lo"))
	})

	It("should decode 52080901", func() {
		buf := []byte{0x01, 0x09, 0x08, 0x52}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_sub_co_ci_u32_e32 v4, vcc_lo, v1, v4, vcc_lo"))
	})

	It("should decode D4C20000 00020805", func() {
		buf := []byte{0x00, 0x00, 0xC2, 0xD4, 0x05, 0x08, 0x02, 0x00}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cmp_eq_u32_e64 s0, s5, v4"))
	})

	It("should decode 7D861210", func() {
		buf := []byte{0x10, 0x12, 0x86, 0x7D}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cmp_le_u32_e32 vcc_lo, s16, v9"))
	})

	It("should decode D5010009 01A98280", func() {
		buf := []byte{0x09, 0x00, 0x01, 0xD5, 0x80, 0x82, 0xA9, 0x01}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cndmask_b32_e64 v9, 0, -1, vcc_lo"))
	})

	It("should decode 7D840A05", func() {
		buf := []byte{0x05, 0x0A, 0x84, 0x7D}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cmp_eq_u32_e32 vcc_lo, s5, v5"))
	})

	It("should decode 7D8A0A80", func() {
		buf := []byte{0x80, 0x0A, 0x8A, 0x7D}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cmp_ne_u32_e32 vcc_lo, 0, v5"))
	})

	It("should decode 7E045702", func() {
		buf := []byte{0x02, 0x57, 0x04, 0x7E}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_rcp_iflag_f32_e32 v2, v2"))
	})

	It("should decode 4E0A0610", func() {
		buf := []byte{0x10, 0x06, 0x0A, 0x4E}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_subrev_nc_u32_e32 v5, s16, v3"))
	})

	It("should decode 240C0680", func() {
		buf := []byte{0x80, 0x06, 0x0C, 0x24}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_max_i32_e32 v6, 0, v3"))
	})

	It("should decode 220E0A0D", func() {
		buf := []byte{0x0D, 0x0A, 0x0E, 0x22}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_min_i32_e32 v7, s13, v5"))
	})

	It("should decode 7D020F06", func() {
		buf := []byte{0x06, 0x0F, 0x02, 0x7D}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cmp_lt_i32_e32 vcc_lo, v6, v7"))
	})

	It("should decode 7E0A02F3", func() {
		buf := []byte{0xF3, 0x02, 0x0A, 0x7E}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_mov_b32_e32 v5, -1.0"))
	})

	It("should decode BF800000", func() {
		buf := []byte{0x00, 0x00, 0x80, 0xBF}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("s_nop 0x0"))
	})

	It("should decode 7D060E11", func() {
		buf := []byte{0x11, 0x0E, 0x06, 0x7D}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cmp_le_i32_e32 vcc_lo, s17, v7"))
	})

	It("should decode 7D8C111A", func() {
		buf := []byte{0x1A, 0x11, 0x8C, 0x7D}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cmp_ge_u32_e32 vcc_lo, v26, v8"))
	})

	It("should decode 062E1D17", func() {
		buf := []byte{0x17, 0x1D, 0x2E, 0x06}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_add_f32_e32 v23, v23, v14"))
	})
})
