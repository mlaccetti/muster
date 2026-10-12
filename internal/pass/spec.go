package pass

type Spec struct {
	ID    string `yaml:"id"`
	Units []Unit `yaml:"units"`
	Gates []Gate `yaml:"gates"`
}

type Unit struct {
	Name    string `yaml:"name"`
	Brief   string `yaml:"brief"`
	Harness string `yaml:"harness,omitempty"`
	Tier    string `yaml:"tier,omitempty"`
}

type Gate struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
	Timeout int    `yaml:"timeout_seconds"`
}
