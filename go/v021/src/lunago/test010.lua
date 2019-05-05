-- @Author: konyka
-- @Date:   2019-05-02 11:12:54
-- @Last Modified by:   konyka
-- @Last Modified time: 2019-05-02 11:23:01
function newCounter ()
    -- body
    local count = 0
    return function ()
        -- body
            count = count + 1
            return count
    end
end


c1 = newCounter()
print(c1())

print(c1())

c2 = newCounter()
print(c2())
print(c1())
print(c2())
