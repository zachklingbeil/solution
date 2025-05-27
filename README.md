# Solution

## install

```sh
# binary from "go build -o solution" included in repo
git clone https://github.com/zachklingbeil/solution.git ~/solution && cd ~/solution
```

```sh
# input
./solution input1.txt

# output
0 100
1 0
2 75
```

```sh
# input
./solution input2.txt

# output
0 66
1 100
```

# Introduction

### Before coding a single character, and certainly, only after reading README.md...

I noticed while examining input1.txt that the files contents had two distinct sections - [Stations] and [Charger Availability Reports] - and made the following observations:

### [Stations]

-   The challenge wants me to determine station uptime based on the uptime of chargers at a station.
-   The first values in each line are 0, 1, and 2. Did someone copy line numbers from their text editor?
-   Preconditions guaranteed that StationID's would be unique amongst ChargerID's.
-   Line 0 has data in two columns... a station with 2 chargers?
-   Lines 1,2 have data in only one column... stations with 1 charger?
-   1001, 1002, 1003, 1004

### [Charger Availability Reports]

-   The data in this section is tabular.
-   Each line has four values separated by spaces (columns).
-   The first value in each line matches values from Stations (1001, 1002, 1003, 1004).
-   This time, however, there are duplicate keys which are indicative of data recorded over time.

Based on these observations and assuming the structure of data received from the source cannot be modified and alternative data sources are either unavailable and/or unfeasible…

I concluded:

-   The first value in [Stations] are StationID's.
-   There are three stations and a total of four chargers from input1.txt.
-   input2.txt does not have unique values for ChargerID's under the [Stations] section. Based on the description of the challenge, the data does not conform to the preconditions and can't be used.
-   Even though input2.txt does not conform, the application still needs to be able to return the expected output.

### **Please refer to main.go and fx/main.go for commented code detailing this solution.**

## Final Thoughts

Given the provided runtime environment, the binary generated using "go build -o solution" along with files input1.txt and input2.txt have been included in the repo. This matches the example cli execution "./your_submission relative/path/to/input/file" from the challenge description.

The compiled binary (solution) can be moved to your system's PATH, using "sudo mv solution /usr/local/bin". Once moved, the program will execute globally, without needing to be in the project directory.

For example:
solution input1.txt
or
solution ./relative/path/to/input2.txt
