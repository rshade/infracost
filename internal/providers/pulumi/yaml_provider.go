package pulumi

import (
	"os"

	"encoding/json"

	"github.com/infracost/infracost/internal/config"
	"github.com/infracost/infracost/internal/schema"
	"github.com/pkg/errors"
	pyaml "github.com/pulumi/pulumi-yaml/pkg/pulumiyaml"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

type YamlProvider struct {
	ctx                  *config.ProjectContext
	Path                 string
	includePastResources bool
}

func NewYamlProvider(ctx *config.ProjectContext, includePastResources bool) schema.Provider {
	return &YamlProvider{
		ctx:                  ctx,
		Path:                 ctx.ProjectConfig.Path,
		includePastResources: includePastResources,
	}
}

func (p *YamlProvider) Type() string {
	return "pulumi_yaml"
}

func (p *YamlProvider) DisplayType() string {
	return "Pulumi Yaml File"
}

func (p *YamlProvider) AddMetadata(metadata *schema.ProjectMetadata) {
	// no op
}

func (p *YamlProvider) LoadResources(usage map[string]*schema.UsageData) ([]*schema.Project, error) {
	b, err := os.ReadFile(p.Path)
	if err != nil {
		return []*schema.Project{}, errors.Wrap(err, "Error reading Pulumi Yamlfile")
	}
	var yamlTemplate pyaml.Template
	err = yaml.Unmarshal(b, &yamlTemplate)

	if err != nil {
		return []*schema.Project{}, errors.Wrap(err, "Error reading Pulumi Yamlfile")
	}
	bjson, err := json.Marshal(yamlTemplate)
	gjsonResult := gjson.ParseBytes(bjson)

	if err != nil {
		return []*schema.Project{}, errors.Wrap(err, "Error reading Pulumi preview JSON file")
	}

	metadata := config.DetectProjectMetadata(p.ctx.ProjectConfig.Path)
	metadata.Type = p.Type()
	p.AddMetadata(metadata)
	name := schema.GenerateProjectName(metadata, p.ctx.ProjectConfig.Name, p.ctx.RunContext.IsCloudEnabled())

	project := schema.NewProject(name, metadata)
	parser := NewParser(p.ctx)
	pastResources, resources, err := parser.parsePulumiYaml(yamlTemplate, usage, gjsonResult)
	if err != nil {
		return []*schema.Project{project}, errors.Wrap(err, "Error parsing Pulumi preview JSON file")
	}

	project.PastResources = pastResources
	project.Resources = resources

	if !p.includePastResources {
		project.PastResources = nil
	}

	return []*schema.Project{project}, nil
}
