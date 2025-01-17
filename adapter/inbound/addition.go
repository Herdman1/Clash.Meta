package inbound

import (
	"net"

	C "github.com/Dreamacro/clash/constant"
)

type Addition func(metadata *C.Metadata)

func (a Addition) Apply(metadata *C.Metadata) {
	a(metadata)
}

func WithInName(name string) Addition {
	return func(metadata *C.Metadata) {
		metadata.InName = name
	}
}

func WithInUser(user string) Addition {
	return func(metadata *C.Metadata) {
		metadata.InUser = user
	}
}

func WithSpecialRules(specialRules string) Addition {
	return func(metadata *C.Metadata) {
		metadata.SpecialRules = specialRules
	}
}

func WithSpecialProxy(specialProxy string) Addition {
	return func(metadata *C.Metadata) {
		metadata.SpecialProxy = specialProxy
	}
}

func WithSrcAddr(addr net.Addr) Addition {
	return func(metadata *C.Metadata) {
		if ip, port, err := parseAddr(addr); err == nil {
			metadata.SrcIP = ip
			metadata.SrcPort = port
		}
	}
}
