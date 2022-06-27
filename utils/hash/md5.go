package hash

import (
	"crypto/md5"
	"fmt"
	"io"
)

func MD5(s string) (m string) {
	h := md5.New()
	_, _ = io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
