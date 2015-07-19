package main

import (
	"golang.org/x/net/context"
)

import (
	. "proto"
)

const (
	SERVICE = "[AUTH]"
)
const (
	_port = ":50006"
)

type server struct {
}

func (s *server) init() {
}

func (s *server) Auth(context.Context, *Auth_Certificate) (*Auth_Result, error) {
	return nil, nil
}
