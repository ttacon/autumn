package config

type Config struct {
	Controller ControllerConfig
	Router     RouterConfig
	Service    ServiceConfig
}

type ControllerConfig struct {
	Framework  string
	ModulePath string
}

type RouterConfig struct {
	Framework  string
	ModulePath string
}

type ServiceConfig struct {
	Framework           string
	ModulePath          string
	TemplatesToGenerate []string
}

// We need to be able to specify (with sane defaults):
//
//  - API controller framework
//   - e.g. < CreateModelRequest, CreateModelResponse >
//  - API router framework
//   - router.Post('/', fn (CreateReq, CreateResp))
//  - Service framework
//   - Mongo storage, Postgresql storage, etc
//  - Support Hooks
//   - naming for functions to be picked up:
//    - `Hook$Stage$Resource$Action`, e.g. `HookPreResourceCreate`
//

//
//  ***** PLUGINS *****
//
// We'll want to suport plugins for controllers, routers and service
// frameworks. Our config will validate plugins and pull templates from
// a global container that is loaded at boot time. This will allow third
// party plugins to load template files.
