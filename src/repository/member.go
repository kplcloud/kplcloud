package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"gopkg.in/guregu/null.v3"
)

type MemberRepository interface {
	Find(email string) (res *types.Member, err error)
	Login(email, password string) (*types.Member, error)
	Create(username, email, password string) error
	Update(param *types.Member) error
	GetInfoById(id int64) (member *types.Member, err error)
	GetMembers(emailLike string, nsName string) (members []types.Member, err error)
	FindById(mId int64) (*types.Member, error)
	GetNssByMemberId(mId int64) ([]string, error)
	GetRolesByMemberId(mId int64) ([]*types.Role, error)
	BindWechat(email, openid string) error
	UnBindWechat(mId int64) error
	GetMembersByIds(ids []int64) (members []*types.Member, err error)
	GetMembersByNss(nss []string) (members []types.Member, err error)
	GetMembersByEmails(emails []string) (members []types.Member, err error)
	GetMembersAll() (members []types.Member, err error)
	CreateMember(m *types.Member) error
	UpdateMember(m *types.Member) error
	FindOffsetLimit(offset, limit int, email string) (res []*types.Member, err error)
	Count(email string) (int64, error)
}

type member struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &member{db: db}
}

func (c *member) CreateMember(m *types.Member) error {

	err := c.db.Save(m).Error
	if err != nil {
		return err
	}

	tx := c.db.Begin()

	err = tx.Model(&m).Association("Namespaces").Append(m.Namespaces).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&m).Association("Roles").Append(m.Roles).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (c *member) UpdateMember(m *types.Member) error {

	err := c.db.Model(m).Where("id = ?", m.ID).Updates(m).Error
	if err != nil {
		return err
	}

	tx := c.db.Begin()

	err = tx.Model(&m).Association("Namespaces").Replace(m.Namespaces).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&m).Association("Roles").Replace(m.Roles).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (c *member) Find(email string) (res *types.Member, err error) {
	var m types.Member
	if err = c.db.Where("email = ?", email).
		Preload("Namespaces").Preload("Groups").Preload("Roles").Find(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (c *member) FindById(mId int64) (*types.Member, error) {
	var m types.Member
	err := c.db.Preload("Namespaces").First(&m, mId).Error
	return &m, err
}

func (c *member) Login(email, password string) (*types.Member, error) {
	var m types.Member
	var err error
	if err = c.db.Select("id,email,username,state,created_at,updated_at").Where("email = ? AND password = ?", email, password).
		Preload("Namespaces").Preload("Groups").Preload("Roles").Find(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (c *member) Create(username, email, password string) error {
	return c.db.Save(&types.Member{
		Email:    email,
		State:    int64(1),
		Username: username,
		Password: null.StringFrom(password),
	}).Error
}

func (c *member) Update(param *types.Member) error {
	var m types.Member
	return c.db.Model(&m).Updates(param).Error
}

func (c *member) GetInfoById(id int64) (member *types.Member, err error) {
	m := new(types.Member)
	err = c.db.Attrs(types.Member{ID: 0}).FirstOrInit(&m, id).Error
	return m, err
}

func (c *member) GetMembers(emailLike string, nsName string) (members []types.Member, err error) {

	//SELECT T0.`id`, T0.`email`, T0.`username`, T0.`state`, T0.`openid`, T0.`phone`, T0.`created_at`, T0.`updated_at` FROM `members` T0
	// INNER JOIN `namespaces_memberss` T1 ON
	// T1.`members_id` = T0.`id`
	// INNER JOIN `namespaces` T2 ON T2.`id` = T1.`namespaces_id`
	// WHERE T0.`email` LIKE BINARY ?
	// AND T0.`state` != ?
	// AND T2.`name_en` = ?
	// LIMIT 1000] - `%z%`, `2`, `operations`
	var nm types.NamespacesMembers
	var m types.Member
	var n types.Namespace
	var ms []types.Member

	err = c.db.Joins("inner join "+nm.TableName()+" t1 on t1.members_id = "+m.TableName()+".id").
		Joins("inner join "+n.TableName()+" t2 on t2.id = t1.namespaces_id").
		Where(m.TableName()+".email like ?", "%"+emailLike+"%").
		Where(m.TableName()+".state != ?", 2).
		Where("t2.name = ?", nsName).
		//Preload("Namespaces","name_en = ?",nsNameEn).
		Find(&ms).Error

	return ms, err
}

func (c *member) GetNssByMemberId(mId int64) ([]string, error) {
	var m types.Member
	var ns []types.Namespace
	err := c.db.First(&m, mId).Related(&ns, "Namespaces").Error
	if err != nil {
		return nil, err
	}
	var res []string
	for _, v := range ns {
		res = append(res, v.Name)
	}
	return res, nil
}

func (c *member) GetRolesByMemberId(mId int64) (roles []*types.Role, err error) {
	var m types.Member
	err = c.db.First(&m, mId).Related(&roles, "roles").Error
	if err != nil {
		return nil, err
	}
	return
}

func (c *member) BindWechat(email, openid string) error {
	var m types.Member
	n := c.db.Model(&m).Where("email = ?", email).Update("openid", openid).RowsAffected
	if n == 0 {
		return errors.New("0 rows affected or returned ")
	} else {
		return nil
	}
}

func (c *member) UnBindWechat(mId int64) error {
	var m types.Member
	return c.db.Model(&m).Where("id = ?", mId).Update("openid", "").Error
}

func (c *member) GetMembersByIds(ids []int64) (members []*types.Member, err error) {
	var ms []*types.Member
	err = c.db.Where("id in (?)", ids).Find(&ms).Error
	return ms, err
}

func (c *member) GetMembersByNss(nss []string) (members []types.Member, err error) {
	var ms []types.Member
	if err = c.db.Where("namespace in (?)", nss).
		Preload("Namespaces").Find(&ms).Error; err != nil {
		return nil, err
	}
	return ms, nil
}

func (c *member) GetMembersByEmails(emails []string) (members []types.Member, err error) {
	var ms []types.Member
	err = c.db.Where("email in (?)", emails).Find(&ms).Error
	return ms, err
}

func (c *member) GetMembersAll() (members []types.Member, err error) {
	var ms []types.Member
	err = c.db.Select("id,email,username,state,openid").Find(&ms).Error
	return ms, err
}

func (c *member) FindOffsetLimit(offset, limit int, email string) (res []*types.Member, err error) {
	query := c.db.Select("id,email,username,state,city,department,phone,created_at,updated_at")
	if email != "" {
		query = query.Where("email like ?", "%"+email+"%")
	}

	err = query.Preload("Namespaces").Preload("Roles").Order(gorm.Expr("id DESC")).Offset(offset).Limit(limit).Find(&res).Error
	return
}

func (c *member) Count(email string) (int64, error) {
	var count int64
	var m types.Member
	query := c.db.Model(&m)
	if email != "" {
		query = query.Where("email like ?", "%"+email+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}
