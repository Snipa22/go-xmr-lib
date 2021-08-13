package Base58

import (
	"reflect"
	"testing"
)

func TestDecodeFromString(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"", args{""}, []byte{}, false},
		{"5Q", args{"5Q"}, []byte{0xFF}, false},
		{"LUv", args{"LUv"}, []byte{0xFF, 0xFF}, false},
		{"2UzHL", args{"2UzHL"}, []byte{0xFF, 0xFF, 0xFF}, false},
		{"7YXq9G", args{"7YXq9G"}, []byte{0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"VtB5VXc", args{"VtB5VXc"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"3CUsUpv9t", args{"3CUsUpv9t"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"Ahg1opVcGW", args{"Ahg1opVcGW"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"jpXCZedGfVQ", args{"jpXCZedGfVQ"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"jpXCZedGfVQ5Q", args{"jpXCZedGfVQ5Q"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"jpXCZedGfVQLUv", args{"jpXCZedGfVQLUv"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"jpXCZedGfVQ2UzHL", args{"jpXCZedGfVQ2UzHL"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"jpXCZedGfVQ7YXq9G", args{"jpXCZedGfVQ7YXq9G"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"jpXCZedGfVQVtB5VXc", args{"jpXCZedGfVQVtB5VXc"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"jpXCZedGfVQ3CUsUpv9t", args{"jpXCZedGfVQ3CUsUpv9t"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"jpXCZedGfVQAhg1opVcGW", args{"jpXCZedGfVQAhg1opVcGW"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"jpXCZedGfVQjpXCZedGfVQ", args{"jpXCZedGfVQjpXCZedGfVQ"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeFromString(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeToString(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"No Data", args{[]byte{}}, "", false},
		{"1 bytes of 00", args{[]byte{0x00}}, "11", false},
		{"2 bytes of 00", args{[]byte{0x00, 0x00}}, "111", false},
		{"3 bytes of 00", args{[]byte{0x00, 0x00, 0x00}}, "11111", false},
		{"4 bytes of 00", args{[]byte{0x00, 0x00, 0x00, 0x00}}, "111111", false},
		{"5 bytes of 00", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00}}, "1111111", false},
		{"6 bytes of 00", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, "111111111", false},
		{"7 bytes of 00", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, "1111111111", false},
		{"8 bytes of 00", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, "11111111111", false},
		{"9 bytes of 00", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, "1111111111111", false},
		{"10 bytes of 00", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, "11111111111111", false},
		{"11 bytes of 00", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, "1111111111111111", false},
		{"12 bytes of 00", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, "11111111111111111", false},
		{"13 bytes of 00", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, "111111111111111111", false},
		{"14 bytes of 00", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, "11111111111111111111", false},
		{"15 bytes of 00", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, "111111111111111111111", false},
		{"16 bytes of 00", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, "1111111111111111111111", false},
		{"0x06156013762879F7FFFFFFFFFF", args{[]byte{0x06, 0x15, 0x60, 0x13, 0x76, 0x28, 0x79, 0xF7, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}}, "22222222222VtB5VXc", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeToString(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncodeToString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeBlock(t *testing.T) {
	type args struct {
		toDecode string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{

		// Bad Block Sizes
		{"Too Small", args{""}, nil, true},
		{"Too large", args{"111111111111"}, nil, true},
		// 1-byte block
		{"1", args{"1"}, nil, true},
		{"z", args{"z"}, nil, true},
		// 2-bytes block
		{"11", args{"11"}, []byte{0x00}, false},
		{"5Q", args{"5Q"}, []byte{0xFF}, false},
		{"5R", args{"5R"}, nil, true},
		{"zz", args{"zz"}, nil, true},
		// 3-bytes block
		{"111", args{"111"}, []byte{0x00, 0x00}, false},
		{"LUv", args{"LUv"}, []byte{0xFF, 0xFF}, false},
		{"LUw", args{"LUw"}, nil, true},
		{"zzz", args{"zzz"}, nil, true},
		// 4-bytes block
		{"1111", args{"1111"}, nil, true},
		{"zzzz", args{"zzzz"}, nil, true},
		// 5-bytes block
		{"11111", args{"11111"}, []byte{0x00, 0x00, 0x00}, false},
		{"2UzHL", args{"2UzHL"}, []byte{0xFF, 0xFF, 0xFF}, false},
		{"2UzHM", args{"2UzHM"}, nil, true},
		{"zzzzz", args{"zzzzz"}, nil, true},
		// 6-bytes block
		{"111111", args{"111111"}, []byte{0x00, 0x00, 0x00, 0x00}, false},
		{"7YXq9G", args{"7YXq9G"}, []byte{0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"7YXq9H", args{"7YXq9H"}, nil, true},
		{"zzzzzz", args{"zzzzzz"}, nil, true},
		// 7-bytes block
		{"1111111", args{"1111111"}, []byte{0x00, 0x00, 0x00, 0x00, 0x00}, false},
		{"VtB5VXc", args{"VtB5VXc"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"VtB5VXd", args{"VtB5VXd"}, nil, true},
		{"zzzzzzz", args{"zzzzzzz"}, nil, true},
		// 8-bytes block
		{"11111111", args{"11111111"}, nil, true},
		{"zzzzzzzz", args{"zzzzzzzz"}, nil, true},
		// 9-bytes block
		{"111111111", args{"111111111"}, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, false},
		{"3CUsUpv9t", args{"3CUsUpv9t"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"3CUsUpv9u", args{"3CUsUpv9u"}, nil, true},
		{"zzzzzzzzz", args{"zzzzzzzzz"}, nil, true},
		// 10-bytes block
		{"1111111111", args{"1111111111"}, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, false},
		{"Ahg1opVcGW", args{"Ahg1opVcGW"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"Ahg1opVcGX", args{"Ahg1opVcGX"}, nil, true},
		{"zzzzzzzzzz", args{"zzzzzzzzzz"}, nil, true},
		// 11-bytes block
		{"11111111111", args{"11111111111"}, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, false},
		{"jpXCZedGfVQ", args{"jpXCZedGfVQ"}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"jpXCZedGfVR", args{"jpXCZedGfVR"}, nil, true},
		{"zzzzzzzzzzz", args{"zzzzzzzzzzz"}, nil, true},
		// Invalid symbols
		{"01111111111", args{"01111111111"}, nil, true},
		{"11111111110", args{"11111111110"}, nil, true},
		{"11111011111", args{"11111011111"}, nil, true},
		{"I1111111111", args{"I1111111111"}, nil, true},
		{"O1111111111", args{"O1111111111"}, nil, true},
		{"l1111111111", args{"l1111111111"}, nil, true},
		{"_1111111111", args{"_1111111111"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeBlock(tt.args.toDecode)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeBlock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeBlock(t *testing.T) {
	type args struct {
		toEncode []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// Single byte tests
		{"0x00", args{[]byte{0x00}}, "11", false},
		{"0x39", args{[]byte{0x39}}, "1z", false},
		{"0xFF", args{[]byte{0xFF}}, "5Q", false},

		// Two byte tests
		{"0x0000", args{[]byte{0x00, 0x00}}, "111", false},
		{"0x0039", args{[]byte{0x00, 0x39}}, "11z", false},
		{"0x0100", args{[]byte{0x01, 0x00}}, "15R", false},
		{"0xFFFF", args{[]byte{0xFF, 0xFF, 0xFF}}, "2UzHL", false},

		{"0x000000", args{[]byte{0x00, 0x00, 0x00}}, "11111", false},
		{"0x000039", args{[]byte{0x00, 0x00, 0x39}}, "1111z", false},
		{"0x010000", args{[]byte{0x01, 0x00, 0x00}}, "11LUw", false},
		{"0xFFFFFF", args{[]byte{0xFF, 0xFF, 0xFF}}, "2UzHL", false},

		{"0x00000039", args{[]byte{0x00, 0x00, 0x00, 0x39}}, "11111z", false},
		{"0xFFFFFFFF", args{[]byte{0xFF, 0xFF, 0xFF, 0xFF}}, "7YXq9G", false},
		{"0x0000000039", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x39}}, "111111z", false},
		{"0xFFFFFFFFFF", args{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}}, "VtB5VXc", false},
		{"0x000000000039", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x39}}, "11111111z", false},
		{"0xFFFFFFFFFFFF", args{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}}, "3CUsUpv9t", false},
		{"0x00000000000039", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x39}}, "111111111z", false},
		{"0xFFFFFFFFFFFFFF", args{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}}, "Ahg1opVcGW", false},
		{"0x0000000000000039", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x39}}, "1111111111z", false},
		{"0xFFFFFFFFFFFFFFFF", args{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}}, "jpXCZedGfVQ", false},

		{"0x0000000000000000", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, "11111111111", false},
		{"0x0000000000000001", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}}, "11111111112", false},
		{"0x0000000000000008", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08}}, "11111111119", false},
		{"0x0000000000000009", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09}}, "1111111111A", false},
		{"0x000000000000003A", args{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3A}}, "11111111121", false},
		{"0x00FFFFFFFFFFFFFF", args{[]byte{0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}}, "1Ahg1opVcGW", false},
		{"0x06156013762879F7", args{[]byte{0x06, 0x15, 0x60, 0x13, 0x76, 0x28, 0x79, 0xF7}}, "22222222222", false},
		{"0x05E022BA374B2A00", args{[]byte{0x05, 0xE0, 0x22, 0xBA, 0x37, 0x4B, 0x2A, 0x00}}, "1z111111111", false},
		{"No Data", args{[]byte{}}, "", true},
		{"9 Bytes", args{[]byte{0x05, 0xE0, 0x22, 0xBA, 0x37, 0x4B, 0x2A, 0x00, 0x00}}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeBlock(tt.args.toEncode)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("encodeBlock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uint8BeTo64(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"0x0000000000000001", args{[]byte{1}}, 0x01, false},
		{"0x0000000000000102", args{[]byte{1, 2}}, 0x0102, false},
		{"0x0000000000010203", args{[]byte{1, 2, 3}}, 0x010203, false},
		{"0x0000000001020304", args{[]byte{1, 2, 3, 4}}, 0x01020304, false},
		{"0x0000000102030405", args{[]byte{1, 2, 3, 4, 5}}, 0x0102030405, false},
		{"0x0000010203040506", args{[]byte{1, 2, 3, 4, 5, 6}}, 0x010203040506, false},
		{"0x0001020304050607", args{[]byte{1, 2, 3, 4, 5, 6, 7}}, 0x01020304050607, false},
		{"0x0102030405060708", args{[]byte{1, 2, 3, 4, 5, 6, 7, 8}}, 0x0102030405060708, false},
		{"0x010203040506070809", args{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}}, 0x00, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uint8BeTo64(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("uint8BeTo64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("uint8BeTo64() got = %v, want %v", got, tt.want)
			}
		})
	}
}
