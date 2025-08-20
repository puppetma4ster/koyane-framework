package analyzer

import (
	"os"

	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Info    YamlInfo        `yaml:"info"`
	General GeneralAnalyzer `yaml:"general"`
	Content AnalyzerContent `yaml:"content"`
}
type YamlInfo struct {
	HasGeneral bool
	HasContent bool
	HasStats   bool
}

func SaveToYamlAll(general *GeneralAnalyzer, content *AnalyzerContent) error {
	info := YamlInfo{
		HasGeneral: true,
		HasContent: true,
		HasStats:   true,
	}
	cfg := Config{
		Info:    info,
		General: *general,
		Content: *content,
	}
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	absolutePath, err := utils.ResolvePath(utils.AnalyzedWlSaveDir)
	if err != nil {
		return err
	}
	err = os.WriteFile(absolutePath+"/"+general.HashVal+".yaml", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadFromYamlAll(analyzer GeneralAnalyzer) (*GeneralAnalyzer, *AnalyzerContent, error) {
	var fileName string = analyzer.HashVal + ".yaml"
	absolutePath, err := utils.ResolvePath(utils.AnalyzedWlSaveDir)
	if err != nil {
		return nil, nil, err
	}

	data, err := os.ReadFile(absolutePath + "/" + fileName)
	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	return &cfg.General, &cfg.Content, nil
}

func IsSavedObj(general *GeneralAnalyzer) (bool, error) {
	var fileName string = general.HashVal + ".yaml"

	absolutePath, err := utils.ResolvePath(utils.AnalyzedWlSaveDir)
	if err != nil {
		return false, err
	}
	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if fileName == entry.Name() {
			return true, nil
		}
	}
	return false, nil
}
