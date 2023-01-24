vim.api.nvim_create_user_command("ChatGPT", function()
    local api = vim.api
local cmd = api.nvim_command
local width = api.nvim_win_get_width(0) * 0.9 -- 90% of the current window width
local height = api.nvim_win_get_height(0) * 0.9 -- 90% of the current window height
cmd(string.format("FloatermNew --autoclose=1 --width=%d --height=%d --title=ChatGpt --position=center --wintype=float ~/.local/share/nvim/site/pack/packer/start/Term_ChatGPT/bin/chatGPT",width,height))

end, {})
