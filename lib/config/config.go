package config

type Config struct {
	Name string

	Controller ControllerConfig
	Router     RouterConfig
	Service    ServiceConfig
}

type FrameworkGetter interface {
	GetFramework() string
	GetVersion() string
	GetProtocol() string
}

type FrameworkInfo struct {
	Module   string
	Version  string
	Protocol string
}

type ControllerConfig struct {
	FrameworkInfo
	ModulePath string
}

func (c FrameworkInfo) GetFramework() string {
	return c.Module
}

func (c FrameworkInfo) GetVersion() string {
	return c.Version
}

func (c FrameworkInfo) GetProtocol() string {
	return c.Protocol
}

type RouterConfig struct {
	FrameworkInfo
	ModulePath string
}

type ServiceConfig struct {
	FrameworkInfo
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
