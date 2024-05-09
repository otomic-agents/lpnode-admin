var init_data = {
  data: [
    {
      collectionName: "chainList",
      data: [
        {
          chainId: 9006,
          chainName: "BSC",
          name: "bsc smart chain",
          tokenName: "BNB",
          tokenUsd: 300,
          chainType: "evm",
          image: "otmoic/chainclient-evm:latest",
          serviceName: "chain-client-evm-bsc-server-9006",
          deployName: "chain-client-evm-bsc-9006",
          envList: [
            {
              STATUS_KEY: "chain-client-status-report-bsc",
            },
          ],
        },
        {
          chainId: 9000,
          chainName: "AVAX",
          name: "avax chain",
          tokenName: "AVAX",
          tokenUsd: 300,
          chainType: "evm",
          image: "otmoic/chainclient-evm:latest",
          serviceName: "chain-client-evm-avax-server-9000",
          deployName: "chain-client-evm-avax-9000",
          envList: [
            {
              STATUS_KEY: "chain-client-status-report-avax",
            },
          ],
        },
        {
          chainId: 60,
          chainName: "eth",
          name: "eth chain",
          tokenName: "ETH",
          tokenUsd: 300,
          chainType: "evm",
          image: "otmoic/chainclient-evm:latest",
          serviceName: "chain-client-evm-eth-server-60",
          deployName: "chain-client-evm-eth-60",
          envList: [
            {
              STATUS_KEY: "chain-client-status-report-eth",
            },
          ],
        },
        {
          chainId: 966,
          chainName: "polygon",
          name: "polygon chain",
          tokenName: "MATIC",
          tokenUsd: 300,
          chainType: "evm",
          image: "otmoic/chainclient-evm:latest",
          serviceName: "chain-client-evm-polygon-server-966",
          deployName: "chain-client-evm-polygon-966",
          envList: [
            {
              STATUS_KEY: "chain-client-status-report-polygon",
            },
          ],
        },
        {
          chainId: 614,
          chainName: "op",
          name: "OP Mainnet",
          tokenName: "ETH",
          tokenUsd: 300,
          chainType: "evm",
          image: "otmoic/chainclient-evm:latest",
          serviceName: "chain-client-evm-op-server-614",
          deployName: "chain-client-evm-op-614",
          envList: [
            {
              STATUS_KEY: "chain-client-status-report-op",
            },
          ],
        },
      ],
      filter: [
        {
          name: "chainId",
          type: "int",
        },
        {
          name: "chainName",
          type: "string",
        },
      ],
      set: [
        {
          name: "name",
          type: "string",
        },
        {
          name: "tokenName",
          type: "string",
        },
        {
          name: "tokenUsd",
          type: "int",
        },
        {
          name: "chainType",
          type: "string",
        },
        {
          name: "image",
          type: "string",
        },
        {
          name: "serviceName",
          type: "string",
        },
        {
          name: "deployName",
          type: "string",
        },
        {
          name: "envList",
          type: "array",
        },
      ],
    },
    {
      collectionName: "monitor_list",
      data: [
        {
          cron: "0,5,10,15,20,25,30,35,40,45,50,55 * * * *",
          name: "block-timestamp-monit",
          script_path: "./m_s_block_timestamp.js",
          task_type: "system",
          deploy_message: "",
        },
      ],
      filter: [
        {
          name: "name",
          type: "string",
        },
      ],
      set: [
        {
          name: "cron",
          type: "string",
        },
        {
          name: "script_path",
          type: "string",
        },
        {
          name: "task_type",
          type: "string",
        },
      ],
    },
  ],
};
var json_init_data = JSON.stringify(init_data);
