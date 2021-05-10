package conf

import (
	"fmt"
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
	fmt.Printf("")
	sectionList := cfg.GetSectionList()
	fmt.Println("sectionList:",sectionList)
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
		keyList := cfg.GetKeyList(section)
		fmt.Println("keyList:",keyList)
	}

	userListKey, err := cfg.GetValue("","USER_LIST")
	if err != nil {
		return errors.Unwrap(err)
	}
	fmt.Println(userListKey)
	userListKey2,_ := cfg.GetValue(goconfig.DEFAULT_SECTION, "USER_LIST")
	fmt.Println(userListKey2)
	maxCount := cfg.MustInt("","MAX_COUNT")
	fmt.Println(maxCount)
	maxPrice := cfg.MustFloat64("","MAX_PRICE")
	fmt.Println(maxPrice)
	isShow := cfg.MustBool("","IS_SHOW")
	fmt.Println(isShow)

	db := cfg.MustValue("test","dbdns")
	fmt.Println(db)

	dbProd := cfg.MustValue("prod","dbdns")
	fmt.Println(dbProd)

	//set 值
	cfg.SetValue("","MAX_NEW","10")
	maxNew := cfg.MustInt("","MAX_NEW")
	fmt.Println(maxNew)

	maxNew1,err := cfg.Int("","MAX_NEW")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(maxNew1)
	return nil
}

