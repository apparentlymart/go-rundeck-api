package rundeck

import (
	"encoding/xml"
)

// Project represents a project within Rundeck.
type Project struct {
	Name           string            `xml:"name"`
	Description    string            `xml:"description,omitempty"`

	// RawConfigItems is a raw representation of the XML structure of
	// project configuration. Config is a more convenient representation
	// for most cases, so on read this property is emptied and its
	// contents rebuilt in Config.
	RawConfigItems []ProjectConfigProperty  `xml:"config>property,omitempty"`

	// Config is the project configuration.
	//
	// When making requests, Config and RawConfigItems are combined to produce
	// a single set of configuration settings. Thus it isn't necessary and
	// doesn't make sense to duplicate the same properties in both properties.
	Config         map[string]string `xml:"-"`

	// URL is used only to represent server responses. It is ignored when
	// making requests.
	URL            string            `xml:"url,attr"`

	// XMLName is used only in XML unmarshalling and doesn't need to
	// be set when creating a Project to send to the server.
	XMLName        xml.Name          `xml:"project"`
}

type projects struct {
	XMLName  xml.Name  `xml:"projects"`
	Count    int64     `xml:"count,attr"`
	Projects []Project `xml:"project"`
}

type projectConfig struct {
	XMLName        xml.Name         `xml:"config"`
	RawConfigItems []ProjectConfigProperty `xml:"property,omitempty"`
}

type ProjectConfigProperty struct {
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

func (c *Client) CreateProject(project *Project) (*Project, error) {
	p := &Project{}
	deflateProject(project)
	err := c.post([]string{"projects"}, nil, project, p)
	inflateProject(p)
	return p, err
}

func (c *Client) DeleteProject(name string) error {
	return c.delete([]string{"project", name})
}

func (c *Client) SetProjectConfig(projectName string, config map[string]string) error {
	configItemsIn := make([]ProjectConfigProperty, 0, len(config))
	for k, v := range config {
		configItemsIn = append(configItemsIn, ProjectConfigProperty{
			Key:   k,
			Value: v,
		})
	}

	return c.put(
		[]string{"project", projectName, "config"},
		projectConfig{
			RawConfigItems: configItemsIn,
		},
		nil,
	)
}

func inflateProject(project *Project) {
	project.Config = make(map[string]string)
	for _, config := range project.RawConfigItems {
		project.Config[config.Key] = config.Value
	}
	project.RawConfigItems = []ProjectConfigProperty{}
}

func deflateProject(project *Project) {
	// The user is allowed to populate both RawConfigItems and
	// Config, but we assume they won't put the same config
	// item in both places. If they do, the behavior is undefined.
	rawConfigItems := project.RawConfigItems
	niceConfigItems := project.Config
	totalConfigItems := len(rawConfigItems) + len(niceConfigItems)

	// Make a new slice that has the same contents as rawConfigItems
	// but has the capacity to grow to include the niceConfigItems too.
	comboConfigItems := make([]ProjectConfigProperty, len(rawConfigItems), totalConfigItems)
	copy(comboConfigItems, rawConfigItems)

	// Now we can append the niceConfigItems.
	for k, v := range niceConfigItems {
		comboConfigItems = append(comboConfigItems, ProjectConfigProperty{
			Key:   k,
			Value: v,
		})
	}

	project.RawConfigItems = comboConfigItems
	project.Config = map[string]string{}
}

func inflateProjects(projects []Project) {
	for _, project := range projects {
		inflateProject(&project)
	}
}
