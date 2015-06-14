package rundeck

import (
	"encoding/xml"
)

type Project struct {
	XMLName        xml.Name          `xml:"project"`
	Name           string            `xml:"name"`
	Description    string            `xml:"description,omitempty"`
	URL            string            `xml:"url,attr"`
	RawConfigItems []ConfigProperty  `xml:"config,omitempty"`
	Config         map[string]string `xml:"-"`
}

type projects struct {
	XMLName  xml.Name  `xml:"projects"`
	Count    int64     `xml:"count,attr"`
	Projects []Project `xml:"project"`
}

type ConfigProperty struct {
	XMLName xml.Name `xml:"property"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:"value,attr"`
}

func (c *Client) GetAllProjects() ([]Project, error) {
	p := &projects{}
	err := c.get([]string{"projects"}, nil, p)
	inflateProjects(p.Projects)
	return p.Projects, err
}

func (c *Client) GetProject(name string) (*Project, error) {
	p := &Project{}
	err := c.get([]string{"project", name}, nil, p)
	inflateProject(p)
	return p, err
}

func inflateProject(project *Project) {
	project.Config = make(map[string]string)
	for _, config := range project.RawConfigItems {
		project.Config[config.Key] = config.Value
	}
}

func inflateProjects(projects []Project) {
	for _, project := range projects {
		inflateProject(&project)
	}
}
