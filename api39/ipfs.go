package api39

import ipfs "github.com/ipfs/go-ipfs-api"

func IpfsAddDir(url, path string) (string, error) {
	return ipfs.NewShell(url).AddDir(path)
}
