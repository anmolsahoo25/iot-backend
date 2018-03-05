package main

import (
	"log"
	"net/http"
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

var datastore map[string]*Device = make(map[string]*Device)

func (d *Device) RegisterDevice(args *Args, reply *string) error {
	log.Println("Register Request Received")
	datastore[args.Device_id] = &Device{args.Device_id, []float64{0, 0, 0, 0, 0, 0, 0, 0}, []bool{false, false, false, false, false, false, false, false}}
	return nil
}

func (d *Device) RecvData(args *Args, reply *Device) error {
	_, ok := datastore[args.Device_id]
	if ok {
		reply.Device_id = datastore[args.Device_id].Device_id
		reply.Data = datastore[args.Device_id].Data
		reply.State = datastore[args.Device_id].State
		return nil
	}
	return nil
}

func (d *Device) SendData(args *Args, reply *string) error {
	//log.Printf("%v %v", args.Data, args.State)
	_, ok := datastore[args.Device_id]
	if ok {
		datastore[args.Device_id].Data = args.Data
		datastore[args.Device_id].State = args.State
		return nil
	}
	return nil
}

func main() {
	device := new(Device)
	rpc.Register(device)
	rpc.HandleHTTP()

	log.Fatal(http.ListenAndServe(":9000", nil))
}
