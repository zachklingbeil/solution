# Demo

## install

```sh
# binary from "go build -o solution" included in repo
git clone https://github.com/zachklingbeil/solution.git ~/zachklingbeil && cd ~/zachklingbeil
```

### execute

```sh
./solution input1.txt

# output

# 0 100
# 1 0
# 2 75
```

```sh
./solution input2.txt

# output

# 0 66
# 1 100
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

-   [Charger Availability Reports] data is both structured and complete.
-
-   [Charger Availability Reports] data will be the sole source of data.

I'll need a program that can "find a needle in a haystack," undeterred by the the size of and/or quantity of haystacks.

### **Please refer to main.go and fx/main.go for commented code detailing this solution.**

## Final Thoughts

ommissions vs commissions.
