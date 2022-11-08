package settings

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		desc  string
		ping  int
		hours int
		urls  []string
		err   bool
	}{
		{desc: "bad monitoringEveryHours 0", ping: 1, hours: 0, urls: nil, err: true},
		{desc: "bad monitoringEveryHours -1", ping: 1, hours: -1, urls: nil, err: true},
		{desc: "bad pingingEveryHours 0", ping: 0, hours: 1, urls: nil, err: true},
		{desc: "bad pingingEveryHours -1", ping: -1, hours: 1, urls: nil, err: true},
		{desc: "bad serverAddresses nil", ping: 1, hours: 1, urls: nil, err: true},
		{desc: "bad serverAddresses empty", ping: 1, hours: 1, urls: []string{}, err: true},
		{desc: "bad serverAddresses one url", ping: 1, hours: 1, urls: []string{"http://google.ru"}, err: true},
		{desc: "bad serverAddresses bad urls", ping: 1, hours: 1, urls: []string{"jsdj dk", "4jk4l3 3#"}, err: true},
		{desc: "good", ping: 1, hours: 1, urls: []string{"https://aaa.bbb", "http://ttt.rrr"}, err: false},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := New(tC.ping, tC.hours, tC.urls)
			if tC.err && err == nil {
				t.Error("want err!=nil but got err==nil")
			}
			if !tC.err && err != nil {
				t.Errorf("want err==nil but got err=[%s]", err.Error())
			}
		})
	}
}

func TestGetMonitoringDuration(t *testing.T) {
	s, _ := New(1, 1, []string{"https://aaa.bbb", "http://ttt.rrr"})

	if s.GetMonitoringDuration() != time.Hour {
		t.Error("GetMonitoringDuration works bad 1")
	}

	s, _ = New(1, 13, []string{"https://aaa.bbb", "http://ttt.rrr"})

	if s.GetMonitoringDuration() != 13*time.Hour {
		t.Error("GetMonitoringDuration works bad 13")
	}
}

func TestGetServerAddress(t *testing.T) {
	urlStr1 := "https://aaa.bbb"
	urlStr2 := "http://ttt.rrr"

	s, _ := New(1, 1, []string{urlStr1, urlStr2})
	url1 := s.GetServerAddress()
	if url1.String() != urlStr1 {
		t.Errorf("want [%s] but got [%s]", urlStr1, url1.String())
	}
}

func TestRotateServerAddresses(t *testing.T) {
	urlStr1 := "https://aaa.bbb"
	urlStr2 := "http://ttt.rrr"

	s, _ := New(1, 1, []string{urlStr1, urlStr2})
	s.RotateServerAddresses()
	url1 := s.GetServerAddress()
	if url1.String() != urlStr2 {
		t.Errorf("want second url [%s] but got [%s]", urlStr2, url1.String())
	}

	s.RotateServerAddresses()
	url1 = s.GetServerAddress()
	if url1.String() != urlStr1 {
		t.Errorf("want first url [%s] but got [%s]", urlStr1, url1.String())
	}
}

func TestNeedUpdate(t *testing.T) {
	s, _ := New(1, 1, []string{"https://aaa.bbb", "http://ttt.rrr"})
	n, _ := New(1, 2, []string{"https://aaa.bbb", "http://ttt.rrr"})

	if s.NeedUpdate(s) {
		t.Error("want not NeedUpdate but got needs 1")
	}

	if !s.NeedUpdate(n) {
		t.Error("want NeedUpdate but got not need 2")
	}

	s, _ = New(1, 1, []string{"https://aaa.bbb", "http://ttt.rrr"})
	n, _ = New(1, 1, []string{"http://aaa.bbb", "http://ttt.rrr"})

	if !s.NeedUpdate(n) {
		t.Error("want NeedUpdate but got not need (url)")
	}

	s, _ = New(1, 1, []string{"https://aaa.bbb", "http://ttt.rrr", "http://www.eee"})
	n, _ = New(1, 1, []string{"https://aaa.bbb", "http://ttt.rrr"})

	if !s.NeedUpdate(n) {
		t.Error("want NeedUpdate but got not need (url num)")
	}
}
