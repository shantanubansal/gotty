package client

type V1UserMe struct {
	Metadata *V1ObjectMeta `json:"metadata,omitempty"`
	Spec     *V1UserSpec   `json:"spec,omitempty"`
}

type V1UserSpec struct {
	EmailID   string   `json:"emailId,omitempty"`
	FirstName string   `json:"firstName,omitempty"`
	LastName  string   `json:"lastName,omitempty"`
	Roles     []string `json:"roles"`
}

type V1ObjectMeta struct {
	Annotations     map[string]string `json:"annotations,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
	Name            string            `json:"name,omitempty"`
	Namespace       string            `json:"namespace,omitempty"`
	ResourceVersion string            `json:"resourceVersion,omitempty"`
	SelfLink        string            `json:"selfLink,omitempty"`
	UID             string            `json:"uid,omitempty"`
}

type V1Error struct {
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
	Message string      `json:"message,omitempty"`
	Ref     string      `json:"ref,omitempty"`
}
