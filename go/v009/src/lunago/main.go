/*
* @Author: konyka
* @Date:   2019-04-26 10:01:20
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-01 12:24:40
*/
package main

import "io/ioutil"
import "os"
import "lunago/state"

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}

		ls := state.New()
		ls.Load(data, os.Args[1], "b")
		ls.Call(0, 0)
	}
}
