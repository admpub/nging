package sockaddr

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func sockaddrToAny(sa windows.Sockaddr) (*windows.RawSockaddrAny, Socklen, error) {
	if sa == nil {
		return nil, 0, syscall.EINVAL
	}

	switch sa := sa.(type) {
	case *windows.SockaddrInet4:
		if sa.Port < 0 || sa.Port > 0xFFFF {
			return nil, 0, syscall.EINVAL
		}

		raw := new(windows.RawSockaddrAny)
		raw.Addr.Family = windows.AF_INET
		raw4 := (*windows.RawSockaddrInet4)(unsafe.Pointer(raw))
		p := (*[2]byte)(unsafe.Pointer(&raw4.Port))
		p[0] = byte(sa.Port >> 8)
		p[1] = byte(sa.Port)
		for i := 0; i < len(sa.Addr); i++ {
			raw4.Addr[i] = sa.Addr[i]
		}
		return raw, Socklen(unsafe.Sizeof(*raw4)), nil

	case *windows.SockaddrInet6:
		if sa.Port < 0 || sa.Port > 0xFFFF {
			return nil, 0, syscall.EINVAL
		}

		raw := new(windows.RawSockaddrAny)
		raw.Addr.Family = windows.AF_INET6
		raw6 := (*windows.RawSockaddrInet6)(unsafe.Pointer(raw))
		p := (*[2]byte)(unsafe.Pointer(&raw6.Port))
		p[0] = byte(sa.Port >> 8)
		p[1] = byte(sa.Port)
		raw6.Scope_id = sa.ZoneId
		for i := 0; i < len(sa.Addr); i++ {
			raw6.Addr[i] = sa.Addr[i]
		}
		return raw, Socklen(unsafe.Sizeof(*raw6)), nil

	case *windows.SockaddrUnix:
		return nil, 0, syscall.EWINDOWS
	}
	return nil, 0, syscall.EAFNOSUPPORT
}

func anyToSockaddr(rsa *windows.RawSockaddrAny) (windows.Sockaddr, error) {
	if rsa == nil {
		return nil, syscall.EINVAL
	}

	switch rsa.Addr.Family {
	case windows.AF_UNIX:
		return nil, syscall.EWINDOWS

	case windows.AF_INET:
		pp := (*windows.RawSockaddrInet4)(unsafe.Pointer(rsa))
		sa := new(windows.SockaddrInet4)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		for i := 0; i < len(sa.Addr); i++ {
			sa.Addr[i] = pp.Addr[i]
		}
		return sa, nil

	case windows.AF_INET6:
		pp := (*windows.RawSockaddrInet6)(unsafe.Pointer(rsa))
		sa := new(windows.SockaddrInet6)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.ZoneId = pp.Scope_id
		for i := 0; i < len(sa.Addr); i++ {
			sa.Addr[i] = pp.Addr[i]
		}
		return sa, nil
	}
	return nil, syscall.EAFNOSUPPORT
}
