package rtc

import (
	"errors"
	"fmt"
	"github.com/deepch/vdk/av"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
	uuid "github.com/satori/go.uuid"
	"log"
	"math/rand"
	"strings"
	"time"
)

type RtcChannel struct{
	UUID string

	peerConnection *webrtc.PeerConnection
	state webrtc.ICEConnectionState

	onConnected func()
	onDisconnected func()

	RemoteSdp        string
	LocalSdp         string
	start            bool
	videoPayloadType uint8
	Vpre             time.Duration

	videoTrack *webrtc.Track

	Codecs []av.CodecData
}

func New(sdp string,codecs []av.CodecData)(rtc *RtcChannel,err error){
	rc := &RtcChannel{
		UUID : uuid.NewV4().String(),
		RemoteSdp:sdp,
		Codecs :codecs,
	}
	err = rc.initPeerConn()
	if err!=nil{
		return nil,err
	}
	return rc,nil
}

func (r *RtcChannel) ID()string{
	return r.UUID
}

func (r *RtcChannel) Close(){
	if r.peerConnection!=nil{
		_ = r.peerConnection.Close()
	}
}

func (r *RtcChannel) initPeerConn()error{
	//创建webetc
	mediaEngine := webrtc.MediaEngine{}
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  r.RemoteSdp,
	}
	err := mediaEngine.PopulateFromSDP(offer)
	if err != nil {
		log.Println("PopulateFromSDP error", err)
		return err
	}

	for _, videoCodec := range mediaEngine.GetCodecsByKind(webrtc.RTPCodecTypeVideo) {
		if videoCodec.Name == "H264" && strings.Contains(videoCodec.SDPFmtpLine, "packetization-mode=1") {
			r.videoPayloadType = videoCodec.PayloadType
			break
		}
	}
	if r.videoPayloadType == 0 {
		return errors.New("Remote peer does not support H264")
	}
	if r.videoPayloadType != 126 {
		log.Println("Video might not work with codec", r.videoPayloadType)
	}
	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))
	r.peerConnection, err = api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		return err
	}
	err = r.initTrack()
	if err!=nil{
		return err
	}

	//开始连接
	if err := r.peerConnection.SetRemoteDescription(offer); err != nil {
		log.Println("SetRemoteDescription error", err, offer.SDP)
		return err
	}
	answer, err := r.peerConnection.CreateAnswer(nil)
	if err != nil {
		log.Println("CreateAnswer error", err)
		return err
	}

	r.LocalSdp = answer.SDP
	if err = r.peerConnection.SetLocalDescription(answer); err != nil {
		log.Println("SetLocalDescription error", err)
		return err
	}
	//监听事件
	r.peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("Connection State has changed %s \n", connectionState.String())
		r.state = connectionState
		if connectionState != webrtc.ICEConnectionStateConnected {
			log.Println("Client Close Exit")
			err := r.peerConnection.Close()
			if err != nil {
				log.Println("peerConnection Close error", err)
			}
			if r.onDisconnected!=nil{
				r.onDisconnected()
			}
			return
		}
		if connectionState == webrtc.ICEConnectionStateConnected {
			if r.onConnected!=nil{
				r.onConnected()
			}
		}
	})
	return nil
}

func(r *RtcChannel) initTrack() error{
	//ADD KeepAlive Timer
	timer1 := time.NewTimer(time.Second * 2)
	r.peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			//fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
			timer1.Reset(2 * time.Second)
		})
	})

	videoTrack, err := r.peerConnection.NewTrack(r.videoPayloadType,rand.Uint32(),"video","_pion")
	if err != nil {
		log.Println("NewTrack error", err)
		return err
	}
	_, err = r.peerConnection.AddTransceiverFromTrack(videoTrack,
		webrtc.RtpTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionSendonly,
		},
	)
	if err != nil {
		log.Println("AddTransceiverFromTrack error", err)
		return err
	}
	r.videoTrack = videoTrack
	_, err = r.peerConnection.AddTrack(videoTrack)
	if err != nil {
		log.Println("AddTrack error", err)
		return err
	}

	//ADD Audio Track
	var audioTrack *webrtc.Track
	codecs := r.Codecs
	if len(codecs) > 1 && (codecs[1].Type() == av.PCM_ALAW || codecs[1].Type() == av.PCM_MULAW) {
		switch codecs[1].Type() {
		case av.PCM_ALAW:
			audioTrack, err = r.peerConnection.NewTrack(webrtc.DefaultPayloadTypePCMA, rand.Uint32(), "audio", "_audio")
		case av.PCM_MULAW:
			audioTrack, err = r.peerConnection.NewTrack(webrtc.DefaultPayloadTypePCMU, rand.Uint32(), "audio", "_audio")
		}
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = r.peerConnection.AddTransceiverFromTrack(audioTrack,
			webrtc.RtpTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionSendonly,
			},
		)
		if err != nil {
			log.Println("AddTransceiverFromTrack error", err)
			return err
		}
		_, err = r.peerConnection.AddTrack(audioTrack)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil

}


func (r *RtcChannel) Link(onConnected func(),onDisconnected func()){
	r.onConnected = onConnected
	r.onDisconnected = onDisconnected
}

func (r *RtcChannel) Trans(pck av.Packet,sps,pps []byte){
	if pck.IsKeyFrame {
		r.start = true
	}
	if !r.start {
		return
	}
	if pck.IsKeyFrame {
		pck.Data = append([]byte("\000\000\001"+string(sps)+"\000\000\001"+string(pps)+"\000\000\001"), pck.Data[4:]...)

	} else {
		pck.Data = pck.Data[4:]
	}
	var Vts time.Duration
	if pck.Idx == 0 {
		if r.Vpre != 0 {
			Vts = pck.Time - r.Vpre
		}
		samples := uint32(90000 / 1000 * Vts.Milliseconds())
		err := r.videoTrack.WriteSample(media.Sample{Data: pck.Data, Samples: samples})
		if err != nil {
			fmt.Println(err)
			return
		}
		r.Vpre = pck.Time
	}
}