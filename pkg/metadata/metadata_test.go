package metadata

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/types/known/structpb"
)

func TestMetadata_ToStructPB(t *testing.T) {
	tests := []struct {
		name    string
		m       Metadata
		want    *structpb.Struct
		wantErr bool
	}{
		{
			name: "should return metadata map to pb",
			m: Metadata{
				"k1": "v1",
				"k2": 2,
			},
			want: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"k1": structpb.NewStringValue("v1"),
					"k2": structpb.NewNumberValue(2),
				},
			},
		},
		{
			name: "should return error if metadata map is not parsable to pb",
			m: Metadata{
				"k1": func(x chan struct{}) {
					<-x
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.ToStructPB()
			if (err != nil) != tt.wantErr {
				t.Errorf("Metadata.ToStructPB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Metadata.ToStructPB() = %v, want %v", got, tt.want)
			}
		})
	}
}
