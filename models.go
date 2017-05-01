package youtrackgo

import "encoding/xml"

type Contact struct {
	Type     string `json:"type,omitempty"`
	Verified bool   `json:"verified,omitempty"`
	Email    string `json:"email,omitempty"`
	Jabber   string `json:"jabber,omitempty"`
}

type VCSUserName struct {
	Name     string `json:"name,omitempty"`
}

type User struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Login    string `json:"login,omitempty"`
	Banned   bool   `json:"banned,omitempty"`
	Guest    bool   `json:"guest,omitempty"`
	Contacts []Contact `json:"contacts,omitempty"`
	VCSUserNames []VCSUserName `json:"VCSUserNames,omitempty"`
}

type UsersResult struct {
	Users []User `json:"users"`
}

type Sprint struct {
	Id  string `xml:"id"`
	Url string `xml:"url"`
}

type AgileBoard struct {
	Id      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	Sprints []Sprint `xml:"sprints>sprint"`
}

type AgileBoards struct {
	XMLName xml.Name     `xml:"projectAgileSettingss"` // YouTrack bug reported: https://youtrack.jetbrains.com/issue/JT-41309
	Boards  []AgileBoard `xml:"agileSettings"`
}

type IssueField struct {
	Type 	string	`xml:"xsi:type,attr"`
	Name 	string	`xml:"name,attr"`
	Values 	[]string  `xml:"value"`
}

func (f IssueField)GetStringValue() string {
	return f.Values[0]
}

type Issue struct {
	Id 	string	`xml:"id,attr"`
	EntityId 	string	`xml:"entityId,attr"`
	Fields	[]IssueField	`xml:"field"`
	fields  map[string]*IssueField
}

func (i *Issue)GetField(name string) *IssueField {
	if i.fields == nil {
		i.fields = make(map[string]*IssueField)
		for k, _ := range i.Fields {
			i.fields[i.Fields[k].Name] = &i.Fields[k]
		}
	}
	return i.fields[name]
}

type Issues struct {
	XMLName	xml.Name	`xml:"issueCompacts"`
	Issues	[]Issue		`xml:"issue"`
}