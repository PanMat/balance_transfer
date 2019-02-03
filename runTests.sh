#bin/bash
# run this script to test REST API
# set up variables for use
s=1
echo $(export GOPATH="/home/anshu/go")

#set chaincode version here --> increment this if you have made any changes to chaincode
ccVersion="v6"
transactions=200
cat ./utils/welcome
cat ./utils/start
echo -e "\033[33;33mStarting Hyperledger Network Example.... \033[0m"
sleep 0

query_blockheight(){
    echo -e "\033[33;33mQuering block height in the ledger\033[0m"
    curl -s -X GET  http://localhost:4000/channels/mychannel/blocks \
        -H "authorization: $ORG1_TOKEN" \
        -H "content-type: application/json" \
        -d '{
            "peer": "peer0.org1.example.com"
    }'
    echo $'\n'
}

query_block_by_Number(){
    echo -e "\033[33;33mQuering block#'$1' in the ledger\033[0m"
    curl -s -X GET  http://localhost:4000/channels/mychannel/blockByID \
        -H "authorization: $ORG1_TOKEN" \
        -H "content-type: application/json" \
        -d '{
            "peer": "peer0.org1.example.com",
            "blockID": "'$1'"
    }'
    echo $'\n'
}

#POST request or Enroll on Org1
echo -e "\033[33;33mPOST request for Enroll on Org1 \033[0m"
ORG1_TOKEN=$(curl -s -X POST http://localhost:4000/users -H "content-type: application/x-www-form-urlencoded" -d 'username=John&orgName=org1&secret=thisismysecret')
ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".Token" | sed "s/\"//g")
echo $ORG1_TOKEN
echo $'\n'
sleep $s

#POST request or Enroll on Org2
echo -e "\033[33;33mPOST request for Enroll on Org2 \033[0m"

ORG2_TOKEN=$(curl -s -X POST http://localhost:4000/users -H "content-type: application/x-www-form-urlencoded" -d 'username=Jeff&orgName=org2&secret=thisismysecret')
ORG2_TOKEN=$(echo $ORG2_TOKEN | jq ".Token" | sed "s/\"//g")
echo $ORG2_TOKEN
echo $'\n'
sleep $s

#POST request to Create Channel
echo -e "\033[33;33mPOST request to Create Channel \033[0m" 
echo "channel creation takes time so wait for it!"
curl -s -X POST http://localhost:4000/channels -H "content-type: application/json" \
 -H "Authorization: $ORG1_TOKEN" \
 -d '{
     "channelName":"mychannel", 
     "channelConfigPath":"./artifacts/channel/channel.tx"
     }'
sleep 5

#POST request: Org-1 joining channel
echo $'\n'
echo -e "\033[33;33mPOST request: Org-1 joining channel \033[0m" 
curl -s -X POST http://localhost:4000/channels/mychannel \
    -H "Authorization: $ORG1_TOKEN" \
    -H "content-type: application/json" \
    -d '{
        "peers": ["peer0.org1.example.com",
        "peer1.org1.example.com"]
        }'
echo $'\n'
sleep $s

#POST request: Org-2 joining channel
echo -e "\033[33;33mPOST request: Org-2 joining channel \033[0m" 
curl -s -X POST http://localhost:4000/channels/mychannel \
    -H "Authorization: $ORG2_TOKEN" \
    -H "content-type: application/json" \
    -d '{
        "peers": ["peer0.org2.example.com",
        "peer1.org2.example.com"]
        }'
echo $'\n'
sleep $s

# Install Chain Code on peers at Org1
echo -e "\033[33;33mInstall Chaincode on peers at Org1 \033[0m" 
curl -s -X POST http://localhost:4000/chaincodes \
    -H "authorization: $ORG1_TOKEN" \
    -H "content-type: application/json" \
    -d '{
        "peers": ["peer0.org1.example.com","peer1.org1.example.com"],
        "chaincodeName":"mycc",
        "chaincodePath":"github.com/balance-transfer-go/chaincode/chaincode_example02/go",
        "chaincodeType": "golang",
        "chaincodeVersion":"'$ccVersion'"
        }'
echo $'\n'
sleep $s

# Query Block Height in the ledger
query_blockheight
sleep $s

# Install Chain Code on peers at Org2
echo -e "\033[33;33mInstall Chaincode on peers at Org1 \033[0m" 
curl -s -X POST http://localhost:4000/chaincodes \
    -H "authorization: $ORG2_TOKEN" \
    -H "content-type: application/json" \
    -d '{
        "peers": ["peer0.org2.example.com","peer1.org2.example.com"],
        "chaincodeName":"mycc",
        "chaincodePath":"github.com/balance-transfer-go/chaincode/chaincode_example02/go",
        "chaincodeType": "golang",
        "chaincodeVersion":"'$ccVersion'"
        }'
echo $'\n'
sleep $s

#Instantiate Chaincode
echo -e "\033[33;33mInstantiating chaincode on the channel \033[0m" 
curl -s -X POST http://localhost:4000/channels/mychannel/instantiate \
    -H "authorization: $ORG1_TOKEN" \
    -H "content-type: application/json" \
    -d '{
        "peers": ["peer0.org1.example.com","peer1.org1.example.com"],
        "chaincodeName":"mycc",
        "chaincodeVersion":"'$ccVersion'",
        "chaincodeType": "golang", 
        "args": [ "Init", "John","400" ],
        "chaincodePath":"github.com/balance-transfer-go/chaincode/chaincode_example02/go"
        }'
echo $'\n'
sleep $s


#Invoke chaincode
echo -e "\033[33;33mInvoking chaincode 'Add function' on the channel \033[0m" 
curl -s -X POST  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org1.example.com","peer1.org1.example.com"],
	"fcn":"add",
	"args":["Jeff","500"]
}'
echo $'\n'
sleep $s

# Query Block Height in the ledger
query_blockheight
sleep $s

#Move balance between users
echo -e "\033[33;33mMoving some amount from John to Jeff \033[0m" 
curl -s -X POST  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org1.example.com","peer1.org1.example.com"],
	"fcn":"move",
	"args":["John","Jeff","10"]
}'
echo $'\n'
sleep $s

#Query chaincode
echo -e "\033[33;33mQuering chaincode for current balances of John \033[0m" 
curl -s -X GET  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peer": "peer0.org2.example.com",
	"fcn":"query",
	"args":[ "John" ]
}'
echo $'\n'
sleep $s

#Query chaincode
echo -e "\033[33;33mQuering chaincode for current balances of Jeff \033[0m" 
curl -s -X GET  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peer": "peer0.org2.example.com",
	"fcn":"query",
	"args":[ "Jeff" ]
}'
echo $'\n'
sleep $s

#Delete user from ledger
echo -e "\033[33;33mDeleting user John from ledger \033[0m" 
curl -s -X GET  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peer": "peer0.org2.example.com",
	"fcn":"delete",
	"args":[ "John" ]
}'
echo $'\n'
sleep $s

# Query Block Height in the ledger
query_blockheight
sleep $s

# add more users and mock faster transactions
echo -e "\033[33;33mAdding more users and try to mock faster transactions\033[0m" 
user_org1=("Diann" "Silvio" "Keri" "Issiah" "Estell" "Ondrea" "Colleen" "Elena" "Eolanda" "Sonny")
bal_org1=(1227 2120 3863 3063 3343 1964 1873 1103 3880 2170 )
user_org2=("Ilena" "Ly" "Dulsea" "Shepperd" "Caren" "Nikita" "Fannie" "Addison" "Sherlyn" "Ilsa")
bal_org2=(2887 1899 2672 2872 3892 2521 2379 3683 3080 1869)

addUsers(){
  echo "Adding all users and balances to the ledger!!!!!"  
  echo $'\n'
  echo "    Org-1..........."
  for ((i = 0; i < "${#user_org1[@]}"; i++))
    do          
        user=${user_org1[$i]}
        amt=${bal_org1[$i]}
        # adding users from Org1
        msg=$(curl -s -X POST  http://localhost:4000/channels/mychannel/chaincodes/mycc \
            -H "authorization: $ORG1_TOKEN" \
            -H "content-type: application/json" \
            -d '{
                "peers": ["peer0.org1.example.com","peer1.org1.example.com"],
                "fcn":"add",
                "args":["'${user_org1[$i]}'","'${bal_org1[$i]}'"]
            }')
            echo "  $msg"
    done
    echo $'\n'
    echo "    Org-2..........."
  for ((i = 0; i < "${#user_org2[@]}"; i++))
    do        
        user=${user_org2[$i]}
        amt=${bal_org2[$i]}
        # adding users from Org2
        msg=$(curl -s -X POST  http://localhost:4000/channels/mychannel/chaincodes/mycc \
            -H "authorization: $ORG2_TOKEN" \
            -H "content-type: application/json" \
            -d '{
                "peers": ["peer0.org1.example.com","peer1.org1.example.com"],
                "fcn":"add",
                "args":["'${user_org2[$i]}'","'${bal_org2[$i]}'"]
            }')
            echo "  $msg"
    done
  echo $'\n'
  echo -e "\033[33;33mEnrolled all users....whew!!!!\033[0m" 
  echo $'\n'
}

# create few mock transactions
random(){
  x=$(($1 + RANDOM%(1+$2-$1)))
  echo $x
}
mock_transactions(){
    for ((i = 0; i< transactions; i++))
        do
            amount=$(random 5 10)
            sID=$(random 0 ${#user_org1[@]}-1)
            rID=$(random 0 ${#user_org2[@]}-1)
            echo "Performing transaction=$i, moving $amount from ${user_org1[${sID}]} to ${user_org2[${rID}]}"
            msg=$(curl -s -X POST  http://localhost:4000/channels/mychannel/chaincodes/mycc -H "authorization: $ORG1_TOKEN"  -H "content-type: application/json" -d '{"peers": ["peer0.org1.example.com","peer1.org1.example.com"], "fcn":"move","args":["'${user_org1[${sID}]}'","'${user_org2[${rID}]}'","'$amount'"]}')
            echo "   $msg"
        done
}

query_all_users(){
    echo -e "\033[33;33mFinal query of user balances\033[0m" 
    echo "........Org-1......"
      for ((i = 0; i < "${#user_org1[@]}"; i++))
        do
            msg=$(curl -s -X GET  http://localhost:4000/channels/mychannel/chaincodes/mycc -H "authorization: $ORG1_TOKEN" -H "content-type: application/json" -d '{"peer": "peer0.org1.example.com", "fcn":"query","args":[ "'${user_org1[$i]}'" ]}')
            echo $msg
        done
    echo $'\n'
    echo "........Org-2......"
      for ((i = 0; i < "${#user_org2[@]}"; i++))
        do
            msg=$(curl -s -X GET  http://localhost:4000/channels/mychannel/chaincodes/mycc -H "authorization: $ORG1_TOKEN" -H "content-type: application/json" -d '{"peer": "peer0.org1.example.com", "fcn":"query","args":[ "'${user_org2[$i]}'" ]}')
            echo $msg
        done
}


addUsers
mock_transactions
query_all_users

# Query Block Height in the ledger
query_blockheight
sleep $s

query_block_by_Number 20
sleep $s
