vim.api.nvim_create_user_command("ChatGPT", function()
    local api = vim.api
    local cmd = api.nvim_command
    vim.g.floaterm_width = 0.9
    vim.g.floaterm_height = 0.9
    cmd(string.format("FloatermNew --autoclose=1 --title=ChatGpt --position=center --wintype=float ~/.local/share/nvim/site/pack/packer/start/Term_ChatGPT/bin/chatGPT"))
end, {})


vim.api.nvim_create_user_command("ChatGPTInstall", function()
    os.execute("sh ~/.local/share/nvim/site/pack/packer/start/Term_ChatGPT/install.sh")
end,{})
