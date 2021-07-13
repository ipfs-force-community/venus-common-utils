package apifinfo

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	logging "github.com/ipfs/go-log/v2"
	multiaddr "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

var log = logging.Logger("cliutil")

var (
	infoWithToken = regexp.MustCompile("^[a-zA-Z0-9\\-_]+?\\.[a-zA-Z0-9\\-_]+?\\.([a-zA-Z0-9\\-_]+)?:.+$")
)

type APIInfo struct {
	Addr  string
	Token []byte
}

func ParseApiInfo(s string) APIInfo {
	var tok []byte
	if infoWithToken.Match([]byte(s)) {
		sp := strings.SplitN(s, ":", 2)
		tok = []byte(sp[0])
		s = sp[1]
	}

	return APIInfo{
		Addr:  s,
		Token: tok,
	}
}

//DialArgs parser libp2p address to http/ws protocol, the version argument can be override by address in version
func (a APIInfo) DialArgs(version string) (string, error) {
	ma, err := multiaddr.NewMultiaddr(a.Addr)
	if err == nil {
		_, addr, err := manet.DialArgs(ma)
		if err != nil {
			return "", err
		}

		//override version
		val, err := ma.ValueForProtocol(P_VERSION)
		if err == nil {
			version = val
		}else if err != multiaddr.ErrProtocolNotFound{
			return "", err
		}

		_, err = ma.ValueForProtocol(multiaddr.P_WSS)
		if err == nil {
			return "wss://" + addr + "/rpc/" + version, nil
		}else if err != multiaddr.ErrProtocolNotFound{
			return "", err
		}

		_, err = ma.ValueForProtocol(multiaddr.P_HTTPS)
		if err == nil {
			return "https://" + addr + "/rpc/" + version, nil
		}else if err != multiaddr.ErrProtocolNotFound{
			return "", err
		}

		_, err = ma.ValueForProtocol(multiaddr.P_WS)
		if err == nil {
			return "ws://" + addr + "/rpc/" + version, nil
		}else if err != multiaddr.ErrProtocolNotFound{
			return "", err
		}

		_, err = ma.ValueForProtocol(multiaddr.P_HTTP)
		if err == nil {
			return "http://" + addr + "/rpc/" + version, nil
		}else if err != multiaddr.ErrProtocolNotFound{
			return "", err
		}

		return "ws://" + addr + "/rpc/" + version, nil
	}else{
		log.Warningf("parse libp2p address %s error , plz confirm this error %v", a.Addr, err)
	}

	_, err = url.Parse(a.Addr)
	if err != nil {
		return "", err
	}
	return a.Addr + "/rpc/" + version, nil
}

func (a APIInfo) Host() (string, error) {
	ma, err := multiaddr.NewMultiaddr(a.Addr)
	if err == nil {
		_, addr, err := manet.DialArgs(ma)
		if err != nil {
			return "", err
		}

		return addr, nil
	}

	spec, err := url.Parse(a.Addr)
	if err != nil {
		return "", err
	}
	return spec.Host, nil
}

func (a APIInfo) AuthHeader() http.Header {
	if len(a.Token) != 0 {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer "+string(a.Token))
		return headers
	}
	log.Warn("API Token not set and requested, capabilities might be limited.")
	return nil
}
