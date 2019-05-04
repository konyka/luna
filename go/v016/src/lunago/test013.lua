-- @Author: konyka
-- @Date:   2019-05-03 10:51:48
-- @Last Modified by:   konyka
-- @Last Modified time: 2019-05-03 10:51:51
function div0(a, b)
  if b == 0 then
    error("DIV BY ZERO !")
  else
    return a / b
  end
end

function div1(a, b) return div0(a, b) end
function div2(a, b) return div1(a, b) end

ok, result = pcall(div2, 4, 2); print(ok, result)
ok, err = pcall(div2, 5, 0);    print(ok, err)
ok, err = pcall(div2, {}, {});  print(ok, err)