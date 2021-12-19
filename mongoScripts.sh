mongod --config /usr/local/etc/mongod.conf --fork

ps aux | grep -v grep | grep mongod
