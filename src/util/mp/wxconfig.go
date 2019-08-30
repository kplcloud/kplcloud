/**
 * @Time : 2019-08-26 19:07 
 * @Author : soupzhb@gmail.com
 * @File : wxconfig.go
 * @Software: GoLand
 */

package mp

type resText struct{
	Content	string
}

type resArticle struct{
	Title string
	Description string
	PicURL string
	URL string
}

//文本配置
var TextConfig = map[string] resText{
	"welcomText":{"你好，欢迎关注开普勒云平台"},
	"clickAboutUsText":{"Kplcloud是什么? \n\n kplcloud是一个基于Kubernetes的轻量级PaaS平台，通过可视化的界面对应用进行管理，降低应用容器化的对度，从而减少应用容器化的时间成本。 \n\nKplcloud已在服务于宜人财富多个团队，稳定运行了近两年，目前平台已在生产环境跑着上百个应用，近千个容器。\n\n开普勒云平台地址：https://kplcloud.nsini.com"},
	"clickContactUs":{"技术支持QQ群 \n722578340"},
	}

//文章配置
var ArticleConfig = map[string] resArticle{
	"clickAboutUsArticle":{"你好，这里是xxx系统","快捷、方便、省心……",  "img/wechat/about_us.jpg",""},
}