package client

import (
	"fmt"

	"github.com/go-routeros/routeros/v3"
)

// ContainerConfig represents the global container configuration
// Path: /container/config
//
// This is a singleton resource - there is only one container configuration
// per RouterOS device. It configures the global container settings including
// registry URL, temporary directory, and memory limits.
type ContainerConfig struct {
	RegistryUrl string `mikrotik:"registry-url" codegen:"registry_url"`
	Tmpdir      string `mikrotik:"tmpdir" codegen:"tmpdir"`
	MemoryHigh  int    `mikrotik:"memory-high" codegen:"memory_high"`
	Username    string `mikrotik:"username" codegen:"username"`
	Password    string `mikrotik:"password" codegen:"password"`
}

// Container represents a container instance
// Path: /container
//
// Containers allow running OCI-compatible container images on RouterOS.
// Compatible with Docker Hub, GCR, Quay and other registries.
type Container struct {
	Id                  string `mikrotik:".id" codegen:"id,mikrotikID,read"`
	Name                string `mikrotik:"name" codegen:"name,required,unique_key"`
	RemoteImage         string `mikrotik:"remote-image" codegen:"remote_image"`
	Tag                 string `mikrotik:"tag" codegen:"tag,read"`
	Digest              string `mikrotik:"digest" codegen:"digest,read"`
	File                string `mikrotik:"file" codegen:"file"`
	Interface           string `mikrotik:"interface" codegen:"interface,required"`
	RootDir             string `mikrotik:"root-dir" codegen:"root_dir,required"`
	Cmd                 string `mikrotik:"cmd" codegen:"cmd"`
	Entrypoint          string `mikrotik:"entrypoint" codegen:"entrypoint"`
	Workdir             string `mikrotik:"workdir" codegen:"workdir"`
	Mounts              string `mikrotik:"mounts" codegen:"mounts"`
	Envlist             string `mikrotik:"envlist" codegen:"envlist"`
	Dns                 string `mikrotik:"dns" codegen:"dns"`
	DomainName          string `mikrotik:"domain-name" codegen:"domain_name"`
	Hostname            string `mikrotik:"hostname" codegen:"hostname"`
	Logging             bool   `mikrotik:"logging" codegen:"logging"`
	StartOnBoot         bool   `mikrotik:"start-on-boot" codegen:"start_on_boot"`
	AutoRestartInterval string `mikrotik:"auto-restart-interval" codegen:"auto_restart_interval"`
	StopSignal          int    `mikrotik:"stop-signal" codegen:"stop_signal"`
	Devices             string `mikrotik:"devices" codegen:"devices"`
	CpuList             string `mikrotik:"cpu-list" codegen:"cpu_list"`
	User                string `mikrotik:"user" codegen:"user"`
	MemoryHigh          int    `mikrotik:"memory-high" codegen:"memory_high"`
	Status              string `mikrotik:"status" codegen:"status,read"`
	Comment             string `mikrotik:"comment" codegen:"comment"`
}

// ContainerEnv represents an environment variable list
// Path: /container/envs
//
// Environment variables are organized in lists that can be referenced
// by containers through the envlist attribute.
type ContainerEnv struct {
	Id      string `mikrotik:".id" codegen:"id,mikrotikID,read"`
	List    string `mikrotik:"list" codegen:"list,required"`
	Key     string `mikrotik:"key" codegen:"key,required"`
	Value   string `mikrotik:"value" codegen:"value,required"`
	Comment string `mikrotik:"comment" codegen:"comment"`
}

// ContainerMount represents a volume mount
// Path: /container/mounts
//
// Mounts bind RouterOS filesystem paths to container paths.
// Multiple mounts can be referenced by containers through the mounts attribute.
type ContainerMount struct {
	Id      string `mikrotik:".id" codegen:"id,mikrotikID,read"`
	Name    string `mikrotik:"name" codegen:"name,required,unique_key"`
	Src     string `mikrotik:"src" codegen:"src,required"`
	Dst     string `mikrotik:"dst" codegen:"dst,required"`
	Comment string `mikrotik:"comment" codegen:"comment"`
}

// Implement Resource interface for Container
var _ Resource = (*Container)(nil)

func (c *Container) IDField() string {
	return ".id"
}

func (c *Container) ID() string {
	return c.Id
}

func (c *Container) SetID(id string) {
	c.Id = id
}

func (c *Container) AfterAddHook(reply *routeros.Reply) {
	c.Id = reply.Done.Map["ret"]
}

func (c *Container) FindField() string {
	return "name"
}

func (c *Container) FindFieldValue() string {
	return c.Name
}

func (c *Container) DeleteField() string {
	return "name"
}

func (c *Container) DeleteFieldValue() string {
	return c.Name
}

func (c *Container) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/container/add",
		Find:   "/container/print",
		Update: "/container/set",
		Delete: "/container/remove",
	}[action]
}

// Implement Resource interface for ContainerEnv
var _ Resource = (*ContainerEnv)(nil)

func (ce *ContainerEnv) IDField() string {
	return ".id"
}

func (ce *ContainerEnv) ID() string {
	return ce.Id
}

func (ce *ContainerEnv) SetID(id string) {
	ce.Id = id
}

func (ce *ContainerEnv) AfterAddHook(reply *routeros.Reply) {
	ce.Id = reply.Done.Map["ret"]
}

func (ce *ContainerEnv) FindField() string {
	return ".id"
}

func (ce *ContainerEnv) FindFieldValue() string {
	return ce.Id
}

func (ce *ContainerEnv) DeleteField() string {
	return ".id"
}

func (ce *ContainerEnv) DeleteFieldValue() string {
	return ce.Id
}

func (ce *ContainerEnv) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/container/envs/add",
		Find:   "/container/envs/print",
		Update: "/container/envs/set",
		Delete: "/container/envs/remove",
	}[action]
}

// Implement Resource interface for ContainerMount
var _ Resource = (*ContainerMount)(nil)

func (cm *ContainerMount) IDField() string {
	return ".id"
}

func (cm *ContainerMount) ID() string {
	return cm.Id
}

func (cm *ContainerMount) SetID(id string) {
	cm.Id = id
}

func (cm *ContainerMount) AfterAddHook(reply *routeros.Reply) {
	cm.Id = reply.Done.Map["ret"]
}

func (cm *ContainerMount) FindField() string {
	return "name"
}

func (cm *ContainerMount) FindFieldValue() string {
	return cm.Name
}

func (cm *ContainerMount) DeleteField() string {
	return "name"
}

func (cm *ContainerMount) DeleteFieldValue() string {
	return cm.Name
}

func (cm *ContainerMount) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/container/mounts/add",
		Find:   "/container/mounts/print",
		Update: "/container/mounts/set",
		Delete: "/container/mounts/remove",
	}[action]
}

// CRUD wrapper functions for Container

func (c *Mikrotik) AddContainer(container *Container) (*Container, error) {
	res, err := c.Add(container)
	if err != nil {
		return nil, err
	}
	return res.(*Container), nil
}

func (c *Mikrotik) FindContainer(name string) (*Container, error) {
	container := &Container{Name: name}
	res, err := c.Find(container)
	if err != nil {
		return nil, err
	}
	return res.(*Container), nil
}

func (c *Mikrotik) FindContainerById(id string) (*Container, error) {
	container := &Container{Id: id}
	cmd := []string{container.ActionToCommand(Find), "?.id=" + id}
	res, err := c.connection.Run(cmd...)
	if err != nil {
		return nil, fmt.Errorf("failed to find container by id: %w", err)
	}

	if len(res.Re) == 0 {
		return nil, fmt.Errorf("container not found with id: %s", id)
	}

	if err := Unmarshal(*res, container); err != nil {
		return nil, fmt.Errorf("failed to unmarshal container: %w", err)
	}

	return container, nil
}

func (c *Mikrotik) UpdateContainer(container *Container) (*Container, error) {
	res, err := c.Update(container)
	if err != nil {
		return nil, err
	}
	return res.(*Container), nil
}

func (c *Mikrotik) DeleteContainer(name string) error {
	container := &Container{Name: name}
	return c.Delete(container)
}

// Container lifecycle operations

func (c *Mikrotik) StartContainer(name string) error {
	cmd := []string{"/container/start", "=numbers=" + name}
	_, err := c.connection.Run(cmd...)
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}
	return nil
}

func (c *Mikrotik) StopContainer(name string) error {
	cmd := []string{"/container/stop", "=numbers=" + name}
	_, err := c.connection.Run(cmd...)
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}
	return nil
}

// ContainerConfig operations (singleton - uses /container/config/set)

func (c *Mikrotik) GetContainerConfig() (*ContainerConfig, error) {
	cmd := []string{"/container/config/print"}
	res, err := c.connection.Run(cmd...)
	if err != nil {
		return nil, fmt.Errorf("failed to get container config: %w", err)
	}

	if len(res.Re) == 0 {
		return nil, fmt.Errorf("container config not found")
	}

	var config ContainerConfig
	if err := Unmarshal(*res, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal container config: %w", err)
	}

	return &config, nil
}

func (c *Mikrotik) UpdateContainerConfig(config *ContainerConfig) (*ContainerConfig, error) {
	cmd := Marshal("/container/config/set", config)
	_, err := c.connection.Run(cmd...)
	if err != nil {
		return nil, fmt.Errorf("failed to update container config: %w", err)
	}

	return c.GetContainerConfig()
}

// CRUD wrapper functions for ContainerEnv

func (c *Mikrotik) AddContainerEnv(env *ContainerEnv) (*ContainerEnv, error) {
	res, err := c.Add(env)
	if err != nil {
		return nil, err
	}
	return res.(*ContainerEnv), nil
}

func (c *Mikrotik) FindContainerEnvById(id string) (*ContainerEnv, error) {
	env := &ContainerEnv{Id: id}
	cmd := []string{env.ActionToCommand(Find), "?.id=" + id}
	res, err := c.connection.Run(cmd...)
	if err != nil {
		return nil, fmt.Errorf("failed to find container env by id: %w", err)
	}

	if len(res.Re) == 0 {
		return nil, fmt.Errorf("container env not found with id: %s", id)
	}

	if err := Unmarshal(*res, env); err != nil {
		return nil, fmt.Errorf("failed to unmarshal container env: %w", err)
	}

	return env, nil
}

func (c *Mikrotik) UpdateContainerEnv(env *ContainerEnv) (*ContainerEnv, error) {
	res, err := c.Update(env)
	if err != nil {
		return nil, err
	}
	return res.(*ContainerEnv), nil
}

func (c *Mikrotik) DeleteContainerEnv(id string) error {
	env := &ContainerEnv{Id: id}
	return c.Delete(env)
}

// CRUD wrapper functions for ContainerMount

func (c *Mikrotik) AddContainerMount(mount *ContainerMount) (*ContainerMount, error) {
	res, err := c.Add(mount)
	if err != nil {
		return nil, err
	}
	return res.(*ContainerMount), nil
}

func (c *Mikrotik) FindContainerMount(name string) (*ContainerMount, error) {
	mount := &ContainerMount{Name: name}
	res, err := c.Find(mount)
	if err != nil {
		return nil, err
	}
	return res.(*ContainerMount), nil
}

func (c *Mikrotik) UpdateContainerMount(mount *ContainerMount) (*ContainerMount, error) {
	res, err := c.Update(mount)
	if err != nil {
		return nil, err
	}
	return res.(*ContainerMount), nil
}

func (c *Mikrotik) DeleteContainerMount(name string) error {
	mount := &ContainerMount{Name: name}
	return c.Delete(mount)
}
