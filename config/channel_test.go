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

func TestCrawlWithChannel(t *testing.T) {
	type args struct {
		url string
	}

	want := make(chan []string)
	data := []string{"localhost"}
	want <- data

	want1 := make(chan string)
	want1 <- "localhost"

	tests := []struct {
		name  string
		args  args
		want  chan []string
		want1 chan string
	}{
		{
			name: "channel",
			args: args{ url: "localhost"},
			want: want,
			want1: want1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CrawlWithChannel(tt.args.url)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CrawlWithChannel() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CrawlWithChannel() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
