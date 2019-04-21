/*
 * @Author: rayou
 * @Date: 2019-04-03 17:54:11
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-03 17:59:01
 */

package pipeline

// 要收集的数据
type CollectData struct {
	Type  int                    // 数据类型
	Title string                 // 数据标题, 是文件的时候将会作为文件名
	Cell  map[string]interface{} // 具体数据内容
}

func NewCollectData(dataType int) *CollectData {
	cd := &CollectData{
		Type: dataType,
		Cell: make(map[string]interface{}),
	}
	return cd
}
