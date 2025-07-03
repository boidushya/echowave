#!/bin/bash
TYPE=${1:-patch}
LATEST=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
case $TYPE in
  patch) NEW=$(echo $LATEST | awk -F. -v OFS=. "{print \$1,\$2,\$3+1}") ;;
  minor) NEW=$(echo $LATEST | awk -F. -v OFS=. "{print \$1,\$2+1,0}") ;;
  major) NEW=$(echo $LATEST | awk -F. -v OFS=. "{print \$1+1,0,0}") ;;
esac
git tag $NEW && git push origin $NEW