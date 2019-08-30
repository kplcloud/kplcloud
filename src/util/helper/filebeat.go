package helper

import (
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"strings"
)

type FileBeat struct {
	LogPath   string `json:"log_path"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

//获取confmap中filebeat参数内容
func FileBeatYaml(beat *FileBeat, template *types.Template) (res string, err error) {
	if beat.Namespace == "" || beat.Name == "" {
		return
	}
	beat.LogPath = logPath(beat.LogPath)
	//template, err := models.GetTemplateByKind("FileBeat")
	//if err != nil {
	//	logs.Error("get template FileBeat err ", err)
	//	return
	//}
	bat := map[string]interface{}{
		"log_path":  beat.LogPath + "*.log",
		"name":      beat.Name,
		"namespace": beat.Namespace,
	}
	res, err = encode.EncodeTemplate("FileBeat", template.Detail, bat)
	return
}

//匹配正确的日志路径
func logPath(path string) (res string) {
	if path == "" {
		res = "/var/log/"
		return
	}
	pathData := strings.Split(path, "/")
	for _, v := range pathData {
		if v != "" {
			res = res + "/" + v
		}
	}
	res = res + "/"
	return
}
