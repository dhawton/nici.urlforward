#!/bin/sh

: ${ENV_SECRETS_DIR:=/run/secrets}

function env_secret_debug()
{
    if [ ! -z "$ENV_SECRETS_DEBUG" ]; then
        echo -e "\033[1m$@\033[0m"
    fi
}

# usage: env_secret_expand VAR
#    ie: env_secret_expand 'XYZ_DB_PASSWORD'
# (will check for "$XYZ_DB_PASSWORD" variable value for a placeholder that defines the
#  name of the docker secret to use instead of the original value. For example:
# XYZ_DB_PASSWORD={{DOCKER-SECRET:my-db.secret}}
env_secret_expand() {
  source "${ENV_SECRETS_DIR}/${1}"
}

env_secrets_expand() {
  for f in /run/secrets/*; do
    source $f
  done
}

env_secrets_expand
./app