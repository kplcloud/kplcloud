/**
 * @Time : 2019/7/10 6:27 PM 
 * @Author : yuntinghu1003@gmail.com
 * @File : cpu
 * @Software: GoLand
 */

package pods

import (
	"strings"
)

type CpuInfo struct {
	Cpu       string
	Memory    string
	MaxCpu    string
	MaxMemory string
}

// CPU:200m/500m/1/2,内存:256Mi/512Mi/2G
func CreateCpuData(str string) *CpuInfo {
	cpuMap := strings.Split(str, "/")
	cpu := new(CpuInfo)
	for k, v := range cpuMap {
		if k == 0 {
			cpu.Cpu = v
			cpu.MaxCpu = v
		}
		if k == 1 {
			cpu.Memory = v
			cpu.MaxMemory = v
		}
	}
	//@todo 计算最大CPU及memory
	return cpu
}
