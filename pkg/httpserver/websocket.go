package httpserver

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"newproxy/pkg/capture"
	"strconv"
	"strings"

	"github.com/google/gopacket/layers"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type PacketSocket struct {
	Timestamp string `json:"timestamp"`
	Src       string `json:"src,omitempty"`
	Dst       string `json:"dst,omitempty"`
	Srcport   uint16 `json:"srcport,omitempty"`
	Dstport   uint16 `json:"dstport,omitempty"`
	Protocol  string `json:"protocol,omitempty"`
	Len       int    `json:"len,omitempty"`
	Message   string `json:"message,omitempty"`
	Code      uint16 `json:"httpcode,omitempty"`
}

func (s *KptureServer) hello(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		s.logger.Error(err)
		return &websocket.CloseError{Code: websocket.CloseProtocolError, Text: "uuid is required"}
	}
	s.logger.Info("websocket connected")
	m, _ := url.ParseQuery(c.Request().URL.RawQuery)
	uuid := m.Get("uuid")

	hm, err := s.getHostMap()
	if err != nil {
		s.logger.Errorf("error getting host map: %v", err)
		return &websocket.CloseError{Code: websocket.CloseProtocolError, Text: "uuid is required"}
	}

	fmt.Println(hm)

	s.logger.Info(fmt.Sprintf("Hello %s", uuid))
	if uuid == "" {
		return &websocket.CloseError{Code: websocket.CloseProtocolError, Text: "uuid is required"}
	}

	kpture := s.CaptureManager.GetCapture(uuid)

	if kpture == nil {
		return &websocket.CloseError{Code: websocket.CloseProtocolError, Text: "uuid is required"}
	}

	if kpture.Status.CaptureState != capture.CaptureStatusStarted {
		return &websocket.CloseError{Code: websocket.CloseProtocolError, Text: "uuid is required"}
	}

	s.logger.Info("hello", "uuid", uuid)

	defer ws.Close()
	for packet := range kpture.Weboscket() {
		pk := PacketSocket{}
		pk.Timestamp = packet.Metadata().Timestamp.String()

		if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)

			if dstName, ok := hm[ip.DstIP.String()]; ok {
				pk.Dst = dstName
			} else {
				pk.Dst = ip.DstIP.String()
			}

			if srcName, ok := hm[ip.SrcIP.String()]; ok {
				pk.Src = srcName
			} else {
				pk.Src = ip.SrcIP.String()
			}
			pk.Protocol = "Ipv4"
		}
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			pk.Dstport = uint16(tcp.DstPort)
			pk.Srcport = uint16(tcp.SrcPort)
			pk.Protocol = "TCP"
		}
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			udp, _ := udpLayer.(*layers.UDP)
			pk.Dstport = uint16(udp.DstPort)
			pk.Srcport = uint16(udp.SrcPort)
			pk.Protocol = "UDP"
			if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
				pk.Protocol = "DNS"
				dns, _ := dnsLayer.(*layers.DNS)
				if len(dns.Answers) > 0 {
					pk.Message += " " + dns.Answers[0].String()
				} else if len(dns.Questions) > 0 {
					pk.Message += " " + string(dns.Questions[0].Name)
				}
			}
		}
		if imcp := packet.Layer(layers.LayerTypeICMPv4); imcp != nil {
			pk.Protocol = "Icmpv4"
			imcppacket, _ := imcp.(*layers.ICMPv4)
			pk.Message = imcppacket.TypeCode.String()
		}

		applicationLayer := packet.ApplicationLayer()
		if applicationLayer != nil {
			if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
				pk.Protocol = "HTTP"
				payloadReader := bytes.NewReader(applicationLayer.Payload())
				bufferedPayloadReader := bufio.NewReader(payloadReader)
				request, _ := http.ReadRequest(bufferedPayloadReader)
				if request != nil {
					pk.Message = request.Method + "" + request.URL.String()
				} else {
					bytesReader := bytes.NewReader(applicationLayer.Payload())
					bufReader := bufio.NewReader(bytesReader)
					value1, _, _ := bufReader.ReadLine()
					if strings.Contains(string(value1), "HTTP") {
						pk.Message = string(value1)
						splited := strings.Split(string(value1), " ")
						code, _ := strconv.ParseUint(splited[1], 10, 16)
						pk.Code = uint16(code)
					}
				}

				// response, _ := http.ReadResponse(bufferedPayloadReader, nil)
				// if response != nil {
				// 	pk.Message = fmt.Sprint(response.Status)
				// 	// fmt.Println(response.StatusCode)
				// }
			}
		}
		pk.Len = len(packet.Data())

		if pk.Protocol != "" {
			packetRaw, _ := json.Marshal(pk)
			err := ws.WriteMessage(websocket.TextMessage, packetRaw)
			if err != nil {
				s.logger.Errorf("error write message to websocket: %v", err)
				return err
			}
		}

	}

	return nil
}
