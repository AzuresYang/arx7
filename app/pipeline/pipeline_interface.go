/*
 * @Author: rayou
 * @Date: 2019-04-03 18:00:56
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-05 22:07:31
 */

package pipeline

type Pipeline interface {
	Init() error
	CollectData(*CollectData)
	Stop()
	GetName() string
}
