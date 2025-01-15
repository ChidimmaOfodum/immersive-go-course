package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
)


func main() {

	var HOST = "localhost"
	var mrPort = flag.String("mcrouter", "11211", "port for mc router")
	var mcPort = flag.String("memcacheds", "11212", "server ports for memcacheds \n more than one port is required")

	flag.Parse()

	mcPorts := strings.Split(*mcPort, ",")
	numberOfPorts := len(mcPorts)

	if numberOfPorts <= 1 {
		fmt.Fprint(os.Stderr, "more than one cache is required\n")
		flag.Usage()
		os.Exit(1)
	}
	routerServer := fmt.Sprintf("%s:%s", HOST, *mrPort)

	routerClient := memcache.New(routerServer)
	key := uuid.NewString()
	err := routerClient.Set(&memcache.Item{Key: key, Value: []byte("my value")})

	if err != nil {
		panic(err)
	}
	var absentCount int

	for _, cache := range mcPorts {
		cache = fmt.Sprintf("%s:%s", HOST, cache)
		mc := memcache.New(cache)
		_, err := mc.Get(key)

		if err != nil {
			if errors.Is(err, memcache.ErrCacheMiss) {
				absentCount++
			} else {
				panic(err)
			}
		}

	}

	switch absentCount {
	case numberOfPorts - 1:
		fmt.Println("Cache is sharded")
	case 0:
		fmt.Println("Cache is replicated")
	default:
		fmt.Print("Cannot determine if cache is sharded or replicated")
	}
}
