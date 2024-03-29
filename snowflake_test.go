package idgenerator

import (
	"fmt"
	"testing"
	"time"
)

func TestNewSnowflakeID(t *testing.T) {
	type args struct {
		opts []option
	}
	tests := []struct {
		name       string
		args       args
		want       int64
		wantErr    bool
		outputBits string
	}{
		{
			"WithTimestamp:2024-02-01 WithDatacenterID:31 WithMachineID:15 WithSequenceNumber:1",
			args{[]option{
				WithTimestamp(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
				WithDatacenterID(31),
				WithMachineID(15),
				WithSequenceNumber(1),
			}},
			11234023837724673,
			false,
			"0000000000100111111010010100100100000000001111101111000000000001",
		},
		{
			"Error invalid datacenter ID",
			args{[]option{
				WithDatacenterID(32),
			}},
			0,
			true,
			"0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			"Error invalid machine ID",
			args{[]option{
				WithMachineID(32),
			}},
			0,
			true,
			"0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			"Error invalid sequence number",
			args{[]option{
				WithSequenceNumber(4096),
			}},
			0,
			true,
			"0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			"Error over the maximum lifetime",
			args{[]option{
				WithTimestamp(time.Date(2050, 1, 1, 0, 0, 0, 0, time.UTC)),
				WithBaseTime(time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)),
			}},
			0,
			true,
			"0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			"Error invalid timestamp",
			args{[]option{
				WithTimestamp(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
				WithBaseTime(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			}},
			0,
			true,
			"0000000000000000000000000000000000000000000000000000000000000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSnowflakeID(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSnowflakeID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewSnowflakeID() int64 = %v, want %v", got, tt.want)
			}
			gotBits := fmt.Sprintf("%064b", got)
			if gotBits != tt.outputBits {
				t.Errorf("NewSnowflakeID() bits = %v, want %v", gotBits, tt.outputBits)
			}
		})
	}
}
