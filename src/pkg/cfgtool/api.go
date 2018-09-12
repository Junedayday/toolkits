package cfgtool

// LoadCfgFromYamlFile : get config to struct from Yaml file
func LoadCfgFromYamlFile(path string, cfg interface{}) error {
	return loadCfgFromYamlFile(path, cfg)
}
