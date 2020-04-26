package resolver

import (
	"github.com/BASChain/go-bmail-account"
	"net"
)

type NameResolver interface {
	DomainA(string) net.IP
	DomainMX(string) net.IP
	EmailBCA(string) bmail.Address
}
