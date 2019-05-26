package cube

func reachable(root *BasicBlock, blocks []*BasicBlock) []*BasicBlock {
	visited := map[*BasicBlock]struct{}{}
	var result []*BasicBlock
	var recurse func(*BasicBlock)
	recurse = func(blk *BasicBlock) {
		if _, hasvisited := visited[blk]; !hasvisited {
			visited[blk] = struct{}{}
			result = append(result, blk)
			for _, succ := range blk.successors {
				if succ != nil {
					recurse(succ)
				}
			}
		}
	}
	recurse(root)
	return result
}

func predecessors(blocks []*BasicBlock) []*BasicBlock {
	for _, blk := range blocks {
		blk.predecessors = nil
	}

	for _, blk := range blocks {
		s0 := blk.successors[0]
		s1 := blk.successors[1]

		if s0 != nil {
			s0.predecessors = append(s0.predecessors, blk)
		}

		if s1 != nil && s0 != s1 {
			s1.predecessors = append(s1.predecessors, blk)
		}
	}

	return blocks
}

// tarjan's strongly connected components algorithm
func topologicalSort(blocks []*BasicBlock) []*BasicBlock {
	var stack []*BasicBlock
	var result []*BasicBlock

	index := 0
	onstack := map[*BasicBlock]struct{}{}
	lowlinks := map[*BasicBlock]int{}
	indices := map[*BasicBlock]int{}

	var strongconnect func(*BasicBlock)
	strongconnect = func(blk *BasicBlock) {
		indices[blk] = index
		lowlinks[blk] = index
		onstack[blk] = struct{}{}
		stack = append(stack, blk)
		index += 1

		for _, succ := range blk.successors {
			if succ == nil {
				continue
			} else if _, visited := indices[succ]; !visited {
				strongconnect(succ)
				if lowlinks[succ] < lowlinks[blk] {
					lowlinks[blk] = lowlinks[succ]
				}
			} else if _, isonstack := onstack[succ]; isonstack {
				if indices[succ] < lowlinks[blk] {
					lowlinks[blk] = indices[succ]
				}
			}
		}

		sccomponent := 0
		if lowlinks[blk] == indices[blk] {
			var item *BasicBlock
			end := len(stack)
			for ok := true; ok; ok = item != blk {
				end -= 1
				item = stack[end]
				item.sccomponent = sccomponent
				result = append(result, item)
				delete(onstack, item)
			}
			stack = stack[:end]
			sccomponent += 1
		}
	}

	for _, blk := range blocks {
		if _, visited := indices[blk]; !visited {
			strongconnect(blk)
		}
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}

func Pass_BuildCFG(proc *Procedure) *Procedure {
	blks1 := reachable(proc.entryPoint, proc.blocks)
	blks2 := predecessors(blks1)
	proc.blocks = topologicalSort(blks2)
	return proc
}
