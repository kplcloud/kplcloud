/**
 * @Time : 2019/6/26 4:07 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : paginator
 * @Software: GoLand
 */

package paginator

type Paginator struct {
	perPageNums int // 每页数量
	pageTotal   int // 页数
	page        int //当前是第几页
	total       int // 数据总数
}

func (p *Paginator) PerPageNums() int {
	return p.perPageNums
}
func (p *Paginator) PageTotal() int {
	return p.pageTotal
}

func (p *Paginator) Page() int {
	return p.page
}

func (p *Paginator) Nums() int {
	return p.total
}

func (p *Paginator) SetPerPageNums(perPageNums int) {
	p.perPageNums = perPageNums
}

func (p *Paginator) SetPageTotal(per, total int) {
	p.pageTotal = total/per + 1
}

func (p *Paginator) SetPage(page int) {
	p.page = page
}

func (p *Paginator) SetNums(nums int) {
	p.total = nums
}

func (p *Paginator) Offset() int {
	return (p.page - 1) * p.perPageNums
}

func (p *Paginator) Result() map[string]interface{} {
	return map[string]interface{}{
		"total":     p.total,
		"pageTotal": p.pageTotal,
		"pageSize":  p.perPageNums,
		"page":      p.page,
	}
}

func NewPaginator(page, per, total int) *Paginator {
	p := Paginator{}
	if per <= 0 {
		per = 10
	}
	if page <= 0 {
		page = 1
	}
	p.SetPerPageNums(per)
	p.SetPageTotal(per, total)
	p.SetPage(page)
	p.SetNums(total)
	return &p
}
