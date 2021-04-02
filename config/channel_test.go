package channel

import (
	"reflect"
	"testing"
)

func Test_Crawl(t *testing.T) {
	// types
	type args struct {
		url string
	}
	type Test struct {
		name string
		args args
		want []string
	}
	// init test case
	tests := []Test{
		{
			name: "channel",
			args: args{url: "localhost"},
			want: []string{"localhost"},
		},
	}
	// running test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Crawl(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("crawl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSendInt(t *testing.T) {
	type args struct {
		c chan int
	}
	values := make(chan int)
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "value",
			args: args{values},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go SendInt(tt.args.c, 1)
			got := <-values
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendValues() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSendString(t *testing.T) {
	type args struct {
		c chan string
	}
	value := make(chan string)
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "send string",
			args: args{value},
			want: "go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go SendString(tt.args.c, "go")
			got := <-value
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Got %v, but want %v", got, tt.want)
			}
		})
	}
}
