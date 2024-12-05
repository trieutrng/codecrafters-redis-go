package main

import (
	"net"
	"strings"
	"time"
)

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

func InitReplication(procesor *Processor, opts serverOption) error {
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
		ReplicationServerInfo.Role = "slave"

		masterAddr := strings.Split(opts.replicaOf, " ")
		conn, err := connectMaster(masterAddr[0], masterAddr[1])
		if err != nil {
			return err
		}

		pingResp := &RESP{
			Type: Arrays,
			Nested: []*RESP{
				{
					Type: BulkString,
					Data: []byte("PING"),
				},
			},
		}
		pingMsg := procesor.parser.Serialize(pingResp)
		conn.Write(pingMsg)
	}
	return nil
}

func connectMaster(host, port string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 5*time.Second)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
