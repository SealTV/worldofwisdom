package server

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"
)

func Test_client_ReadWithTimeout(t *testing.T) {
	type args struct {
		ctx     context.Context
		timeout time.Duration
	}
	tests := []struct {
		name    string
		rw      io.ReadWriter
		args    args
		want    string
		wantErr bool
	}{
		{
			"1. success read",
			bytes.NewBuffer([]byte("some message")),
			args{
				context.Background(),
				1 * time.Second,
			},
			"some message",
			false,
		},
		{
			"2. timeout",
			&longtimeReadWriter{
				result:   []byte("some message"),
				duration: 100 * time.Millisecond,
			},
			args{
				context.Background(),
				1 * time.Microsecond,
			},
			"",
			true,
		},
		{
			"3. error read",
			&longtimeReadWriter{
				err: errors.New("some error"),
			},
			args{
				context.Background(),
				1 * time.Second,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.rw)
			got, err := c.ReadWithTimeout(tt.args.ctx, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.ReadWithTimeout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("client.ReadWithTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_Write(t *testing.T) {
	tests := []struct {
		name    string
		rw      io.ReadWriter
		msg     string
		wantErr bool
	}{
		{
			"1. success write",
			bytes.NewBuffer([]byte{}),
			"some message",
			false,
		},
		{
			"2. error write",
			&longtimeReadWriter{
				err: errors.New("some error"),
			},
			"some message",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.rw)

			if err := c.Write(tt.msg); (err != nil) != tt.wantErr {
				t.Errorf("client.Write() error = %v, wantErr %v\n", err, tt.wantErr)
			}

			if !tt.wantErr {
				str, err := c.Read(context.TODO())
				if err != nil {
					t.Errorf("read error: %v", err)
					return
				}

				t.Log(str, tt.msg)

				if str != tt.msg {
					t.Errorf("client.Write() msg check error, got: %s, want: %s\n", str, tt.msg)
				}
			}
		})
	}
}

type longtimeReadWriter struct {
	err      error
	result   []byte
	duration time.Duration
}

func (rw *longtimeReadWriter) Read(p []byte) (int, error) {
	if rw.err != nil {
		return 0, rw.err
	}

	time.Sleep(rw.duration)

	copy(p, rw.result)
	return len(rw.result), nil
}

func (rw *longtimeReadWriter) Write(p []byte) (int, error) {
	if rw.err != nil {
		return 0, rw.err
	}

	copy(rw.result, p)

	return len(p), nil
}
