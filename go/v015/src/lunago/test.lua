-- @Author: konyka
-- @Date:   2019-04-30 15:39:10
-- @Last Modified by:   konyka
-- @Last Modified time: 2019-04-30 15:41:33


local test = {"a", "b", "c"}
test[2] = "B"
test["foo"] = "Bar"
local str = test[3] .. test[2] .. test[1] .. test["foo"] .. #test





