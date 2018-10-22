package label

import (
	"reflect"
	"testing"
)

func TestLabel_String(t *testing.T) {
	tests := []struct {
		name string
		l    Label
		want string
	}{
		{"Normal", Label{"test", "example", "com"}, "test.example.com."},
		{"Root", Label{""}, "."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("Label.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLabelsFromMessage(t *testing.T) {
	type args struct {
		data  []byte
		start int
	}
	tests := []struct {
		name    string
		args    args
		want    Label
		want1   int
		wantErr bool
	}{
		{"Normal", args{[]byte{0x07, 0x65, 0x78, 0x61, 0x6D, 0x70, 0x6C, 0x65, 0x03, 0x63, 0x6F, 0x6D, 0x00}, 0}, Label{"example", "com"}, 13, false},
		{"Redirect", args{[]byte{0x07, 0x65, 0x78, 0x61, 0x6D, 0x70, 0x6C, 0x65, 0x03, 0x63, 0x6F, 0x6D, 0x00, 0xC0, 0x00}, 13}, Label{"example", "com"}, 15, false},
		{"Redirect outside message bounds", args{[]byte{0xC0, 0xFF}, 0}, nil, 0, true},
		{"Start Outside Message Range", args{[]byte{0x0, 0x0}, 12}, nil, 0, true},
		{"Label length too short for redirect", args{[]byte{0xC0}, 0}, nil, 0, true},
		{"Label length > 63", args{[]byte{0x40}, 0}, nil, 0, true},
		{"Label length exceeds message length", args{[]byte{0x20}, 0}, nil, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetLabelsFromMessage(tt.args.data, tt.args.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLabelsFromMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLabelsFromMessage() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetLabelsFromMessage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLabel_Encode(t *testing.T) {
	tests := []struct {
		name string
		l    Label
		want []byte
	}{
		{"Normal", Label{"example", "com"}, []byte{0x07, 0x65, 0x78, 0x61, 0x6D, 0x70, 0x6C, 0x65, 0x03, 0x63, 0x6F, 0x6D, 0x00}},
		{"Empty", Label{}, []byte{0x0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Label.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Label
		wantErr bool
	}{
		{"Normal", "example.com", Label{"example", "com"}, false},
		{"With Root", "example.com.", Label{"example", "com"}, false},
		{"Null", string(0x0), nil, true},
		{"Empty", "", nil, true},
		{"Invalid characters", "²³.test.example.com", nil, true},
		{"Empty segment", "test..example.com", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
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
