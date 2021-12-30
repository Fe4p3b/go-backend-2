math.randomseed(os.time())
request = function()
  local k = math.random(0, 1000)
  local t
  if k > 990 then
     t = "incorrect_admin_token"
  else
     t = "admin_secret_token"
  end

  local url = "/entity?token="..t.."&id="..k.."&data="..k
  return wrk.format("POST", url)
end
