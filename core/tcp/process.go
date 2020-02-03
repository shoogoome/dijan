package tcp

import (
	"bufio"
	"dijan/core/cluster"
	"dijan/utils"
	"io"
	"log"
	"net"
	"strings"
)

type result struct {
	v []byte
	e error
}

func (s *Server) get(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, e := s.readKey(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		v, e := s.Get(k)
		c <- &result{v, e}
	}()
}

func (s *Server) set(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, v, t, e := s.readKeyAndValue(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		c <- &result{nil, s.Set(k, v, t)}
	}()
}

func (s *Server) del(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, e := s.readKey(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		c <- &result{nil, s.Del(k)}
	}()
}

func (s *Server) member(ch chan chan*result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	go func() {
		nodes := make([]string, cluster.Member.NumMembers())
		for index, node := range cluster.Member.Members() {
			nodes[index] = node.Name
		}
		c <- &result{[]byte(strings.Join(nodes, " ")), nil}
	}()
}

func reply(conn net.Conn, resultCh chan chan *result) {
	defer conn.Close()
	for {
		c, open := <-resultCh
		if !open {
			return
		}
		r := <-c
		e := sendResponse(r.v, r.e, conn)
		if e != nil {
			log.Println("close connection due to error:", e)
			return
		}
	}
}

func (s *Server) process(conn net.Conn) {
	r := bufio.NewReader(conn)
	resultCh := make(chan chan *result, utils.GlobalSystemConfig.Rocksdb.AsynchronousNumber)
	defer close(resultCh)
	go reply(conn, resultCh)
	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("close connection due to error:", e)
			}
			return
		}

		switch op {
		case 'S':
			s.set(resultCh, r)
		case 'G':
			s.get(resultCh, r)
		case 'D':
			s.del(resultCh, r)
		case 'M':
			s.member(resultCh, r)
		default:
			log.Println("close connection due to invalid operation:", string(op))
			return
		}
	}
}
