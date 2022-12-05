#!/bin/sh          

CRR_HOME="/continuous-release-radar/"
GO_PATH="/usr/local/go/bin/go"

(crontab -l ; echo "0 12 * * 2 cd $CRR_HOME && $GO_PATH run . 2>&1") | crontab -