package config

import "time"

// Identity is the top-level struct representing a developer's identity.toml.
type Identity struct {
	Meta        Meta            `toml:"meta"`
	Identity    IdentitySection `toml:"identity"`
	Stack       Stack           `toml:"stack"`
	Conventions Conventions     `toml:"conventions"`
	Projects    []Project       `toml:"projects"`
	AI          AI              `toml:"ai"`
	Learned     Learned         `toml:"learned"`
	Private     map[string]any  `toml:"private"`
}

// WithoutPrivate returns a copy of the Identity with the Private section cleared.
// All renderers must use this to prevent private data from leaking into output.
func (id *Identity) WithoutPrivate() Identity {
	cp := *id
	cp.Private = nil
	return cp
}

type Meta struct {
	Version   string    `toml:"version"`
	UpdatedAt time.Time `toml:"updated_at"`
}

type IdentitySection struct {
	Name      string `toml:"name"`
	Tone      string `toml:"tone"`
	Comments  string `toml:"comments"`
	Responses string `toml:"responses"`
	Pace      string `toml:"pace"`
}

type Stack struct {
	Primary   []string   `toml:"primary"`
	Secondary []string   `toml:"secondary"`
	Data      []string   `toml:"data"`
	Infra     []string   `toml:"infra"`
	Messaging []string   `toml:"messaging"`
	Cloud     []string   `toml:"cloud"`
	Testing   []string   `toml:"testing"`
	Avoid     StackAvoid `toml:"avoid"`
}

type StackAvoid struct {
	Items   []string          `toml:"items"`
	Reasons map[string]string `toml:"reasons"`
}

type Conventions struct {
	Formatting    string `toml:"formatting"`
	PRStyle       string `toml:"pr_style"`
	CommitStyle   string `toml:"commit_style"`
	ErrorHandling string `toml:"error_handling"`
	Naming        string `toml:"naming"`
}

type Project struct {
	Name     string   `toml:"name"`
	Repo     string   `toml:"repo"`
	Stack    []string `toml:"stack"`
	Infra    []string `toml:"infra"`
	Context  string   `toml:"context"`
	Patterns []string `toml:"patterns"`
}

type AI struct {
	Verbosity    string `toml:"verbosity"`
	Confirmation string `toml:"confirmation"`
	Suggestions  string `toml:"suggestions"`
	CodeComments string `toml:"code_comments"`
	Tests        string `toml:"tests"`
}

type Learned struct {
	Entries []string `toml:"entries"`
}
