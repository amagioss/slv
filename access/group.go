package access

type Group struct {
	Id          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Email       string   `yaml:"email"`
	Access      []string `yaml:"access"`
}
