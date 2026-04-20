// Package config provides loading, saving, and validation of the jira-cli
// configuration file.
//
// The configuration is stored as a YAML file. By default the file is read
// from $HOME/.jira-cli.yml.  The path can be overridden by setting the
// JIRA_CLI_CONFIG environment variable to an absolute path.
//
// Typical usage:
//
//	cfg, err := config.Load()
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := cfg.Validate(); err != nil {
//		log.Fatalf("invalid config: %v", err)
//	}
//
// To persist configuration changes:
//
//	if err := config.Save(cfg, "/path/to/.jira-cli.yml"); err != nil {
//		log.Fatal(err)
//	}
package config
