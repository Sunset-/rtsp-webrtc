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

var streamMapLock sync.Mutex
var streamMap = make(map[string]*Stream)

type Stream struct {
	sync.Mutex

	Codecs []av.CodecData

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

func GetStream(url string) (s *Stream, err error) {
	var stream *Stream
	var ok bool
	streamMapLock.Lock()
	stream, ok = streamMap[url]
	streamMapLock.Unlock()
	if ok {
		return stream, nil
	}
	stream = &Stream{Url: url}
	stream.ctx, stream.cancel = context.WithCancel(context.Background())
	err = stream.open()
	if err != nil {
		return nil, err
	}
	streamMapLock.Lock()
	cStream, ok := streamMap[url]
	if ok {
		return cStream, nil
	}
	streamMap[url] = stream
	streamMapLock.Unlock()
	return stream, nil
}

func AttachStream(stream *Stream, c TransChannel) {
	stream.Lock()
	stream.transChannels = append(stream.transChannels, c)
	stream.Unlock()
	fmt.Println("新增转发信道，当前信道数：", len(stream.transChannels))
}
func DettachStream(stream *Stream, c TransChannel) {
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

func (s *Stream) open() error {
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
				return err
			}
			continue
		}
		break
	}
	session.RtpKeepAliveTimeout = 10 * time.Second
	codec, err := session.Streams()
	s.Codecs = codec
	if err != nil {
		fmt.Println("rtsp获取码流失败，关闭rtsp：", err)
		s.Close()
		return err
	}
	fmt.Println("RTSP开流成功：", s.Url)
	sps := codec[0].(h264parser.CodecData).SPS()
	pps := codec[0].(h264parser.CodecData).PPS()
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
	return nil
}

func (s *Stream) Close() {
	s.cancel()
	for _, t := range s.transChannels {
		t.Close()
	}
	streamMapLock.Lock()
	delete(streamMap, s.Url)
	streamMapLock.Unlock()
}
