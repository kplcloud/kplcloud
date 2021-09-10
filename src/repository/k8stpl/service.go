/**
 * @Time : 2021/9/6 2:38 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package k8stpl

import (
	"bytes"
	"context"
	"github.com/ghodss/yaml"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"html/template"
	"k8s.io/apimachinery/pkg/util/json"
)

type Middleware func(Service) Service

type Service interface {
	EncodeTemplate(ctx context.Context, kind types.Kind, paramContent map[string]interface{}, data interface{}) (tpl []byte, err error)
	FindByKind(ctx context.Context, kind types.Kind) (tpl types.K8sTemplate, err error)
	Save(ctx context.Context, tpl *types.K8sTemplate) (err error)
	Delete(ctx context.Context, kind types.Kind) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) Delete(ctx context.Context, kind types.Kind) (err error) {

	return
}

func (s *service) Save(ctx context.Context, tpl *types.K8sTemplate) (err error) {
	return s.db.Model(tpl).Save(tpl).Error
}

func (s *service) FindByKind(ctx context.Context, kind types.Kind) (tpl types.K8sTemplate, err error) {
	err = s.db.Model(&types.K8sTemplate{}).Where("kind = ?", kind.String()).First(&tpl).Error
	return
}

func (s *service) EncodeTemplate(ctx context.Context, kind types.Kind, paramContent map[string]interface{}, data interface{}) (tpl []byte, err error) {
	t, err := s.FindByKind(ctx, kind)
	if err != nil {
		return
	}
	tmpl, err := template.New(kind.String()).Parse(t.Content)
	if err != nil {
		return
	}
	var w bytes.Buffer
	err = tmpl.Execute(&w, paramContent)
	if err != nil {
		return
	}
	paramContentJson, err := json.Marshal(paramContent)
	var p = make([]byte, (len(t.Content)*2)+(len(string(paramContentJson))*2))
	n, err := w.Read(p)
	tpl = p[:n]
	if data != nil {
		err = yaml.Unmarshal(tpl, &data)
	}
	return tpl, err
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
