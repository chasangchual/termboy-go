package main

import (
	"log"
)

type Register byte
type Word uint16

type Registers struct {
	A byte
	B byte
	C byte
	D byte
	E byte
	H byte
	L byte
	F byte // Flags Register
}

//See ZILOG z80 cpu manual p.80  (http://www.zilog.com/docs/z80/um0080.pdf)
type Clock struct {
	m Word
	t Word
}

func (c *Clock) Set(m, t Word) {
	c.m, c.t = m, t
}

func (c *Clock) Reset() {
	c.m, c.t = 0, 0
}

type Z80 struct {
	PC             Word // Program Counter
	SP             Word // Stack Pointer
	R              Registers
	MachineCycles  Clock
	LastInstrCycle Clock
	mmu            MMU
}

func NewCPU(m MMU) *Z80 {
	cpu := new(Z80)
	cpu.mmu = m
	cpu.Reset()

	//TODO: startup. additional setup here

	return cpu
}

func (cpu *Z80) Reset() {
	cpu.PC = 0
	cpu.SP = 0
	cpu.R.A = 0
	cpu.R.B = 0
	cpu.R.C = 0
	cpu.R.D = 0
	cpu.R.E = 0
	cpu.R.F = 0
	cpu.R.H = 0
	cpu.R.L = 0
	cpu.MachineCycles.Reset()
	cpu.LastInstrCycle.Reset()
}

func (cpu *Z80) ResetFlags() {
	cpu.R.F = 0x00
}

func (cpu *Z80) IncrementPC(by Word) {
	cpu.PC += by
}

func (cpu *Z80) ToWord(hob, lob byte) Word {
	return (Word(hob) << 8) ^ Word(lob)
}

func (cpu *Z80) Dispatch(Opcode byte) {
	switch Opcode {
	case 0x3E: //LD A, n
		cpu.LDrn(&cpu.R.A)
	case 0x06: //LD B,n
		cpu.LDrn(&cpu.R.B)
	case 0x0E: //LD C,n
		cpu.LDrn(&cpu.R.C)
	case 0x16: //LD D,n
		cpu.LDrn(&cpu.R.D)
	case 0x1E: //LD E,n
		cpu.LDrn(&cpu.R.E)
	case 0x26: //LD H,n
		cpu.LDrn(&cpu.R.H)
	case 0x2E: //LD L,n
		cpu.LDrn(&cpu.R.L)

	case 0x7E: //LD A, (HL)
		cpu.LDr_hl(&cpu.R.A)
	case 0x0A: //LD A, (BC)
		cpu.LDr_bc(&cpu.R.A)
	case 0x1A: //LD A, (DE)
		cpu.LDr_de(&cpu.R.A)
	case 0xFA: //LD A, (nn)
		cpu.LDr_nn(&cpu.R.A)

	case 0x7F: //LD A, A
		cpu.LDrr(&cpu.R.A, &cpu.R.A)
	case 0x78: //LD A, B
		cpu.LDrr(&cpu.R.A, &cpu.R.B)
	case 0x79: //LD A, C
		cpu.LDrr(&cpu.R.A, &cpu.R.C)
	case 0x7A: //LD A, D
		cpu.LDrr(&cpu.R.A, &cpu.R.D)
	case 0x7B: //LD A, E
		cpu.LDrr(&cpu.R.A, &cpu.R.E)
	case 0x7C: //LD A, H
		cpu.LDrr(&cpu.R.A, &cpu.R.H)
	case 0x7D: //LD A, L
		cpu.LDrr(&cpu.R.A, &cpu.R.L)

	case 0x47: //LD B, A
		cpu.LDrr(&cpu.R.B, &cpu.R.A)
	case 0x40: //LD B, B
		cpu.LDrr(&cpu.R.B, &cpu.R.B)
	case 0x41: //LD B, C
		cpu.LDrr(&cpu.R.B, &cpu.R.C)
	case 0x42: //LD B, D
		cpu.LDrr(&cpu.R.B, &cpu.R.D)
	case 0x43: //LD B, E
		cpu.LDrr(&cpu.R.B, &cpu.R.E)
	case 0x44: //LD B, H
		cpu.LDrr(&cpu.R.B, &cpu.R.H)
	case 0x45: //LD B, L
		cpu.LDrr(&cpu.R.B, &cpu.R.L)
	case 0x46: //LD B, (HL)
		cpu.LDr_hl(&cpu.R.B)

	case 0x4F: //LD C, A
		cpu.LDrr(&cpu.R.C, &cpu.R.A)
	case 0x48: //LD C, B
		cpu.LDrr(&cpu.R.C, &cpu.R.B)
	case 0x49: //LD C, C
		cpu.LDrr(&cpu.R.C, &cpu.R.C)
	case 0x4A: //LD C, D
		cpu.LDrr(&cpu.R.C, &cpu.R.D)
	case 0x4B: //LD C, E
		cpu.LDrr(&cpu.R.C, &cpu.R.E)
	case 0x4C: //LD C, H
		cpu.LDrr(&cpu.R.C, &cpu.R.H)
	case 0x4D: //LD C, L
		cpu.LDrr(&cpu.R.C, &cpu.R.L)
	case 0x4E: //LD C, (HL)
		cpu.LDr_hl(&cpu.R.C)

	case 0x57: //LD D, A
		cpu.LDrr(&cpu.R.D, &cpu.R.A)
	case 0x50: //LD D, B
		cpu.LDrr(&cpu.R.D, &cpu.R.B)
	case 0x51: //LD D, C
		cpu.LDrr(&cpu.R.D, &cpu.R.C)
	case 0x52: //LD D, D
		cpu.LDrr(&cpu.R.D, &cpu.R.D)
	case 0x53: //LD D, E
		cpu.LDrr(&cpu.R.D, &cpu.R.E)
	case 0x54: //LD D, H
		cpu.LDrr(&cpu.R.D, &cpu.R.H)
	case 0x55: //LD D, L
		cpu.LDrr(&cpu.R.D, &cpu.R.L)
	case 0x56: //LD D, (HL)
		cpu.LDr_hl(&cpu.R.D)

	case 0x5F: //LD E, A
		cpu.LDrr(&cpu.R.E, &cpu.R.A)
	case 0x58: //LD E, B
		cpu.LDrr(&cpu.R.E, &cpu.R.B)
	case 0x59: //LD E, C
		cpu.LDrr(&cpu.R.E, &cpu.R.C)
	case 0x5A: //LD E, D
		cpu.LDrr(&cpu.R.E, &cpu.R.D)
	case 0x5B: //LD E, E
		cpu.LDrr(&cpu.R.E, &cpu.R.E)
	case 0x5C: //LD E, H
		cpu.LDrr(&cpu.R.E, &cpu.R.H)
	case 0x5D: //LD E, L
		cpu.LDrr(&cpu.R.E, &cpu.R.L)
	case 0x5E: //LD E, (HL)
		cpu.LDr_hl(&cpu.R.E)

	case 0x67: //LD H, A
		cpu.LDrr(&cpu.R.H, &cpu.R.A)
	case 0x60: //LD H, B
		cpu.LDrr(&cpu.R.H, &cpu.R.B)
	case 0x61: //LD H, C
		cpu.LDrr(&cpu.R.H, &cpu.R.C)
	case 0x62: //LD H, D
		cpu.LDrr(&cpu.R.H, &cpu.R.D)
	case 0x63: //LD H, E
		cpu.LDrr(&cpu.R.H, &cpu.R.E)
	case 0x64: //LD H, H
		cpu.LDrr(&cpu.R.H, &cpu.R.H)
	case 0x65: //LD H, L
		cpu.LDrr(&cpu.R.H, &cpu.R.L)
	case 0x66: //LD H, (HL)
		cpu.LDr_hl(&cpu.R.H)

	case 0x6F: //LD L, A
		cpu.LDrr(&cpu.R.L, &cpu.R.A)
	case 0x68: //LD L, B
		cpu.LDrr(&cpu.R.L, &cpu.R.B)
	case 0x69: //LD L, C
		cpu.LDrr(&cpu.R.L, &cpu.R.C)
	case 0x6A: //LD L, D
		cpu.LDrr(&cpu.R.L, &cpu.R.D)
	case 0x6B: //LD L, E
		cpu.LDrr(&cpu.R.L, &cpu.R.E)
	case 0x6C: //LD L, H
		cpu.LDrr(&cpu.R.L, &cpu.R.H)
	case 0x6D: //LD L, L
		cpu.LDrr(&cpu.R.L, &cpu.R.L)
	case 0x6E: //LD L, (HL)
		cpu.LDr_hl(&cpu.R.L)

	case 0x77: //LD (HL), A
		cpu.LDhl_r(&cpu.R.A)
	case 0x70: //LD (HL), B
		cpu.LDhl_r(&cpu.R.B)
	case 0x71: //LD (HL), C
		cpu.LDhl_r(&cpu.R.C)
	case 0x72: //LD (HL), D
		cpu.LDhl_r(&cpu.R.D)
	case 0x73: //LD (HL), E
		cpu.LDhl_r(&cpu.R.E)
	case 0x74: //LD (HL), H
		cpu.LDhl_r(&cpu.R.H)
	case 0x75: //LD (HL), L
		cpu.LDhl_r(&cpu.R.L)

	case 0x02: //LD (BC), A
		cpu.LDbc_r(&cpu.R.A)
	case 0x12: //LD (DE), A
		cpu.LDde_r(&cpu.R.A)
	case 0xEA: //LD (nn), A
		cpu.LDnn_r(&cpu.R.A)

	case 0x36: //LD (HL), n
		cpu.LDhl_n()

	case 0xF2: //LD A,(C)
		cpu.LDr_ffplusc(&cpu.R.A)
	case 0xE2: //LD (C),A
		cpu.LDffplusc_r(&cpu.R.A)

	case 0x87: //ADD A, A
		cpu.AddA_r(&cpu.R.A)
	case 0x80: //ADD A, B
		cpu.AddA_r(&cpu.R.B)
	case 0x81: //ADD A, C
		cpu.AddA_r(&cpu.R.C)
	case 0x82: //ADD A, D
		cpu.AddA_r(&cpu.R.D)
	case 0x83: //ADD A, E
		cpu.AddA_r(&cpu.R.E)
	case 0x84: //ADD A, H
		cpu.AddA_r(&cpu.R.H)
	case 0x85: //ADD A, L
		cpu.AddA_r(&cpu.R.L)
	case 0x86: //ADD A,(HL)
		cpu.AddA_hl()
	case 0xC6: //ADD A,#
		cpu.AddA_n()
	default:
		log.Fatalf("Invalid/Unknown instruction %X", Opcode)
	}
}

func (cpu *Z80) Step() {
	var Opcode byte = cpu.mmu.ReadByte(cpu.PC)
	cpu.IncrementPC(1)

	cpu.Dispatch(Opcode)

	cpu.MachineCycles.m += cpu.LastInstrCycle.m
	cpu.MachineCycles.t += cpu.LastInstrCycle.t
	cpu.LastInstrCycle.Reset()
}

// INSTRUCTIONS START
//-----------------------------------------------------------------------
//LD r,n
func (cpu *Z80) LDrn(r *byte) {
	log.Println("LD r,n")
	var value byte = cpu.mmu.ReadByte(cpu.PC)
	cpu.IncrementPC(1)

	*r = value

	//set clock values
	cpu.LastInstrCycle.Set(2, 8)
}

//LD r,r
func (cpu *Z80) LDrr(r1 *byte, r2 *byte) {
	log.Println("LD r,r")
	*r1 = *r2

	//set clock values
	cpu.LastInstrCycle.Set(1, 4)
}

//LD r,(HL)
func (cpu *Z80) LDr_hl(r *byte) {
	log.Println("LD r,(HL)")

	var HL Word = cpu.ToWord(cpu.R.H, cpu.R.L)
	var value byte = cpu.mmu.ReadByte(HL)

	*r = value

	//set clock values
	cpu.LastInstrCycle.Set(2, 8)
}

//LD (HL),r
func (cpu *Z80) LDhl_r(r *byte) {
	log.Println("LD (HL),r")
	var HL Word = cpu.ToWord(cpu.R.H, cpu.R.L)
	var value byte = *r

	cpu.mmu.WriteByte(HL, value)

	//set clock values
	cpu.LastInstrCycle.Set(2, 8)
}

//LD (BC),r
func (cpu *Z80) LDbc_r(r *byte) {
	log.Println("LD (BC),r")

	var BC Word = cpu.ToWord(cpu.R.B, cpu.R.C)
	var value byte = *r

	cpu.mmu.WriteByte(BC, value)

	//set clock values
	cpu.LastInstrCycle.Set(2, 8)
}

//LD (DE),r
func (cpu *Z80) LDde_r(r *byte) {
	log.Println("LD (DE),r")

	var DE Word = cpu.ToWord(cpu.R.D, cpu.R.E)
	var value byte = *r

	cpu.mmu.WriteByte(DE, value)

	//set clock values
	cpu.LastInstrCycle.Set(2, 8)
}

//LD nn,r
func (cpu *Z80) LDnn_r(r *byte) {
	log.Println("LD nn,r")
	var resultAddr Word = cpu.mmu.ReadWord(cpu.PC)
	cpu.IncrementPC(2)

	cpu.mmu.WriteByte(resultAddr, *r)

	cpu.LastInstrCycle.Set(4, 16)
}

//LD (HL),n
func (cpu *Z80) LDhl_n() {
	log.Println("LD (HL),n")
	var HL Word = cpu.ToWord(cpu.R.H, cpu.R.L)
	var value byte = cpu.mmu.ReadByte(cpu.PC)
	cpu.IncrementPC(1)

	cpu.mmu.WriteByte(HL, value)

	//set clock values
	cpu.LastInstrCycle.Set(3, 12)
}

//LD r, (BC)
func (cpu *Z80) LDr_bc(r *byte) {
	log.Println("LD r,(BC)")

	var BC Word = cpu.ToWord(cpu.R.B, cpu.R.C)
	var value byte = cpu.mmu.ReadByte(BC)

	*r = value

	//set clock values
	cpu.LastInstrCycle.Set(2, 8)
}

//LD r, (DE)
func (cpu *Z80) LDr_de(r *byte) {
	log.Println("LD r,(DE)")

	var DE Word = cpu.ToWord(cpu.R.D, cpu.R.E)
	var value byte = cpu.mmu.ReadByte(DE)

	*r = value

	//set clock values
	cpu.LastInstrCycle.Set(2, 8)
}

//LD r, nn
func (cpu *Z80) LDr_nn(r *byte) {
	log.Println("LD r,(nn)")

	//read 2 bytes from PC
	var nn Word = cpu.mmu.ReadWord(cpu.PC)
	cpu.IncrementPC(2)

	var value byte = cpu.mmu.ReadByte(nn)
	*r = value

	//set clock values
	cpu.LastInstrCycle.Set(4, 16)
}

//LD r,(C)
func (cpu *Z80) LDr_ffplusc(r *byte) {
	log.Println("LD r,(C)")
	var valueAddr Word = 0xFF00 + Word(cpu.R.C)
	*r = cpu.mmu.ReadByte(valueAddr)

	//set clock values
	cpu.LastInstrCycle.Set(2, 8)
}

//LD (C),r
func (cpu *Z80) LDffplusc_r(r *byte) {
	log.Println("LD (C),r")
	var valueAddr Word = 0xFF00 + Word(cpu.R.C)
	cpu.mmu.WriteByte(valueAddr, *r)

	//set clock values
	cpu.LastInstrCycle.Set(2, 8)
}

//ADD A,r
func (cpu *Z80) AddA_r(r *byte) {
	log.Println("ADD A,r")
	var oldA byte = cpu.R.A
	cpu.R.A += *r

	cpu.ResetFlags()

	//set carry flag
	if (oldA + *r) < oldA {
		cpu.R.F = cpu.R.F ^ 0x40
	}

	//set zero flag
	if cpu.R.A == 0x00 {
		cpu.R.F = cpu.R.F ^ 0x80
	}

	//set clock values
	cpu.LastInstrCycle.Set(1, 4)
}

//ADD A,(HL)
func (cpu *Z80) AddA_hl() {
	log.Println("ADD A,(HL)")
	var HL Word = cpu.ToWord(cpu.R.H, cpu.R.L)
	var value byte = cpu.mmu.ReadByte(HL)

	var oldA byte = cpu.R.A
	cpu.R.A += value

	cpu.ResetFlags()

	//set carry flag
	if (oldA + value) < oldA {
		cpu.R.F = cpu.R.F ^ 0x40
	}

	//set zero flag
	if cpu.R.A == 0x00 {
		cpu.R.F = cpu.R.F ^ 0x80
	}

	//set clock values
	cpu.LastInstrCycle.Set(2, 8)
}

//ADD A,n
func (cpu *Z80) AddA_n() {
	log.Println("ADD A,n")
	var value byte = cpu.mmu.ReadByte(cpu.PC)
	cpu.IncrementPC(1)

	var oldA byte = cpu.R.A
	cpu.R.A += value

	cpu.ResetFlags()

	//set carry flag
	if (oldA + value) < oldA {
		cpu.R.F = cpu.R.F ^ 0x40
	}

	//set zero flag
	if cpu.R.A == 0x00 {
		cpu.R.F = cpu.R.F ^ 0x80
	}

	//set clock values
	cpu.LastInstrCycle.Set(2, 8)
}

//-----------------------------------------------------------------------
//INSTRUCTIONS END