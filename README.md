# Dig

Dig is embed client for service discovery, helping nodes discover each other. 



**Backends:**

- [x] redis
- [ ] etcd
- [ ] zookeeper
- [ ] consul



## Getting started

```golang
package main

import (
	"github.com/Sunmxt/dig"
	"time"
)

func main() {
    reg, _ := dig.Connect("redis", "127.0.0.1:6379", "my-services", 10, 10) // Connect to center.
	svc, _ := reg.Service("svc1") // Open service "svc1"
	// Publish node information.
	svc.Publish(&dig.Node{
		Name: "svc-node1",
		Metadata: map[string]string{
			"rpc-endpoint": "172.17.0.9:7890",
			"remark": "node of svc1",
			"rpc-proto": "gob",
			"rpc-weight": "8",
		},
		Timeout: 300,
	})
	// Start discovering other services and nodes.
	for {
        reg.Poll()
        time.Sleep(time.Second)
	}
}
```



