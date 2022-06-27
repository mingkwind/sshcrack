package hash

import "sshpwd_crack/vars"

func MakeTaskHash(k string) string {
	return MD5(k)
}

func CheckTaskHash(hash string) bool {
	_, ok := vars.SuccessHash.Load(hash)
	return ok
}

func SetTaskHash(hash string) {
	vars.SuccessHash.Store(hash, true)
}
