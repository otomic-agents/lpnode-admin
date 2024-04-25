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
          image: "kldtks/edge:otmoic-chainclient-evm-6ef74826",
          serviceName: "chain-client-evm-bsc-server-9006",
          envList: [
            {
              STATUS_KEY: "chain-client-status-report-bsc",
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
