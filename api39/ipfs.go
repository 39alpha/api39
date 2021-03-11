package api39

import (
	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/kataras/iris/v12"
	"log"
	"regexp"
)

func IpfsAddDir(url, path string) (string, error) {
	return ipfs.NewShell(url).AddDir(path)
}

func IpfsGetAddr(url string) ([]string, error) {
	id, err := ipfs.NewShell(url).ID()
	if err != nil {
		return nil, err
	}
	return id.Addresses, err
}

func Addr(ctx iris.Context) {
	cfg, ok := ctx.Values().Get("config").(*Config)
	if !ok {
		log.Println("failed to retrieve configuration from context")
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"message": "failed to retrieve configuration from context",
		})
		return
	}

	id, err := ipfs.NewShell(cfg.Ipfs.Url).ID()
	if err != nil {
		log.Printf("failed to get IPFS address: %v\n", err)
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"message": "the IPFS gateway does not appear to be operational",
			"error":   err,
		})
		return
	}

	re := regexp.MustCompile(`^(/ip6|/ip4/(127\.0\.0\.1|192\.\d+\.\d+\.\d+))`)

	addresses := make([]string, 0, len(id.Addresses))
	for _, addr := range id.Addresses {
		if !re.MatchString(addr) {
			addresses = append(addresses, addr)
		}
	}

	if len(id.Addresses) == 0 {
		log.Println("failed to get IPFS address: no addresses found")
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"message": "the IPFS gateway does not appear to be operational",
			"error":   "no addresses found for the IPFS gateway",
		})
		return
	}

	_, _ = ctx.JSON(iris.Map{
		"message":   "success",
		"addresses": addresses,
	})
}
