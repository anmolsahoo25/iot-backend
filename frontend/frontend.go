package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"net/rpc"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

var BACKEND_SERVICE_SERVICE_HOST string = os.Getenv("BACKEND_SERVICE_SERVICE_HOST") + ":9000"

func main() {

	// Create the router and the routes
	router := mux.NewRouter()
	router.HandleFunc("/status", ServerStatus).Methods("GET")
	router.HandleFunc("/register", RegisterDevice).Methods("POST")
	router.HandleFunc("/recv/{device_id}", GetDeviceData).Methods("GET")
	router.HandleFunc("/send/{device_id}", SendDeviceData).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func ServerStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("Server status checked")
	log.Println(BACKEND_SERVICE_SERVICE_HOST)
	config, err := rest.InClusterConfig()
	clientset, err := kubernetes.NewForConfig(config)
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	var PodNames []string
	var CreationTimes []metav1.Time

	for _, pod := range pods.Items {
		value, ok := pod.Labels["app"]
		if ok == true && value == "iot-frontend" {
			log.Printf("%v %v\n", value, ok)
			PodNames = append(PodNames, pod.Name)
			CreationTimes = append(CreationTimes, pod.CreationTimestamp)
		}
	}
	var senddata = struct {
		Message       string
		PodNames      []string
		CreationTimes []metav1.Time
	}{
		"Success",
		PodNames,
		CreationTimes,
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(senddata)

}

func RegisterDevice(w http.ResponseWriter, r *http.Request) {
	device_id := GenerateRandomId()
	client, _ := rpc.DialHTTP("tcp", BACKEND_SERVICE_SERVICE_HOST)
	args := Args{device_id, []float64{0, 0, 0, 0, 0, 0, 0, 0}, []bool{false, false, false, false, false, false, false, false}}
	var reply string
	client.Call("Device.RegisterDevice", args, &reply)
	//log.Printf("%v", struct{ Device_id string }{device_id})
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct{ Device_id string }{device_id})
}

func GetDeviceData(w http.ResponseWriter, r *http.Request) {
	//log.Println("testing")
	device_id := mux.Vars(r)["device_id"]
	var t Device
	args := Args{device_id, []float64{0, 0, 0, 0, 0, 0, 0, 0}, []bool{false, false, false, false, false, false, false, false}}
	client, _ := rpc.DialHTTP("tcp", BACKEND_SERVICE_SERVICE_HOST)
	client.Call("Device.RecvData", args, &t)
	//log.Printf("Recv Call %v %v", t.Data, t.State)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)

}
func SendDeviceData(w http.ResponseWriter, r *http.Request) {
	device_id := mux.Vars(r)["device_id"]
	decoder := json.NewDecoder(r.Body)
	var t Device
	decoder.Decode(&t)
	data := t.Data
	state := t.State
	client, _ := rpc.DialHTTP("tcp", BACKEND_SERVICE_SERVICE_HOST)
	args := Args{device_id, data, state}
	var reply string
	client.Call("Device.SendData", args, &reply)
	//log.Printf("%v %v %v", device_id, data, state)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(struct{ Device_id string }{device_id})
}

func GenerateRandomId() string {
	const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "1234567890"

	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}
