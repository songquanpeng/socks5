package handler

import (
	"io"
	"net"
)

func forward(src, dst net.Conn) {
	defer src.Close()
	defer dst.Close()
	io.Copy(src, dst)
}
