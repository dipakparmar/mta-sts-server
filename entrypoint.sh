#!/bin/sh

echo "Checking configuration file ..."

# check the config file exists

if [ ! -f "$HOME/.config/mta-sts-server/config.yaml" ]; then
    echo "Config file does not exist"
    echo "Will attempt to create it"
else
    echo "Config file exists"
    echo "Will attempt to use it"
    cat $HOME/.config/mta-sts-server/config.yaml
    echo "Starting server"
    /app/server start
    exit 0
fi

# check the port is set
if [ -z "$PORT" ]; then
    echo "PORT is not set"
    # use default port 8080
    PORT=8080
    echo "Using default port $PORT"
fi

# variable to use to check mx check - if domain is set then skip the mx check 
MX_CHECK=0

# check the domain is set
if [ -z "$DOMAIN" ]; then
    echo "DOMAIN is not set"
    echo "Will attempt using the domain from the mx"
    MX_CHECK=1
fi



# if domain is set then get the mx using dig and store it in the mx variable and if mx has multiple values then use space as the delimiter and store it in the mx variable

if [ -n "$DOMAIN" ]; then
    MX=$(dig +short mx $DOMAIN | awk '{print $2}' | tr ' ' ',') || true 
    echo "Finding MX for $DOMAIN"
    echo "MX is $MX"
    MX_CHECK=0
fi

# if MX_CHECK is set to 1 then check mx is set
if [ "$MX_CHECK" -eq 1 ]; then
    if [ -z "$MX" ]; then
        echo "MX is not set"
        exit 1
    fi
fi


# check the mode is set
if [ -z "$MODE" ]; then
    echo "MODE is not set"
    echo "Using default mode testing"
    MODE="testing"
fi

# check the max_age is set

if [ -z "$MAX_AGE" ]; then
    echo "MAX_AGE is not set"
    echo "Using default max_age 86400"
    MAX_AGE=86400
fi

mkdir -p $HOME/.config/mta-sts-server

# write the config file
cat <<EOF > $HOME/.config/mta-sts-server/config.yaml
port: $PORT
domain: $DOMAIN
mx: $MX
mode: $MODE
max_age: $MAX_AGE
verbose: true
EOF

echo "Config file created"

/app/server config view

# run the server
/app/server start