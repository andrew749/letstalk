# AWS jobs

Jobs in this folder should exist on lambda so that we can invoke them
ad hoc.

Each background job should be independent and contain all the code needed.

Currently the function building process is very manual. It will likely remain this was while we have a relatively small number of functions.

## Building a function

To build a function, run the following:
```
./build_function.sh <function_name>
```

This script takes care of uploading the function to aws.
