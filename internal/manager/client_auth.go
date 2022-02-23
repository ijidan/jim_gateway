package manager

import "github.com/spf13/cast"

func ParseToken(token string) uint64  {
	return cast.ToUint64(token)
}