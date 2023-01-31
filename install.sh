#!/bin/sh

FILE=~/.local/share/nvim/site/pack/packer/start/Term_ChatGPT/bin/chatGPT
DIRFILE=~/.local/share/nvim/site/pack/packer/start/Term_ChatGPT
if [ $# -gt 0 ]; then
    if [ ! -f "$FILE" ]; then
        cd $DIRFILE
        go build . 
        mkdir bin
        mv Term_ChatGPT ./bin/chatGPT
    fi
else
    go build .
    mkdir bin
    mv Term_ChatGPT ./bin/chatGPT
fi

