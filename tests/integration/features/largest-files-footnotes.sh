
# Test: Largest files table footnote for long file path
# set -e

TEST_REPO="test-long-path-repo"
LONG_PATH="a/very/long/path/that/exceeds/the/limit/for/display/in/the/table/and/should/be/truncated/by/the/tool/verylongfilename.txt"
FILE_CONTENT="This is a test file with a very long path."

# Clean up any previous test
rm -rf "$TEST_REPO"
mkdir "$TEST_REPO"
cd "$TEST_REPO"
git init -q

# Create the long path file
mkdir -p "$(dirname "$LONG_PATH")"
echo "$FILE_CONTENT" > "$LONG_PATH"
git add "$LONG_PATH"
git commit -m "Add file with long path" -q

# Run git-metrics and capture output
OUTPUT=$(../../../../git-metrics -r . 2>/dev/null)


# macOS compatible substring extraction for TRUNCATED
LONG_PATH_LENGTH=${#LONG_PATH}
FIRST_PART=$(echo "$LONG_PATH" | cut -c1-20)
LAST_PART=$(echo "$LONG_PATH" | awk -v len="$LONG_PATH_LENGTH" 'BEGIN{print substr("'$LONG_PATH'", len-19, 20)}')
TRUNCATED="${FIRST_PART}...${LAST_PART}"

# echo "$OUTPUT"

echo "TRUNCATED: $TRUNCATED"

# echo "$OUTPUT" | grep -F "${TRUNCATED} [1]"

# echo "hey"

# exit 100

TABLE_MATCH=$(echo "$OUTPUT" | grep -E "${TRUNCATED} \[1\]" 2>&1)
GREP_EXIT_CODE=$?
echo "GREP_EXIT_CODE: $GREP_EXIT_CODE"
if [ $GREP_EXIT_CODE -ne 0 ] || [ -z "$TABLE_MATCH" ]; then
  echo "FAIL: Truncated path with footnote marker not found in largest files table."
  exit 2
fi

# Check for footnote with full path at the bottom
FOOTNOTE_MATCH=$(echo "$OUTPUT" | grep "\[1\] $LONG_PATH")
if [ -z "$FOOTNOTE_MATCH" ]; then
  echo "FAIL: Footnote with full path not found."
  exit 3
fi

echo "PASS: Largest files table footnote for long file path works as expected."

# Cleanup
cd ..
rm -rf "$TEST_REPO"
