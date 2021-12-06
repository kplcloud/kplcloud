/**
 * @Time : 2021/12/6 10:19 AM
 * @Author : solacowa@gmail.com
 * @File : pty_handler
 * @Software: GoLand
 */

package terminal

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"io"
	"k8s.io/client-go/tools/remotecommand"
)

const END_OF_TRANSMISSION = "\u0004"

type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

type Session struct {
	id            string
	sockJSSession sockjs.Session
	sizeChan      chan remotecommand.TerminalSize
	logger        log.Logger
	bound         chan error
	doneChan      chan struct{}
}

type Message struct {
	Op        string `json:"op"`
	SessionID string `json:"sessionId"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
	Data      string `json:"data"`
}

type Result struct {
	SessionId string `json:"sessionId,omitempty"`
	Token     string `json:"token,omitempty"`
	Cluster   string `json:"cluster,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	PodName   string `json:"podName,omitempty"`
	Container string `json:"container,omitempty"`
	Cmd       string `json:"cmd,omitempty"`
}

func (t Session) Read(p []byte) (int, error) {
	m, err := t.sockJSSession.Recv()
	if err != nil {
		err = errors.Wrap(err, "sockJSSession.Recv")
		return 0, err
	}

	var msg Message
	if err := json.Unmarshal([]byte(m), &msg); err != nil {
		fmt.Println(fmt.Sprintf("read msg (%s) form client error.%v", string(p), err))
		return 0, err
	}
	switch msg.Op {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{msg.Cols, msg.Rows}
		return 0, nil
	default:
		return copy(p, END_OF_TRANSMISSION), fmt.Errorf("unknown message type '%s'", msg.Op)
	}
}

func (t Session) Write(p []byte) (int, error) {
	msg, err := json.Marshal(Message{
		Op:   "stdout",
		Data: string(p),
	})
	if err != nil {
		return 0, err
	}

	if err = t.sockJSSession.Send(string(msg)); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (t Session) Close(status uint32, reason string) {
	err := t.sockJSSession.Close(status, reason)
	fmt.Println(fmt.Sprintf("close socket (%s). %d, %s, %v", t.id, status, reason, err))
}

// Next TerminalSize handles pty->process resize events
// Called in a loop from remotecommand as long as the process is running
func (t Session) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

// Toast can be used to send the user any OOB messages
// hterm puts these in the center of the terminal
func (t Session) Toast(p string) error {
	msg, err := json.Marshal(Message{
		Op:   "toast",
		Data: p,
	})
	if err != nil {
		return err
	}

	if err = t.sockJSSession.Send(string(msg)); err != nil {
		return err
	}
	return nil
}
