package httpext

import (
	"fmt"
	"io"
	"log"
)

func ToAddr(p int) string {
	return fmt.Sprintf(":%d", p)
}

func HandleClose(closer io.Closer) {
	if closer != nil {
		err := closer.Close()
		if err != nil {
			log.Printf("error on close %T: %s\n", closer, err.Error())
		}
	}
}
