# json-stream-split
 
Simple library that will split json objects out of a io.reader

It was developed to handle processing of aws firehose s3 objects where record delimiter was not set.


## Usage

```go
package main

import jsonstreamsplit "github.com/nodefortytwo/json-stream-split"import (
    os

    github.com/nodefortytwo/json-stream-split
)

func main(){
    reader, err := os.Open("someHugeFile")
    if err != nil {
        panic(err)  
    }
    
    events, err := jsonstreamsplit.Split(reader)
    ...

}
``` 