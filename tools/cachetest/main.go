package main

import (
	"fmt"
	"kola/cache"
	"strconv"
)

func main() {
	c := cache.NewCache("cachetest", "test")
	if err := c.Start(); err != nil {
		panic(err)
	}

	val, err := c.Get("testval")
	if err != nil {
		panic(err)
	}
	fmt.Printf("old value: %s\n", val)

	var ival int
	if val != nil {
		ival, err = strconv.Atoi(string(val))
		ival++
	} else {
		ival = 0
	}

	val = []byte(fmt.Sprintf("%d", ival))
	if err := c.Put("testval", val); err != nil {
		panic(err)
	}

	fmt.Printf("put new value: %s\n", val)
}
