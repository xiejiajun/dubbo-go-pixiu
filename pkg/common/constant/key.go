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

package constant

const (
	HTTPConnectManagerFilter = "dgp.filters.http_connect_manager"
	HTTPAuthorityFilter      = "dgp.filters.http.authority_filter"
	HTTPRouterFilter         = "dgp.filters.http.router"
	HTTPApiFilter            = "dgp.filters.http.api"
	HTTPDomainFilter         = "dgp.filters.http.domain"
	RemoteCallFilter         = "dgp.filters.remote_call"
	// TODO timeout控制Filter
	TimeoutFilter = "dgp.filters.timeout"
	// TODO 用于打印请求链路统计日志的Filter
	LoggerFilter = "dgp.filters.logger"
	// TODO 用于恢复panic的Filter
	RecoveryFilter = "dgp.filters.recovery"
	// TODO 用于将响应结果返回给调用方
	ResponseFilter = "dgp.filters.response"
	// TODO 用于打印访问日志的Filter（Nginx的accessLog也是这么打印的吧)
	AccessLogFilter = "dgp.filters.access_log"
)

const (
	LocalMemoryApiDiscoveryService = "api.ds.local_memory"
)
