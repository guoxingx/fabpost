#!/bin/sh
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

CHANNEL_NAME=$1
if [ ! -n "$CHANNEL_NAME" ] ;then
	CHANNEL_NAME=mychannel
	echo CHANNEL_NAME not specified, use default.
fi
echo use CHANNEL_NAME \"$CHANNEL_NAME\"

export PATH=$GOPATH/src/github.com/hyperledger/fabric/build/bin:${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}

# remove previous crypto material and config transactions
rm -fr config
rm -fr crypto-config
mkdir config
mkdir crypto-config

# generate crypto material
cryptogen generate --config=./crypto-config.yaml
if [ "$?" -ne 0 ]; then
  echo "Failed to generate crypto material..."
  exit 1
fi

# generate genesis block for orderer
configtxgen -profile OneOrgOrdererGenesis -outputBlock ./config/genesis.block
if [ "$?" -ne 0 ]; then
  echo "Failed to generate orderer genesis block..."
  exit 1
fi

# generate channel configuration transaction
configtxgen -profile OneOrgChannel -outputCreateChannelTx ./config/channel.tx -channelID $CHANNEL_NAME
if [ "$?" -ne 0 ]; then
  echo "Failed to generate channel configuration transaction..."
  exit 1
fi

# generate anchor peer transaction
configtxgen -profile OneOrgChannel -outputAnchorPeersUpdate ./config/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for Org1MSP..."
  exit 1
fi

# replace ca file in docker-compose.yaml
CaFile=$(ls crypto-config/peerOrganizations/org1.example.com/ca | awk '{print}'|grep -v .pem)
echo replace ca file in docker-compose.yaml with \"$CaFile\"
find . -name 'docker-compose.yml' | xargs perl -pi -e "s|FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/.*|FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/$CaFile|g"
find . -name 'docker-compose.yml' | xargs perl -pi -e "s|FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server-config/.*|FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server-config/$CaFile|g"

