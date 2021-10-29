/**
 * @Time: 2019-06-29 09:27
 * @Author: solacowa@gmail.com
 * @File: prometheus
 * @Software: GoLand
 */

package public

import (
	"encoding/json"
	"strings"
)

type PrometheusAlerts interface {
	Get() (prom Prom)
	String() string
	Alerts() []Alert
	GetName(label *Label) string
	GetNamespace(label *Label) string
	GetAlertName() string
	GetDesc() string
	From(alert *Alert) string
	To(alert *Alert) string
}

type prometheusAlerts struct {
	body []byte
	prom Prom
}

type Prom struct {
	Receiver    string  `json:"receiver"`
	Status      string  `json:"status"`
	Alerts      []Alert `json:"alerts"`
	GroupLabels struct {
		AlertName string `json:"alertname"`
	} `json:"groupLabels"`
	CommonLabels      Label `json:"commonLabels"`
	CommonAnnotations struct {
		Summary     string `json:"summary"`
		Description string `json:"description"`
	} `json:"common_annotations"`
	Labels      string `json:"labels"`
	ExternalURL string `json:"externalURL"`
	StartsAt    string `json:"startsAt"`
	EndsAt      string `json:"endsAt"`
	GroupKey    string `json:"group_key"`
	Version     string `json:"version"`
}

type Alert struct {
	Status       string      `json:"status"`
	Labels       Label       `json:"labels"`
	Annotations  Annotations `json:"annotations"`
	StartsAt     string      `json:"startsAt"`
	EndsAt       string      `json:"endsAt"`
	GeneratorURL string      `json:"generatorURL"`
}

type Label struct {
	AlertName           string `json:"alertname"`
	Container           string `json:"container"`
	ContainerName       string `json:"container_name"`
	Deployment          string `json:"deployment"`
	Instance            string `json:"instance"`
	Job                 string `json:"job"`
	K8sApp              string `json:"k8s_app"`
	KubernetesName      string `json:"kubernetes_name"`
	KubernetesNamespace string `json:"kubernetes_namespace"`
	Namespace           string `json:"namespace"`
	Pod                 string `json:"pod"`
	Severity            string `json:"severity"`
	SearchID            string `json:"searchID"`
	DestinationService  string `json:"destination_service"`
	SourceService       string `json:"source_service"`
}

type Annotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
	From        string `json:"from"`
	To          string `json:"to"`
}

func NewPrometheusAlerts(body []byte) (p PrometheusAlerts, err error) {
	var prom Prom
	err = json.Unmarshal(body, &prom)
	if err != nil {
		return
	}
	return &prometheusAlerts{body, prom}, nil
}

func (c *prometheusAlerts) Get() (prom Prom) {
	return c.prom
}

func (c *prometheusAlerts) String() string {
	return string(c.body)
}

func (c *prometheusAlerts) Alerts() []Alert {
	return c.prom.Alerts
}

func (c *prometheusAlerts) GetAlertName() string {
	return c.prom.GroupLabels.AlertName
}

func (c *prometheusAlerts) GetNamespace(label *Label) string {
	if label.Namespace != "" {
		return label.Namespace
	}
	if label.KubernetesNamespace != "" {
		return label.KubernetesNamespace
	}
	if label.DestinationService != "" {
		destinationServices := strings.Split(label.DestinationService, ".")
		return destinationServices[1]
	}
	return ""
}

func (c *prometheusAlerts) GetName(label *Label) string {
	if label.Deployment != "" {
		return label.Deployment
	}
	if label.Container != "" {
		return label.Container
	}
	if label.ContainerName != "" {
		return label.ContainerName
	}

	if label.DestinationService != "" {
		destinationServices := strings.Split(label.DestinationService, ".")
		return destinationServices[0]
	}

	return ""
}

func (c *prometheusAlerts) GetDesc() string {
	var desc string
	alerts := c.prom.Alerts

	if len(alerts) > 0 {
		for _, v := range alerts {
			desc = v.Annotations.Description
			if desc != "" {
				return desc //只返回一个desc即可，太长徽信模板消息也展示不了
			}
		}
	}

	if c.prom.CommonAnnotations.Description != "" {
		desc = c.prom.CommonAnnotations.Description
		return desc
	}

	desc = "no capture description"
	return desc
}

func (c *prometheusAlerts) From(alert *Alert) string {
	return alert.Annotations.From
}

func (c *prometheusAlerts) To(alert *Alert) string {
	return alert.Annotations.To
}
