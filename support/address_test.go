package support

import "testing"

func TestValidateAddress(t *testing.T) {
	type args struct {
		address string
		tag     []byte
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"Valid Main Address", args{"491JTRtHaWdYLfgEkkkWgV6tmURjx8W2nRcZ1QUJ9mYSEmqzg5CDfTZCQquNcheYuwDphLbYw8HjYhes9dKtsVBS7aPb6LC", []byte{0x12}}, true, false},
		{"Valid SubAddress", args{"87i7kA61fNvMboXiYWHVygPAggKJPETFqLXXcdH4mQTrECvrTxZMtt6e6owj1k8jUVjNR11eBuBMWHFBtxAwEVcm9dcSUxr", []byte{0x2A}}, true, false},
		{"Valid Main Address, invalid tag", args{"491JTRtHaWdYLfgEkkkWgV6tmURjx8W2nRcZ1QUJ9mYSEmqzg5CDfTZCQquNcheYuwDphLbYw8HjYhes9dKtsVBS7aPb6LC", []byte{0x13}}, false, false},
		{"Valid SubAddress, invalid tag", args{"87i7kA61fNvMboXiYWHVygPAggKJPETFqLXXcdH4mQTrECvrTxZMtt6e6owj1k8jUVjNR11eBuBMWHFBtxAwEVcm9dcSUxr", []byte{0x2B}}, false, false},
		{"Invalid Main Address - Bad Checksum", args{"491JTRtHaWdYLfgEkkkWgV6tmURjx8W2nRcZ1QUJ9mYSEmqzg5CDfTZCQquNcheYuwDphLbYw8HjYhes9dKtsVBS7aPb6LB", []byte{0x12}}, false, true},
		{"Invalid SubAddress - Bad Checksum", args{"87i7kA61fNvMboXiYWHVygPAggKJPETFqLXXcdH4mQTrECvrTxZMtt6e6owj1k8jUVjNR11eBuBMWHFBtxAwEVcm9dcSUxc", []byte{0x2A}}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateAddress(tt.args.address, tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddressDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddressDecode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
