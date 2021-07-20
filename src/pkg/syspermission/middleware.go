/**
 * @Time : 5/25/21 3:12 PM
 * @Author : solacowa@gmail.com
 * @File : middleware
 * @Software: GoLand
 */

package syspermission

type Middleware func(Service) Service
