package rundeck

import (
	"fmt"
	"testing"
)

func TestMarshalJobCommandPlugin(t *testing.T) {
	testMarshalXml(t, []marshalTest{
		marshalTest{
			"with-config",
			JobCommandPlugin{
				Type: "foo-plugin",
				Config: map[string]string{
					"woo": "foo",
					"bar": "baz",
				},
			},
			`<JobCommandPlugin type="foo-plugin"><configuration><entry key="bar" value="baz"></entry><entry key="woo" value="foo"></entry></configuration></JobCommandPlugin>`,
		},
		marshalTest{
			"with-empty-config",
			JobCommandPlugin{
				Type: "foo-plugin",
				Config: map[string]string{},
			},
			`<JobCommandPlugin type="foo-plugin"></JobCommandPlugin>`,
		},
		marshalTest{
			"with-zero-value-config",
			JobCommandPlugin{
				Type: "foo-plugin",
			},
			`<JobCommandPlugin type="foo-plugin"></JobCommandPlugin>`,
		},
	})
}

func TestUnmarshalJobCommandPlugin(t *testing.T) {
	testUnmarshalXml(t, []unmarshalTest{
		unmarshalTest{
			"with-config",
			`<JobCommandPlugin type="foo-plugin"><configuration><entry key="woo" value="foo"/><entry key="bar" value="baz"/></configuration></JobCommandPlugin>`,
			&JobCommandPlugin{},
			func (rv interface {}) error {
				v := rv.(*JobCommandPlugin)
				if v.Type != "foo-plugin" {
					return fmt.Errorf("got Type %s, but expecting foo-plugin", v.Type)
				}
				if len(v.Config) != 2 {
					return fmt.Errorf("got %i Config values, but expecting 2", len(v.Config))
				}
				if v.Config["woo"] != "foo" {
					return fmt.Errorf("Config[\"woo\"] = \"%s\", but expecting \"foo\"", v.Config["woo"])
				}
				if v.Config["bar"] != "baz" {
					return fmt.Errorf("Config[\"bar\"] = \"%s\", but expecting \"baz\"", v.Config["bar"])
				}
				return nil
			},
		},
		unmarshalTest{
			"with-empty-config",
			`<JobCommandPlugin type="foo-plugin"><configuration/></JobCommandPlugin>`,
			&JobCommandPlugin{},
			func (rv interface {}) error {
				v := rv.(*JobCommandPlugin)
				if v.Type != "foo-plugin" {
					return fmt.Errorf("got Type %s, but expecting foo-plugin", v.Type)
				}
				if len(v.Config) != 0 {
					return fmt.Errorf("got %i Config values, but expecting 0", len(v.Config))
				}
				return nil
			},
		},
	})
}
