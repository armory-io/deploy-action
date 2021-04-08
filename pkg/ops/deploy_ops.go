package ops

type deployTask struct {
	CloudProvider string `json:"cloudProvider"`
	Account       string `json:"account"`
	Type          string `json:"type"`
}

type deployManifestTask struct {
	deployTask
	Manifests []interface{} `json:"manifests"`
	Source    string        `json:"source"`
	Moniker   Moniker       `json:"moniker"`
}

type Moniker struct {
	App string `json:"app"`
}

func NewKubernetesDeployTask(account, application string, manifestProvider ManifestProvider) (*deployManifestTask, error) {
	manifests, err := manifestProvider()
	if err != nil {
		return nil, err
	}

	return &deployManifestTask{
		deployTask: deployTask{
			CloudProvider: "kubernetes",
			Account:       account,
			Type:          "deployManifest",
		},
		Manifests: manifests,
		Source:    "text",
		Moniker:   Moniker{App: application},
	}, nil
}

type Task struct {
	Job         []interface{} `json:"job"`
	Application string        `json:"application"`
	Description string        `json:"description"`
}
