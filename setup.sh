#!/bin/bash

# 実行場所
EXECUTE_DIR_ABSPATH=`pwd`
# このファイルの絶対パス
THIS_SCRIPT_FILE_ABSPATH=`readlink -f ${BASH_SOURCE}`
# ディレクトリの絶対パス
THIS_SCRIPT_DIR_ABSPATH=`dirname ${THIS_SCRIPT_FILE_ABSPATH}`
# このスクリプトのディレクトリに移動
cd $THIS_SCRIPT_DIR_ABSPATH

# 設定ファイルのパス
CONF_PATH="${HOME}/.app/mvsc"

if [ ! -d $CONF_PATH ]; then
  # 設定格納先のディレクトリを生成
  mkdir -p $CONF_PATH
  # 設定ファイルをコピー
  cp "config.yml" $CONF_PATH
else
  echo "config dir already exists"
fi

# 実行場所に戻る
cd $EXECUTE_DIR_ABSPATH
