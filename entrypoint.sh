#!/usr/bin/env bash
echo "$INPUT_CONFIG" > ${HOME}/.spinconfig
INPUT_CONFIGPATH=${HOME}/.spinconfig /usr/bin/action