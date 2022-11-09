#!/bin/bash          

CRR_HOME="/YOUR_PATH/continuous-release-radar/"
GO_PATH="/YOUR_PATH/go"

(crontab -l ; echo "0 12 * * 2 cd $CRR_HOME && GO_PATH run . 2>&1") | crontab -