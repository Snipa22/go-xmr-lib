package wallet

import (
	"testing"
)

func TestSendXMR(t *testing.T) {
	type args struct {
		inXfer XMRWalletTransfer
	}
	tests := []struct {
		name    string
		args    args
		want    XMRTransferReceipt
		wantErr bool
	}{
		{
			name: "Should transfer text XMR",
			args: args{inXfer: XMRWalletTransfer{
				Destinations: []XMRPayment{
					{Amount: 100,
						Address: "9ttGj2BUaDibLAjE8AegzMFnCtYd7qefSYKYWNEh8F7DFzpGYcG9gUfdwzzaPMoWiDfGddoBm8nk253BaS7KV9iHJrWRfnt"},
				},
				Priority: 1,
				RingSize: 16,
				GetTxKey: true,
			}},
			want:    XMRTransferReceipt{},
			wantErr: false,
		},
		{
			name: "Should not transfer text XMR",
			args: args{inXfer: XMRWalletTransfer{
				Destinations: []XMRPayment{
					{Amount: 9223372036854775807,
						Address: "9ttGj2BUaDibLAjE8AegzMFnCtYd7qefSYKYWNEh8F7DFzpGYcG9gUfdwzzaPMoWiDfGddoBm8nk253BaS7KV9iHJrWRfnt"},
				},
				Priority: 1,
				RingSize: 16,
				GetTxKey: true,
			}},
			want:    XMRTransferReceipt{},
			wantErr: true,
		},
		{
			name: "Should not send to mainnet from testnet",
			args: args{inXfer: XMRWalletTransfer{
				Destinations: []XMRPayment{
					{Amount: 100,
						Address: "4A9WRKaBgwd2ZuchFodWXgL2uKYu8h1J4VPVHRiSpDaTDvWpQhVGJ4p5ZvYWF2XzniW1ZDe9p3suZgiUHnEuf6vtSb88gWk"},
				},
				Priority: 1,
				RingSize: 16,
				GetTxKey: true,
			}},
			want:    XMRTransferReceipt{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SendXMR(tt.args.inXfer)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendXMR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
