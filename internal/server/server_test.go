package server

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sealtv/worldofwisdom/internal/server/mocks"
)

//go:generate mockgen -destination=mocks/server.go -package=mocks -source=server.go

func TestServer_Run(t *testing.T) {
	canceledCtx, _ := context.WithTimeout(context.Background(), time.Nanosecond)

	tests := []struct {
		name    string
		prepare func(svc *mocks.MockService)
		ctx     context.Context
		wantErr bool
	}{
		{
			"1. success",
			func(svc *mocks.MockService) {
				svc.EXPECT().ProcessClient(gomock.Any()).Return(nil)
			},
			context.Background(),
			false,
		},
		{
			"2. error on handle client",
			func(svc *mocks.MockService) {
				svc.EXPECT().ProcessClient(gomock.Any()).Return(errors.New("unexpected error"))
			},
			context.Background(),
			false,
		},
		{
			"3. canceled context",
			func(svc *mocks.MockService) {},
			canceledCtx,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mocks.NewMockService(ctrl)

			tt.prepare(svc)

			i := 0
			l := &mockListener{
				accept: func() (net.Conn, error) {
					if i == 0 {
						i++
						return &mockConn{
							read:  func(b []byte) (n int, err error) { return 0, nil },
							close: func() error { return nil },
						}, nil
					}
					return nil, &net.OpError{Err: net.ErrClosed}
				},
			}
			s := New(l, svc)

			if err := s.Run(tt.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Server.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_handleIncomingConnections_AcceptConn(t *testing.T) {
	i := 0
	listener := &mockListener{
		accept: func() (net.Conn, error) {
			if i == 0 {
				i++
				return &mockConn{}, nil
			}

			return nil, &net.OpError{Err: net.ErrClosed}
		},
	}

	conns := make(chan net.Conn)
	errs := make(chan error)

	go handleIncomingConnections(listener, conns, errs)

	select {
	case c := <-conns:
		if c == nil {
			t.Errorf("handleIncomingConnections() got = %v, want non-nil", c)
		}
	case err := <-errs:
		if err != nil {
			t.Errorf("handleIncomingConnections() got = %v, want nil", err)
		}
	}
}

func Test_handleIncomingConnections_UnexpectedErr(t *testing.T) {
	i := 0
	listener := &mockListener{
		accept: func() (net.Conn, error) {
			if i == 0 {
				i++
				return nil, errors.New("unexpected error")
			}

			return nil, &net.OpError{Err: net.ErrClosed}
		},
	}
	conns := make(chan net.Conn)
	errs := make(chan error)
	go handleIncomingConnections(listener, conns, errs)

	select {
	case c := <-conns:
		if c != nil {
			t.Errorf("handleIncomingConnections() got = %v, want non-nil", c)
		}
	case err := <-errs:
		if err == nil {
			t.Errorf("handleIncomingConnections() got = %v, want nil", err)
		}
	}
}

// mock for net.Listener
type mockListener struct {
	accept func() (net.Conn, error)
	close  func() error
	addr   func() net.Addr
}

func (m *mockListener) Accept() (net.Conn, error) {
	return m.accept()
}

func (m *mockListener) Close() error {
	return m.close()
}

func (m *mockListener) Addr() net.Addr {
	return m.addr()
}

// mock for net.Conn
type mockConn struct {
	read   func(b []byte) (n int, err error)
	write  func(b []byte) (n int, err error)
	close  func() error
	local  func() net.Addr
	remote func() net.Addr
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.read(b)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.write(b)
}

func (m *mockConn) Close() error {
	return m.close()
}

func (m *mockConn) LocalAddr() net.Addr {
	return m.local()
}

func (m *mockConn) RemoteAddr() net.Addr {
	return m.remote()
}

func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}
