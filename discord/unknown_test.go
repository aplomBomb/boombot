package discord

import (
	"reflect"
	"testing"

	"github.com/andersfylling/disgord"
)

func Test_unknown(t *testing.T) {
	type args struct {
		data *disgord.Message
	}
	tests := []struct {
		name    string
		args    args
		want    *disgord.Message
		wantErr bool
	}{
		{
			name: "unknown | create message success",
			args: args{
				
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unknown(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("unknown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unknown() = %v, want %v", got, tt.want)
			}
		})
	}
}
