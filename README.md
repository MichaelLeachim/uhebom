This is a library for unsupervised data extraction from HTML pages

It has the following name:

* **U**nsupervised
* **H**TML 
* **E**extraction
* **B**ased
* **0**n
* **M**ining Data Records

In short, **Uhebom**. 

It consists of two parts:

* MDR algorithm for extracting data regions from a HTML web page
* [Needleman–Wunsch](https://en.wikipedia.org/wiki/Needleman%E2%80%93Wunsch_algorithm) algorithm for alignment of data records

The MDR algorithm based on [Mining Data Records](https://dl.acm.org/citation.cfm?id=956826) paper. 
The implementation is heavily inspired by [this](https://github.com/scrapinghub/pydepta) library. 

The alignment part uses this [Needleman–Wunsch](https://github.com/MichaelLeachim/wunsch) implementation. 

## The purpose of this work

Is to provide a fast and portable way to extract 
repeating data in tabular form from HTML pages. 
This implementation also aims to work in JS 
environment. 


## Installation 
```golang
go get -u github.com/MichaelLeachim/uhebom
```


## Usage 

```golang

import (
  extractor "github.com/MichaelLeachim/uhebom"
  log
)

func main(){
  datum_extracted := extractor.Extract([]byte("<html><div>Hello world</div></html>))
  log.Println(datum_extracted)
}

```

## Demo 

You should check out the result of the system
TODO: implement the HTML example of this library usage.

