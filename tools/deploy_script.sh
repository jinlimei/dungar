#!/usr/bin/env bash
set -euo pipefail

SSH_DIR="$HOME/.ssh"
IPADDR="10.0.0.90"

umask 077

mkdir -p "$SSH_DIR"
chmod 700 "$SSH_DIR"
echo "SSH dir chmod'd correctly"

touch "$SSH_DIR/deploy.key" "$SSH_DIR/known_hosts"
chmod 600 "$SSH_DIR/deploy.key" "$SSH_DIR/known_hosts"
echo "SSH-related files chmod'd correctly"

cleanup() {
  rm -f "$SSH_DIR/deploy.key" "$SSH_DIR/known_hosts"
}

trap cleanup EXIT

if [[ -z "${PROD_SSH_KEY:-}" || -z "${PROD_KNOWN_HOSTS:-}" ]]; then
  echo "PROD_SSH_KEY AND/OR PROD_KNOWN_HOSTS ARE EMPTY"
  exit 1
fi

printf '%s\n' "$PROD_SSH_KEY" > "$SSH_DIR/deploy.key"
printf '%s\n' "$PROD_KNOWN_HOSTS" > "$SSH_DIR/known_hosts"
echo "SSH Key and Known Hosts set"

NOW=$(date +%s)
echo "Setting Deploy Date: $NOW"

CMD_ARGS=(-o BatchMode=yes -o IdentitiesOnly=yes -i "$SSH_DIR/deploy.key")

echo "Transferring Dungar..."
scp "${CMD_ARGS[@]}" ./bin/dungar dungar@"$IPADDR":"dungar_$NOW"

echo "Dungar Transferred. Updating Watcher"
ssh "${CMD_ARGS[@]}" dungar@"$IPADDR" "echo $NOW > watcher.tmp && mv watcher.tmp watcher"
echo "Completed push to production. Supervisor on production will now deploy"
