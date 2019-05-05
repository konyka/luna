/*
* @Author: konyka
* @Date:   2019-05-05 17:08:44
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-06 07:43:38
*/

package stdlib

import . "lunago/api"

var coFuncs = map[string]GoFunction{
    "create":      coCreate,
    "resume":      coResume,
    "yield":       coYield,
    "status":      coStatus,
    "isyieldable": coYieldable,
    "running":     coRunning,
    "wrap":        coWrap,
}


















