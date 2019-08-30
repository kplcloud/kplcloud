/**
 * @Time : 2019-06-25 17:52
 * @Author : solacowa@gmail.com
 * @File : jenkins
 * @Software: GoLand
 */

/**
 * 参考了: https://github.com/bndr/gojenkins
 * 由于他这个没对老版本的jenkins进行兼容，所以会有一些调整
 */

package jenkins

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Auth struct {
	Username string
	ApiToken string
}

type Jenkins interface {
	SetHTTPClient(client *http.Client)
	SetApiVersion(apiVersion string)
	GetJobs() ([]Job, error)
	GetJob(name string) (job Job, err error)
	GetJobConfig(name string) (job MavenJobItem, err error)
	GetBuild(job Job, number int) (build Build, err error)
	GetLastBuild(job Job) (build Build, err error)
	GetBuildHistory(name string) (build []*Build, err error)
	CreateJob(mavenJobItem MavenJobItem, jobName string) error
	CreateFreeJob(jobName string, freeJobItem []byte) error
	DeleteJob(job Job) error
	AddJobToView(viewName string, job Job) error
	CreateView(listView ListView) error
	GetView(name string) (view ListView, err error)
	AbortJob(name string, buildId int) error
	Build(name string, params url.Values) error
	GetBuildConsoleOutput(build Build) ([]byte, error)
	GetQueue() (queue Queue, err error)
	GetArtifact(build Build, artifact Artifact) ([]byte, error)
	SetBuildDescription(build Build, description string) error
	GetComputerObject() (co ComputerObject, err error)
	GetComputers() ([]Computer, error)
	GetComputer(name string) (computer Computer, err error)
	CreateJobParams(params Params) error
	CreateJavaJobParams(params Params) error
}

type service struct {
	auth       *Auth
	host       string
	client     *http.Client
	apiVersion string
}

func NewService(auth *Auth, host string, apiVersion string) (Jenkins, error) {
	return &service{auth: auth, host: host, apiVersion: apiVersion, client: http.DefaultClient}, nil
}

func (c *service) SetHTTPClient(client *http.Client) {
	c.client = client
}

func (c *service) SetApiVersion(apiVersion string) {
	c.apiVersion = apiVersion
}

func (c *service) GetView(name string) (view ListView, err error) {
	err = c.get(fmt.Sprintf("/view/%s", name), nil, &view)
	return
}

func (c *service) GetJobs() ([]Job, error) {
	var payload = struct {
		Jobs []Job `json:"jobs"`
	}{}
	err := c.get("", nil, &payload)
	return payload.Jobs, err
}

func (c *service) GetJob(name string) (job Job, err error) {
	err = c.get(fmt.Sprintf("/job/%s", name), nil, &job)
	return
}

func (c *service) GetJobConfig(name string) (job MavenJobItem, err error) {
	err = c.getXml(fmt.Sprintf("/job/%s/config.xml", name), nil, &job)
	return
}

func (c *service) GetBuild(job Job, number int) (build Build, err error) {
	err = c.get(fmt.Sprintf("/job/%s/%d", job.Name, number), nil, &build)
	return
}

func (c *service) GetLastBuild(job Job) (build Build, err error) {
	err = c.get(fmt.Sprintf("/job/%s/lastBuild", job.Name), nil, &build)
	return
}

func (c *service) GetBuildHistory(name string) (build []*Build, err error) {
	reqUrl := fmt.Sprintf("/%s/job/%s/api/json?tree=builds[id,number,url,fullDisplayName,description,timestamp,duration,estimatedDuration,building,keepLog,result,artifacts,actions,changeSet]", c.host, name)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("http status: %d", resp.StatusCode))
		return
	}

	var body struct {
		Class  string   `json:"_class"`
		Builds []*Build `json:"builds"`
	}
	err = c.parseResponse(resp, &body)

	build = body.Builds
	return
}

func (c *service) CreateJob(mavenJobItem MavenJobItem, jobName string) error {
	mavenJobItemXml, err := xml.Marshal(mavenJobItem)
	if err != nil {
		return err
	}
	return c.postXml("/createItem", url.Values{"name": []string{jobName}}, bytes.NewReader(mavenJobItemXml), nil)
}

func (c *service) CreateFreeJob(jobName string, freeJobItem []byte) error {
	return c.postXml("/createItem", url.Values{"name": []string{jobName}}, bytes.NewReader(freeJobItem), nil)
}

func (c *service) DeleteJob(job Job) error {
	return c.post(fmt.Sprintf("/job/%s/doDelete", job.Name), nil, nil)
}

func (c *service) AddJobToView(viewName string, job Job) error {
	return c.post(fmt.Sprintf("/view/%s/addJobToView", viewName), url.Values{"name": []string{job.Name}}, map[string]string{"name": job.Name})
}
func (c *service) CreateView(listView ListView) error {
	xmlListView, err := xml.Marshal(listView)
	if err != nil {
		return err
	}
	return c.postXml("/createView", url.Values{"name": []string{listView.Name}}, bytes.NewReader(xmlListView), nil)
}

func (c *service) AbortJob(name string, buildId int) error {
	return c.post(fmt.Sprintf("/job/%s/%d/stop", name, buildId), nil, nil)
}

func (c *service) Build(name string, params url.Values) error {
	job, err := c.GetJob(name)
	if err != nil {
		return err
	}

	if hasParams(job) {
		return c.post(fmt.Sprintf("/job/%s/buildWithParameters", job.Name), params, nil)
	} else {
		return c.post(fmt.Sprintf("/job/%s/build", job.Name), params, nil)
	}
}

func (c *service) GetBuildConsoleOutput(build Build) ([]byte, error) {
	requestUrl := fmt.Sprintf("%s/consoleText", build.URL)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()
	return ioutil.ReadAll(res.Body)
}

func (c *service) GetQueue() (queue Queue, err error) {
	err = c.get(fmt.Sprintf("/queue"), nil, &queue)
	return
}

func (c *service) GetArtifact(build Build, artifact Artifact) ([]byte, error) {
	requestUrl := fmt.Sprintf("/%s/artifact/%s", build.URL, artifact.RelativePath)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()
	return ioutil.ReadAll(res.Body)
}

func (c *service) SetBuildDescription(build Build, description string) error {
	requestUrl := fmt.Sprintf("%ssubmitDescription?description=%s", build.URL, url.QueryEscape(description))
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return err
	}

	res, err := c.sendRequest(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != 200 {
		return fmt.Errorf("Unexpected response: expected '200' but received '%d'", res.StatusCode)
	}

	return nil
}

func (c *service) GetComputerObject() (co ComputerObject, err error) {
	err = c.get(fmt.Sprintf("/computer"), nil, &co)
	return
}

func (c *service) GetComputers() ([]Computer, error) {
	var payload = struct {
		Computers []Computer `json:"computer"`
	}{}
	err := c.get("/computer", nil, &payload)
	return payload.Computers, err
}

func (c *service) GetComputer(name string) (computer Computer, err error) {
	err = c.get(fmt.Sprintf("/computer/%s", name), nil, &computer)
	return
}

func (c *service) getXml(path string, params url.Values, body interface{}) (err error) {
	reqUrl := c.buildUrl(path, params)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return
	}

	return c.parseXmlResponse(resp, body)
}

func (c *service) get(path string, params url.Values, body interface{}) (err error) {
	req, err := http.NewRequest("GET", c.buildUrl(path, params), nil)
	if err != nil {
		return
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return
	}

	if !(200 <= resp.StatusCode && resp.StatusCode <= 299) {
		return errors.New(fmt.Sprintf("error: HTTP GET returned status code %d (expected 2xx)", resp.StatusCode))
	}

	return c.parseResponse(resp, body)
}

func (c *service) post(path string, params url.Values, body interface{}) (err error) {
	requestUrl := c.buildUrl(path, params)
	req, err := http.NewRequest("POST", requestUrl, nil)
	if err != nil {
		return
	}

	if _, err := c.checkCrumb(req); err != nil {
		return err
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return
	}
	if !(200 <= resp.StatusCode && resp.StatusCode <= 299) {
		return errors.New(fmt.Sprintf("error: HTTP POST returned status code %d (expected 2xx)", resp.StatusCode))
	}

	return c.parseResponse(resp, body)
}

func (c *service) postXml(path string, params url.Values, xmlBody io.Reader, body interface{}) (err error) {
	reqUrl := c.host + path

	if params != nil {
		queryStr := params.Encode()
		if queryStr != "" {
			reqUrl += "?" + queryStr
		}
	}

	req, err := http.NewRequest("POST", reqUrl, xmlBody)
	if err != nil {
		return
	}

	if _, err := c.checkCrumb(req); err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/xml")
	resp, err := c.sendRequest(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("error: HTTP POST returned status code returned: %d", resp.StatusCode))
	}

	return c.parseXmlResponse(resp, body)

}

func (c *service) buildUrl(path string, params url.Values) (reqUrl string) {
	reqUrl = c.host + path + "/api/json"
	if params != nil {
		queryString := params.Encode()
		if queryString != "" {
			reqUrl = reqUrl + "?" + queryString
		}
	}
	return
}

func (c *service) sendRequest(req *http.Request) (*http.Response, error) {
	if c.auth != nil {
		req.SetBasicAuth(c.auth.Username, c.auth.ApiToken)
	}
	return c.client.Do(req)
}

func (c *service) parseResponse(resp *http.Response, body interface{}) (err error) {
	defer func() {
		_ = resp.Body.Close()
	}()

	if body == nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(data, body)
}

func (c *service) parseXmlResponse(resp *http.Response, body interface{}) (err error) {
	defer func() {
		_ = resp.Body.Close()
	}()

	if body == nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	xmlData := string(data)

	xmlData = strings.Replace(xmlData, "version='1.1'", "version='1.0'", 1)

	return xml.Unmarshal([]byte(xmlData), body)
}

func (c *service) checkCrumb(req *http.Request) (*http.Request, error) {

	// api - store jenkins api useCrumbs response
	api := struct {
		UseCrumbs bool `json:"useCrumbs"`
	}{}

	err := c.get("/api/json", url.Values{"tree": []string{"useCrumbs"}}, &api)
	if err != nil {
		return req, err
	}

	if !api.UseCrumbs {
		// CSRF Protection is not enabled
		return req, nil
	}

	// get crumb field and value
	crumb := struct {
		Crumb             string `json:"crumb"`
		CrumbRequestField string `json:"crumbRequestField"`
	}{}

	err = c.get("/crumbIssuer", nil, &crumb)
	if err != nil {
		return req, err
	}

	// update header
	req.Header.Set(crumb.CrumbRequestField, crumb.Crumb)

	return req, nil
}

// hasParams returns a boolean value indicating if the job is parameterized
func hasParams(job Job) bool {
	for _, action := range job.Actions {
		if len(action.ParameterDefinitions) > 0 {
			return true
		}
	}
	return false
}
