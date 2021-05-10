package conf

import (
	"github.com/dusbin/gosupervisor/ulog"
	"testing"
)

func TestParse(t *testing.T) {
	ulog.InitLog("conf_test","debug")
	err := Parse(ReadFromString,configData)
	t.Log(err)
	t.Log(Config)
	err = Parse(ReadFromFile,notExistPath)
	t.Log(err)
}
var notExistPath=`notExistPath`
var configData = `
; supervisor config file

[unix_http_server]
; (the path to the socket file)
file=/var/run/supervisor.sock   
; sockef file mode (default 0700)
chmod=0700                       

[program:dus1]
directory=/root/src/dus
command=/root/src/dus/dus
stderr_log=/root/src/dus/err.log
autorestart= true

[group:dus]
programs=dus1


[program:test1]
directory=/root/src/dus
command=/root/src/dus/test1
stderr_log=/root/src/dus/err_test1.log
autorestart=true

[group:test]
programs=test1,test2

[supervisord]
; (main log file;default $CWD/supervisord.log)
logfile=/var/log/supervisor/supervisord.log 
; (supervisord pidfile;default supervisord.pid)
pidfile=/var/run/supervisord.pid 
; ('AUTO' child log dir, default $TEMP)
childlogdir=/var/log/supervisor            

; the below section must remain in the config file for RPC
; (supervisorctl/web interface) to work, additional interfaces may be
; added by defining them in separate rpcinterface: sections
[rpcinterface:supervisor]
supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface

[supervisorctl]
; use a unix:// URL  for a unix socket
serverurl=unix:///var/run/supervisor.sock 

; The [include] section can just contain the "files" setting.  This
; setting can list multiple files (separated by whitespace or
; newlines).  It can also contain wildcards.  The filenames are
; interpreted as relative to this file.  Included files *cannot*
; include files themselves.

[include]
files = /etc/supervisor/conf.d/*.conf
`
