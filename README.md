# go-config parsed toml files in specific order to generate a validated struct

## Install

```shell
go get -u github.com/beckend/go-config
```

### Description and usage

- a directory with config files
- the directory of files with be parsed in order: `base.toml` -> `[env].toml` -> `local.toml` and the later one will overwrite the values of the previous, all files are optionally existing.
- `[env].toml` is calulated by providing in options struct with key `EnvKeyRunEnv`, so set for example to `RUN_ENV` it will read environment variable `RUN_ENV` and if the value is for example `staging` it becomes `staging.toml`.
- the config structs which are to be validated can be annotated correctly according to https://github.com/go-playground/validator - which is being used in this library.

Example:

when `/path/to/directory/with-configs` contains `base.toml` with contents:

```toml
APIKeyGithub = 'secret-key'
```

`local.toml` with contents:

```toml
APIKeyGithub = 'secret-key-local-dev'
```

```golang
package mypackage

import (
  fmt "fmt"
	configuration "github.com/beckend/go-config"
)


// See https://github.com/go-playground/validator
type MyConfig struct {
	APIKeyGithub string `validate:"required"`
}

func main() {
  ...

  myConfig := configuration.GetConfig(configuration.GetConfigOptions{
    CreateConfig: func(options configuration.CallbackGetConfigOptions) interface{} {
      returned := MyConfig{}

      options.FailOnError(options.Config.BindStruct("", &returned))
      options.Validate(returned)

      return returned
    },
    EnvKeyRunEnv: "RUN_ENV",
    PathConfigs:  "/path/to/directory/with-configs",
  }).(MyConfig)

  // prints secret-key-local-dev since local.toml is the last parsed in priority chain
  fmt.Println(myConfig.APIKeyGithub)
  ...
}
```
