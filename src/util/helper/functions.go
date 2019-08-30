/**
 * Created by GoLand.
 * Email: xzghua@gmail.com
 * Date: 2019-07-04
 * Time: 15:34
 */
package helper

import (
	"github.com/go-kit/kit/endpoint"
	"strings"
)

func EpsMapMerge(map1 map[string]endpoint.Endpoint, map2 map[string]endpoint.Endpoint) map[string]endpoint.Endpoint {
	if len(map1) > len(map2) {
		map1, map2 = map2, map1
	}
	for k, v := range map1 {
		map2[k] = v
	}
	return map2
}

// merge arr2 to arr1
func EmsArrMerge(arr1 []endpoint.Middleware, arr2 []endpoint.Middleware) []endpoint.Middleware {
	for _, v := range arr2 {
		arr1 = append(arr1, v)
	}
	return arr1
}

func GitName(gitUrl string) string {
	var gitName string

	if strings.HasPrefix(gitUrl, "git@") {
		gitUrl = strings.Split(gitUrl, ":")[1]
	}

	name := strings.Split(gitUrl, "/")
	gitName = strings.Replace(name[len(name)-2]+"/"+name[len(name)-1], ".git", "", -1)
	return gitName
}

func GitUrl(gitUrl string) string {
	var git = "git@github.com:"
	if gitUrl == "" {
		return git
	}
	if index := strings.Index(gitUrl, "http://"); index != -1 {
		git = strings.Trim("git@"+gitUrl[7:], "/")
		git += ":"
	}
	if index := strings.Index(gitUrl, "https://"); index != -1 {
		git = strings.Trim("git@"+gitUrl[8:], "/")
		git += ":"
	}
	return git
}

func FormatBuildPath(buildPath string) (path string) {
	var compareStr string
	if strings.Contains(buildPath, "Dockerfile") {
		compareStr = strings.Replace(buildPath, "Dockerfile", "", -1)
	} else {
		compareStr = buildPath
	}
	path = strings.Trim(compareStr, "/")
	if buildPath == "" || buildPath == "/" {
		return "./"
	}
	return
}

func Convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = Convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = Convert(v)
		}
	}
	return i

}
