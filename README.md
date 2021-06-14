# HashGraph BFT

   Implementation of HashGraph Byzantine Fault tolerance algorithm in Golang programming language with Zero MQ.

## Contents
  - [About](#about)
  - [Installation](#installation)
  - [Execution](#execution)

## About
Implementation of Hashgraph Byzantine Fault tolerance algorithm in Go programming language with ZeroMQ framework.

<p>Based on: <br>
"The Hashgraph Protocol: Efficient Asynchronous BFT for High-Throughput Distributed Ledgers" <br>
by Leemon Baird and Atul Luykx.<p>

## Installation
### Golang
If you have not already installed **Golang** follow the instructions [here](https://golang.org/doc/install).

### Clone Repository
```bash
cd ~/go/src/
git clone https://github.com/Giorgos-P/HashGraphBFT.git
```
### Setup
<p>In bin folder create the following paths:<br>
bin/keys<br>
bin/logs/client<br>
bin/logs/error<br>
bin/logs/out<br>
</p>

## Execution
### Hashgraph Client
You can find the instructions for our BFT client [here](https://github.com/Giorgos-P/HashGraph_Client).

### Manually
To install the program and generate the keys run:
```bash
go install HashgraphBFT
HashgraphBFT generate_keys <N>      // For key generation
```
Open <N> different terminals and in each terminal run:
```bash
HashgraphBFT <ID> <N> <Clients> <Scenario> <Remote>
```

### Script
Adjust the script (HashgraphBFT/scripts/run.sh) and run:
```bash
bash ~/go/src/HashgraphBFT/scripts/run.sh
```
When you are done and want to kill the processes run:
```bash
bash ~/go/src/HashgraphBFT/scripts/killHash.sh
```

## References
- [The Hashgraph Protocol: Efficient Asynchronous BFT for High-Throughput Distributed Ledgers](https://hedera.com/hh-ieee_coins_paper-200516.pdf).