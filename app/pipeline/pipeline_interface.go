/*
 * @Author: rayou
 * @Date: 2019-04-03 18:00:56
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-03 18:10:13
 */

package pipeline

type Pipeline interface {
	Init() error
	CollectData(CollectData)
	Stop() error
	GetName() string
}
