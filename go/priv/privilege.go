package priv

import (
	"bytes"
	"fmt"
	"strings"
)

// Privilege is action type that can be granted to user.
type Privilege int

const (
	ReadPrivilege     Privilege = 1 << 0
	WritePrivilege    Privilege = 1 << 1
	CreateCQPrivilege Privilege = 1 << 2
	InsertPrivilege   Privilege = 1 << 3
	SelectPrivilege   Privilege = 1 << 4
	DeletePrivilege   Privilege = 1 << 5
	DropPrivilege     Privilege = 1 << 6

	ShowUsersPrivilege      Privilege = 1 << 16
	CreateUserPrivilege     Privilege = 1 << 17
	ShowRolesPrivilege      Privilege = 1 << 18
	CreateRolePrivilege     Privilege = 1 << 19
	GrantPrivilege          Privilege = 1 << 20
	ShowDatabasesPrivilege  Privilege = 1 << 21
	CreateDatabasePrivilege Privilege = 1 << 22
	ShowSysInfoPrivilege    Privilege = 1 << 23
	SetSysInfoPrivilege     Privilege = 1 << 24
	AuditPrivilege          Privilege = 1 << 25
	ShowQueriesPrivilege    Privilege = 1 << 26

	NoPrivilege           Privilege = 0
	AllGlobalPrivileges   Privilege = 0xFFFFFFFF
	AllResourcePrivileges Privilege = 0xFFFF
)

// String returns a string representation of a Privilege.
func (p Privilege) String() string {
	if p == AllGlobalPrivileges || p == AllResourcePrivileges {
		return privilege2name[p]
	}

	privMask := Privilege(1)
	privs := make([]string, 0, 16)
	for privMask != 0 {
		if privMask&p == privMask {
			if name, ok := privilege2name[privMask]; ok {
				privs = append(privs, name)
			}
		}
		privMask <<= 1
	}

	return strings.Join(privs, ", ")
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
	ShowQueriesPrivilege:    "SHOW QUERIES",

	AllGlobalPrivileges:   "ALL PRIVILEGES",
	AllResourcePrivileges: "ALL PRIVILEGES",
}

var name2privilege = map[string]Privilege{
	"READ":      ReadPrivilege,
	"WRITE":     WritePrivilege,
	"CREATE CQ": CreateCQPrivilege,
	"INSERT":    InsertPrivilege,
	"SELECT":    SelectPrivilege,
	"DELETE":    DeletePrivilege,
	"DROP":      DropPrivilege,

	"SHOW USERS":      ShowUsersPrivilege,
	"CREATE USER":     CreateUserPrivilege,
	"SHOW ROLES":      ShowRolesPrivilege,
	"CREATE ROLE":     CreateRolePrivilege,
	"GRANT":           GrantPrivilege,
	"SHOW DATABASES":  ShowDatabasesPrivilege,
	"CREATE DATABASE": CreateDatabasePrivilege,
	"SHOW SYSINFO":    ShowSysInfoPrivilege,
	"SET SYSINFO":     SetSysInfoPrivilege,
	"AUDIT":           AuditPrivilege,
	"SHOW QUERIES":    ShowQueriesPrivilege,

	"ALL":            AllGlobalPrivileges,
	"ALL PRIVILEGES": AllGlobalPrivileges,
}

// ResourcePath represent dot seperated path, eg `"my.db".autogen.cpu`
type ResourcePath string

func (r ResourcePath) Idents() ([]string, error) {
	qr := strings.NewReader(string(r))
	p := NewParser(qr)
	idents, err := p.parseSegmentedIdents()
	if err != nil {
		return nil, err
	}
	for _, ident := range idents {
		if strings.TrimSpace(ident) == "" {
			return nil, fmt.Errorf("invalid resource path %s, expect a full path", r)
		}
	}
	return idents, nil
}

// NewPrivilegeTree create an empty privilege tree
func NewPrivilegeTree() *PrivilegeTree {
	return &PrivilegeTree{}
}

// PrivilegeTree is a collection of privileges on some resources.
type PrivilegeTree struct {
	Privilege Privilege
	Tree      map[string]*PrivilegeTree
}

// SetAll set full privileges to privilege tree.
func (t *PrivilegeTree) SetAll() {
	t.Privilege = AllGlobalPrivileges
	t.Tree = nil
}

// ClearAll clear all privileges from privilege tree.
func (t *PrivilegeTree) ClearAll() {
	t.Privilege = NoPrivilege
	t.Tree = nil
}

// AddGlobal add some privileges to global resource.
func (t *PrivilegeTree) AddGlobal(privilege Privilege) {
	t.Privilege |= privilege
	if t.Tree == nil {
		return
	}
	for _, v := range t.Tree {
		if v != nil {
			v.DeleteGlobal(privilege)
		}
	}
}

// DeleteGlobal delete some privileges from all resources.
func (t *PrivilegeTree) DeleteGlobal(privilege Privilege) {
	t.Privilege &^= privilege
	if t.Tree == nil {
		return
	}
	for _, v := range t.Tree {
		if v != nil {
			v.DeleteGlobal(privilege)
		}
	}
	t.prune()
}

// Add add some privileges to given resource
func (t *PrivilegeTree) Add(resource string, privilege Privilege) error {
	privilege &= AllResourcePrivileges

	sum := NoPrivilege
	idents, err := ResourcePath(resource).Idents()
	if err != nil {
		return err
	}
	for _, ident := range idents {
		sum ^= t.Privilege
		if t.Tree == nil {
			t.Tree = make(map[string]*PrivilegeTree)
		}
		if t.Tree[ident] == nil {
			t.Tree[ident] = NewPrivilegeTree()
		}
		t = t.Tree[ident]
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

	return nil
}

// Delete delete some privielges from all resources under the given resource name.
func (t *PrivilegeTree) Delete(resource string, privilege Privilege) error {
	privilege &= AllResourcePrivileges

	sum := NoPrivilege
	idents, err := ResourcePath(resource).Idents()
	if err != nil {
		return err
	}
	for _, ident := range idents {
		sum ^= t.Privilege
		if t.Tree == nil {
			t.Tree = make(map[string]*PrivilegeTree)
		}
		if t.Tree[ident] == nil {
			t.Tree[ident] = NewPrivilegeTree()
		}
		t = t.Tree[ident]
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

	return nil
}

// UnionWith combine all privileges of 2 privilege trees.
func (t *PrivilegeTree) UnionWith(s *PrivilegeTree) {
	t.union(NoPrivilege, NoPrivilege, NoPrivilege, s, true)
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

	if s.Tree == nil {
		return
	}
	for k, v := range s.Tree {
		if t.Tree == nil {
			t.Tree = make(map[string]*PrivilegeTree)
		}
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
func (t *PrivilegeTree) DifferentWith(s *PrivilegeTree) {
	t.sub(NoPrivilege, NoPrivilege, NoPrivilege, s, true)
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

	if s.Tree == nil {
		return
	}
	for k, v := range s.Tree {
		if t.Tree == nil {
			t.Tree = make(map[string]*PrivilegeTree)
		}
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
	privilege &^= t.Privilege
	return privilege == NoPrivilege
}

// Contain checks if privileges set contains privileges on the given resource.
func (t *PrivilegeTree) Contain(resource string, privilege Privilege) (bool, error) {
	idents, err := ResourcePath(resource).Idents()
	if err != nil {
		return false, err
	}

	sum := NoPrivilege
	sum ^= t.Privilege
	for _, ident := range idents {
		if t.Tree == nil {
			break
		}
		if t.Tree[ident] == nil {
			break
		}
		t = t.Tree[ident]
		sum ^= t.Privilege
	}
	sum &= privilege

	return sum == privilege, nil
}

// Contains checks if the privilege set contains all privileges from another set.
func (t *PrivilegeTree) Contains(s *PrivilegeTree) bool {
	s.sub(NoPrivilege, NoPrivilege, NoPrivilege, t, true)
	return s.powerless()
}

func (t *PrivilegeTree) containSub(tsum, ssum, newsum Privilege, s *PrivilegeTree, root bool) bool {
	if t == nil || s == nil {
		return true
	}

	tsum ^= t.Privilege
	ssum ^= s.Privilege
	current := tsum ^ ssum ^ newsum
	if !root {
		current &= AllResourcePrivileges
	}
	newsum ^= current
	t.Privilege = current

	if t.Privilege != NoPrivilege {
		return false
	}

	if s.Tree == nil {
		return true
	}
	for k, v := range s.Tree {
		if t.Tree == nil {
			t.Tree = make(map[string]*PrivilegeTree)
		}
		tt := t.Tree[k]
		if tt == nil && v != nil {
			t.Tree[k] = NewPrivilegeTree()
		}

		ret := t.Tree[k].containSub(tsum, ssum, newsum, v, false)
		if !ret {
			return ret
		}
	}

	return true
}

// tidy free useless memory.
func (t *PrivilegeTree) prune() {
	if t.powerless() {
		t.ClearAll()
		return
	}
	if t.Tree == nil {
		return
	}
	if len(t.Tree) == 0 {
		t.Tree = nil
		return
	}
	for _, v := range t.Tree {
		if v != nil {
			v.prune()
		}
	}
}

// if PrivilegeTree don't has any privilege
func (t *PrivilegeTree) powerless() bool {
	if t.Privilege != NoPrivilege {
		return false
	}
	if t.Tree == nil {
		return true
	}

	for _, v := range t.Tree {
		if v != nil {
			if !v.powerless() {
				return false
			}
		}
	}

	return true
}

// String serialize PrivilegeTree to string, just for debug.
func (t *PrivilegeTree) String() string {
	buf := &bytes.Buffer{}
	t.string(buf, "")
	return buf.String()
}

func (t *PrivilegeTree) string(buf *bytes.Buffer, indent string) {
	buf.WriteString(fmt.Sprintf("%s%s\n", indent, t.Privilege))
	indent += "    "
	for k, v := range t.Tree {
		buf.WriteString(fmt.Sprintf("%s%s\n", indent, k))
		v.string(buf, indent)
	}
}
