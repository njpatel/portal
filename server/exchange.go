package server

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	api "github.com/njpatel/portal/api"
)

func init() {
	rand.Seed(int64(os.Getpid())) // This is going to blow up in my face
}

type session struct {
	started    time.Time
	senderIP   string
	receiverIP string
	token      string

	senderChan   chan string
	receiverChan chan *api.Frame
}

type exchange struct {
	sync.Mutex
	sessions map[string]*session
}

func newExchange() *exchange {
	return &exchange{
		sessions: make(map[string]*session),
	}
}

// newSession creates a new session and returns the token and waiting chan
func (e *exchange) newSession(senderIP string) (string, <-chan string, chan<- *api.Frame) {
	t := time.Now()
	sess := &session{
		started:      t,
		token:        e.generateToken(senderIP, t.String()),
		senderIP:     senderIP,
		senderChan:   make(chan string),
		receiverChan: make(chan *api.Frame),
	}

	e.Lock()
	e.sessions[sess.token] = sess
	e.Unlock()

	return sess.token, sess.senderChan, sess.receiverChan
}

func (*exchange) generateToken(ip, time string) string {
	// This'll never clash, lol (FIXME)
	data := fmt.Sprintf("%s%s%f", ip, time, rand.Float64())
	sum := sha1.Sum([]byte(data))
	return hex.EncodeToString(sum[:])[:20]
}

func (e *exchange) delete(token string) {
	e.Lock()
	sess := e.sessions[token]
	if sess != nil {
		close(sess.senderChan)
		close(sess.receiverChan)
		delete(e.sessions, token)
		logger.Tracef("Active sessions: %d", len(e.sessions))
	}
	e.Unlock()
}

func (e *exchange) connect(token, receiverIP string) (<-chan *api.Frame, error) {
	e.Lock()
	defer e.Unlock()

	sess := e.sessions[token]
	if sess == nil {
		return nil, errors.New("Invalid token")
	}

	if sess.receiverIP != "" {
		return nil, errors.New("Session in progress")
	}

	sess.receiverIP = receiverIP

	// Unlock the pending put request
	sess.senderChan <- receiverIP

	return sess.receiverChan, nil
}
