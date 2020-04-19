package main

import (
	"fmt"
	"webservice/consistenthash"
)

func main() {
	c := consistenthash.New(10)
	c.Add("host1")
	c.Add("host2")
	c.Add("host3")

	datas := []string{"data1", "data2", "data3", "data4", "data5"}
	for _, u := range datas {
		server := c.Get(u)
		fmt.Printf("%s => %s\n", u, server)
	}

	c.Add("host4")
	c.Add("host5")
	for _, u := range datas {
		server := c.Get(u)
		fmt.Printf("%s => %s\n", u, server)
	}

}
