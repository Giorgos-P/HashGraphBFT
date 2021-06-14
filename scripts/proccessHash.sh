# Find how many hashgraph processes are running
# subtract one from that number

ps aux | grep 'Hash*' | wc -l
