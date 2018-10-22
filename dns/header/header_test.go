package header

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		args    []byte
		want    *Header
		wantErr bool
	}{
		{"Empty input", []byte(nil), nil, true},
		{"Input too short", []byte("Hello world"), nil, true},
		{"Normal", []byte("Testing header"), &Header{[2]byte{'T', 'e'}, [2]byte{'s', 't'}, 26990, 26400, 26725, 24932}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeader_IsResponse(t *testing.T) {
	tests := []struct {
		name string
		h    Header
		want bool
	}{
		{"Query", Header{[2]byte{0x0, 0x0}, [2]byte{0x0, 0x0}, 0, 0, 0, 0}, false},
		{"Response", Header{[2]byte{0x0, 0x0}, [2]byte{0x80, 0x0}, 0, 0, 0, 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.IsResponse(); got != tt.want {
				t.Errorf("Header.IsResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeader_OpCode(t *testing.T) {
	tests := []struct {
		name string
		h    Header
		want uint8
	}{
		{"OpCode 0", Header{[2]byte{0x0, 0x0}, [2]byte{0x0, 0x0}, 0, 0, 0, 0}, 0},
		{"OpCode 15", Header{[2]byte{0x0, 0x0}, [2]byte{0x78, 0x0}, 0, 0, 0, 0}, 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.OpCode(); got != tt.want {
				t.Errorf("Header.OpCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeader_AuthoritativeAnswer(t *testing.T) {
	tests := []struct {
		name string
		h    Header
		want bool
	}{
		{"AA unset", Header{[2]byte{0x0, 0x0}, [2]byte{0x0, 0x0}, 0, 0, 0, 0}, false},
		{"AA set", Header{[2]byte{0x0, 0x0}, [2]byte{0x4, 0x0}, 0, 0, 0, 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.AuthoritativeAnswer(); got != tt.want {
				t.Errorf("Header.AuthoritativeAnswer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeader_Truncated(t *testing.T) {
	tests := []struct {
		name string
		h    Header
		want bool
	}{
		{"TC unset", Header{[2]byte{0x0, 0x0}, [2]byte{0x0, 0x0}, 0, 0, 0, 0}, false},
		{"TC set", Header{[2]byte{0x0, 0x0}, [2]byte{0x2, 0x0}, 0, 0, 0, 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Truncated(); got != tt.want {
				t.Errorf("Header.Truncated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeader_RecursionDesired(t *testing.T) {
	tests := []struct {
		name string
		h    Header
		want bool
	}{
		{"RD unset", Header{[2]byte{0x0, 0x0}, [2]byte{0x0, 0x0}, 0, 0, 0, 0}, false},
		{"RD set", Header{[2]byte{0x0, 0x0}, [2]byte{0x1, 0x0}, 0, 0, 0, 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.RecursionDesired(); got != tt.want {
				t.Errorf("Header.RecursionDesired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeader_RecursionAvailable(t *testing.T) {
	tests := []struct {
		name string
		h    Header
		want bool
	}{
		{"RA unset", Header{[2]byte{0x0, 0x0}, [2]byte{0x0, 0x0}, 0, 0, 0, 0}, false},
		{"RA set", Header{[2]byte{0x0, 0x0}, [2]byte{0x0, 0x80}, 0, 0, 0, 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.RecursionAvailable(); got != tt.want {
				t.Errorf("Header.RecursionAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeader_ZeroBits(t *testing.T) {
	tests := []struct {
		name string
		h    Header
		want bool
	}{
		{"Z unset", Header{[2]byte{0x0, 0x0}, [2]byte{0x0, 0x0}, 0, 0, 0, 0}, true},
		{"Z set", Header{[2]byte{0x0, 0x0}, [2]byte{0x0, 0x70}, 0, 0, 0, 0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.ZeroBits(); got != tt.want {
				t.Errorf("Header.ZeroBits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeader_ResponseCode(t *testing.T) {
	tests := []struct {
		name string
		h    Header
		want uint8
	}{
		{"RCode 0", Header{[2]byte{0x0, 0x0}, [2]byte{0x0, 0x0}, 0, 0, 0, 0}, 0},
		{"RCode 15", Header{[2]byte{0x0, 0x0}, [2]byte{0x0, 0xF}, 0, 0, 0, 0}, 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.ResponseCode(); got != tt.want {
				t.Errorf("Header.ResponseCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeader_String(t *testing.T) {
	tests := []struct {
		name string
		h    Header
		want string
	}{
		{"All active", Header{[2]byte{0xC0, 0xDE}, [2]byte{0xE7, 0x8D}, 1, 1, 1, 1}, "Message with ID C0DE.\nThe flags have the following values: QR: true, OpCode: 12, AA: true, TC: true, RD: true, RA: true, Zero bits are actually zero: true, RCode: 13.\nThe message contains 1 questions, 1 answers, 1 authority pointers and 1 additional records."},
		{"Some active", Header{[2]byte{0xC0, 0xDE}, [2]byte{0x10, 0x73}, 1, 0, 0, 0}, "Message with ID C0DE.\nThe flags have the following values: QR: false, OpCode: 2, AA: false, TC: false, RD: false, RA: false, Zero bits are actually zero: false, RCode: 3.\nThe message contains 1 questions, 0 answers, 0 authority pointers and 0 additional records."},
		{"None active", Header{[2]byte{0xC0, 0xDE}, [2]byte{0x0, 0x20}, 0, 0, 0, 0}, "Message with ID C0DE.\nThe flags have the following values: QR: false, OpCode: 0, AA: false, TC: false, RD: false, RA: false, Zero bits are actually zero: false, RCode: 0.\nThe message contains 0 questions, 0 answers, 0 authority pointers and 0 additional records."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.String(); got != tt.want {
				t.Errorf("Header.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeader_Encode(t *testing.T) {
	tests := []struct {
		name string
		h    Header
		want []byte
	}{
		{"Normal", Header{[2]byte{0x12, 0x34}, [2]byte{0x56, 0x78}, 37035, 52719, 4660, 22136}, []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF, 0x12, 0x34, 0x56, 0x78}},
		{"Null", Header{[2]byte{0x0, 0x0}, [2]byte{0x0, 0x0}, 0, 0, 0, 0}, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Header.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryHeader(t *testing.T) {
	tests := []struct {
		name string
		rd   bool
		want *Header
	}{
		{"Normal", false, &Header{[2]byte{0x0, 0x0}, [2]byte{0x0, 0x0}, 0, 0, 0, 0}},
		{"RD", true, &Header{[2]byte{0x0, 0x0}, [2]byte{0x01, 0x0}, 0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewQueryHeader(tt.rd)
			if got.ID == [2]byte{0x0, 0x0} {
				t.Error("NewQueryHeader() should never have a null ID")
			}
			got.ID = [2]byte{0x0, 0x0}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewErrorHeader(t *testing.T) {
	type args struct {
		id    [2]byte
		aa    bool
		rcode uint8
	}
	tests := []struct {
		name string
		args args
		want *Header
	}{
		{"AA set", args{[2]byte{0x12, 0x34}, true, 9}, &Header{[2]byte{0x12, 0x34}, [2]byte{0x84, 0x89}, 0, 0, 0, 0}},
		{"AA unset", args{[2]byte{0x12, 0x34}, false, 5}, &Header{[2]byte{0x12, 0x34}, [2]byte{0x80, 0x85}, 0, 0, 0, 0}},
		{"Null", args{[2]byte{0x0, 0x0}, false, 0}, &Header{[2]byte{0x0, 0x0}, [2]byte{0x80, 0x80}, 0, 0, 0, 0}},
		{"RCode > 10", args{[2]byte{0x12, 0x34}, false, 15}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewErrorHeader(tt.args.id, tt.args.aa, tt.args.rcode); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewErrorHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAnswerHeader(t *testing.T) {
	type args struct {
		id [2]byte
		aa bool
		rd bool
	}
	tests := []struct {
		name string
		args args
		want *Header
	}{
		{"Normal", args{[2]byte{0x12, 0x34}, false, false}, &Header{[2]byte{0x12, 0x34}, [2]byte{0x80, 0x80}, 0, 0, 0, 0}},
		{"AA", args{[2]byte{0x12, 0x34}, true, false}, &Header{[2]byte{0x12, 0x34}, [2]byte{0x84, 0x80}, 0, 0, 0, 0}},
		{"RD", args{[2]byte{0x12, 0x34}, false, true}, &Header{[2]byte{0x12, 0x34}, [2]byte{0x81, 0x80}, 0, 0, 0, 0}},
		{"AA and RD", args{[2]byte{0x12, 0x34}, true, true}, &Header{[2]byte{0x12, 0x34}, [2]byte{0x85, 0x80}, 0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAnswerHeader(tt.args.id, tt.args.aa, tt.args.rd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAnswerHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
