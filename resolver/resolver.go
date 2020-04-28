package resolver

import (
	"github.com/BASChain/go-bmail-account"
	"github.com/ethereum/go-ethereum/crypto"
	"net"
)

type NameResolver interface {
	DomainA(domain string) net.IP
	DomainMX(domain string) net.IP
	BMailBCA(mailHash string) (address bmail.Address, cname string)
}

func BMailNameHash(mailName string) string {
	hash := crypto.Keccak256Hash([]byte(mailName))
	return hash.String()
}
