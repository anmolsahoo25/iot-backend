package main

import (
	"fmt"
	"net/rpc"
)

type Device struct {
	Device_id string
	Data      []float64
	State     []bool
}

type Args struct {
	Device_id string
	Data      []float64
	State     []bool
}

func main() {
	client, _ := rpc.DialHTTP("tcp", "localhost:9000")

	args := Args{"bNU8gjRvYZiDlFRk", []float64{1, 2, 3, 4, 5, 6, 7, 8}, []bool{false, true, false, true, false, true, false, true}}

	var string_reply string
	var reply Device

	client.Call("Device.SendData", args, &string_reply)
	client.Call("Device.RecvData", args, &reply)
	fmt.Printf("%v %v\n", reply.Data, reply.State)
}
