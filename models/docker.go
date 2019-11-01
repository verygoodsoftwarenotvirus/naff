package models

type DockerComposeFile struct {
	Version  string                          `json:"version"`
	Services map[string]DockerComposeService `json:"services"`
}

type (
	DockerComposeBuild struct {
		Context    string `json:"context"`
		Dockerfile string `json:"dockerfile"`
	}

	DockerComposeService struct {
		Image       string                `json:"image,omitempty"`
		Environment map[string]string     `json:"environment,omitempty"`
		Logging     *DockerComposeLogging `json:"logging,omitempty"`
		Ports       []string              `json:"ports,omitempty"`
		Links       []string              `json:"links,omitempty"`
		Volumes     []DockerVolume        `json:"volumes,omitempty"`
		Command     string                `json:"command,omitempty"`
		Build       *DockerComposeBuild   `json:"build,omitempty"`
		DependsOn   []string              `json:"depends_on,omitempty"`
	}

	DockerVolume struct {
		Source string `json:"source"`
		Target string `json:"target"`
		Type   string `json:"type"`
	}

	DockerComposeLogging struct {
		Driver string `json:"driver"`
	}
)
