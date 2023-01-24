vim.api.nvim_create_user_command("ChatGPT", function()
  require("chatgpt").openChat()
end, {})
