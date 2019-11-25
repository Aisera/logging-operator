/*
 * Copyright © 2019 Banzai Cloud
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package plugins

// ForwardOutput plugin name
const ForwardOutput = "forward"

// ForwardOutputDefaultValues default values for the plugin
var ForwardOutputDefaultValues = map[string]string{
	"name":               "target",
	"bufferPath":         "/buffers/forward",
	"chunkLimit":         "2M",
	"queueLimit":         "8",
	"timekey":            "1h",
	"timekey_wait":       "10m",
	"timekey_use_utc":    "true",
	"retry_max_interval": "30",
	"flush_interval":     "5s",
	"flush_thread_count": "2",
	"retry_forever":      "true",
	"forward_tag":        "", // tag messages before sending them with forward.XYZ.
	"tlsSharedKey":       "", // enables tls and must match with the shared key on the remote side
	"tlsCACertFile":      "/fluentd/tls/caCert",
	"tlsCertFile":        "/fluentd/tls/clientCert",
	"tlsKeyFile":         "/fluentd/tls/clientKey",
	"clientHostname":     "fluentd.client", // this must be different from the hostname on the remote side
}

// ForwardOutputTemplate for the ForwardOutput plugin
const ForwardOutputTemplate = `
{{ if .forward_tag -}}
<filter {{ .pattern }}.** >
  @type record_transformer
  <record>
    forward_tag "{{ .forward_tag }}"
    original_tag ${tag}
  </record>
</filter>

<match {{ .pattern }}.** >
  @type rewrite_tag_filter
  <rule>
    key forward_tag
    pattern ^(.+)$
    tag forward.$1.${tag}
  </rule>
</match>

<match forward.{{ .forward_tag }}.{{ .pattern }}.** >
{{ else -}}
<match {{ .pattern }}.** >
{{ end -}}
  @type forward

  {{ if not (eq .tlsSharedKey "") -}}
  transport tls
  tls_version TLSv1_2
  tls_cert_path                {{ .tlsCACertFile }}
  tls_client_cert_path         {{ .tlsCertFile }}
  tls_client_private_key_path  {{ .tlsKeyFile }}
  <security>
    self_hostname           {{ .clientHostname }}
    shared_key              {{ .tlsSharedKey }}
  </security>
  {{ end -}}

  <server>
    name {{ .name }}
    host {{ .host }}
    port {{ .port }}
  </server>

  <buffer tag, time>
    @type file
    path {{ .bufferPath }}
    timekey {{ .timekey }}
    timekey_wait {{ .timekey_wait }}
    timekey_use_utc {{ .timekey_use_utc }}
    flush_mode interval
    retry_type exponential_backoff
    flush_thread_count {{ .flush_thread_count }}
    flush_interval {{ .flush_interval }}
    retry_forever {{ .retry_forever }}
    retry_max_interval {{ .retry_max_interval }}
    chunk_limit_size {{ .chunkLimit }}
    queue_limit_length {{ .queueLimit }}
    overflow_action block
  </buffer>
</match>`