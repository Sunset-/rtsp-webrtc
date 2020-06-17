package main

import (
	"encoding/base64"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"webrtc/rtc"
	"webrtc/stream"
)

type Stream struct {
}

func main() {
	go serveHTTP()
	//device, err := goonvif.NewDevice("172.16.133.207:80")//门口相机
	//device, err := goonvif.NewDevice("172.16.133.159:80")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//device.Authenticate("admin", "DFwl123456")

	//capabilities := Device.GetCapabilities{Category:"Media"}
	//获取通道列表
	//res, err := device.CallMethod(Media.GetProfiles{})
	//if err != nil {
	//	return
	//}
	//gp := &Media.GetProfilesResponse{}
	//err = Decode(res, gp)
	//if err != nil {
	//	return
	//}
	//if len(gp.Profiles) > 0 {
	//	for _, p := range gp.Profiles {
	//		fmt.Println("通道：", p.Name, ",", p.Token)
	//	}
	//	res, err := device.CallMethod(Media.GetStreamUri{
	//		ProfileToken: gp.Profiles[0].Token,
	//	})
	//	if err != nil {
	//		return
	//	}
	//	gp := &Media.GetStreamUriResponse{}
	//	err = Decode(res, gp)
	//	fmt.Println("开流地址：", gp.MediaUri.Uri)
	//	//openStream(prependUsername(string(gp.MediaUri.Uri),"admin", "DFwl123456"))
	//}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	select {
	case <-c:
		break
	}
}

func prependUsername(uri, username, password string) string {
	index := strings.Index(uri, "//")
	return uri[:index+2] + username + ":" + password + "@" + uri[index+2:]
}

type Envelope struct {
	XMLName      xml.Name
	EnvelopeBody EnvelopeBody `xml:"Body"`
}

type EnvelopeBody struct {
	XMLName  xml.Name
	Response []byte `xml:",innerxml"`
}

func Decode(res *http.Response, ptr interface{}) error {
	evp := &Envelope{
		EnvelopeBody: EnvelopeBody{},
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(data, evp)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(evp.EnvelopeBody.Response, ptr)
	if err != nil {
		return err
	}
	return nil
}

func serveHTTP() {
	router := gin.Default()
	router.Static("/", "./web")
	router.POST("/api/recive", reciver)
	err := router.Run(":8200")
	if err != nil {
		log.Fatalln("Start HTTP Server error", err)
	}
}

func reciver(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	data := c.PostForm("data")
	url := c.PostForm("url")

	//解析对端的SDP
	sdp, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Println("DecodeString error", err)
		return
	}

	s, err := stream.GetStream(url)
	if err != nil {
		log.Println("rtsp开流失败", err)
		return
	}
	rc, err := rtc.New(string(sdp), s.Codecs)
	if err != nil {
		_, err = c.Writer.Write([]byte(err.Error()))
	} else {
		_, err = c.Writer.Write([]byte(base64.StdEncoding.EncodeToString([]byte(rc.LocalSdp))))
	}

	rc.Link(func() {
		stream.AttachStream(s, rc)
	}, func() {
		stream.DettachStream(s, rc)
	})
}
