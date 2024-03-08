#!/bin/bash

mkdir -p ~/.ssh/

IPADDR="10.0.0.90"

echo "Host $IPADDR" > ~/.ssh/ssh_config
echo "  StrictHostKeyChecking no" >> ~/.ssh/ssh_config
echo "SSH Config Set"

echo "$PROD_SSH_KEY" > ~/.ssh/deploy.key
echo "$PROD_KNOWN_HOSTS" > ~/.ssh/known_hosts
echo "SSH Key and Known Hosts set"

chmod 600 ~/.ssh/deploy.key
chmod 640 ~/.ssh/known_hosts
chmod 640 ~/.ssh/ssh_config
echo "SSH Settings CHModded correctly"

NOW=$(date +%s)
echo "Setting Deploy Date: $NOW"

cp ./bin/dungar "dungar_$NOW"
echo "Transferring Dungar..."
scp -i ~/.ssh/deploy.key "dungar_$NOW" dungar@$IPADDR:~

echo "Dungar Transferred. Updating Watcher"
ssh -i ~/.ssh/deploy.key -t dungar@$IPADDR "echo $NOW > watcher"
echo "Completed push to production. Supervisor on production will now deploy"
