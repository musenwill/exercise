package game24

import "sort"

type Node struct {
	Left, Right *Node
	Value       float64
	Operator    string
}

func allTrees(n1, n2, n3, n4 float64) []*Node {
	opCombs := AllOpCombines()
	permutations := AllPermutations([]float64{n1, n2, n3, n4})

	trees := make([]*Node, len(opCombs)*len(permutations)*5)
	for _, ops := range opCombs {
		for _, per := range permutations {
			t1 := TreeType1(ops[0], ops[1], ops[2], per[0], per[1], per[2], per[3])
			t2 := TreeType2(ops[0], ops[1], ops[2], per[0], per[1], per[2], per[3])
			t3 := TreeType3(ops[0], ops[1], ops[2], per[0], per[1], per[2], per[3])
			t4 := TreeType4(ops[0], ops[1], ops[2], per[0], per[1], per[2], per[3])
			t5 := TreeType5(ops[0], ops[1], ops[2], per[0], per[1], per[2], per[3])
			trees = append(trees, t1, t2, t3, t4, t5)
		}
	}

	return trees
}

func AllPermutations(nums []float64) [][]float64 {
	if len(nums) == 0 {
		return nil
	}
	result := [][]float64{}
	sort.Float64s(nums)
	permutationHelper(nums, 0, &result)

	return result
}

func permutationHelper(nums []float64, i int, result *[][]float64) {
	n := len(nums)
	if i == n-1 {
		tmp := make([]float64, n)
		copy(tmp, nums)
		*result = append(*result, tmp)
		return
	}
	for k := i; k < n; k++ {
		if k != i && nums[k] == nums[i] {
			continue
		}
		nums[k], nums[i] = nums[i], nums[k]
		permutationHelper(nums, i+1, result)
	}
	for k := n - 1; k > i; k-- {
		nums[i], nums[k] = nums[k], nums[i]
	}
}

func AllOpCombines() [][]string {
	results := make([][]string, 0)
	ops := []string{"+", "-", "*", "/"}
	for _, o1 := range ops {
		for _, o2 := range ops {
			for _, o3 := range ops {
				results = append(results, []string{o1, o2, o3})
			}
		}
	}
	return results
}

func TreeType1(o1, o2, o3 string, n1, n2, n3, n4 float64) *Node {
	node1, node2, node3, node4 := &Node{Value: n1}, &Node{Value: n2}, &Node{Value: n3}, &Node{Value: n4}
	node5 := &Node{Left: node1, Right: node2, Operator: o1}
	node6 := &Node{Left: node5, Right: node3, Operator: o2}
	node7 := &Node{Left: node6, Right: node4, Operator: o3}
	return node7
}

func TreeType2(o1, o2, o3 string, n1, n2, n3, n4 float64) *Node {
	node1, node2, node3, node4 := &Node{Value: n1}, &Node{Value: n2}, &Node{Value: n3}, &Node{Value: n4}
	node5 := &Node{Left: node2, Right: node3, Operator: o1}
	node6 := &Node{Left: node1, Right: node5, Operator: o2}
	node7 := &Node{Left: node6, Right: node4, Operator: o3}
	return node7
}

func TreeType3(o1, o2, o3 string, n1, n2, n3, n4 float64) *Node {
	node1, node2, node3, node4 := &Node{Value: n1}, &Node{Value: n2}, &Node{Value: n3}, &Node{Value: n4}
	node5 := &Node{Left: node1, Right: node2, Operator: o1}
	node6 := &Node{Left: node3, Right: node4, Operator: o2}
	node7 := &Node{Left: node5, Right: node6, Operator: o3}
	return node7
}

func TreeType4(o1, o2, o3 string, n1, n2, n3, n4 float64) *Node {
	node1, node2, node3, node4 := &Node{Value: n1}, &Node{Value: n2}, &Node{Value: n3}, &Node{Value: n4}
	node5 := &Node{Left: node3, Right: node4, Operator: o1}
	node6 := &Node{Left: node2, Right: node5, Operator: o2}
	node7 := &Node{Left: node1, Right: node6, Operator: o3}
	return node7
}

func TreeType5(o1, o2, o3 string, n1, n2, n3, n4 float64) *Node {
	node1, node2, node3, node4 := &Node{Value: n1}, &Node{Value: n2}, &Node{Value: n3}, &Node{Value: n4}
	node5 := &Node{Left: node1, Right: node2, Operator: o1}
	node6 := &Node{Left: node5, Right: node3, Operator: o2}
	node7 := &Node{Left: node4, Right: node6, Operator: o3}
	return node7
}
