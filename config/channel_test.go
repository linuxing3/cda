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
			name: "测试多类型通道",
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
	// 测试无缓冲通道
	values := make(chan int)
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "测试Int无缓冲通道",
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
	// 测试有缓冲通道
	bufferedValues := make(chan int, 2)
	tests = []struct {
		name string
		args args
		want int
	}{
		{
			name: "测试Int有缓冲通道",
			args: args{bufferedValues},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go SendInt(tt.args.c, 1)
			go SendInt(tt.args.c, 1)
			got1 := <-bufferedValues
			got2 := <-bufferedValues
			if !reflect.DeepEqual(got1, tt.want) {
				t.Errorf("SendValues() got = %v, want %v", got1, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want) {
				t.Errorf("SendValues() got = %v, want %v", got2, tt.want)
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
			name: "测试String无缓冲通道",
			args: args{value},
			want: "go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go SendString(tt.args.c, "go")
			got := <-value // FIXED 这里不能用tt.args.c
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Got %v, but want %v", got, tt.want)
			}
		})
	}
}

// TODO: 无法测试多重通道的异步
func TestTransformStringToStringArray(t *testing.T) {
	type args struct {
		worklist    chan []string
		unseenLinks chan string
	}

	workList := make(chan []string)

	unseenList := make(chan string)
	unseenList <- "link3"


	tests := []struct {
		name string
		args args
	}{
		{
			name: "将字符串改成字符串数组",
			args: args{
				worklist: workList,
				unseenLinks: unseenList,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go TransformStringToStringArray(tt.args.worklist, tt.args.unseenLinks)
			// after transformation
			link3 := <- workList

			reflect.DeepEqual(link3, []string{"link3"})
		})
	}
}
