/*
 * @Author: rayou
 * @Date: 2019-04-03 17:54:11
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-03 17:59:01
 */

package pipeline
// 要收集的数据
type CollectData Struct{
	Type int	// 数据类型
	Cell map[string]interface{}		// 具体数据内容
}

func NewCollectData(dataType int) *CollectData{
	cd := &CollectData{
		Type: dataType,
		Cell: make(map[string]interface{})
	}
}