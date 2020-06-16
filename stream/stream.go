package stream

import (
	"context"
	"fmt"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/codec/h264parser"
	"github.com/deepch/vdk/format/rtsp"
	"sync"
	"time"
)

var streamMap = make(map[string]*Stream)

type Stream struct {
	sync.Mutex

	Url    string
	ctx    context.Context
	cancel context.CancelFunc

	transChannels []TransChannel
}

type TransChannel interface {
	ID() string
	Trans(pck av.Packet, sps, pps []byte)
	Close()
}

func AttachStream(url string, c TransChannel) {
	var stream *Stream
	stream, ok := streamMap[url]
	if !ok {
		stream = &Stream{Url: url}
		stream.ctx, stream.cancel = context.WithCancel(context.Background())
		streamMap[url] = stream
		stream.open()
	}
	stream.Lock()
	stream.transChannels = append(stream.transChannels, c)
	stream.Unlock()
	fmt.Println("新增转发信道，当前信道数：", len(stream.transChannels))
}
func DettachStream(url string, c TransChannel) {
	var stream *Stream
	stream, ok := streamMap[url]
	if !ok {
		return
	}
	stream.Lock()
	var after []TransChannel
	for _, lc := range stream.transChannels {
		if lc.ID() != c.ID() {
			after = append(after, lc)
		}
	}
	stream.transChannels = after
	fmt.Println("移除转发信道，当前信道数：", len(stream.transChannels))
	if len(stream.transChannels) == 0 {
		fmt.Println("当前信道数为0，关闭输入流。")
		stream.Close()
	}
	stream.Unlock()
}

func (s *Stream) open() {
	var session *rtsp.Client
	var retryCount int
	for {
		var err error
		session, err = rtsp.DialTimeout(s.Url, 10*time.Second)
		if err != nil {
			fmt.Println("rtsp连接失败，稍后重试：", err)
			time.Sleep(5 * time.Second)
			retryCount++
			if retryCount == 3 {
				fmt.Println("重试3次未成功，取消连接")
				s.Close()
				return
			}
			continue
		}
		break
	}
	session.RtpKeepAliveTimeout = 10 * time.Second
	codec, err := session.Streams()
	if err != nil {
		fmt.Println("rtsp获取码流失败，关闭rtsp：", err)
		s.Close()
		return
	}
	sps := codec[0].(h264parser.CodecData).SPS()
	pps := codec[0].(h264parser.CodecData).PPS()
	fmt.Println(sps, pps)
	go func() {
	LOOP:
		for {
			select {
			case <-s.ctx.Done():
				break LOOP
			default:
			}
			pkt, err := session.ReadPacket()
			if err != nil {
				break
			}
			var transChannels []TransChannel
			s.Lock()
			transChannels = s.transChannels
			s.Unlock()
			for _, s := range transChannels {
				s.Trans(pkt, sps, pps)
			}
			//fmt.Println(pkt.Idx, ":", len(pkt.Data))
		}
		err = session.Close()
		if err != nil {
			fmt.Println("RTSP关流失败：", s.Url)
		} else {
			fmt.Println("RTSP关流成功：", s.Url)
		}
	}()
}

func (s *Stream) Close() {
	s.cancel()
	for _,t := range s.transChannels{
		t.Close()
	}
	delete(streamMap, s.Url)
}
