# hmm
hmm is a tool for instant charting your data

# Why?

Imagine you have a script that gives you some numbers. Benchmark script that gives you an amount of requests per seconds, for example. It would be great to see these numbers on a chart! hmm can help you with this. You can send an output for your script to hmm's stdin and hmm will open your browser and show these numbers to you as a beautiful atimatically updating chart.

# Example

Suppose you want to observe load average numbers for your operating system. You can print them with a small bash command:

```
echo "Time 1min 5min 15min"; while true; do echo -n $(date +%H:%M:%S); echo -n " "; uptime | cut -d " " -f 9-11; sleep 1; done;
Time 1min 5min 15min
19:32:33 2.38 2.12 2.04
19:32:34 2.38 2.12 2.04
19:32:35 2.40 2.12 2.04
19:32:36 2.38 2.12 2.04
19:32:37 4.01 2.13 2.04
^C
```

Cool. We can use hmm to see this numbers as a chart by sending same data to hmm's stdin:

```
cat <(echo "Time 1min 5min 15min"; while true; do echo -n $(date +%H:%M:%S); echo -n " "; uptime | cut -d " " -f 9-11; sleep 1; done;) | hmm
```

It will look like this:

![hmm gif](hmm.gif)

# Installation

Similar to most of the go packages, just use go get:

```
go get -u githib.com/mkevac/hmm
```
