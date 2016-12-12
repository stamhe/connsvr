package tests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/simplejia/connsvr"
	"github.com/simplejia/connsvr/comm"
	"github.com/simplejia/connsvr/conf"
	"github.com/simplejia/connsvr/proto"
	"github.com/simplejia/utils"

	"net"
	"testing"
)

func TestMsgsTcp(t *testing.T) {
	cmd := comm.PUSH
	rid := "r1"
	uid := "u_TestMsgsTcp"
	text := "hello world"
	msgId := ""

	conn, err := net.Dial(
		"udp",
		fmt.Sprintf("%s:%d", utils.LocalIp, conf.C.App.Bport),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	msg := proto.NewMsg(comm.UDP)
	msg.SetCmd(cmd)
	msg.SetRid(rid)
	msg.SetUid(uid)
	msg.SetBody(text)
	msg.SetExt(`{"msgid": "1"}`)
	data, ok := msg.Encode()
	if !ok {
		t.Fatal("msg.Encode() error")
	}

	_, err = conn.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond)

	conn, err = net.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", utils.LocalIp, conf.C.App.Tport),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	msg = proto.NewMsg(comm.TCP)
	msg.SetCmd(comm.MSGS)
	msg.SetUid("")
	msg.SetRid(rid)
	msg.SetBody(msgId)
	data, ok = msg.Encode()
	if !ok {
		t.Fatal("msg.Encode() error")
	}

	_, err = conn.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	_msg := new(proto.MsgTcp)
	ok = _msg.Decode(bufio.NewReader(conn))
	if !ok {
		t.Fatal("_msg.DecodeHeader() error")
	}

	if _msg.Cmd() == comm.ERR {
		t.Errorf("get: %v, expected: %v", _msg.Cmd(), msg.Cmd())
	}
	if _msg.Rid() != rid {
		t.Errorf("get: %s, expected: %s", _msg.Rid(), rid)
	}

	expect_body, _ := json.Marshal([]string{text})
	if body := _msg.Body(); body != string(expect_body) {
		t.Errorf("get: %s, expected: %s", _msg.Body(), expect_body)
	}
}