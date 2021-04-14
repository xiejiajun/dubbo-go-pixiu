/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package remote

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/client/dubbo"
	clienthttp "github.com/apache/dubbo-go-pixiu/pkg/client/http"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// nolint
func Init() {
	extension.SetFilterFunc(constant.RemoteCallFilter, remoteFilterFunc())
}

func remoteFilterFunc() fc.FilterFunc {
	return New(defaultNewParams()).Do()
}

func defaultNewParams() mockLevel {
	mock := 1
	mockStr := os.Getenv(constant.EnvMock)
	if len(mockStr) > 0 {
		i, err := strconv.Atoi(mockStr)
		if err == nil {
			mock = i
		}
	}

	return mockLevel(mock)
}

type mockLevel int8

const (
	open = iota
	close
	all
)

// clientFilter is a filter for recover.
type clientFilter struct {
	level mockLevel
}

// New create timeout filter.
func New(level mockLevel) filter.Filter {
	if level < 0 || level > 2 {
		level = close
	}
	return &clientFilter{
		level: level,
	}
}

// Do execute clientFilter filter logic
// support: 1 http 2 dubbo 2 http 2 http
func (f clientFilter) Do() fc.FilterFunc {
	return func(c fc.Context) {
		f.doRemoteCall(c.(*contexthttp.HttpContext))
	}
}

func (f clientFilter) doRemoteCall(c *contexthttp.HttpContext) {
	api := c.GetAPI()

	if (f.level == open && api.Mock) || (f.level == all) {
		c.SourceResp = &filter.ErrResponse{
			Message: "mock success",
		}
		c.Next()
		return
	}

	typ := api.Method.IntegrationRequest.RequestType

	// TODO 根据协议类型获取对应的客户端
	cli, err := matchClient(typ)
	if err != nil {
		c.Err = err
		return
	}

	// TODO 通过范化客户端调用目标Dubbo服务的dubbo接口 / 通过http客户端调用目标Http服务的http接口
	resp, err := cli.Call(client.NewReq(c.Ctx, c.Request, *api))
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] client call err:%v!", err)
		c.Err = err
		return
	}

	logger.Debugf("[dubbo-go-pixiu] client call resp:%v", resp)

	c.SourceResp = resp
	// response write in response filter.
	c.Next()
}

func matchClient(typ config.RequestType) (client.Client, error) {
	switch strings.ToLower(string(typ)) {
	case string(config.DubboRequest):
		// TODO Dubbo协议用dubbo Client处理，这里获取到的是PX.beforeStart初始化好的Client
		return dubbo.SingletonDubboClient(), nil
	case string(config.HTTPRequest):
		// TODO http协议使用HttpClient处理, http协议的Client在这里初始化，也是单例的，这里的Http协议用于调用Http接口，
		//  不能调用dubbo 2.7.x及以下版本的dubbo接口(dubbo 3.x兼容Http2, 是否可以通过Http协议调用Dubbo3.x待观望)
		// TODO 这里也说明了dubbo-go-pixiu不止是一个dubbo接口网关，还可以做Http网关
		return clienthttp.SingletonHTTPClient(), nil
	default:
		return nil, errors.New("not support")
	}
}
