-- @Author: konyka
-- @Date:   2019-05-03 09:21:20
-- @Last Modified by:   konyka
-- @Last Modified time: 2019-05-03 09:21:22
t = {a = 1, b = 2, c = 3}
for k, v in pairs(t) do
  print(k, v)
end

t = {"a", "b", "c"}
for k, v in ipairs(t) do
  print(k, v)
end