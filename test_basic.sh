#!/bin/sh

# Clear the system before starting
slv system prune -y

# Create a new profile
slv profile new -n testorg

# Create a new environment
SLV_SECRET_KEY=$(slv env new -n testenv --add | grep "Secret Key")
SLV_SECRET_KEY=$(echo $SLV_SECRET_KEY | cut -d ':' -f 2 | sed 's/^[[:space:]]*//' | sed 's/[[:space:]]*$//')

# Generate a random vault file name
VAULT_FILE="$(mktemp).slv"

# Create a new vault
slv vault new -v $VAULT_FILE -s testenv --enable-hash

# Add a secret to the vault
slv secret put -v $VAULT_FILE -n testsecret -s testvalue

# Export the secret key
export SLV_SECRET_KEY

# Get the secret value
SECRET_VALUE=$(slv secret get -v $VAULT_FILE -n testsecret)

# Check if the secret value matches
if [ "$SECRET_VALUE" != "testvalue" ]; then
    echo "Secret value does not match"
    exit 1
else
    echo "Test Successful!!"
fi

# Delete the vault
rm $VAULT_FILE