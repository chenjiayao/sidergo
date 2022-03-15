package cluster

import (
	"bytes"
	"testing"
)

func Test_client_parseMulti(t *testing.T) {
	c := makeClient("localhost:3101")

	cmd := "3\r\nget\r\n" // 注意这里传递进去的是没有 $
	want := "$3\r\nget\r\n"

	var buf bytes.Buffer
	buf.WriteString(cmd)

	got, _ := c.parseMulti(&buf)
	gotString := string(got)

	if want != gotString {
		t.Errorf("test failed, got: %s", string(got))
	}
}

func Test_client_parseNumber(t *testing.T) {
	c := makeClient("localhost:3101")

	cmd := "10\r\n" // 注意这里传递进去的是没有 :
	want := ":10\r\n"

	var buf bytes.Buffer
	buf.WriteString(cmd)

	got, _ := c.parseNumber(&buf)
	gotString := string(got)

	if want != gotString {
		t.Errorf("test failed, got: %s", string(got))
	}
}

func Test_client_parseError(t *testing.T) {
	c := makeClient("localhost:3101")

	cmd := "ERR value is not a valid float\r\n" //
	want := "-ERR value is not a valid float\r\n"

	var buf bytes.Buffer
	buf.WriteString(cmd)

	got, _ := c.parseError(&buf)
	gotString := string(got)

	if want != gotString {
		t.Errorf("test failed, got: %s", string(got))
	}
}

func Test_client_parseArray(t *testing.T) {
	c := makeClient("localhost:3101")

	cmd := "3\r\n$3\r\nget\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	want := "*3\r\n$3\r\nget\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"

	var buf bytes.Buffer
	buf.WriteString(cmd)

	got, _ := c.parseArray(&buf)
	gotString := string(got)

	if want != gotString {
		t.Errorf("test failed, got: %s", string(got))
	}
}
