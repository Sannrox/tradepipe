#!/usr/bin/env sh
#
# Build a static binary for the host OS/ARCH
#

set -eu

VERSION=${VERSION:-$(git describe --tags --abbrev=0)}
BUILDTIME=${BUILDTIME:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}
LASTMODIFIEDDATE=${LASTMODIFIEDDATE:-$(ls -l CHANGELOG.md | awk '{print $6,  $7, $8}')}

echo "${LASTMODIFIEDDATE}"

echo "### [${VERSION}] - ${BUILDTIME}" >> CHANGELOG.md
echo "<details>" >> CHANGELOG.md
echo "<summary> Details </summary>" >> CHANGELOG.md
git log --since="${LASTMODIFIEDDATE}" --format="  * %s" --no-merges --grep=#resolve --grep=#close >> CHANGELOG.md
echo "\n" >> CHANGELOG.md
echo "</details>" >> CHANGELOG.md
exit

