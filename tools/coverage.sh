#!/bin/bash
#
# Code coverage generation

COVERAGE_DIR="${COVERAGE_DIR:-coverage}"
PKG_LIST=$(go list ./... | grep -v /vendor/ | grep -v /cmd)

# Create the coverage files directory
mkdir -p "$COVERAGE_DIR";
echo "[coverage] made $COVERAGE_DIR"

# Create a coverage file for each package
for package in ${PKG_LIST}; do
    go test -covermode=count -coverprofile "${COVERAGE_DIR}/${package##*/}.cov" "$package" ;
done ;

echo "[coverage] built coverage for files"

echo "[coverage] merging coverage profiles"
# Merge the coverage profile files
echo 'mode: count' > "${COVERAGE_DIR}"/coverage.out ;
tail -q -n +2 "${COVERAGE_DIR}"/*.cov >> "${COVERAGE_DIR}"/coverage.out ;
echo "[coverage] completed merge process"

# Display the global code coverage
echo "[coverage] running global code coverage"
go tool cover -func="${COVERAGE_DIR}"/coverage.out ;
echo "[coverage] completed global code coverage"

# If needed, generate HTML report
if [ "$1" == "html" ]; then
	echo "[coverage] generating html report"
	go tool cover -html="${COVERAGE_DIR}"/coverage.out -o coverage.html ;
	echo "[coverage] generated html report"
fi

# Remove the coverage files directory
rm -rf "$COVERAGE_DIR";
