function say_hello()
  local api = vim.api
  local cmd = api.nvim_command
  cmd("echom 'Hello, World!'")
end

