package priv_test

import (
	"testing"

	"github.com/musenwill/exercise/priv"
)

func TestPrivilegeOf(t *testing.T) {
	var tests = []struct {
		s      string
		hasErr bool
	}{
		{
			s:      "all",
			hasErr: false,
		},
		{
			s:      "all privileges",
			hasErr: false,
		},
		{
			s:      "ALL privileges",
			hasErr: false,
		},
		{
			s:      "create cq",
			hasErr: false,
		},
		{
			s:      "whatever",
			hasErr: true,
		},
	}

	for _, test := range tests {
		_, err := priv.PrivilegeOf(test.s)
		if act, exp := err != nil, test.hasErr; act != exp {
			t.Fatalf("privilege of %s got %v expect %v", test.s, act, exp)
		}
	}
}

func TestPrivilegeGroup(t *testing.T) {
	var tests = []struct {
		privilege priv.Privilege
		group     priv.Privilege
	}{
		{
			privilege: priv.ReadPrivilege,
			group:     priv.AllResourcePrivileges,
		},
		{
			privilege: priv.DropPrivilege,
			group:     priv.AllResourcePrivileges,
		},
		{
			privilege: priv.ReadPrivilege,
			group:     priv.AllGlobalPrivileges,
		},
		{
			privilege: priv.ShowQueriesPrivilege,
			group:     priv.AllGlobalPrivileges,
		},
	}

	for _, test := range tests {
		if test.privilege&test.group != test.privilege {
			t.Fatalf("expect privilege %s be in group %s", test.privilege, test.group)
		}
	}
}

func TestPrivilegeToString(t *testing.T) {
	var tests = []struct {
		privilege priv.Privilege
		str       string
	}{
		{
			privilege: priv.AllGlobalPrivileges,
			str:       "ALL PRIVILEGES",
		},
		{
			privilege: priv.AllResourcePrivileges,
			str:       "ALL PRIVILEGES",
		},
		{
			privilege: priv.ReadPrivilege | priv.DropPrivilege | priv.ShowUsersPrivilege | priv.ShowQueriesPrivilege,
			str:       "READ, DROP, SHOW USERS, SHOW QUERIES",
		},
	}
	for _, test := range tests {
		if test.str != test.privilege.String() {
			t.Fatalf("privilege to string got %s expect %s", test.privilege.String(), test.str)
		}
	}
}

func TestPrivilegeSetAll(t *testing.T) {

	set := priv.NewPrivilegeTree()
	set.SetAll()

	var tests = []struct {
		resource  string
		privilege priv.Privilege
		contains  bool
	}{
		{
			resource:  "notexist",
			privilege: priv.AllGlobalPrivileges,
			contains:  true,
		},
		{
			resource:  "notexist",
			privilege: priv.AllResourcePrivileges,
			contains:  true,
		},
		{
			resource:  "notexist",
			privilege: priv.CreateUserPrivilege | priv.GrantPrivilege | priv.AuditPrivilege | priv.SelectPrivilege,
			contains:  true,
		},
		{
			resource:  "mydb.autogen.cpu",
			privilege: priv.AllResourcePrivileges,
			contains:  true,
		},
		{
			resource:  "mydb.autogen.cpu",
			privilege: priv.SelectPrivilege | priv.InsertPrivilege | priv.DeletePrivilege | priv.DropPrivilege,
			contains:  true,
		},
	}

	for _, test := range tests {
		act, err := set.Contain(test.resource, test.privilege)
		if err != nil {
			t.Fatal(err)
		}
		if act, exp := act, test.contains; act != exp {
			if exp {
				t.Fatalf("expect SetAll() contain privileges [%s] on resource %s", test.privilege, test.resource)
			} else {
				t.Fatalf("expect SetAll() not contain privileges [%s] on resource %s", test.privilege, test.resource)
			}
		}
	}
}

func TestPrivilegeSetClearAll(t *testing.T) {
	set := priv.NewPrivilegeTree()
	set.SetAll()
	set.ClearAll()

	var tests = []struct {
		resource  string
		privilege priv.Privilege
		contains  bool
	}{
		{
			resource:  "notexist",
			privilege: priv.AllGlobalPrivileges,
			contains:  false,
		},
		{
			resource:  "notexist",
			privilege: priv.AllResourcePrivileges,
			contains:  false,
		},
		{
			resource:  "notexist",
			privilege: priv.CreateUserPrivilege | priv.GrantPrivilege | priv.AuditPrivilege,
			contains:  false,
		},
		{
			resource:  "mydb.autogen.cpu",
			privilege: priv.AllResourcePrivileges,
			contains:  false,
		},
		{
			resource:  "mydb.autogen.cpu",
			privilege: priv.SelectPrivilege | priv.InsertPrivilege | priv.DeletePrivilege | priv.DropPrivilege,
			contains:  false,
		},
	}

	for _, test := range tests {
		act, err := set.Contain(test.resource, test.privilege)
		if err != nil {
			t.Fatal(err)
		}
		if act, exp := act, test.contains; act != exp {
			if exp {
				t.Fatalf("expect ClearAll() contain privileges [%s] on resource %s", test.privilege, test.resource)
			} else {
				t.Fatalf("expect ClearAll() not contain privileges [%s] on resource %s", test.privilege, test.resource)
			}
		}
	}
}

func TestIdentsOk(t *testing.T) {
	var tests = []struct {
		s string
		t []string
	}{
		{
			s: `a`,
			t: []string{"a"},
		},
		{
			s: ` a `,
			t: []string{"a"},
		},
		{
			s: `a.b`,
			t: []string{"a", "b"},
		},
		{
			s: `a.b.c`,
			t: []string{"a", "b", "c"},
		},
		{
			s: `a.b.c`,
			t: []string{"a", "b", "c"},
		},
		{
			s: `"a.b.c".b.c`,
			t: []string{"a.b.c", "b", "c"},
		},
		{
			s: `"a.b.c"."b.c".c`,
			t: []string{"a.b.c", "b.c", "c"},
		},
	}
	for _, test := range tests {
		act, err := priv.ResourcePath(test.s).Idents()
		if err != nil {
			t.Fatalf("idents of resource path %s expect %v got error '%v'", test.s, test.t, err)
		}
		if act, exp := act, test.t; !compare(act, exp) {
			t.Fatalf("idents of resource path %s got %v expect %v", test.s, act, exp)
		}
	}
}

func TestIdentsErr(t *testing.T) {
	var tests = []struct {
		s string
		e string
	}{
		{
			s: ``,
			e: `found EOF, expected identifier at line 1, char 1`,
		},
		{
			s: `on`,
			e: `found ON, expected identifier at line 1, char 1`,
		},
		{
			s: `.`,
			e: `found ., expected identifier at line 1, char 1`,
		},
		{
			s: `a.`,
			e: `found EOF, expected identifier at line 1, char 4`,
		},
		{
			s: `a..c`,
			e: `invalid resource path a..c, expect a full path`,
		},
		{
			s: `a. .c`,
			e: `found ., expected identifier at line 1, char 4`,
		},
	}
	for _, test := range tests {
		if _, err := priv.ResourcePath(test.s).Idents(); err == nil || err.Error() != test.e {
			t.Fatalf("idents of resource path %s got error '%v' expect error '%v'", test.s, err, test.e)
		}
	}
}

func TestPrivilegeSetAddDelete(t *testing.T) {
	set := priv.NewPrivilegeTree()
	set.AddGlobal(priv.GrantPrivilege | priv.InsertPrivilege)
	set.Add("mydb.autogen", priv.SelectPrivilege)
	set.Add("mydb.autogen.cpu", priv.DeletePrivilege|priv.DropPrivilege)
	set.Add("mydb.autogen.mem", priv.DeletePrivilege|priv.DropPrivilege)
	set.Add("yourdb.daily", priv.SelectPrivilege)

	privilege := priv.GrantPrivilege | priv.InsertPrivilege | priv.SelectPrivilege | priv.DeletePrivilege | priv.DropPrivilege
	act, _ := set.Contain("mydb.autogen.cpu", privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect mydb.autogen.cpu resource contain privileges [%s]", privilege)
	}

	privilege = priv.SelectPrivilege
	act, _ = set.Contain("mydb", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb resource not contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("notexist", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect notexist resource not contain privileges [%s]", privilege)
	}

	privilege = priv.DeletePrivilege
	act, _ = set.Contain("mydb", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb resource not contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("notexist", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect resource notexist not contain privileges [%s]", privilege)
	}
	set.DeleteGlobal(privilege)
	act, _ = set.Contain("notexist", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect resource notexist not contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("mydb.autogen", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb.autogen resource not contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("mydb.autogen.cpu", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb.autogen.cpu resource not contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("mydb.autogen.mem", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb.autogen.mem resource not contain privileges [%s]", privilege)
	}

	privilege = priv.GrantPrivilege
	set.Delete("mydb.autogen.cpu", privilege)
	act, _ = set.Contain("mydb.autogen.cpu", privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect mydb.autogen.cpu resource contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("notexist", privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect notexist resource contain privileges [%s]", privilege)
	}
	set.DeleteGlobal(privilege)
	act = set.GlobalContain(privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect global resource not contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("mydb.autogen.cpu", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb.autogen.cpu resource not contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("notexist", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect notexist resource not contain privileges [%s]", privilege)
	}

	privilege = priv.DropPrivilege
	set.Delete("mydb.autogen.mem", privilege)
	act, _ = set.Contain("notexist", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect notexist resource not contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("mydb", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb resource not contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("mydb.autogen.mem", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb.autogen.mem resource contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("mydb.autogen.cpu", privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect mydb.autogen.cpu resource contain privileges [%s]", privilege)
	}

	privilege = priv.InsertPrivilege
	set.Delete("mydb.autogen.cpu", privilege)
	act, _ = set.Contain("notexist", privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect notexist resource contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("mydb", privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect mydb resource contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("mydb.autogen.mem", privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect mydb.autogen.mem resource contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("yourdb.daily", privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect yourdb resource contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("mydb.autogen.cpu", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb.autogen.cpu resource not contain privileges [%s]", privilege)
	}
	set.Delete("yourdb.autogen.speed", privilege)
	act, _ = set.Contain("yourdb.autogen.speed", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect yourdb.autogen.speed resource no contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("yourdb.autogen", privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect yourdb.autogen resource contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("yourdb", privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect yourdb resource contain privileges [%s]", privilege)
	}
	set.Delete("yourdb.daily", privilege)
	act, _ = set.Contain("yourdb.daily", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect yourdb.daily resource not contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("yourdb.daily.cpu", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect yourdb.daily.cpu resource not contain privileges [%s]", privilege)
	}
	set.Add("yourdb.daily.mem", privilege)
	act, _ = set.Contain("yourdb.daily.cpu", privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect yourdb.daily.cpu resource not contain privileges [%s]", privilege)
	}
	act, _ = set.Contain("yourdb.daily.mem", privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect yourdb.daily.mem resource contain privileges [%s]", privilege)
	}
}

func TestPrivilegeUnion(t *testing.T) {
	var testExists = []struct {
		r string
		p priv.Privilege
		t bool
	}{
		{
			r: "global",
			p: priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege,
			t: true,
		},
		{
			r: "noexists",
			p: priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege,
			t: true,
		},
		{
			r: "mydb",
			p: priv.SelectPrivilege | priv.InsertPrivilege,
			t: true,
		},
		{
			r: "mydb.autogen.cpu",
			p: priv.DeletePrivilege | priv.DropPrivilege,
			t: true,
		},
		{
			r: "yourdb.autogen.cpu",
			p: priv.InsertPrivilege,
			t: true,
		},
	}
	var testNotExists = []struct {
		r string
		p priv.Privilege
		t bool
	}{
		{
			r: "mydb.daily",
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: "yourdb.autogen",
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: "yourdb.autogen.mem",
			p: priv.InsertPrivilege,
			t: false,
		},
	}

	setA := priv.NewPrivilegeTree()
	setA.AddGlobal(priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege)
	setA.Add("mydb", priv.SelectPrivilege)
	setA.Add("mydb.autogen.cpu", priv.DeletePrivilege|priv.DropPrivilege)
	setA.Add("mydb.autogen.mem", priv.DeletePrivilege|priv.DropPrivilege)
	setA.Delete("mydb.daily", priv.SelectPrivilege)
	setA.Delete("yourdb.autogen", priv.InsertPrivilege)
	setA.Add("yourdb.autogen.cpu", priv.InsertPrivilege)

	setB := priv.NewPrivilegeTree()
	setB.Add("mydb.daily.cpu", priv.SelectPrivilege)
	setB.Delete("yourdb.autogen.cpu", priv.SelectPrivilege)
	setB.AddGlobal(priv.InsertPrivilege)
	setB.Delete("yourdb.daily", priv.InsertPrivilege)
	setB.Add("yourdb.daily.cpu", priv.InsertPrivilege)
	setB.Delete("yourdb.autogen.mem", priv.InsertPrivilege)
	setB.DeleteGlobal(priv.AuditPrivilege)

	runCases(t, setA, testExists)
	runCases(t, setA, testNotExists)
	setA.UnionWith(setB)
	runCases(t, setA, testExists)

	tests := []struct {
		r string
		p priv.Privilege
		t bool
	}{
		{
			r: "mydb.daily.cpu",
			p: priv.SelectPrivilege,
			t: true,
		},
		{
			r: "mydb.daily.mem",
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: "mydb.daily",
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: "yourdb.autogen.cpu",
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: "yourdb",
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: "global",
			p: priv.InsertPrivilege,
			t: true,
		},
		{
			r: "yourdb.autogen",
			p: priv.InsertPrivilege,
			t: true,
		},
		{
			r: "yourdb.daily",
			p: priv.InsertPrivilege,
			t: true,
		},
		{
			r: "yourdb.daily.cpu",
			p: priv.InsertPrivilege,
			t: true,
		},
		{
			r: "yourdb.autogen.cpu",
			p: priv.InsertPrivilege,
			t: true,
		},
		{
			r: "yourdb.autogen.mem",
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: "global",
			p: priv.AuditPrivilege,
			t: true,
		},
	}

	runCases(t, setA, tests)
}

func TestPrivilegeDiff(t *testing.T) {
	var testExists = []struct {
		r string
		p priv.Privilege
		t bool
	}{
		{
			r: "global",
			p: priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege,
			t: true,
		},
		{
			r: "noexists",
			p: priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege,
			t: true,
		},
		{
			r: "mydb",
			p: priv.SelectPrivilege | priv.InsertPrivilege,
			t: true,
		},
		{
			r: "mydb.autogen.cpu",
			p: priv.DeletePrivilege | priv.DropPrivilege,
			t: true,
		},
		{
			r: "yourdb.autogen.cpu",
			p: priv.InsertPrivilege,
			t: true,
		},
	}
	var testNotExists = []struct {
		r string
		p priv.Privilege
		t bool
	}{
		{
			r: "global",
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: "global",
			p: priv.ShowUsersPrivilege,
			t: false,
		},
		{
			r: "mydb.daily",
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: "yourdb.autogen",
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: "yourdb.autogen.mem",
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: "mydb.daily.cpu",
			p: priv.SelectPrivilege,
			t: false,
		},
	}

	setA := priv.NewPrivilegeTree()
	setA.AddGlobal(priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege)
	setA.Add("mydb", priv.SelectPrivilege)
	setA.Add("mydb.autogen.cpu", priv.DeletePrivilege|priv.DropPrivilege)
	setA.Add("mydb.autogen.mem", priv.DeletePrivilege|priv.DropPrivilege)
	setA.Delete("mydb.daily", priv.SelectPrivilege)
	setA.Delete("yourdb.autogen", priv.InsertPrivilege)
	setA.Add("yourdb.autogen.cpu", priv.InsertPrivilege)
	setA.Add("yourdb.autogen.disk", priv.InsertPrivilege)

	setB := priv.NewPrivilegeTree()
	setB.Add("mydb.daily.cpu", priv.SelectPrivilege)
	setB.Delete("yourdb.autogen.cpu", priv.SelectPrivilege)
	setB.AddGlobal(priv.InsertPrivilege)
	setB.Delete("yourdb.daily", priv.InsertPrivilege)
	setB.Add("yourdb.daily.cpu", priv.InsertPrivilege)
	setB.Delete("yourdb.autogen.mem", priv.InsertPrivilege)
	setB.Delete("yourdb.autogen.disk", priv.InsertPrivilege)
	setB.DeleteGlobal(priv.AuditPrivilege)

	runCases(t, setA, testExists)
	runCases(t, setA, testNotExists)
	setA.DifferentWith(setB)

	var tests = []struct {
		r string
		p priv.Privilege
		t bool
	}{
		{
			r: "global",
			p: priv.AuditPrivilege,
			t: true,
		},
		{
			r: "noexists",
			p: priv.AuditPrivilege,
			t: true,
		},
		{
			r: "mydb.daily",
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: "mydb.daily.cpu",
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: "global",
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: "noexists",
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: "yourdb.autogen",
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: "yourdb.autogen.cpu",
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: "yourdb.daily",
			p: priv.InsertPrivilege,
			t: true,
		},
		{
			r: "yourdb.daily.cpu",
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: "yourdb.autogen.mem",
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: "yourdb.autogen.mem",
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: "yourdb.autogen.disk",
			p: priv.InsertPrivilege,
			t: true,
		},
	}

	runCases(t, setA, tests)
}

func runCases(t *testing.T, set *priv.PrivilegeTree, tests []struct {
	r string
	p priv.Privilege
	t bool
}) {
	for _, test := range tests {
		var act bool
		if test.r == "global" {
			act = set.GlobalContain(test.p)
		} else {
			act, _ = set.Contain(test.r, test.p)
		}
		if act, exp := act, test.t; act != exp {
			if exp {
				t.Fatalf("expect %s resource contain privileges [%s]", test.r, test.p)
			} else {
				t.Fatalf("expect %s resource not contain privileges [%s]", test.r, test.p)
			}
		}
	}
}

func compare(a, b []string) bool {
	lenA, lenB := len(a), len(b)
	if lenA != lenB {
		return false
	}

	for i := 0; i < lenA; i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
