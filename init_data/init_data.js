var init_data = {
  data: [
    {
      "collectionName": "chainList",
      "data": [
        {
          "chainId": 9000,
          "chainName": "AVAX",
          "name": "Avax c chain",
          "tokenName": "AVAX",
          "tokenUsd": 150,
          "chainType": "evm"
        },
        {

          "chainId": 9006,
          "chainName": "BSC",
          "name": "bsc smart chain",
          "tokenName": "BNB",
          "tokenUsd": 300,
          "chainType": "evm"
        },
        {

          "chainId": 397,
          "chainName": "NEAR",
          "name": "near core",
          "tokenName": "NEAR",
          "tokenUsd": 100,
          "chainType": "near"
        },
        {

          "chainId": 144,
          "chainName": "XRP",
          "name": "xrp chain",
          "tokenName": "XRP",
          "tokenUsd": 100,
          "chainType": "xrp"
        }
      ],
      "filter": [
        {
          "name": "chainId",
          "type": "int"
        },
        {
          "name": "chainName",
          "type": "string"
        }
      ],
      "set": [
        {
          "name": "name",
          "type": "string"
        },
        {
          "name": "tokenName",
          "type": "string"
        },
        {
          "name": "tokenUsd",
          "type": "int"
        },
        {
          "name": "chainType",
          "type": "string"
        }
      ]
    },
    {
      "collectionName": "monitor_list",
      "data": [
        {
          "corn": "0,5,10,15,20,25,30,35,40,45,50,55 * * * *",
          "name": "block-timestamp-monit",
          "script_path": "./m_s_block_timestamp.js",
          "task_type": "system",
          "deploy_message": "",
        }
      ],
      "filter": [
        {
          "name": "name",
          "type": "string"
        }
      ],
      "set": [
        {
          "name": "corn",
          "type": "string"
        },
        {
          "name": "script_path",
          "type": "string"
        },
        {
          "name": "task_type",
          "type": "string"
        }
      ]
    }
  ]
}
var json_init_data = JSON.stringify(init_data)