package priv

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// Privilege is action type that can be granted to user.
type Privilege int

const ( // resource privileges
	ReadPrivilege Privilege = 1 << iota
	WritePrivilege
	CreateCQPrivilege
	InsertPrivilege
	SelectPrivilege
	DeletePrivilege
	DropPrivilege
)
const ( // global privileges
	ShowUsersPrivilege Privilege = 1 << (iota + 16)
	CreateUserPrivilege
	ShowRolesPrivilege
	CreateRolePrivilege
	GrantPrivilege
	ShowDatabasesPrivilege
	CreateDatabasePrivilege
	ShowSysInfoPrivilege
	SetSysInfoPrivilege
	AuditPrivilege
	ShowCQSPrivilege
)
const ( // grouped privileges
	NoPrivilege           Privilege = 0
	AllResourcePrivileges Privilege = 0xFFFF
	AllGlobalPrivileges   Privilege = 1<<31 - 1

	// ReadGroupPrivileges and WriteGroupPrivileges is define for compatible with old version
	ReadGroupPrivileges  Privilege = ShowDatabasesPrivilege | SelectPrivilege | CreateCQPrivilege
	WriteGroupPrivileges Privilege = CreateDatabasePrivilege | DeletePrivilege | DropPrivilege
)

// String returns a string representation of a Privilege.
func (p Privilege) String() string {
	if p == AllGlobalPrivileges || p == AllResourcePrivileges {
		return privilege2name[p]
	}

	privMask := Privilege(1)
	names := make([]string, 0, 16)
	for privMask != 0 {
		if privMask&p == privMask {
			if name, ok := privilege2name[privMask]; ok {
				names = append(names, name)
			}
		}
		privMask <<= 1
	}

	return strings.Join(names, ", ")
}

// PrivilegeOf find privilege of given name.
func PrivilegeOf(name string) (Privilege, error) {
	p, ok := name2privilege[strings.ToUpper(strings.TrimSpace(name))]
	if !ok {
		return NoPrivilege, fmt.Errorf("unknown privilege '%s'", name)
	}
	return p, nil
}

var privilege2name = map[Privilege]string{
	ReadPrivilege:     "READ",
	WritePrivilege:    "WRITE",
	CreateCQPrivilege: "CREATE CQ",
	InsertPrivilege:   "INSERT",
	SelectPrivilege:   "SELECT",
	DeletePrivilege:   "DELETE",
	DropPrivilege:     "DROP",

	ShowUsersPrivilege:      "SHOW USERS",
	CreateUserPrivilege:     "CREATE USER",
	ShowRolesPrivilege:      "SHOW ROLES",
	CreateRolePrivilege:     "CREATE ROLE",
	GrantPrivilege:          "GRANT",
	ShowDatabasesPrivilege:  "SHOW DATABASES",
	CreateDatabasePrivilege: "CREATE DATABASE",
	ShowSysInfoPrivilege:    "SHOW SYSINFO",
	SetSysInfoPrivilege:     "SET SYSINFO",
	AuditPrivilege:          "AUDIT",
	ShowCQSPrivilege:        "SHOW CQS",

	AllGlobalPrivileges:   "ALL PRIVILEGES",
	AllResourcePrivileges: "ALL PRIVILEGES",
}

var name2privilege = make(map[string]Privilege)

func init() {
	for k, v := range privilege2name {
		name2privilege[v] = k
	}
	name2privilege["ALL"] = AllGlobalPrivileges
	name2privilege["ALL PRIVILEGES"] = AllGlobalPrivileges
}

// ResourcePath represent dot seperated path, eg `"my.db".autogen.cpu`
type ResourcePath struct {
	Segs []string
}

var GlobalResource = NewResourcePath()

func CreateResourcePathUnsafe(resource string) *ResourcePath {
	r, _ := CreateResourcePath(resource)
	return r
}

func CreateResourcePath(resource string) (*ResourcePath, error) {
	if strings.TrimSpace(resource) == "" {
		return NewResourcePath(), nil
	}

	qr := strings.NewReader(string(resource))
	p := NewParser(qr)
	segs, err := p.parseSegmentedIdents()
	if err != nil {
		return nil, err
	}
	return NewResourcePath(segs...), nil
}

func NewResourcePath(segs ...string) *ResourcePath {
	if len(segs) > 0 && segs[0] == "" { // first part shall not be empty
		return &ResourcePath{}
	}
	if len(segs) > 2 && segs[1] == "" {
		segs[1] = "autogen" // autogen may have omited
	}

	return &ResourcePath{segs}
}

func (r *ResourcePath) String() string {
	var buf bytes.Buffer

	for i, seg := range r.Segs {
		if i != 0 {
			buf.WriteString(".")
		}
		buf.WriteString(QuoteIdent(seg))
	}

	return buf.String()
}

// PrivilegeSet manipulates privilges on resources.
type PrivilegeSet interface {
	// SetAll set full privileges to privilege set.
	SetAll()
	// ClearAll clear all privileges from privilege set.
	ClearAll()
	// AddGlobal add some privileges to global resource.
	AddGlobal(privilege Privilege)
	// DeleteGlobal delete some privileges from all resources.
	DeleteGlobal(privilege Privilege)
	// Add some privileges to given resource.
	Add(resource *ResourcePath, privilege Privilege)
	// Delete some privielges from all resources under the given resource name.
	Delete(resource *ResourcePath, privilege Privilege)
	// UnionWith combine all privileges of 2 privilege sets.
	UnionWith(s PrivilegeSet)
	// DifferentWith delete all privileges from the given privilege set.
	DifferentWith(s PrivilegeSet)
	// GlobalContain checks if root node have the given privileges.
	// It should be noticed that this does not mean have privileges on every resources.
	GlobalContain(privilege Privilege) bool
	// Contain checks if privileges set contains privileges on the given resource.
	Contain(resource *ResourcePath, privilege Privilege) bool
	// Contains checks if the privilege set contains all privileges from another set.
	Contains(s PrivilegeSet) bool
	// Powerless check set if don't has any privilege.
	Powerless() bool
}

// NewPrivilegeTree create an empty privilege tree
func NewPrivilegeTree() *PrivilegeTree {
	return &PrivilegeTree{Tree: make(map[string]*PrivilegeTree)}
}

// PrivilegeTree is an implementation of PrivilegeSet interface.
type PrivilegeTree struct {
	Privilege Privilege
	Tree      map[string]*PrivilegeTree
}

func (t *PrivilegeTree) implPrivilegeSet() {
	var _ PrivilegeSet = (*PrivilegeTree)(nil)
}

// SetAll set full privileges to privilege tree.
func (t *PrivilegeTree) SetAll() {
	t.Privilege = AllGlobalPrivileges
	t.Tree = make(map[string]*PrivilegeTree)
}

// ClearAll clear all privileges from privilege tree.
func (t *PrivilegeTree) ClearAll() {
	t.Privilege = NoPrivilege
	t.Tree = make(map[string]*PrivilegeTree)
}

// AddGlobal add some privileges to global resource.
func (t *PrivilegeTree) AddGlobal(privilege Privilege) {
	t.Privilege |= privilege
	for _, v := range t.Tree {
		if v != nil {
			v.DeleteGlobal(privilege)
		}
	}
}

// DeleteGlobal delete some privileges from all resources.
func (t *PrivilegeTree) DeleteGlobal(privilege Privilege) {
	t.Privilege &^= privilege
	for _, v := range t.Tree {
		if v != nil {
			v.DeleteGlobal(privilege)
		}
	}
	t.prune()
}

// Add some privileges to given resource
func (t *PrivilegeTree) Add(resource *ResourcePath, privilege Privilege) {
	privilege &= AllResourcePrivileges

	sum := NoPrivilege
	for _, seg := range resource.Segs {
		sum ^= t.Privilege
		if t.Tree[seg] == nil {
			t.Tree[seg] = NewPrivilegeTree()
		}
		t = t.Tree[seg]
	}
	result := ^sum            // make sure (result ^ sum) & privilege = privilege
	result &= privilege       // clear all bits unrelated with incoming privilege
	t.Privilege &^= privilege // reset related bits to zero
	t.Privilege |= result     // set with new value

	for _, v := range t.Tree {
		if v != nil {
			v.DeleteGlobal(privilege)
		}
	}
	t.prune()
}

// Delete some privielges from all resources under the given resource name.
func (t *PrivilegeTree) Delete(resource *ResourcePath, privilege Privilege) {
	privilege &= AllResourcePrivileges

	sum := NoPrivilege
	for _, seg := range resource.Segs {
		sum ^= t.Privilege
		if t.Tree[seg] == nil {
			t.Tree[seg] = NewPrivilegeTree()
		}
		t = t.Tree[seg]
	}
	sum &= privilege          // clear all bits unrelated with incoming privilege
	t.Privilege &^= privilege // reset related bits to zero
	t.Privilege |= sum        // set with new value, (t.Privilege & privilege) ^ sum = 0

	for _, v := range t.Tree {
		if v != nil {
			v.DeleteGlobal(privilege)
		}
	}

	t.prune()
}

// UnionWith combine all privileges of 2 privilege trees.
func (t *PrivilegeTree) UnionWith(s PrivilegeSet) {
	t.union(NoPrivilege, NoPrivilege, NoPrivilege, s.(*PrivilegeTree), true)
	t.prune()
}

func (t *PrivilegeTree) union(tsum, ssum, newsum Privilege, s *PrivilegeTree, root bool) {
	if t == nil || s == nil {
		return
	}

	tsum ^= t.Privilege
	ssum ^= s.Privilege
	current := (tsum | ssum) ^ newsum
	if !root {
		current &= AllResourcePrivileges
	}
	newsum ^= current
	t.Privilege = current

	for k, v := range s.Tree {
		if tt := t.Tree[k]; tt == nil && v != nil {
			t.Tree[k] = NewPrivilegeTree()
		}
		t.Tree[k].union(tsum, ssum, newsum, v, false)
	}
	for k, v := range t.Tree {
		if _, exist := s.Tree[k]; !exist {
			v.union(tsum, ssum, newsum, NewPrivilegeTree(), false)
		}
	}
}

// DifferentWith delete all privileges from the given privilege tree.
func (t *PrivilegeTree) DifferentWith(s PrivilegeSet) {
	t.sub(NoPrivilege, NoPrivilege, NoPrivilege, s.(*PrivilegeTree), true)
	t.prune()
}

func (t *PrivilegeTree) sub(tsum, ssum, newsum Privilege, s *PrivilegeTree, root bool) {
	if t == nil || s == nil {
		return
	}

	tsum ^= t.Privilege
	ssum ^= s.Privilege
	current := tsum&^ssum ^ newsum
	if !root {
		current &= AllResourcePrivileges
	}
	newsum ^= current
	t.Privilege = current

	for k, v := range s.Tree {
		tt := t.Tree[k]
		if tt == nil && v != nil {
			t.Tree[k] = NewPrivilegeTree()
		}

		t.Tree[k].sub(tsum, ssum, newsum, v, false)
	}
	for k, v := range t.Tree {
		if _, exist := s.Tree[k]; !exist {
			v.sub(tsum, ssum, newsum, NewPrivilegeTree(), false)
		}
	}
}

// GlobalContain checks if root node have the given privileges.
// It should be noticed that this does not mean have privileges on every resources.
func (t *PrivilegeTree) GlobalContain(privilege Privilege) bool {
	return privilege&^t.compatibleWithReadWrite(t.Privilege) == NoPrivilege
}

// Contain checks if privileges set contains privileges on the given resource.
func (t *PrivilegeTree) Contain(resource *ResourcePath, privilege Privilege) bool {
	sum := NoPrivilege
	sum ^= t.Privilege
	for _, seg := range resource.Segs {
		if t.Tree[seg] == nil {
			break
		}
		t = t.Tree[seg]
		sum ^= t.Privilege
	}

	return t.compatibleWithReadWrite(sum)&privilege == privilege
}

// Contains checks if the privilege set contains all privileges from another set.
func (t *PrivilegeTree) Contains(s PrivilegeSet) bool {
	s.(*PrivilegeTree).sub(NoPrivilege, NoPrivilege, NoPrivilege, t, true)
	return s.Powerless()
}

// tidy free useless memory.
func (t *PrivilegeTree) prune() {
	if t.Powerless() {
		t.ClearAll()
		return
	}
	for _, v := range t.Tree {
		if v != nil {
			v.prune()
		}
	}
}

// read privilege or write privilege in old version equals a group of privileges
// in current version, so should handle read and write privilege especially
func (t *PrivilegeTree) compatibleWithReadWrite(privilege Privilege) Privilege {
	if privilege&ReadPrivilege == ReadPrivilege {
		privilege |= ReadGroupPrivileges
	}
	if privilege&WritePrivilege == WritePrivilege {
		privilege |= WriteGroupPrivileges
	}
	return privilege
}

// if PrivilegeTree don't has any privilege
func (t *PrivilegeTree) Powerless() bool {
	if t.Privilege != NoPrivilege {
		return false
	}

	for _, v := range t.Tree {
		if v != nil {
			if !v.Powerless() {
				return false
			}
		}
	}

	return true
}

/*
String serialize PrivilegeTree to string, use BFS strategy.
One node can be serialized as (name, privilege), an empty node can be
serialized as ().
All children of a node can be bracketed in [], and all nodes of one floor can
be bracketed in {}

Example, a tree as follow:

        ┌- (yourdb, 3)
        |
(, 1) --|              ┌- (daily, 5)
		|              |
		└- (mydb, 2) --|                 ┌- (mem, 7)
					   |                 |
					   └- (autogen, 4) --|
										 |
										 └- (cpu, 6)
will be serialized as:
{[(, 1)]} {[(yourdb, 3)(mydb, 2)]} {[][(daily, 5)(autogen, 4)]} {[][(mem, 7)(cpu, 6)]}
*/

func (t *PrivilegeTree) String() string {
	buf := &bytes.Buffer{}

	floor := []map[string]*PrivilegeTree{{"": t}}
	for len(floor) > 0 {
		newFloor := make([]map[string]*PrivilegeTree, 0)

		buf.WriteString("{")
		for _, f := range floor {
			buf.WriteString("[")
			for k, v := range f {
				if v != nil {
					buf.WriteString(fmt.Sprintf("(%s,%d)", QuoteIdent(k), v.Privilege))
					newFloor = append(newFloor, v.Tree)
				} else {
					buf.WriteString("()")
				}
			}
			buf.WriteString("]")
		}
		buf.WriteString("}")

		floor = newFloor
	}

	return buf.String()
}

// LoadPrivilegeTree unserialize PrivilegeTree from string.
// Uncomplete version, can not handle situation which character '{}[]()' contained
// in node name.
func LoadPrivilegeTree(s string) (*PrivilegeTree, error) {
	s = strings.TrimSpace(s)
	runes := []rune(s)

	floors, err := parseSections(runes, '{', '}')
	if err != nil {
		return nil, err
	}

	for _, floor := range floors {
		children, err := parseSections(floor, '[', ']')
		if err != nil {
			return nil, err
		}
		for _, f := range children {
			fmt.Println(string(f))
		}
	}

	return nil, nil
}

func parseSections(s []rune, start, end rune) ([][]rune, error) {
	var sections [][]rune

	len := len(s)
	var left int
	for left < len-1 {
		if s[left] != start {
			left++
			continue
		}
		var right int
		for right = left + 1; right < len; right++ {
			if s[right] == end {
				break
			}
		}
		if right >= len {
			return nil, errors.New("")
		}

		sections = append(sections, s[left+1:right])
		left = right + 1
	}

	return sections, nil
}
