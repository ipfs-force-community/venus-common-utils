package apiinfo

import (
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/xerrors"
	"strconv"
	"strings"
)

const P_VERSION = multiaddr.P_WSS + 1

func init() {
	multiaddr.AddProtocol(multiaddr.Protocol{
		Name:  "version",
		Code:  P_VERSION,
		VCode: multiaddr.CodeToVarint(P_VERSION),
		Size:  multiaddr.LengthPrefixedVarSize,
		Transcoder: multiaddr.NewTranscoderFromFunctions(func(s string) ([]byte, error) {
			if !strings.HasPrefix(s, "v") {
				return nil, xerrors.New("version must start with version prefix v")
			}
			if len(s) < 2 {
				return nil, xerrors.New("must give a specify version such as v0")
			}
			_, err := strconv.Atoi(s[1:])
			if err != nil {
				return nil, xerrors.New("version part must be number")
			}
			return []byte(s), nil
		}, func(bytes []byte) (string, error) {
			vStr := string(bytes)
			if !strings.HasPrefix(vStr, "v") {
				return "", xerrors.New("version must start with version prefix v")
			}
			if len(vStr) < 2 {
				return "", xerrors.New("must give a specify version such as v0")
			}
			_, err := strconv.Atoi(vStr[1:])
			if err != nil {
				return "", xerrors.New("version part must be number")
			}
			return vStr, nil
		}, nil),
	})
}
