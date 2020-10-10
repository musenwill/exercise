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
			privilege: priv.ShowCQSPrivilege,
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
			privilege: priv.ReadPrivilege | priv.DropPrivilege | priv.ShowUsersPrivilege | priv.ShowCQSPrivilege,
			str:       "READ, DROP, SHOW USERS, SHOW CQS",
		},
	}
	for _, test := range tests {
		if test.str != test.privilege.String() {
			t.Fatalf("privilege to string got %s expect %s", test.privilege.String(), test.str)
		}
	}
}

// 8   	1000000000	         0.532 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPrivilegeGlobalContain(b *testing.B) {
	set := priv.NewPrivilegeTree()
	set.AddGlobal(priv.GrantPrivilege | priv.InsertPrivilege)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.GlobalContain(priv.GrantPrivilege | priv.InsertPrivilege)
	}
}

// 8   	20796274	        53.5 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPrivilegeContainDeepPath(b *testing.B) {
	set := priv.NewPrivilegeTree()
	set.AddGlobal(priv.GrantPrivilege | priv.InsertPrivilege)
	set.Add(priv.CreateResourcePathUnsafe("mydb"), priv.SelectPrivilege)
	set.Delete(priv.CreateResourcePathUnsafe("mydb.autogen"), priv.SelectPrivilege)
	set.Add(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), priv.DeletePrivilege|priv.DropPrivilege|priv.SelectPrivilege)
	set.Add(priv.CreateResourcePathUnsafe("mydb.autogen.mem"), priv.DeletePrivilege|priv.DropPrivilege|priv.SelectPrivilege)
	set.Add(priv.CreateResourcePathUnsafe("yourdb.daily"), priv.SelectPrivilege)

	resource := priv.CreateResourcePathUnsafe("mydb.autogen.cpu")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Contain(resource, priv.DeletePrivilege|priv.DropPrivilege|priv.SelectPrivilege)
	}
}

// 8   	371420978	         3.20 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPrivilegeContainShallowPath(b *testing.B) {
	set := priv.NewPrivilegeTree()
	set.AddGlobal(priv.GrantPrivilege | priv.InsertPrivilege)
	resource := priv.CreateResourcePathUnsafe("")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Contain(resource, priv.GrantPrivilege|priv.InsertPrivilege)
	}
}

// 8   	 1000000	      1053 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPrivilegeContains(b *testing.B) {
	setA := priv.NewPrivilegeTree()
	setA.AddGlobal(priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("mydb"), priv.SelectPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), priv.DeletePrivilege|priv.DropPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("mydb.autogen.mem"), priv.DeletePrivilege|priv.DropPrivilege)
	setA.Delete(priv.CreateResourcePathUnsafe("mydb.daily"), priv.SelectPrivilege)
	setA.Delete(priv.CreateResourcePathUnsafe("yourdb.autogen"), priv.InsertPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"), priv.InsertPrivilege)

	setB := priv.NewPrivilegeTree()
	setB.Add(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), priv.DeletePrivilege|priv.DropPrivilege)
	setB.Add(priv.CreateResourcePathUnsafe("mydb.autogen.mem"), priv.DeletePrivilege|priv.DropPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"), priv.InsertPrivilege)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setA.Contains(setB)
	}
}

func TestPrivilegeSetAll(t *testing.T) {

	set := priv.NewPrivilegeTree()
	set.SetAll()

	var tests = []struct {
		resource  *priv.ResourcePath
		privilege priv.Privilege
		contains  bool
	}{
		{
			resource:  priv.CreateResourcePathUnsafe("notexist"),
			privilege: priv.AllGlobalPrivileges,
			contains:  true,
		},
		{
			resource:  priv.CreateResourcePathUnsafe("notexist"),
			privilege: priv.AllResourcePrivileges,
			contains:  true,
		},
		{
			resource:  priv.CreateResourcePathUnsafe("notexist"),
			privilege: priv.CreateUserPrivilege | priv.GrantPrivilege | priv.AuditPrivilege | priv.SelectPrivilege,
			contains:  true,
		},
		{
			resource:  priv.CreateResourcePathUnsafe("mydb.autogen.cpu"),
			privilege: priv.AllResourcePrivileges,
			contains:  true,
		},
		{
			resource:  priv.CreateResourcePathUnsafe("mydb.autogen.cpu"),
			privilege: priv.SelectPrivilege | priv.InsertPrivilege | priv.DeletePrivilege | priv.DropPrivilege,
			contains:  true,
		},
	}

	for _, test := range tests {
		act := set.Contain(test.resource, test.privilege)
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
		resource  *priv.ResourcePath
		privilege priv.Privilege
		contains  bool
	}{
		{
			resource:  priv.CreateResourcePathUnsafe("notexist"),
			privilege: priv.AllGlobalPrivileges,
			contains:  false,
		},
		{
			resource:  priv.CreateResourcePathUnsafe("notexist"),
			privilege: priv.AllResourcePrivileges,
			contains:  false,
		},
		{
			resource:  priv.CreateResourcePathUnsafe("notexist"),
			privilege: priv.CreateUserPrivilege | priv.GrantPrivilege | priv.AuditPrivilege,
			contains:  false,
		},
		{
			resource:  priv.CreateResourcePathUnsafe("mydb.autogen.cpu"),
			privilege: priv.AllResourcePrivileges,
			contains:  false,
		},
		{
			resource:  priv.CreateResourcePathUnsafe("mydb.autogen.cpu"),
			privilege: priv.SelectPrivilege | priv.InsertPrivilege | priv.DeletePrivilege | priv.DropPrivilege,
			contains:  false,
		},
	}

	for _, test := range tests {
		act := set.Contain(test.resource, test.privilege)
		if act, exp := act, test.contains; act != exp {
			if exp {
				t.Fatalf("expect ClearAll() contain privileges [%s] on resource %s", test.privilege, test.resource)
			} else {
				t.Fatalf("expect ClearAll() not contain privileges [%s] on resource %s", test.privilege, test.resource)
			}
		}
	}
}

func TestResourcePathOk(t *testing.T) {
	var tests = []struct {
		s string
		t []string
	}{
		{
			s: ``,
			t: []string{},
		},
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
		{
			s: `a..c`,
			t: []string{"a", "autogen", "c"},
		},
	}
	for _, test := range tests {
		act, err := priv.CreateResourcePath(test.s)
		if err != nil {
			t.Fatalf("idents of resource path %s expect %v got error '%v'", test.s, test.t, err)
		}
		if act, exp := act.Segs, test.t; !compare(act, exp) {
			t.Fatalf("idents of resource path %s got %v expect %v", test.s, act, exp)
		}
	}
}

func TestResourcePathErr(t *testing.T) {
	var tests = []struct {
		s string
		e string
	}{
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
			s: `a. .c`,
			e: `found ., expected identifier at line 1, char 4`,
		},
	}
	for _, test := range tests {
		if _, err := priv.CreateResourcePath(test.s); err == nil || err.Error() != test.e {
			t.Fatalf("idents of resource path %s got error '%v' expect error '%v'", test.s, err, test.e)
		}
	}
}

func TestPrivilegeSetAddDelete(t *testing.T) {
	set := priv.NewPrivilegeTree()
	set.AddGlobal(priv.GrantPrivilege | priv.InsertPrivilege)
	set.Add(priv.CreateResourcePathUnsafe("mydb.autogen"), priv.SelectPrivilege)
	set.Add(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), priv.DeletePrivilege|priv.DropPrivilege)
	set.Add(priv.CreateResourcePathUnsafe("mydb.autogen.mem"), priv.DeletePrivilege|priv.DropPrivilege)
	set.Add(priv.CreateResourcePathUnsafe("yourdb.daily"), priv.SelectPrivilege)

	privilege := priv.GrantPrivilege | priv.InsertPrivilege | priv.SelectPrivilege | priv.DeletePrivilege | priv.DropPrivilege
	act := set.Contain(priv.CreateResourcePathUnsafe("mydb..cpu"), privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect mydb.autogen.cpu resource contain privileges [%s]", privilege)
	}

	privilege = priv.SelectPrivilege
	act = set.Contain(priv.CreateResourcePathUnsafe("mydb"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb resource not contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("notexist"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect notexist resource not contain privileges [%s]", privilege)
	}

	privilege = priv.DeletePrivilege
	act = set.Contain(priv.CreateResourcePathUnsafe("mydb"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb resource not contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("notexist"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect resource notexist not contain privileges [%s]", privilege)
	}
	set.DeleteGlobal(privilege)
	act = set.Contain(priv.CreateResourcePathUnsafe("notexist"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect resource notexist not contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("mydb.autogen"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb.autogen resource not contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb.autogen.cpu resource not contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("mydb.autogen.mem"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb.autogen.mem resource not contain privileges [%s]", privilege)
	}

	privilege = priv.GrantPrivilege
	set.Delete(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), privilege)
	act = set.Contain(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect mydb.autogen.cpu resource contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("notexist"), privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect notexist resource contain privileges [%s]", privilege)
	}
	set.DeleteGlobal(privilege)
	act = set.GlobalContain(privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect global resource not contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb.autogen.cpu resource not contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("notexist"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect notexist resource not contain privileges [%s]", privilege)
	}

	privilege = priv.DropPrivilege
	set.Delete(priv.CreateResourcePathUnsafe("mydb.autogen.mem"), privilege)
	act = set.Contain(priv.CreateResourcePathUnsafe("notexist"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect notexist resource not contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("mydb"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb resource not contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("mydb.autogen.mem"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb.autogen.mem resource contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect mydb.autogen.cpu resource contain privileges [%s]", privilege)
	}

	privilege = priv.InsertPrivilege
	set.Delete(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), privilege)
	act = set.Contain(priv.CreateResourcePathUnsafe("notexist"), privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect notexist resource contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("mydb"), privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect mydb resource contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("mydb.autogen.mem"), privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect mydb.autogen.mem resource contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("yourdb.daily"), privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect yourdb resource contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect mydb.autogen.cpu resource not contain privileges [%s]", privilege)
	}
	set.Delete(priv.CreateResourcePathUnsafe("yourdb.autogen.speed"), privilege)
	act = set.Contain(priv.CreateResourcePathUnsafe("yourdb.autogen.speed"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect yourdb.autogen.speed resource no contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("yourdb.autogen"), privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect yourdb.autogen resource contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("yourdb"), privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect yourdb resource contain privileges [%s]", privilege)
	}
	set.Delete(priv.CreateResourcePathUnsafe("yourdb.daily"), privilege)
	act = set.Contain(priv.CreateResourcePathUnsafe("yourdb.daily"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect yourdb.daily resource not contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("yourdb.daily.cpu"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect yourdb.daily.cpu resource not contain privileges [%s]", privilege)
	}
	set.Add(priv.CreateResourcePathUnsafe("yourdb.daily.mem"), privilege)
	act = set.Contain(priv.CreateResourcePathUnsafe("yourdb.daily.cpu"), privilege)
	if act, exp := act, false; act != exp {
		t.Fatalf("expect yourdb.daily.cpu resource not contain privileges [%s]", privilege)
	}
	act = set.Contain(priv.CreateResourcePathUnsafe("yourdb.daily.mem"), privilege)
	if act, exp := act, true; act != exp {
		t.Fatalf("expect yourdb.daily.mem resource contain privileges [%s]", privilege)
	}
}

func TestPrivilegeUnion(t *testing.T) {
	var testExists = []struct {
		r *priv.ResourcePath
		p priv.Privilege
		t bool
	}{
		{
			r: priv.CreateResourcePathUnsafe("global"),
			p: priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("noexists"),
			p: priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("mydb"),
			p: priv.SelectPrivilege | priv.InsertPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("mydb.autogen.cpu"),
			p: priv.DeletePrivilege | priv.DropPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"),
			p: priv.InsertPrivilege,
			t: true,
		},
	}
	var testNotExists = []struct {
		r *priv.ResourcePath
		p priv.Privilege
		t bool
	}{
		{
			r: priv.CreateResourcePathUnsafe("mydb.daily"),
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen"),
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen.mem"),
			p: priv.InsertPrivilege,
			t: false,
		},
	}

	setA := priv.NewPrivilegeTree()
	setA.AddGlobal(priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("mydb"), priv.SelectPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), priv.DeletePrivilege|priv.DropPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("mydb.autogen.mem"), priv.DeletePrivilege|priv.DropPrivilege)
	setA.Delete(priv.CreateResourcePathUnsafe("mydb.daily"), priv.SelectPrivilege)
	setA.Delete(priv.CreateResourcePathUnsafe("yourdb.autogen"), priv.InsertPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"), priv.InsertPrivilege)

	setB := priv.NewPrivilegeTree()
	setB.Add(priv.CreateResourcePathUnsafe("mydb.daily.cpu"), priv.SelectPrivilege)
	setB.Delete(priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"), priv.SelectPrivilege)
	setB.AddGlobal(priv.InsertPrivilege)
	setB.Delete(priv.CreateResourcePathUnsafe("yourdb.daily"), priv.InsertPrivilege)
	setB.Add(priv.CreateResourcePathUnsafe("yourdb.daily.cpu"), priv.InsertPrivilege)
	setB.Delete(priv.CreateResourcePathUnsafe("yourdb.autogen.mem"), priv.InsertPrivilege)
	setB.DeleteGlobal(priv.AuditPrivilege)

	runCases(t, setA, testExists)
	runCases(t, setA, testNotExists)
	setA.UnionWith(setB)
	runCases(t, setA, testExists)

	tests := []struct {
		r *priv.ResourcePath
		p priv.Privilege
		t bool
	}{
		{
			r: priv.CreateResourcePathUnsafe("mydb.daily.cpu"),
			p: priv.SelectPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("mydb.daily.mem"),
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("mydb.daily"),
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"),
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb"),
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("global"),
			p: priv.InsertPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen"),
			p: priv.InsertPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.daily"),
			p: priv.InsertPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.daily.cpu"),
			p: priv.InsertPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"),
			p: priv.InsertPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen.mem"),
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("global"),
			p: priv.AuditPrivilege,
			t: true,
		},
	}

	runCases(t, setA, tests)
}

func TestPrivilegeDiff(t *testing.T) {
	var testExists = []struct {
		r *priv.ResourcePath
		p priv.Privilege
		t bool
	}{
		{
			r: priv.CreateResourcePathUnsafe("global"),
			p: priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("noexists"),
			p: priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("mydb"),
			p: priv.SelectPrivilege | priv.InsertPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("mydb.autogen.cpu"),
			p: priv.DeletePrivilege | priv.DropPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"),
			p: priv.InsertPrivilege,
			t: true,
		},
	}
	var testNotExists = []struct {
		r *priv.ResourcePath
		p priv.Privilege
		t bool
	}{
		{
			r: priv.CreateResourcePathUnsafe("global"),
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("global"),
			p: priv.ShowUsersPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("mydb.daily"),
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen"),
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen.mem"),
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("mydb.daily.cpu"),
			p: priv.SelectPrivilege,
			t: false,
		},
	}

	setA := priv.NewPrivilegeTree()
	setA.AddGlobal(priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("mydb"), priv.SelectPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), priv.DeletePrivilege|priv.DropPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("mydb.autogen.mem"), priv.DeletePrivilege|priv.DropPrivilege)
	setA.Delete(priv.CreateResourcePathUnsafe("mydb.daily"), priv.SelectPrivilege)
	setA.Delete(priv.CreateResourcePathUnsafe("yourdb.autogen"), priv.InsertPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"), priv.InsertPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("yourdb.autogen.disk"), priv.InsertPrivilege)

	setB := priv.NewPrivilegeTree()
	setB.Add(priv.CreateResourcePathUnsafe("mydb.daily.cpu"), priv.SelectPrivilege)
	setB.Delete(priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"), priv.SelectPrivilege)
	setB.AddGlobal(priv.InsertPrivilege)
	setB.Delete(priv.CreateResourcePathUnsafe("yourdb.daily"), priv.InsertPrivilege)
	setB.Add(priv.CreateResourcePathUnsafe("yourdb.daily.cpu"), priv.InsertPrivilege)
	setB.Delete(priv.CreateResourcePathUnsafe("yourdb.autogen.mem"), priv.InsertPrivilege)
	setB.Delete(priv.CreateResourcePathUnsafe("yourdb.autogen.disk"), priv.InsertPrivilege)
	setB.DeleteGlobal(priv.AuditPrivilege)

	runCases(t, setA, testExists)
	runCases(t, setA, testNotExists)
	setA.DifferentWith(setB)

	var tests = []struct {
		r *priv.ResourcePath
		p priv.Privilege
		t bool
	}{
		{
			r: priv.CreateResourcePathUnsafe("global"),
			p: priv.AuditPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("noexists"),
			p: priv.AuditPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("mydb.daily"),
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("mydb.daily.cpu"),
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("global"),
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("noexists"),
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen"),
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"),
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.daily"),
			p: priv.InsertPrivilege,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.daily.cpu"),
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen.mem"),
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen.mem"),
			p: priv.InsertPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen.disk"),
			p: priv.InsertPrivilege,
			t: true,
		},
	}

	runCases(t, setA, tests)
}

func TestDiffSelf(t *testing.T) {
	setA := priv.NewPrivilegeTree()
	setA.AddGlobal(priv.GrantPrivilege | priv.AuditPrivilege | priv.InsertPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("mydb"), priv.SelectPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), priv.DeletePrivilege|priv.DropPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("mydb.autogen.mem"), priv.DeletePrivilege|priv.DropPrivilege)
	setA.Delete(priv.CreateResourcePathUnsafe("mydb.daily"), priv.SelectPrivilege)
	setA.Delete(priv.CreateResourcePathUnsafe("yourdb.autogen"), priv.InsertPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"), priv.InsertPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("yourdb.autogen.disk"), priv.InsertPrivilege)

	setA.DifferentWith(setA)
	if !setA.Powerless() {
		t.Fatalf("expect self diff to be powerless")
	}
}

func TestReadWritecompatibility(t *testing.T) {
	setA := priv.NewPrivilegeTree()
	setA.AddGlobal(priv.ReadPrivilege)
	setA.Delete(priv.CreateResourcePathUnsafe("mydb"), priv.ReadPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("mydb.autogen"), priv.ReadPrivilege)
	setA.Delete(priv.CreateResourcePathUnsafe("mydb.autogen.cpu"), priv.ReadPrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("yourdb"), priv.WritePrivilege)
	setA.Delete(priv.CreateResourcePathUnsafe("yourdb.autogen"), priv.WritePrivilege)
	setA.Add(priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"), priv.WritePrivilege)

	var tests = []struct {
		r *priv.ResourcePath
		p priv.Privilege
		t bool
	}{
		{
			r: priv.CreateResourcePathUnsafe("global"),
			p: priv.ReadGroupPrivileges,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("noexists"),
			p: priv.ReadGroupPrivileges,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("global"),
			p: priv.DeletePrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("mydb"),
			p: priv.SelectPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("mydb.autogen"),
			p: priv.ReadGroupPrivileges,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("mydb.autogen.cpu"),
			p: priv.ShowDatabasesPrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb"),
			p: priv.WriteGroupPrivileges,
			t: true,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen"),
			p: priv.CreateDatabasePrivilege,
			t: false,
		},
		{
			r: priv.CreateResourcePathUnsafe("yourdb.autogen.cpu"),
			p: priv.WriteGroupPrivileges,
			t: true,
		},
	}
	runCases(t, setA, tests)
}

func runCases(t *testing.T, set priv.PrivilegeSet, tests []struct {
	r *priv.ResourcePath
	p priv.Privilege
	t bool
}) {
	for _, test := range tests {
		var act bool
		if test.r.String() == "global" {
			act = set.GlobalContain(test.p)
		} else {
			act = set.Contain(test.r, test.p)
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
