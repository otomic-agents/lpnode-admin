# otmoic-lpnode-admin
## Introduction
otmoic-lpnode-admin is the backend application of the otmoic application, which implements the backend API that the frontend depends on. It is installed together with the Lpnode and Dashborad programs as the basic system, and completes the initialization of basic data.

## Main Functions
* Manage users' wallets and private keys based on secretVault;
* Manage which trading pairs provide liquidity;
* Manage the installation of Amm and Exchange-adapter programs;
* Set cross-chain transaction fees and hedging logic;
* Manage the configuration of Amm and Exchange-adapter;
* Upload your custom scripts to complete scheduled functions such as system monitoring and data analysis;
* Maintain and manage the data status of Lpnode and ChainClient, where ChainClient is used for interacting with the blockchain;
* Register your Lp Account with Relay;
* View the running status of various components of the system.