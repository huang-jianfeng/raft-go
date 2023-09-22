package mr

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"

	"6.5840/mr"
)

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.

	// uncomment to send the Example RPC to the coordinator.
	// CallExample()
	running := true

	for running {
		request := &TaskRequestMsg{}
		result := &TaskRequestResultMsg{}
		ok := call("Coordinator.AllocTask", request, result)
		if ok {
			switch result.taskType {
			case Finish:
				running = false
			case MapTask:

			case ReduceTask:

			}
		} else {

		}
	}

}

func ExecMapTask(filenames []string, mapf func(string, string) []KeyValue) {

	cache := make(map[int][]mr.KeyValue)
	for _, filename := range filenames {
		file, err := os.Open(filename)

		if err != nil {
			log.Fatal("cannot open %v", filename)
		}

		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("cannot read %v", filename)
		}
		file.Close()

		kva := mapf(filename, string(content))
		for _, kv := range kva {
			cache[ihash(kv.Key)] = append(cache[ihash(kv.Key)], kv)
		}
	}
}

func CallAllocTask(request *TaskRequestMsg, result *TaskRequestResultMsg) {
	ok := call("Coordinator.AllocTask", request, result)
	if ok {

	} else {

	}
}

func CallSubmitTaskResult() {

}

// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
