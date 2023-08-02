package app

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sealtv/worldofwisdom/internal/app/mocks"
)

//go:generate mockgen -destination=mocks/app.go -package=mocks -source=app.go

func TestApp_ProcessClient(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(cli *mocks.MockClienter, pow *mocks.MockPoWer, wb *mocks.MockWisdomBooker)
		wantErr bool
	}{
		{
			"1. success",
			func(cli *mocks.MockClienter, pow *mocks.MockPoWer, wb *mocks.MockWisdomBooker) {
				rnd := "random"
				pow.EXPECT().GetChallenge().Return(rnd)

				call := cli.EXPECT().Write(rnd).Return(nil)
				call = cli.EXPECT().ReadWithTimeout(gomock.Any(), gomock.Any()).Return("response", nil).After(call)
				call = pow.EXPECT().IsValid(rnd + "response").Return(true).After(call)

				quote := "quote"
				call = wb.EXPECT().GetRandomQuote().Return(quote).After(call)
				cli.EXPECT().Write(quote).Return(nil).After(call)
			},
			false,
		},
		{
			"2. error on write resp",
			func(cli *mocks.MockClienter, pow *mocks.MockPoWer, wb *mocks.MockWisdomBooker) {
				rnd := "random"
				pow.EXPECT().GetChallenge().Return(rnd)

				call := cli.EXPECT().Write(rnd).Return(nil)
				call = cli.EXPECT().ReadWithTimeout(gomock.Any(), gomock.Any()).Return("response", nil).After(call)
				call = pow.EXPECT().IsValid(rnd + "response").Return(true).After(call)

				quote := "quote"
				call = wb.EXPECT().GetRandomQuote().Return(quote).After(call)
				cli.EXPECT().Write(quote).Return(errors.New("some error")).After(call)
			},
			false,
		},
		{
			"3. invalid pow",
			func(cli *mocks.MockClienter, pow *mocks.MockPoWer, wb *mocks.MockWisdomBooker) {
				rnd := "random"
				call := pow.EXPECT().GetChallenge().Return(rnd)
				call = cli.EXPECT().Write(rnd).Return(nil).After(call)
				call = cli.EXPECT().ReadWithTimeout(gomock.Any(), gomock.Any()).Return("response", nil).After(call)
				call = pow.EXPECT().IsValid(rnd + "response").Return(false).After(call)
				cli.EXPECT().Write("Invalid PoW response").Return(nil).After(call)
			},
			false,
		},
		{
			"4. error on write invalid pow resp",
			func(cli *mocks.MockClienter, pow *mocks.MockPoWer, wb *mocks.MockWisdomBooker) {
				rnd := "random"
				call := pow.EXPECT().GetChallenge().Return(rnd)
				call = cli.EXPECT().Write(rnd).Return(nil).After(call)
				call = cli.EXPECT().ReadWithTimeout(gomock.Any(), gomock.Any()).Return("response", nil).After(call)
				call = pow.EXPECT().IsValid(rnd + "response").Return(false).After(call)
				cli.EXPECT().Write("Invalid PoW response").Return(errors.New("some error")).After(call)
			},
			false,
		},
		{
			"5. error on read resp",
			func(cli *mocks.MockClienter, pow *mocks.MockPoWer, wb *mocks.MockWisdomBooker) {
				rnd := "random"
				call := pow.EXPECT().GetChallenge().Return(rnd)
				call = cli.EXPECT().Write(rnd).Return(nil).After(call)
				cli.EXPECT().ReadWithTimeout(gomock.Any(), gomock.Any()).Return("", errors.New("some error")).After(call)
			},
			false,
		},
		{
			"6. error on write challange",
			func(cli *mocks.MockClienter, pow *mocks.MockPoWer, wb *mocks.MockWisdomBooker) {
				rnd := "random"
				call := pow.EXPECT().GetChallenge().Return(rnd)
				cli.EXPECT().Write(rnd).Return(errors.New("some error")).After(call)
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cli := mocks.NewMockClienter(ctrl)
			pow := mocks.NewMockPoWer(ctrl)
			wb := mocks.NewMockWisdomBooker(ctrl)

			tt.prepare(cli, pow, wb)

			a := NewApp(pow, wb)
			if err := a.ProcessClient(cli); (err != nil) != tt.wantErr {
				t.Errorf("App.ProcessClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
