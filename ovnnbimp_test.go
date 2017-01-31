package libovndb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	SOCKET = "/var/run/openvswitch/ovnnb_db.sock"

	LSW   = "TEST_LSW"
	LSP   = "TEST_LSP"
	LSP_SECOND   = "TEST_LSP_SECOND "
	ADDR  = "36:46:56:76:86:96 127.0.0.1"
	MATCH = "outport == \"96d44061-1823-428b-a7ce-f473d10eb3d0\" && ip && ip.dst == 10.97.183.61"
	MATCH_SECOND = "outport == \"96d44061-1823-428b-a7ce-f473d10eb3d0\" && ip && ip.dst == 10.97.183.62"
)

var ovndbapi OVNDBApi

func init() {
	ovndbapi = GetInstance(SOCKET, UNIX, "", 0)
}

func TestACLs(t *testing.T) {
	var c []*OvnCommand = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.LSWAdd(LSW))
	c = append(c, ovndbapi.LSPAdd(LSW, LSP))

	c = append(c, ovndbapi.LSPSetAddress(LSP, ADDR))
	c = append(c, ovndbapi.ACLAdd(LSW, "to-lport", MATCH, "drop", 1001, nil, false))
	ovndbapi.Execute(c...)

	lsps := ovndbapi.GetLogicPortsBySwitch(LSW)
	assert.Equal(t, true, len(lsps) == 1 && lsps[0].Name == LSP, "test[%s]: %v", "added port", lsps)
	assert.Equal(t, true, len(lsps) == 1 && lsps[0].Addresses[0] == ADDR, "test[%s]", "setted port address")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.LSPAdd(LSW, LSP_SECOND))
	ovndbapi.Execute(c...)
	lsps = ovndbapi.GetLogicPortsBySwitch(LSW)
	assert.Equal(t, true, len(lsps) == 2, "test[%s]: %+v", "added 2 ports", lsps)

	acls := ovndbapi.GetACLsBySwitch(LSW)
	assert.Equal(t, true, len(acls) == 1 && acls[0].Match == MATCH &&
		acls[0].Action == "drop" && acls[0].Priority == 1001, "test[%s]", "add acl")

	assert.Equal(t, true, nil == ovndbapi.ACLAdd(LSW, "to-lport", MATCH, "drop", 1001, nil, false), "test[%s]", "add same acl twice, should only one added.")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ACLAdd(LSW, "to-lport", MATCH_SECOND, "drop", 1001, nil, false))
	ovndbapi.Execute(c...)
	acls = ovndbapi.GetACLsBySwitch(LSW)
	assert.Equal(t, true, len(acls) == 2, "test[%s]", "add second acl")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ACLDel(LSW, "to-lport", MATCH, 1001))
	ovndbapi.Execute(c...)
	acls = ovndbapi.GetACLsBySwitch(LSW)
	assert.Equal(t, true, len(acls) == 1, "test[%s]", "acl remove")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ACLDel(LSW, "to-lport", MATCH_SECOND, 1001))
	ovndbapi.Execute(c...)
	acls = ovndbapi.GetACLsBySwitch(LSW)
	assert.Equal(t, true, len(acls) == 0, "test[%s]", "acl remove")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.LSPDel(LSP))
	ovndbapi.Execute(c...)
	lsps = ovndbapi.GetLogicPortsBySwitch(LSW)
	assert.Equal(t, true, len(lsps) == 1, "test[%s]", "one port remove")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.LSPDel(LSP_SECOND))
	ovndbapi.Execute(c...)
	lsps = ovndbapi.GetLogicPortsBySwitch(LSW)
	assert.Equal(t, true, len(lsps) == 0, "test[%s]", "one port remove")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.LSWDel(LSW))
	ovndbapi.Execute(c...)

}

func findAS(name string) bool {
	as := ovndbapi.GetAddressSets()
	for _, a := range as {
		if a.Name == name {
			return true
		}
	}
	return false
}

func TestAddressSet(t *testing.T) {
	var c []*OvnCommand = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ASAdd("AS1", []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}))
	ovndbapi.Execute(c...)
	as := ovndbapi.GetAddressSets()
	assert.Equal(t, true, findAS("AS1"), "test[%s] and value[%v]", "address set added.", as[0].Addresses)

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ASAdd("AS2", []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}))
	ovndbapi.Execute(c...)
	as = ovndbapi.GetAddressSets()
	assert.Equal(t, true, findAS("AS2"), "test[%s] and value[%v]", "address set added.", as[0].Addresses)

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ASDel("AS1"))
	ovndbapi.Execute(c...)
	assert.Equal(t, false, findAS("AS1"), "test AS remove")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ASDel("AS2"))
	ovndbapi.Execute(c...)
	assert.Equal(t, false, findAS("AS2"), "test AS remove")
}