package cfgtool

import (
	"testing"
)

func TestLoadYamlFile(t *testing.T) {
	type yamlCfg struct {
		ID        int
		Name      string
		Locations []string `yaml:",flow"`
	}
	testCfg := &yamlCfg{}
	err := loadCfgFromYamlFile("../../configs/testing/testing.yaml", testCfg)
	if err != nil {
		t.Error("Load yaml file failed!")
		return
	}
	if testCfg.ID != 1 || testCfg.Name != "Panjun" || len(testCfg.Locations) != 3 {
		t.Error("read yaml file wrong!")
	}
}
