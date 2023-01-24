local M = {}

M.openChat = function()
    local api = vim.api
  local cmd = api.nvim_command
  local pwd = api.nvim_call_function("getcwd",{})
  cmd("echom 'Current working directory: "..pwd.."'")
end
