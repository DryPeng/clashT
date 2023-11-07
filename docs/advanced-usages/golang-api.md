---
sidebarTitle: Integrating ClashT in Golang Programs
sidebarOrder: 3
---

# Integrating ClashT in Golang Programs

If clash does not fit your own usage, you can use ClashTTT in your own Golang code.

There is already basic support:

```go
package main

import (
	"context"
	"fmt"
	"io"
	"net"

"github.com/DryPeng/clashT/adapter/outbound"

"github.com/DryPeng/clashT/constant"
	"github.com/DryPeng/clashT/listener/socks"
)

func main() {
	in := make(chan constant.ConnContext, 100)
	defer close(in)

	l, err := socks.New("127.0.0.1:10000", in)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	println("listen at:", l.Address())

	direct := outbound.NewDirect()

	for c := range in {
		conn := c
		metadata := conn.Metadata()
		fmt.Printf("request incoming from %s to %s\n", metadata.SourceAddress(), metadata.RemoteAddress())
		go func () {
			remote, err := direct.DialContext(context.Background(), metadata)
			if err != nil {
				fmt.Printf("dial error: %s\n", err.Error())
				return
			}
			relay(remote, conn.Conn())
		}()
	}
}

func relay(l, r net.Conn) {
	go io.Copy(l, r)
	io.Copy(r, l)
}
```
