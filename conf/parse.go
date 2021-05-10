package conf

import (
	"github.com/Unknwon/goconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
)
type readType uint32

const (
	_ readType = iota
	ReadFromString
	ReadFromFile
)
var Config = &ConfigInfo{}
func Parse(t readType,data string) error {
	var cfg = &goconfig.ConfigFile{}
	var err error
	switch t {
	case ReadFromFile :
		cfg, err = goconfig.LoadConfigFile(data)
	case ReadFromString :
		cfg, err = goconfig.LoadFromData([]byte(data))
	}
	if err != nil {
		return errors.WithMessage(err,"Parse failed")
	}
	sectionList := cfg.GetSectionList()
	for _,section := range sectionList {
		switch section{
		case unixHttpServer :
			Config.parseUnixHttpServer(cfg)
		case supervisorCtl :
			Config.parseSupervisorCtl(cfg)
		case supervisorD :
			Config.parseSupervisorD(cfg)
		case rpcInterfaceSupervisor :
			Config.parseRpcInterfaceSupervisor(cfg)
		case include :
			Config.parseInclude(cfg)
		default:
			sections := strings.Split(section,":")
			if len(sections) != 2 {
				logrus.Warn("not support section:",section)
				continue
			}
			switch  sections[0] {
			case program :
				Config.parsePrograms(cfg,section,sections[1])
			case group :
				Config.parseGroups(cfg,section,sections[1])
			default:
				logrus.Warn("not support section:",section)
			}
		}
	}
	return nil
}

