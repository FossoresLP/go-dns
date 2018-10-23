package query

import (
	"reflect"
	"testing"

	"github.com/fossoreslp/go-dns/dns/label"
	"github.com/fossoreslp/go-dns/dns/record-names"
)

func TestQuery_String(t *testing.T) {
	tests := []struct {
		name string
		q    Query
		want string
	}{
		{"Normal", Query{label.Label{"test", "example", "com"}, names.QTYPE_ANY, names.QCLASS_ANY}, "Query for \"test.example.com.\", Type: 255, Class: 255"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.String(); got != tt.want {
				t.Errorf("Query.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		name label.Label
		t    names.QTYPE
	}
	tests := []struct {
		name string
		args args
		want Query
	}{
		{"Normal", args{label.Label{"test", "example", "com"}, names.QTYPE(names.A)}, Query{label.Label{"test", "example", "com"}, names.QTYPE(names.A), names.QCLASS(names.IN)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.name, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Encode(t *testing.T) {
	tests := []struct {
		name string
		q    Query
		want []byte
	}{
		{"Normal", Query{label.Label{"test", "example", "com"}, names.QTYPE(names.A), names.QCLASS(names.IN)}, []byte{0x04, 0x74, 0x65, 0x73, 0x74, 0x07, 0x65, 0x78, 0x61, 0x6D, 0x70, 0x6C, 0x65, 0x03, 0x63, 0x6F, 0x6D, 0x00, 0x0, 0x01, 0x0, 0x01}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		message []byte
		number  uint16
	}
	tests := []struct {
		name    string
		args    args
		want    []Query
		want1   int
		wantErr bool
	}{
		{"Normal", args{[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x04, 0x74, 0x65, 0x73, 0x74, 0x07, 0x65, 0x78, 0x61, 0x6D, 0x70, 0x6C, 0x65, 0x03, 0x63, 0x6F, 0x6D, 0x00, 0x0, 0x01, 0x0, 0x01}, 1}, []Query{Query{label.Label{"test", "example", "com"}, names.QTYPE(names.A), names.QCLASS(names.IN)}}, 34, false},
		{"Null", args{nil, 1}, nil, 0, true},
		{"Header Only", args{[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, 0}, []Query{}, 12, false},
		{"Label outside message bounds", args{[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, 1}, nil, 0, true},
		{"Message too short for Type and Class", args{[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, 1}, nil, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Parse(tt.args.message, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
