package main

import "fmt"

const (
	Ping = "PING"
	Echo = "ECHO"
)

type Executor func(resp *RESP) (*RESP, error)

type Processor struct {
	parser    RespParser
	memory    *Memory
	executors map[string]Executor
}

func NewProcessor(respParser RespParser, memory *Memory) *Processor {
	return &Processor{
		parser:    respParser,
		memory:    memory,
		executors: initExecutors(memory),
	}
}

func (p *Processor) Accept(cmd []byte) ([]byte, error) {
	resp, err := p.parser.Deserialize(cmd)
	if err != nil {
		return nil, err
	}
	if resp.Type != Arrays {
		return nil, fmt.Errorf("expected Arrays type for command, but received: %v", RespTypeString(resp.Type))
	}
	if len(resp.Nested) == 0 {
		return nil, fmt.Errorf("invalid command: command empty")
	}

	executor, ok := p.executors[string(resp.Nested[0].Data)]
	if !ok {
		return nil, fmt.Errorf("command not supported")
	}

	output, err := executor(resp)
	if err != nil {
		return nil, err
	}

	return p.parser.Serialize(output), nil
}

func initExecutors(memory *Memory) map[string]Executor {
	return map[string]Executor{
		"PING": ping(),
		"ECHO": echo(),
		"GET":  get(memory),
		"SET":  set(memory),
	}
}

func ping() Executor {
	return func(resp *RESP) (*RESP, error) {
		return &RESP{
			Type: SimpleString,
			Data: []byte("PONG"),
		}, nil
	}
}

func echo() Executor {
	return func(resp *RESP) (*RESP, error) {
		if len(resp.Nested) < 2 {
			return nil, fmt.Errorf("ECHO command error: input insufficient")
		}
		return &RESP{
			Type: BulkString,
			Data: resp.Nested[1].Data,
		}, nil
	}
}

func set(memory *Memory) Executor {
	return func(resp *RESP) (*RESP, error) {
		if len(resp.Nested) < 3 {
			return nil, fmt.Errorf("insufficient arguments for SET")
		}
		argKey, argVal := resp.Nested[1], resp.Nested[2]
		key, val := string(argKey.Data), string(argVal.Data)

		memory.Put(key, val)

		return &RESP{
			Type: SimpleString,
			Data: []byte("OK"),
		}, nil
	}
}

func get(memory *Memory) Executor {
	return func(resp *RESP) (*RESP, error) {
		if len(resp.Nested) < 2 {
			return nil, fmt.Errorf("insufficient arguments for GET")
		}
		argKey := resp.Nested[1]
		key := string(argKey.Data)

		val := memory.Get(key)

		return &RESP{
			Type: SimpleString,
			Data: []byte(val),
		}, nil
	}
}
