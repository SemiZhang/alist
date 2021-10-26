package drivers

import (
	"fmt"
	"github.com/Xhofe/alist/conf"
	"github.com/Xhofe/alist/model"
	"github.com/Xhofe/alist/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Native struct {

}

func (n Native) Save(account model.Account) {
	log.Debugf("save a account: [%s]",account.Name)
}

func (n Native) Path(path string, account *model.Account) (*model.File, []*model.File, error) {
	fullPath := filepath.Join(account.RootFolder,path)
	log.Debugf("%s-%s-%s",account.RootFolder,path,fullPath)
	if !utils.Exists(fullPath) {
		return nil,nil,fmt.Errorf("path not found")
	}
	if utils.IsDir(fullPath) {
		result := make([]*model.File,0)
		files, err := ioutil.ReadDir(fullPath)
		if err != nil {
			return nil, nil, err
		}
		for _,f := range files {
			if strings.HasPrefix(f.Name(),".") {
				continue
			}
			time := f.ModTime()
			file := &model.File{
				Name:      f.Name(),
				Size:      f.Size(),
				Type:      0,
				UpdatedAt: &time,
			}
			if f.IsDir() {
				file.Type = conf.FOLDER
			}else {
				file.Type = utils.GetFileType(filepath.Ext(f.Name()))
			}
			result = append(result, file)
		}
		return nil, result, nil
	}
	f,err := os.Stat(fullPath)
	if err != nil {
		return nil, nil, err
	}
	time := f.ModTime()
	file := &model.File{
		Name:      f.Name(),
		Size:      f.Size(),
		Type:      utils.GetFileType(filepath.Ext(f.Name())),
		UpdatedAt: &time,
	}
	return file, nil, nil
}

func (n Native) Link(path string, account *model.Account) (string,error) {
	fullPath := filepath.Join(account.RootFolder,path)
	s, err := os.Stat(fullPath)
	if err != nil {
		return "", err
	}
	if s.IsDir() {
		return "", fmt.Errorf("can't down folder")
	}
	return fullPath,nil
}

var _ Driver = (*Native)(nil)

func init() {
	RegisterDriver("native", &Native{})
}