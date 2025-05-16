# Demo

requires: go v1.24.3

```
# execute
go version

# output
go version go1.24.3 linux/amd64
```

handle dir

```
zip -r electric.zip .

mkdir -p electric
unzip electric.zip -d electric
```

```
mkdir -p electric && cd electric
unzip path/to/electric.zip

go run main.go input1.txt
go run main.go input2.txt
```

```
# build, execute within repo
mkdir -p electric && cd electric
unzip path/to/electric.zip
go build -o electric

./electric input1.txt
```

```
# execute from $PATH
# chmod +x ./electric
# sudo mv electric /usr/local/bin/
# electric path/to/input.txt
```

# Introduction

### Before coding a single character, and certainly, only after reading README.md...

I examined the initial data and noticed the files contents had two distinct sections - [Stations] and [Charger Availability Reports] - and made the following observations:

### [Stations]

-   data seems to be formatted as a table with three rows
-   the first values in each row are 0, 1, and 2. Did someone copy line numbers from their text editor?
-   either way, they are all ints, they can be used as keys for parsing data in the row
-   the first row has data in two columns, the other two only have data in one column.

### [Charger Availability Reports]

-   Each row has data in four columns, values separated by spaces.
-   the first values in each row are ints which can also be used as keys.
-   this time, however, there are duplicate keys which means i’ll need to map the data.
-   Duplicate keys are indicative of data recorded over time.

Based on these observations and assuming the structure of data received from the source cannot be modified and alternative data sources are either unavailable and/or unfeasible…

I concluded:

-   [Stations] data would require significant assumptions to be useful in calculations. [Stations] data should not be used in its current state.

-   [Charger Availability Reports] data is both structured and complete. [Charger Availability Reports] data will be the sole source of data.

I'll need a program that can "find a needle in a haystack," undeterred by the the size of and/or quantity of haystacks. Input haystack in the form of a .txt file, output needle(s).

## Solution

I've written a program for calculating uptime for nodes in a network.

-- introduce package structure.
-- struct
-- methods

ommissions vs commissions.
