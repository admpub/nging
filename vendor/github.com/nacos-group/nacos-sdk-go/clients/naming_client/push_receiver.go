/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package naming_client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/nacos-group/nacos-sdk-go/util"
)

type PushReceiver struct {
	port        int
	host        string
	hostReactor *HostReactor
}

type PushData struct {
	PushType    string `json:"type"`
	Data        string `json:"data"`
	LastRefTime int64  `json:"lastRefTime"`
}

var (
	GZIP_MAGIC = []byte("\x1F\x8B")
)

func NewPushReceiver(hostReactor *HostReactor) *PushReceiver {
	pr := PushReceiver{
		hostReactor: hostReactor,
	}
	pr.startServer()
	return &pr
}

func (us *PushReceiver) tryListen() (*net.UDPConn, bool) {
	addr, err := net.ResolveUDPAddr("udp", us.host+":"+strconv.Itoa(us.port))
	if err != nil {
		logger.Errorf("can't resolve address,err: %+v", err)
		return nil, false
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logger.Errorf("error listening %s:%d,err:%+v", us.host, us.port, err)
		return nil, false
	}

	return conn, true
}

func (us *PushReceiver) startServer() {
	var conn *net.UDPConn

	for i := 0; i < 3; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		port := r.Intn(1000) + 54951
		us.port = port
		conn1, ok := us.tryListen()

		if ok {
			conn = conn1
			logger.Infof("udp server start, port: " + strconv.Itoa(port))
			break
		}

		if !ok && i == 2 {
			logger.Errorf("failed to start udp server after trying 3 times.")
		}
	}

	go func() {
		defer conn.Close()
		for {
			us.handleClient(conn)
		}
	}()
}

func (us *PushReceiver) handleClient(conn *net.UDPConn) {
	data := make([]byte, 4024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		logger.Errorf("failed to read UDP msg because of %+v", err)
		return
	}

	s := TryDecompressData(data[:n])
	logger.Info("receive push: "+s+" from: ", remoteAddr)

	var pushData PushData
	err1 := json.Unmarshal([]byte(s), &pushData)
	if err1 != nil {
		logger.Infof("failed to process push data.err:%+v", err1)
		return
	}
	ack := make(map[string]string)

	if pushData.PushType == "dom" || pushData.PushType == "service" {
		us.hostReactor.ProcessServiceJson(pushData.Data)

		ack["type"] = "push-ack"
		ack["lastRefTime"] = strconv.FormatInt(pushData.LastRefTime, 10)
		ack["data"] = ""

	} else if pushData.PushType == "dump" {
		ack["type"] = "dump-ack"
		ack["lastRefTime"] = strconv.FormatInt(pushData.LastRefTime, 10)
		ack["data"] = util.ToJsonString(us.hostReactor.serviceInfoMap)
	} else {
		ack["type"] = "unknow-ack"
		ack["lastRefTime"] = strconv.FormatInt(pushData.LastRefTime, 10)
		ack["data"] = ""
	}

	bs, _ := json.Marshal(ack)
	c, err := conn.WriteToUDP(bs, remoteAddr)
	if err != nil {
		logger.Errorf("WriteToUDP failed,return:%d,err:%+v", c, err)
	}
}

func TryDecompressData(data []byte) string {

	if !IsGzipFile(data) {
		return string(data)
	}
	reader, err := gzip.NewReader(bytes.NewReader(data))

	if err != nil {
		logger.Errorf("failed to decompress gzip data,err:%+v", err)
		return ""
	}

	defer reader.Close()
	bs, err := ioutil.ReadAll(reader)

	if err != nil {
		logger.Errorf("failed to decompress gzip data,err:%+v", err)
		return ""
	}

	return string(bs)
}

func IsGzipFile(data []byte) bool {
	if len(data) < 2 {
		return false
	}

	return bytes.HasPrefix(data, GZIP_MAGIC)
}
