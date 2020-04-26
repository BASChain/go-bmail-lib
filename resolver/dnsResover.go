package resolver

import (
	"github.com/BASChain/go-bmail-account"
	"net"
)

type DNSResolver struct {
}

func (ds *DNSResolver) DomainA(domain string) net.IP {
	panic("implement me")
}

func (ds *DNSResolver) DomainMX(domainMX string) net.IP {
	panic("implement me")
}

func (ds *DNSResolver) EmailBCA(mailAddress string) bmail.Address {
	panic("implement me")
}

func NewDnsResolver() NameResolver {
	obj := &DNSResolver{}

	return obj
}
