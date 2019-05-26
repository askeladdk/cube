package cube

import "errors"

func reallycrudessa(proc *Procedure) (*Procedure, error) {
	ssadef := func(localidx int) int {
		local := &proc.locals[localidx]
		ssaregidx := len(proc.ssaregs)
		proc.ssaregs = append(proc.ssaregs, SSAReg{
			local:      local,
			generation: local.generations,
		})
		local.generations += 1
		local.lastssareg = ssaregidx
		return ssaregidx
	}

	ssause := func(localidx int) int {
		local := proc.locals[localidx]
		return local.lastssareg
	}

	for _, blk := range proc.blocks {
		for localidx, _ := range proc.locals {
			ssaregidx := ssadef(localidx)
			blk.ssaparams = append(blk.ssaparams, ssaregidx)
		}

		for i, _ := range blk.instructions {
			insr := &blk.instructions[i]

			op0 := insr.operands[0]
			op1 := insr.operands[1]
			op2 := insr.operands[2]

			if otype, val := op2.unpack(); otype == operandType_LOC {
				insr.operands[2] = operandReg(ssause(val))
			}

			if otype, val := op1.unpack(); otype == operandType_LOC {
				insr.operands[1] = operandReg(ssause(val))
			}

			if otype, val := op0.unpack(); otype == operandType_LOC {
				insr.operands[0] = operandReg(ssadef(val))
			} else if otype != operandType_NIL {
				return nil, errors.New("invalid destination type")
			}
		}

		for localidx, _ := range proc.locals {
			ssaregidx := ssause(localidx)
			for i, succ := range blk.successors {
				if succ != nil {
					blk.jmpargs[i] = append(blk.jmpargs[i], ssaregidx)
				}
			}
		}

		if otype, val := blk.jmpretval.unpack(); otype == operandType_LOC {
			blk.jmpretval = operandReg(ssause(val))
		}
	}

	return proc, nil
}
