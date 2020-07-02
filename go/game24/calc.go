package game24

import (
	"fmt"
	"sort"
	"strconv"
)

const precision float64 = 0.001

type Node struct {
	Left, Right *Node
	Value       float64
	Operator    *Operator
	flatted     bool
}

type Operator struct {
	Level       int
	Value       string
	Commutative bool
}

var PlusOp *Operator = &Operator{Level: 1, Value: "+", Commutative: true}
var SubOp *Operator = &Operator{Level: 1, Value: "-", Commutative: false}
var MultiOp *Operator = &Operator{Level: 4, Value: "*", Commutative: true}
var DivOp *Operator = &Operator{Level: 4, Value: "/", Commutative: false}

func Calc(n1, n2, n3, n4, target float64) {
	trees := allTrees(n1, n2, n3, n4)
	trees = calcFilter(trees, target)
	trees = revsFilter(trees)
	for _, tree := range trees {
		fmt.Println(inOrderFormatTree(tree))
	}
}

func revsFilter(trees []*Node) []*Node {
	dict := make(map[string]*Node)
	for _, tree := range trees {
		// flatTree(tree)
		revsTree(tree)
		// unflatTree(tree)
		dict[inOrderFormatTree(tree)] = tree
	}

	result := make([]*Node, 0, len(dict))
	for _, tree := range dict {
		result = append(result, tree)
	}

	return result
}

func calcFilter(trees []*Node, target float64) []*Node {
	results := make([]*Node, 0, len(trees))
	for _, tree := range trees {
		val, err := calcTree(tree)
		if err == nil && target-precision <= val && target+precision >= val {
			results = append(results, tree)
		}
	}
	return results
}

func allTrees(n1, n2, n3, n4 float64) []*Node {
	opCombs := AllOpCombines()
	permutations := AllPermutations([]float64{n1, n2, n3, n4})
	trees := make([]*Node, 0, len(opCombs)*len(permutations)*5)
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

func revsTree(tree *Node) {
	if !tree.IsLeaf() {
		revsTree(tree.Left)
		revsTree(tree.Right)

		if tree.Operator.Commutative {
			if !tree.Right.IsLeaf() {
				if !tree.Left.IsLeaf() {
					if tree.Left.Operator.Level < tree.Right.Operator.Level {
						tree.Left, tree.Right = tree.Right, tree.Left
					}
				} else {
					if tree.Operator.Level < tree.Right.Operator.Level {
						tree.Left, tree.Right = tree.Right, tree.Left
					}
				}
			} else {
				if !tree.Left.IsLeaf() {
					if tree.Left.Operator.Level < tree.Operator.Level {
						tree.Left, tree.Right = tree.Right, tree.Left
					}
				} else {
					if tree.Left.Value < tree.Right.Value {
						tree.Left, tree.Right = tree.Right, tree.Left
					}
				}
			}
		}
	}
}

func calcTree(tree *Node) (float64, error) {
	if !tree.IsLeaf() {
		leftVal, err := calcTree(tree.Left)
		if err != nil {
			return 0, err
		}
		rightVal, err := calcTree(tree.Right)
		if err != nil {
			return 0, err
		}
		switch tree.Operator.Value {
		case "+":
			tree.Value = leftVal + rightVal
		case "-":
			tree.Value = leftVal - rightVal
		case "*":
			tree.Value = leftVal * rightVal
		case "/":
			if rightVal == 0 {
				return 0, fmt.Errorf("divied zero")
			} else {
				tree.Value = leftVal / rightVal
			}
		default:
			return 0, fmt.Errorf("unsupported operator %s", tree.Operator.Value)
		}
	}
	return tree.Value, nil
}

func flatTree(tree *Node) {
	if !tree.IsLeaf() {
		flatTree(tree.Left)
		flatTree(tree.Right)
		if tree.Operator.Value == "-" {
			tree.Right.Value *= -1
			tree.Operator, _ = GetOperator("+")
			tree.Right.flatted = true
		} else if tree.Operator.Value == "/" {
			tree.Right.Value = 1.0 / tree.Right.Value
			tree.Operator, _ = GetOperator("*")
			tree.Right.flatted = true
		}
	}
}

func unflatTree(tree *Node) {
	if !tree.IsLeaf() {
		unflatTree(tree.Left)
		unflatTree(tree.Right)
		if tree.Operator.Value == "+" {
			if tree.Right.flatted {
				tree.Right.Value *= -1
				tree.Right.flatted = false
				tree.Operator, _ = GetOperator("-")
			}
		} else if tree.Operator.Value == "*" && tree.Right.flatted {
			if tree.Right.flatted {
				tree.Right.Value = 1.0 / tree.Right.Value
				tree.Right.flatted = false
				tree.Operator, _ = GetOperator("/")
			}
		}
	}
}

func inOrderFormatTree(tree *Node) string {
	if !tree.IsLeaf() {
		left, right := inOrderFormatTree(tree.Left), inOrderFormatTree(tree.Right)
		if !tree.Left.IsLeaf() && tree.Left.Operator.Level < tree.Operator.Level {
			left = fmt.Sprintf("(%s)", left)
		}
		if !tree.Right.IsLeaf() {
			if tree.Right.Operator.Level < tree.Operator.Level {
				right = fmt.Sprintf("(%s)", right)
			} else if tree.Right.Operator.Level == tree.Operator.Level {
				if !tree.Operator.Commutative {
					right = fmt.Sprintf("(%s)", right)
				}
			}
		}
		return fmt.Sprintf("%s %s %s", left, tree.Operator.Value, right)
	}

	return strconv.Itoa(int(tree.Value))
}

func printTree(tree *Node) string {
	if !tree.IsLeaf() {
		left, right := printTree(tree.Left), printTree(tree.Right)
		if !tree.Left.IsLeaf() && tree.Left.Operator.Level < tree.Operator.Level {
			left = fmt.Sprintf("(%s)", left)
		}
		if !tree.Right.IsLeaf() {
			if tree.Right.Operator.Level < tree.Operator.Level {
				right = fmt.Sprintf("(%s)", right)
			} else if tree.Right.Operator.Level == tree.Operator.Level {
				if !tree.Operator.Commutative {
					right = fmt.Sprintf("(%s)", right)
				}
			}
		}
		return fmt.Sprintf("%s %s %s", left, tree.Operator.Value, right)
	}

	return fmt.Sprintf("%f", tree.Value)
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

func AllOpCombines() [][]*Operator {
	results := make([][]*Operator, 0)
	ops := []string{"+", "-", "*", "/"}
	for _, o1 := range ops {
		for _, o2 := range ops {
			for _, o3 := range ops {
				op1, _ := GetOperator(o1)
				op2, _ := GetOperator(o2)
				op3, _ := GetOperator(o3)
				results = append(results, []*Operator{op1, op2, op3})
			}
		}
	}
	return results
}

func TreeType1(o1, o2, o3 *Operator, n1, n2, n3, n4 float64) *Node {
	node1, node2, node3, node4 := &Node{Value: n1}, &Node{Value: n2}, &Node{Value: n3}, &Node{Value: n4}
	node5 := &Node{Left: node1, Right: node2, Operator: o1}
	node6 := &Node{Left: node5, Right: node3, Operator: o2}
	node7 := &Node{Left: node6, Right: node4, Operator: o3}
	return node7
}

func TreeType2(o1, o2, o3 *Operator, n1, n2, n3, n4 float64) *Node {
	node1, node2, node3, node4 := &Node{Value: n1}, &Node{Value: n2}, &Node{Value: n3}, &Node{Value: n4}
	node5 := &Node{Left: node2, Right: node3, Operator: o1}
	node6 := &Node{Left: node1, Right: node5, Operator: o2}
	node7 := &Node{Left: node6, Right: node4, Operator: o3}
	return node7
}

func TreeType3(o1, o2, o3 *Operator, n1, n2, n3, n4 float64) *Node {
	node1, node2, node3, node4 := &Node{Value: n1}, &Node{Value: n2}, &Node{Value: n3}, &Node{Value: n4}
	node5 := &Node{Left: node1, Right: node2, Operator: o1}
	node6 := &Node{Left: node3, Right: node4, Operator: o2}
	node7 := &Node{Left: node5, Right: node6, Operator: o3}
	return node7
}

func TreeType4(o1, o2, o3 *Operator, n1, n2, n3, n4 float64) *Node {
	node1, node2, node3, node4 := &Node{Value: n1}, &Node{Value: n2}, &Node{Value: n3}, &Node{Value: n4}
	node5 := &Node{Left: node3, Right: node4, Operator: o1}
	node6 := &Node{Left: node2, Right: node5, Operator: o2}
	node7 := &Node{Left: node1, Right: node6, Operator: o3}
	return node7
}

func TreeType5(o1, o2, o3 *Operator, n1, n2, n3, n4 float64) *Node {
	node1, node2, node3, node4 := &Node{Value: n1}, &Node{Value: n2}, &Node{Value: n3}, &Node{Value: n4}
	node5 := &Node{Left: node1, Right: node2, Operator: o1}
	node6 := &Node{Left: node5, Right: node3, Operator: o2}
	node7 := &Node{Left: node4, Right: node6, Operator: o3}
	return node7
}

func GetOperator(op string) (*Operator, error) {
	switch op {
	case "+":
		return PlusOp, nil
	case "-":
		return SubOp, nil
	case "*":
		return MultiOp, nil
	case "/":
		return DivOp, nil
	default:
		return nil, fmt.Errorf("unsupported operator %s", op)
	}
}

func (tree *Node) IsLeaf() bool {
	return tree.Left == nil && tree.Right == nil
}
