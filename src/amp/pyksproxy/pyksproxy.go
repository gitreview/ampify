// Public Domain (-) 2010-2011 The Ampify Authors.
// See the Ampify UNLICENSE file for details.

// The pyksproxy package provides a client library to interact with Scalien's
// Keyspace nodes via the Python ``keyspace-proxy``.
package pyksproxy

import (
	"amp/runtime"
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	SUCCESS         = '\x00'
	RESPONSE_NIL    = '\x01'
	RESPONSE_NUMBER = '\x02'
	RESPONSE_STRING = '\x03'
	RESPONSE_ARRAY  = '\x04'
	RESPONSE_MAP    = '\x05'
	FAILED          = '\xf1'
	NO_SERVICE      = '\xf2'
	UNKNOWN_ERROR   = '\xf3'
	UNEXPECTED_TYPE = '\xf4'
)

var (
	scriptPath = runtime.AmpifyRoot + "/environ/keyspace-proxy"
)

type KeyspaceError struct {
	Message string
}

func (err *KeyspaceError) String() string {
	return "KeyspaceError: " + string(err.Message)
}

type KeyspaceProxy struct {
	Connected bool
	Servers   string
	process   *os.Process
	stdin     *os.File
	stdout    *os.File
	stderr    *os.File
	lock      *sync.Mutex
}

func Keyspace(servers string) (*KeyspaceProxy, os.Error) {
	serverList := strings.Split(servers, " ")
	argv := make([]string, len(serverList)+1)
	argv[0] = scriptPath
	for idx, server := range serverList {
		argv[idx+1] = server
	}
	stdinRead, stdinWrite, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	stdoutRead, stdoutWrite, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	stderrRead, stderrWrite, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	fd := []*os.File{stdinRead, stdoutWrite, stderrWrite}
	process, err := os.StartProcess(scriptPath, argv, &os.ProcAttr{
		Dir:   "",
		Env:   os.Environ(),
		Files: fd,
	})
	if err != nil {
		return nil, err
	}
	return &KeyspaceProxy{
		Servers: servers,
		process: process,
		stdin:   stdinWrite,
		stdout:  stdoutRead,
		stderr:  stderrRead,
		lock:    &sync.Mutex{},
	},
		nil
}

func EncodeSize(value int) (result []byte) {
	for {
		leftBits := value & 127
		value >>= 7
		if value > 0 {
			leftBits += 128
		}
		result = append(result, byte(leftBits))
		if value == 0 {
			break
		}
	}
	return result
}

func PackNumber(value int, buffer *bytes.Buffer) {
	buffer.Write([]byte("\x04"))
	buffer.Write(EncodeSize(value))
}

func PackString(value string, buffer *bytes.Buffer) {
	buffer.Write([]byte("\x05"))
	buffer.Write(EncodeSize(len(value)))
	buffer.Write([]byte(value))
}

func (keyspace *KeyspaceProxy) Send(command string, args ...string) ([]byte, os.Error) {
	buffer := bytes.NewBuffer([]byte(command))
	if len(args) > 0 {
		buffer.Write(EncodeSize(len(args)))
		for _, arg := range args {
			PackString(arg, buffer)
		}
	}
	_, err := keyspace.stdin.Write(buffer.Bytes())
	if err != nil {
		return nil, err
	}
	responseSlice := make([]byte, 1)
	_, err = keyspace.stdout.Read(responseSlice)
	if err != nil {
		return nil, err
	}
	responseByte := responseSlice[0]
	switch responseByte {
	case SUCCESS:
		return nil, nil
	case RESPONSE_NIL:
		return nil, nil
	case RESPONSE_STRING:
		return nil, nil
	case FAILED:
		return nil, &KeyspaceError{"FAILED"}
	case NO_SERVICE:
		return nil, &KeyspaceError{"NO_SERVICE"}
	case UNEXPECTED_TYPE:
		return nil, &KeyspaceError{"UNEXPECTED_TYPE"}
	}
	return nil, nil
}

func (keyspace *KeyspaceProxy) Get(key string) ([]byte, os.Error) {
	return keyspace.Send("\x03", key)
}

func (keyspace *KeyspaceProxy) GetDirty(key string) ([]byte, os.Error) {
	return keyspace.Send("\x04", key)
}

func (keyspace *KeyspaceProxy) Set(key string, value string) ([]byte, os.Error) {
	return keyspace.Send("\x05", key, value)
}

func (keyspace *KeyspaceProxy) Lock() {
	keyspace.lock.Lock()
}

func (keyspace *KeyspaceProxy) Unlock() {
	keyspace.lock.Unlock()
}

func (keyspace *KeyspaceProxy) String() string {
	return fmt.Sprintf("<keyspace: %s>", keyspace.Servers)
}
