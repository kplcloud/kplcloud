/**
 * @Time : 2019-07-19 14:19
 * @Author : soupzhb@gmail.com
 * @File : uploadimg.go
 * @Software: GoLand
 */

package upload

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type SopImg interface {
	Base64ToFile(info ImgInfo) error
}

type ImgInfo struct {
	Type   string
	Base64 string
	Path   string
}

type RepImgInfo struct {
	Name      string
	PathLong  string
	PathSmall string
	Postfix   string
	Width     int
	Height    int
}

func (im *ImgInfo) Base64ToFile() (rep RepImgInfo, err error) {
	s, err := base64.StdEncoding.DecodeString(im.Base64) //转码

	bf := bytes.NewBuffer(s) //写入buffer

	//格式判断
	var img image.Image
	var err2 error
	var postfix string
	postfix = "jpg"
	err2 = errors.New("格式不正确")
	switch strings.ToLower(im.Type) {
	case "png":
		img, err2 = png.Decode(bf)
		postfix = "png"
		break
	case "jpg", "jpeg":
		img, err2 = jpeg.Decode(bf)
		postfix = "jpg"
		break
	case "gif":
		img, err2 = gif.Decode(bf)
		postfix = "gif"
		break
	}

	//获取宽高
	var w, h int
	if err2 == nil {
		c := img.Bounds()
		w = c.Max.X
		h = c.Max.Y
	} else {
		w = 0
		h = 0
	}

	name, p := im.GetName()

	//目录判断
	if im.Path == "" {
		im.Path = "/opt/data/images"
	}

	//创建目录
	err = os.MkdirAll(im.Path+p, 0711)
	if err != nil {
		return
	}

	fileName := name + "." + postfix
	filePath := im.Path + p + fileName

	rep.Name = fileName
	rep.PathLong = filePath
	rep.PathSmall = p + fileName
	rep.Postfix = postfix
	rep.Width = w
	rep.Height = h

	//写入文件
	err = ioutil.WriteFile(filePath, s, 0666)
	if err != nil {
		return
	}

	return
}

func (im *ImgInfo) GetName() (name string, path string) {
	ctx := md5.New()
	t := time.Now().UnixNano()
	ts := strconv.FormatInt(t, 10)
	n := rand.Int63n(100000)
	ns := strconv.FormatInt(n, 10)
	ctx.Write([]byte(im.Base64 + ts + ns))
	name = hex.EncodeToString(ctx.Sum(nil))
	ym := time.Now().Format("200601")
	path = "/" + ym + "/" + string([]rune(name)[10:12]) + "/" + string([]rune(name)[21:23]) + "/"
	return
}
