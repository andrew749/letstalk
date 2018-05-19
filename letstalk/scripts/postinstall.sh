#!/bin/sh
# postinstall.sh

# TL;DR node require() and react-native require() types conflict, so I'm
# commenting out the node type definition since we're in a RN env.
# https://github.com/DefinitelyTyped/DefinitelyTyped/issues/15960
# https://github.com/aws/aws-amplify/issues/281
# https://github.com/aws/aws-sdk-js/issues/1926

case "$OSTYPE" in
  darwin*)
    echo "OSX";
    sed -i '' "s/\(^declare var require: NodeRequire;\)/\/\/\1/g" node_modules/\@types/node/index.d.ts;;
  linux*)
    echo "LINUX";
    sed -i'' "s/\(^declare var require: NodeRequire;\)/\/\/\1/g" node_modules/\@types/node/index.d.ts;;
  *)        echo "unknown: $OSTYPE" ;;
esac
