package conf

import (
	"github.com/Unknwon/goconfig"
	"strings"
)

// default
const (
	unixHttpServer         = "unix_http_server"
	supervisorD            = "supervisord"
	rpcInterfaceSupervisor = "rpcinterface:supervisor"
	supervisorCtl          = "supervisorctl"
	include                = "include"
	program                = "program"
	group                  = "group"
)

type ConfigInfo struct {
	UnixHttpServer         UnixHttpServer
	SupervisorD            SupervisorD
	RpcInterfaceSupervisor RpcInterfaceSupervisor
	SupervisorCtl          SupervisorCtl
	Include                Include
	Programs               []Program
	Groups                 []Group
}
type Group struct {
	Name     string
	Programs []string
}
type Program struct {
	Name        string
	Directory   string
	Command     string
	StdErrLog   string
	AutoRestart bool
}

type UnixHttpServer struct {
	File  string
	Chmod string
}
type SupervisorD struct {
	LogFile     string
	PidFile     string
	ChildLogDir string
}
type RpcInterfaceSupervisor struct {
	SupervisorRpcInterfaceFactory string
}
type SupervisorCtl struct {
	ServerURL string
}
type Include struct {
	Files string
}

/*
	[unix_http_server]
	file=/var/run/supervisor.sock   ; (the path to the socket file)
	chmod=0700                       ; sockef file mode (default 0700)
*/
func (c *ConfigInfo) parseUnixHttpServer(cfg *goconfig.ConfigFile) {
	c.UnixHttpServer.File = cfg.MustValue(unixHttpServer, "file")
	c.UnixHttpServer.Chmod = cfg.MustValue(unixHttpServer, "chmod")
}

/*
	[supervisord]
	logfile=/var/log/supervisor/supervisord.log ; (main log file;default $CWD/supervisord.log)
	pidfile=/var/run/supervisord.pid ; (supervisord pidfile;default supervisord.pid)
	childlogdir=/var/log/supervisor            ; ('AUTO' child log dir, default $TEMP)
*/
func (c *ConfigInfo) parseSupervisorD(cfg *goconfig.ConfigFile) {
	c.SupervisorD.LogFile = cfg.MustValue(supervisorD, "logfile")
	c.SupervisorD.PidFile = cfg.MustValue(supervisorD, "pidfile")
	c.SupervisorD.ChildLogDir = cfg.MustValue(supervisorD, "childlogdir")
}

/*
	[rpcinterface:supervisor]
	supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface
*/
func (c *ConfigInfo) parseRpcInterfaceSupervisor(cfg *goconfig.ConfigFile) {
	c.RpcInterfaceSupervisor.SupervisorRpcInterfaceFactory = cfg.MustValue(rpcInterfaceSupervisor, "supervisor.rpcinterface_factory")
}

/*
	[supervisorctl]
	serverurl=unix:///var/run/supervisor.sock ; use a unix:// URL  for a unix socket
*/
func (c *ConfigInfo) parseSupervisorCtl(cfg *goconfig.ConfigFile) {
	c.SupervisorCtl.ServerURL = cfg.MustValue(supervisorCtl, "serverurl")
}

/*
	[include]
	files = /etc/supervisor/conf.d/*.conf
*/
func (c *ConfigInfo) parseInclude(cfg *goconfig.ConfigFile) {
	c.Include.Files = cfg.MustValue(include, "files")
}

/*
	[program:dus1]
	directory=/root/src/dus
	command=/root/src/dus/dus
	stderr_log=/root/src/dus/err.log
	autorestart= true
*/
func (c *ConfigInfo) parsePrograms(cfg *goconfig.ConfigFile, section, name string) {
	var program Program
	program.Name = name
	program.Directory = cfg.MustValue(section, "directory")
	program.Command = cfg.MustValue(section, "command")
	program.StdErrLog = cfg.MustValue(section, "stderr_log")
	program.AutoRestart = cfg.MustBool(section,"autorestart",true)
	c.Programs = append(c.Programs, program)
}

/*
	[group:test]
	programs=test1,test2
 */
func (c *ConfigInfo) parseGroups(cfg *goconfig.ConfigFile, section, name string) {
	var group Group
	group.Name = name
	group.Programs = strings.Split(cfg.MustValue(section, "programs",""),",")
	c.Groups = append(c.Groups,group)
}
