#!/bin/bash
# Fix CompletableFuture return type issue in MockConfigCatController.java

CONTROLLER_FILE="$1"

if [ ! -f "$CONTROLLER_FILE" ]; then
    echo "File not found: $CONTROLLER_FILE"
    exit 0
fi

# Create a temporary file
TMP_FILE="${CONTROLLER_FILE}.tmp"

# Use awk to fix the method
awk '
/return service\.listCats\(\)/ {
    print "    try {"
    print "      var response = service.listCats().get();"
    print ""
    print "      HttpStatus httpStatus = HttpStatus.valueOf(200);"
    in_thenapply = 1
    next
}
in_thenapply && /\.thenApply\(response -&gt; \{/ {
    print "      httpServletResponse.setStatus(httpStatus.value());"
    print ""
    print "      Map<String, String> responseHeaders = response.getReservedMap();"
    print "      if (responseHeaders != null) {"
    print "        responseHeaders.forEach((k, v) -> httpServletResponse.addHeader(k, v));"
    print "      }"
    print ""
    print "      if (response.getContent() == null) {"
    print "        return ResponseEntity.noContent().build();"
    print "      }"
    in_thenapply = 0
    in_lambda = 1
    next
}
in_lambda && /return new ResponseEntity\(mappingJacksonValue, httpStatus\);/ {
    print "      return new ResponseEntity(mappingJacksonValue, httpStatus);"
    print "    } catch (Exception e) {"
    print "      log.error(\"Error in listCats\", e);"
    print "      throw new RuntimeException(\"Error calling service\", e);"
    print "    }"
    in_lambda = 0
    next
}
{ print }
' "$CONTROLLER_FILE" > "$TMP_FILE"

# Replace original file
mv "$TMP_FILE" "$CONTROLLER_FILE"

echo "✅ Fixed CompletableFuture issue in $CONTROLLER_FILE"

