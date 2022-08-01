package daemon

import (
	"reflect"
	"testing"
)

func TestGetBlockHeaderByHash(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		args    args
		want    BlockHeader
		wantErr bool
	}{
		{
			name: "Known good hash",
			args: args{hash: "e22cf75f39ae720e8b71b3d120a5ac03f0db50bba6379e2850975b4859190bc6"},
			want: BlockHeader{
				BlockSize:                 210,
				BlockWeight:               210,
				CumulativeDifficulty:      754734824984346,
				CumulativeDifficultyTop64: 0,
				Depth:                     1767506,
				Difficulty:                815625611,
				DifficultyTop64:           0,
				Hash:                      "e22cf75f39ae720e8b71b3d120a5ac03f0db50bba6379e2850975b4859190bc6",
				Height:                    912345,
				LongTermWeight:            210,
				MajorVersion:              1,
				MinerTxHash:               "c7da3965f25c19b8eb7dd8db48dcd4e7c885e2491db77e289f0609bf8e08ec30",
				MinorVersion:              2,
				Nonce:                     1646,
				NumTxes:                   0,
				OrphanStatus:              false,
				PowHash:                   "",
				PrevHash:                  "b61c58b2e0be53fad5ef9d9731a55e8a81d972b8d90ed07c04fd37ca6403ff78",
				Reward:                    7388968946286,
				Timestamp:                 1452793716,
				WideCumulativeDifficulty:  "0x2ae6d65248f1a",
				WideDifficulty:            "0x309d758b",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBlockHeaderByHash(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockHeaderByHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBlockHeaderByHash() got = %v, want %v", got, tt.want)
			}
		})
	}
}
