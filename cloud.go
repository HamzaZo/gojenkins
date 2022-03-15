package gojenkins

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"strconv"
	"strings"
)

type KubernetesCloud struct {
	Raw      *[]CloudResponse
	K8sCloud CloudConfig
	Jenkins  *Jenkins
	Base     string
}

type CloudConfig struct {
	CloudName     string
	Namespace     string
	JenkinsURL    string
	JenkinsTunnel string
	Operation     string
}

type CloudResponse struct {
	Name                     string            `json:"name"`
	ServerUrl                string            `json:"serverUrl"`
	ServerCertificate        string            `json:"serverCertificate"`
	SkipTlsVerify            bool              `json:"skipTlsVerify"`
	JenkinsUrl               string            `json:"jenkinsUrl"`
	JenkinsTunnel            string            `json:"jenkinsTunnel"`
	Namespace                string            `json:"namespace"`
	AddMasterProxyEnvVars    bool              `json:"addMasterProxyEnvVars"`
	CapOnlyOnAlivePods       bool              `json:"capOnlyOnAlivePods"`
	UseJenkinsProxy          bool              `json:"useJenkinsProxy"`
	Labels                   map[string]string `json:"labels"`
	UsageRestricted          bool              `json:"usageRestricted"`
	PodRetention             string            `json:"podRetention"`
	ContainerCap             int               `json:"containerCap"`
	CredentialsID            string            `json:"credentialsID"`
	WaitForPodSec            int               `json:"waitForPodSec"`
	DirectConnection         bool              `json:"directConnection"`
	ReadTimeout              int               `json:"readTimeout"`
	MaxRequestsPerHost       int               `json:"maxRequestsPerHost"`
	DefaultsProviderTemplate string            `json:"defaultsProviderTemplate"`
	Templates                []string          `json:"templates"`
	PodLabels                map[string]string `json:"podLabels"`
	RetentionTimeout         int               `json:"retentionTimeout"`
	ConnectionTimeout        int               `json:"connectionTimeout"`
	WebSocket                bool              `json:"webSocket"`
}

var (
	temp        *template.Template
	tpl         bytes.Buffer
	cloudConfig = []byte(`
import org.csanchez.jenkins.plugins.kubernetes.*
import org.csanchez.jenkins.plugins.kubernetes.model.*
import jenkins.model.Jenkins
import groovy.json.JsonOutput

{{ if eq .Operation "create" }}
def addKubernetesCloud(cloudList, config) {
    def cloud = new KubernetesCloud(
            cloudName = config.cloudName ?: 'Kubernetes'
    )
    cloud.serverCertificate = config.serverCertificate ?: ''
    cloud.skipTlsVerify = config.skipTlsVerify ?: false
    cloud.credentialsId = config.credentialsId ?: ''
    cloud.jenkinsTunnel = config.jenkinsTunnel ?: ''
    cloud.usageRestricted = config.usageRestricted ?: true
    cloud.serverUrl = config.serverUrl ?: ''
    cloud.namespace = config.namespace ?: ''
    cloud.jenkinsUrl = config.jenkinsUrl ?: ''
    cloud.containerCapStr = config.containerCapStr ?: '10'
    cloudList.add(cloud)
}
private configure(config) {
    def instance = Jenkins.getInstance()
    def clouds = instance.clouds
    config.each { name, details ->
        Iterator iter = clouds.iterator();
        while (iter.hasNext()) {
            elem = iter.next();
            if (elem.name == details.cloudName) {
               iter.remove();
            }
        }
        addKubernetesCloud(clouds, details)
    }
    def lstClouds = []
    clouds.each { elem ->
        def nKubernetesCloudC4 = new KubernetesCloudC4 (
            defaultsProviderTemplate:elem.defaultsProviderTemplate,
            name:                    elem.name,
            serverUrl:               elem.serverUrl,
            useJenkinsProxy:         elem.useJenkinsProxy,
            serverCertificate:       elem.serverCertificate,
            skipTlsVerify:           elem.skipTlsVerify,
            addMasterProxyEnvVars:   elem.addMasterProxyEnvVars,
            capOnlyOnAlivePods:      elem.capOnlyOnAlivePods,
            namespace:               elem.namespace,
            webSocket:               elem.webSocket,
            directConnection:        elem.directConnection,
            jenkinsUrl:              elem.jenkinsUrl,
            jenkinsTunnel:           elem.jenkinsTunnel,
            credentialsId:           elem.credentialsId,
            containerCap:            elem.containerCap,
            retentionTimeout:        elem.retentionTimeout,
            connectTimeout:          elem.connectTimeout,
            readTimeout :            elem.readTimeout ,
            labels:                  elem.labels,
            usageRestricted:         elem.usageRestricted,
            maxRequestsPerHost:      elem.maxRequestsPerHost,
            waitForPodSec :          elem.waitForPodSec,
            podRetention :           elem.podRetention
        )
        def lstLabels = []
        elem.podLabels.each { podLabel ->
            def nLabel = new PodLabel(key: podLabel.key, value: podLabel.value)
            lstLabels.add(nLabel)
        }
        nKubernetesCloudC4.podLabels = lstLabels
        lstClouds.add(nKubernetesCloudC4)
    }
    def lstCloudsJson = JsonOutput.toJson(lstClouds)
    return lstCloudsJson
}
{{ end }}
{{ if eq .Operation "read" }}
private configure(config) {
    def instance = Jenkins.getInstance()
    def clouds = instance.clouds
    def lstClouds = []
    config.each { name, details ->
        Iterator iter = clouds.iterator();
        while (iter.hasNext()) {
            elem = iter.next();
            if (elem.name == details.cloudName) {
                       def nKubernetesCloudC4 = new KubernetesCloudC4 (
                           defaultsProviderTemplate:elem.defaultsProviderTemplate,
                           name:                    elem.name,
                           serverUrl:               elem.serverUrl,
                           useJenkinsProxy:         elem.useJenkinsProxy,
                           serverCertificate:       elem.serverCertificate,
                           skipTlsVerify:           elem.skipTlsVerify,
                           addMasterProxyEnvVars:   elem.addMasterProxyEnvVars,
                           capOnlyOnAlivePods:      elem.capOnlyOnAlivePods,
                           namespace:               elem.namespace,
                           webSocket:               elem.webSocket,
                           directConnection:        elem.directConnection,
                           jenkinsUrl:              elem.jenkinsUrl,
                           jenkinsTunnel:           elem.jenkinsTunnel,
                           credentialsId:           elem.credentialsId,
                           containerCap:            elem.containerCap,
                           retentionTimeout:        elem.retentionTimeout,
                           connectTimeout:          elem.connectTimeout,
                           readTimeout :            elem.readTimeout ,
                           labels:                  elem.labels,
                           usageRestricted:         elem.usageRestricted,
                           maxRequestsPerHost:      elem.maxRequestsPerHost,
                           waitForPodSec :          elem.waitForPodSec,
                           podRetention :           elem.podRetention
                       )
                       def lstLabels = []
                       elem.podLabels.each { podLabel ->
                           def nLabel = new PodLabel(key: podLabel.key, value: podLabel.value)
                           lstLabels.add(nLabel)
                       }
                       nKubernetesCloudC4.podLabels = lstLabels
                       lstClouds.add(nKubernetesCloudC4)
            }
        }
    }
    def lstCloudsJson = JsonOutput.toJson(lstClouds)
    return lstCloudsJson
}
{{ end }}
{{ if eq .Operation "delete" }}
private configure(config) {
    def instance = Jenkins.getInstance()
    def clouds = instance.clouds
    config.each { name, details ->
        Iterator iter = clouds.iterator();
        while (iter.hasNext()) {
            elem = iter.next();
            if (elem.name == details.cloudName) {
               iter.remove();
            }
        }
    }
}
{{ end }}
configure 'k8s-cloud': [
        {{ if or ( eq .Operation "read") (eq .Operation "delete") }}
        cloudName    : '{{ .CloudName }}'
        {{ else if  eq .Operation "create" }}
        cloudName    : '{{ .CloudName }}',
        namespace    : '{{.Namespace}}',
        jenkinsUrl   : '{{ .JenkinsURL}}',
        jenkinsTunnel: '{{ .JenkinsTunnel}}'
        {{ end }}
]
{{ if or ( eq .Operation "create") (eq .Operation "read") }}
public class PodLabel {
    String key;
    String value;
}
public class KubernetesCloudC4 {
    String defaultsProviderTemplate;
    List<PodTemplate> templates;
    String name;
    String serverUrl;
    boolean useJenkinsProxy;
    String serverCertificate;
    boolean skipTlsVerify;
    boolean addMasterProxyEnvVars;
    boolean capOnlyOnAlivePods;
    String namespace;
    boolean webSocket;
    boolean directConnection;
    String jenkinsUrl;
    String jenkinsTunnel;
    String credentialsId;
    Integer containerCap;
    int retentionTimeout;
    int connectTimeout;
    int readTimeout;
    Map<String, String> labels;
    List<PodLabel> podLabels;
    boolean usageRestricted;
    int maxRequestsPerHost;
    Integer waitForPodSec;
    String podRetention;
}
{{ end }}
`)
)

func (k *KubernetesCloud) Get(ctx context.Context) (*KubernetesCloud, error) {
	output, err := k.renderTemplate()
	if err != nil {
		return nil, err
	}
	data := map[string]string{
		"script": strings.TrimSpace(output),
	}
	r, err := k.Jenkins.Requester.Post(ctx, k.Base, nil, k.Raw, data)

	if err != nil {
		return nil, err
	}
	if r.StatusCode == 200 {
		return k, nil
	}

	return nil, errors.New(strconv.Itoa(r.StatusCode))
}

func (k *KubernetesCloud) Create(ctx context.Context) (*KubernetesCloud, error) {
	output, err := k.renderTemplate()
	if err != nil {
		return nil, err
	}
	data := map[string]string{
		"script": strings.TrimSpace(output),
	}
	r, err := k.Jenkins.Requester.Post(ctx, k.Base, nil, k.Raw, data)

	if err != nil {
		return nil, err
	}
	if r.StatusCode == 200 {
		return k, nil
	}

	return nil, errors.New(strconv.Itoa(r.StatusCode))
}

func (k *KubernetesCloud) Delete(ctx context.Context) (bool, error) {
	output, err := k.renderTemplate()
	if err != nil {
		return false, err
	}
	data := map[string]string{
		"script": strings.TrimSpace(output),
	}
	resp, err := k.Jenkins.Requester.Post(ctx, k.Base, nil, k.Raw, data)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, errors.New(strconv.Itoa(resp.StatusCode))
	}

	return true, nil
}

func (k *KubernetesCloud) renderTemplate() (string, error) {
	temp = template.Must(template.New("cloud").Parse(string(cloudConfig)))
	err := temp.Execute(&tpl, k.K8sCloud)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
