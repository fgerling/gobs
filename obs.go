package obs

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const baseUrl string = "api.suse.de"

type Client struct {
	BaseURL    *url.URL
	username   string
	password   string
	httpClient *http.Client
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := xml.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.username, c.password)
	if body != nil {
		req.Header.Set("Content-Type", "application/xml")
	}
	req.Header.Set("Accept", "application/xml")
	return req, nil
}
func (c *Client) GetReleaseRequests(group string, state string) ([]ReleaseRequest, error) {
	req, err := c.newRequest("GET", "/request", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("view", "collection")
	q.Add("group", group)
	q.Add("states", state)
	req.URL.RawQuery = q.Encode()

	var collection Collection
	_, err = c.do(req, &collection)
	return collection.ReleaseRequests, err
}
func (c *Client) GetPatchinfo(rr ReleaseRequest) (*Patchinfo, error) {
	project := rr.Actions[0].Source.Project
	/* https://api.suse.de/source/SUSE:Maintenance:11688/patchinfo/_patchinfo */
	patchinfo_url := fmt.Sprintf("/source/%v/patchinfo/_patchinfo", project)
	req, err := c.newRequest("GET", patchinfo_url, nil)
	if err != nil {
		return nil, err
	}

	var patchinfo Patchinfo
	_, err = c.do(req, &patchinfo)
	return &patchinfo, err
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("Got status code: %v for %q\n", resp.StatusCode, req.URL))
	}
	defer resp.Body.Close()
	err = xml.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

func NewClient(username string, password string) *Client {
	return &Client{
		BaseURL:    &url.URL{Host: baseUrl, Scheme: "https"},
		username:   username,
		password:   password,
		httpClient: &http.Client{},
	}
}

func GetRepo(rr ReleaseRequest) string {
	trgPrjStr := strings.Replace(rr.Actions[0].Target.Project, ":", "_", -1)
	srcPrjStr := strings.Replace(rr.Actions[0].Source.Project, ":", ":/", -1)
	repo := fmt.Sprintf("http://download.suse.de/ibs/%s/%s/%s.repo", srcPrjStr, trgPrjStr, rr.Actions[0].Source.Project)
	return repo
}
