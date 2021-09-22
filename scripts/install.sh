#!/bin/bash

if [ $INSTALL_PREFIX ];then
    echo "XFW_INSTALL_PREFIX = $INSTALL_PREFIX"
else
    export XFW_INSTALL_PREFIX=".."
fi

cp "../share/config.yaml" "$XFW_INSTALL_PREFIX/share/config.yaml"
