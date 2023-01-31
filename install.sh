#!/bin/sh

FILE=~/.local/share/nvim/site/pack/packer/start/Term_ChatGPT/bin/chatGPT
DIRFILE=~/.local/share/nvim/site/pack/packer/start/Term_ChatGPT
if [ -f "$FILE" ]; then
    echo "$FILE exists."
else 
    cd $DIRFILE
    go build . 
    mv Term_ChatGPT ./bin/chatGPT
fi
