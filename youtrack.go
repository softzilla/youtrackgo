package youtrackgo

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/json"
	"strings"
)

type YouTrack struct {
	YouTrackBase *url.URL
	HubBase      *url.URL
	Token        string
}

func New(youtrackUrl string, hubUrl string, token string) (*YouTrack, error) {
	var err error
	if youtrackUrl == "" || token == "" {
		return nil, errors.New("Required valid YouTrack url and token")
	}
	youtrack := YouTrack{
		Token: token,
	}
	if youtrack.YouTrackBase, err = url.Parse(youtrackUrl); err != nil {
		return nil, err
	}
	if youtrack.HubBase, err = url.Parse(hubUrl); err != nil {
		return nil, err
	}
	return &youtrack, nil
}

func (yt YouTrack) CallYoutrack(endpointUrl string) ([]byte, error) {
	var (
		req      *http.Request
		err      error
		endpoint *url.URL
	)
	if endpoint, err = url.Parse(endpointUrl); err != nil {
		return nil, err
	}
	if req, err = http.NewRequest("GET", yt.YouTrackBase.ResolveReference(endpoint).String(), nil); err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+yt.Token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (yt YouTrack) CallHub(endpointUrl string) ([]byte, error) {
	var (
		req      *http.Request
		err      error
		endpoint *url.URL
	)
	if endpoint, err = url.Parse(endpointUrl); err != nil {
		return nil, err
	}
	Url := yt.HubBase.ResolveReference(endpoint).String()
	if req, err = http.NewRequest("GET", Url, nil); err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+yt.Token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("Error accessing url %s", Url))
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (yt YouTrack) GetUsers(query string) ([]User, error) {
	var (
		err error
		res []byte
	)

	res, err = yt.CallHub(fmt.Sprintf("api/rest/users?query=%s&$top=1000", url.QueryEscape(query)))
	if err != nil {
		return nil, err
	}
	u := new(UsersResult)
	if err = json.Unmarshal(res, u); err != nil {
		return nil, err
	}
	return u.Users, nil
}

func (yt YouTrack) GetAgileBoards() ([]AgileBoard, error) {
	var (
		err error
		res []byte
	)
	res, err = yt.CallYoutrack("rest/admin/agile")
	if err != nil {
		return nil, err
	}

	boards := new(AgileBoards)
	if err = xml.Unmarshal(res, boards); err != nil {
		return nil, err
	}
	return boards.Boards, nil
}


type GetIssuesParameters struct {
	Query 	string
	Start 	int
	Limit  	int
}

func (p GetIssuesParameters)GetQuery() string {
	query := make([]string, 0, 4)
	if p.Query != "" {
		query = append(query, fmt.Sprintf("filter=%s", url.QueryEscape(p.Query)))
	}
	query = append(query, fmt.Sprintf("after=%d", p.Start))
	if p.Limit == 0 {
		p.Limit = 100
	}
	query = append(query, fmt.Sprintf("max=%d", p.Limit))
	return strings.Join(query, "&")
}

func (yt YouTrack) GetIssues(params *GetIssuesParameters) ([]Issue, error) {
	var (
		err error
		res []byte
	)

	if params == nil {
		params = &GetIssuesParameters{}
	}
	reqUrl := fmt.Sprintf("rest/issue?%s", params.GetQuery())
	res, err = yt.CallYoutrack(reqUrl)
	if err != nil {
		return nil, err
	}

	issues := new(Issues)
	if err = xml.Unmarshal(res, issues); err != nil {
		return nil, err
	}
	return issues.Issues, nil
}