package discord

import (
	"reflect"
	"testing"

	"github.com/andersfylling/disgord"
)

func Test_unknown(t *testing.T) {
	type args struct {
		data    *disgord.MessageCreate
		message *disgord.Message
	}
	tests := []struct {
		name string
		args args
		want *disgord.Message
	}{
		{
			name: "unknown | success",
			args: args{
				data: disgord.MessageCreate{
					
				}
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unknown(tt.args.data, tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unknown() = %v, want %v", got, tt.want)
			}
		})
	}
}
