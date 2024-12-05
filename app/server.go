package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

type serverOption struct {
	port      string
	replicaOf string
}

type redisReplicationInfo struct {
	Role                       string `info:"role"`
	ConnectedSlaves            int    `info:"connected_slaves"`
	MasterReplid               string `info:"master_replid"`
	MasterReplOffset           int    `info:"master_repl_offset"`
	SecondReplOffset           int    `info:"second_repl_offset"`
	ReplBacklogActive          int    `info:"repl_backlog_active"`
	ReplBacklogSize            int    `info:"repl_backlog_size"`
	ReplBacklogFirstByteOffset int    `info:"repl_backlog_first_byte_offset"`
	ReplBacklogHistlen         int    `info:"repl_backlog_histlen"`
}

var ReplicationServerInfo redisReplicationInfo

func getServerOptions(args []string) serverOption {
	opts := serverOption{
		port: "6379",
	}
	for i, arg := range args {
		switch arg {
		case "--port":
			if i < len(args)-1 {
				opts.port = args[i+1]
			}
		case "--replicaof":
			opts.replicaOf = args[i+1]
		}
	}
	return opts
}

func initDefaultReplicationInfo(procesor *Processor, opts serverOption) {
	ReplicationServerInfo = redisReplicationInfo{
		Role:                       "master",
		ConnectedSlaves:            0,
		MasterReplid:               "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb", // hard coded
		MasterReplOffset:           0,
		SecondReplOffset:           -1,
		ReplBacklogActive:          0,
		ReplBacklogSize:            1048576,
		ReplBacklogFirstByteOffset: 0,
		ReplBacklogHistlen:         0,
	}

	if len(opts.replicaOf) > 0 {
		// masterAddr := strings.Split(opts.replicaOf, " ")
		ReplicationServerInfo.Role = "slave"
	}
}

func main() {
	opts := getServerOptions(os.Args)

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", opts.port))
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	// init dependencies
	respParser := NewRESP()
	memory := NewMemory()
	processor := NewProcessor(respParser, memory)

	// store replication info
	go initDefaultReplicationInfo(processor, opts)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handle(conn, processor)
	}
}

func handle(conn net.Conn, processor *Processor) {
	buf := make([]byte, 1024)
	for {
		read, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error when parsing command!", err.Error())
			break
		}
		if read == 0 {
			fmt.Println("No data read")
			break
		}
		output, err := processor.Accept(buf)
		if err != nil {
			fmt.Println("Invalid command: ", err)
			break
		}
		conn.Write(output)
	}
}
