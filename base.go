package coprhd

type (
	ResourceLink struct {
		Rel  string `json:"rel,omitempty"`
		Href string `json:"href,omitempty"`
	}

	ResourceId struct {
		Id string `json:"id"`
	}

	Resource struct {
		ResourceId `json:",inline"`
		Link       ResourceLink `json:"link,omitempty"`
	}

	NamedResource struct {
		Resource `json:",inline"`
		Name     string `json:"name"`
	}

	BaseObject struct {
		NamedResource `json:",inline"`
		Inactive      bool     `json:"inactive"`
		Global        bool     `json:"global"`
		Remote        bool     `json:"remote"`
		Vdc           Resource `json:"vdc"`
		Tags          []string `json:"tags"`
		Internal      bool     `json:"internal"`
		Project       Resource `json:"project,omitempty"`
		Tenant        Resource `json:"tenant,omitempty"`
		CreationTime  int64    `json:"creation_time"`
		VArray        Resource `json:"varray,omitempty"`
		Owner         string   `json:"owner,omitempty"`
		Type          string   `json:"type,omitempty"`
	}
)
