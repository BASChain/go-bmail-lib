package resolver

import (
	"github.com/BASChain/go-bmail-account"
	"net"
)

type EthResolver struct {
}

func (er *EthResolver) DomainA(domain string) net.IP {
	panic("implement me")
}

func (er *EthResolver) DomainMX(domainMX string) net.IP {
	panic("implement me")
}

func (er *EthResolver) EmailBCA(emailAddress string) bmail.Address {
	panic("implement me")
}

func NewEthResolver() NameResolver {
	obj := &EthResolver{}

	return obj
}
