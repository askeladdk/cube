package cube

func reachable(root *basicBlock, blocks []*basicBlock) []*basicBlock {
	visited := map[*basicBlock]struct{}{}
	var result []*basicBlock
	var recurse func(*basicBlock)
	recurse = func(blk *basicBlock) {
		if _, hasvisited := visited[blk]; !hasvisited {
			visited[blk] = struct{}{}
			result = append(result, blk)
		}
	}
	recurse(root)
	return result
}

func predecessors(blocks []*basicBlock) []*basicBlock {
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

func topologicalSort(blocks []*basicBlock) []*basicBlock {
	var stack []*basicBlock
	var result []*basicBlock

	index := 0
	onstack := map[*basicBlock]struct{}{}
	lowlinks := map[*basicBlock]int{}
	indices := map[*basicBlock]int{}

	var strongconnect func(*basicBlock)
	strongconnect = func(blk *basicBlock) {
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
			var item *basicBlock
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
