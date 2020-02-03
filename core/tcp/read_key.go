package tcp

import (
	"bufio"
	"dijan/utils"
	"errors"
	"io"
	"strconv"
)

func (s *Server) readKey(r *bufio.Reader) (string, error) {
	klen, e := readLen(r)
	if e != nil {
		return "", e
	}
	k := make([]byte, klen)
	_, e = io.ReadFull(r, k)
	if e != nil {
		return "", e
	}
	key := string(k)
	addr, ok := s.ShouldProcess(key)
	if !ok {
		return "", errors.New(addr + utils.GlobalSystemConfig.Server.TcpListenPort)
	}
	return key, nil
}

func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, int, error) {
	klen, e := readLen(r)
	if e != nil {
		return "", nil, 0, e
	}
	vlen, e := readLen(r)
	if e != nil {
		return "", nil, 0, e
	}
	tlen, e := readLen(r)
	if e != nil {
		return "", nil, 0, e
	}
	k := make([]byte, klen)
	_, e = io.ReadFull(r, k)
	if e != nil {
		return "", nil, 0, e
	}
	key := string(k)
	addr, ok := s.ShouldProcess(key)

	v := make([]byte, vlen)
	_, e = io.ReadFull(r, v)
	if e != nil {
		if !ok {
			return "", nil, 0, errors.New(addr)
		}
		return "", nil, 0, e
	}
	if tlen == 0 {
		if !ok {
			return "", nil, 0, errors.New(addr)
		}
		return key, v, -1, nil
	}
	t := make([]byte, tlen)
	_, e = io.ReadFull(r, t)
	if e != nil {
		if !ok {
			return "", nil, 0, errors.New(addr)
		}
		return "", nil, 0, e
	}
	ttl, e := strconv.Atoi(string(t))
	if e != nil {
		if !ok {
			return "", nil, 0, errors.New(addr)
		}
		return "", nil, 0, e
	}
	if !ok {
		return "", nil, 0, errors.New(addr)
	}
	return key, v, ttl, nil
}
