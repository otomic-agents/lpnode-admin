package main

import (
	adminapiservice "admin-panel"
	accountcex "admin-panel/gen/account_cex"
	accountdex "admin-panel/gen/account_dex"
	ammordercenter "admin-panel/gen/amm_order_center"
	basedata "admin-panel/gen/base_data"
	bridgeconfig "admin-panel/gen/bridge_config"
	chainconfig "admin-panel/gen/chain_config"
	configresource "admin-panel/gen/config_resource"
	dexwallet "admin-panel/gen/dex_wallet"
	hedge "admin-panel/gen/hedge"
	installctrlpanel "admin-panel/gen/install_ctrl_panel"
	lpregister "admin-panel/gen/lp_register"
	lpmonit "admin-panel/gen/lpmonit"
	mainlogic "admin-panel/gen/main_logic"
	ordercenter "admin-panel/gen/order_center"
	relayaccount "admin-panel/gen/relay_account"
	statuslist "admin-panel/gen/status_list"
	taskmanager "admin-panel/gen/task_manager"
	tokenmanager "admin-panel/gen/token_manager"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		hostF     = flag.String("host", "0.0.0.0", "Server host (valid values: 0.0.0.0)")
		domainF   = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		httpPortF = flag.String("http-port", "", "HTTP port (overrides host HTTP port specified in service design)")
		secureF   = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF      = flag.Bool("debug", false, "Log request and response bodies")
	)
	flag.Parse()

	// Setup logger. Replace logger with your own log package of choice.
	var (
		logger *log.Logger
	)
	{
		logger = log.New(os.Stderr, "[adminapiservice] ", log.Ltime)
	}

	// Initialize the services.
	var (
		mainLogicSvc        mainlogic.Service
		accountCexSvc       accountcex.Service
		accountDexSvc       accountdex.Service
		ammOrderCenterSvc   ammordercenter.Service
		baseDataSvc         basedata.Service
		bridgeConfigSvc     bridgeconfig.Service
		chainConfigSvc      chainconfig.Service
		configResourceSvc   configresource.Service
		installCtrlPanelSvc installctrlpanel.Service
		dexWalletSvc        dexwallet.Service
		hedgeSvc            hedge.Service
		lpmonitSvc          lpmonit.Service
		orderCenterSvc      ordercenter.Service
		lpRegisterSvc       lpregister.Service
		relayAccountSvc     relayaccount.Service
		statusListSvc       statuslist.Service
		taskManagerSvc      taskmanager.Service
		tokenManagerSvc     tokenmanager.Service
	)
	{
		mainLogicSvc = adminapiservice.NewMainLogic(logger)
		accountCexSvc = adminapiservice.NewAccountCex(logger)
		accountDexSvc = adminapiservice.NewAccountDex(logger)
		ammOrderCenterSvc = adminapiservice.NewAmmOrderCenter(logger)
		baseDataSvc = adminapiservice.NewBaseData(logger)
		bridgeConfigSvc = adminapiservice.NewBridgeConfig(logger)
		chainConfigSvc = adminapiservice.NewChainConfig(logger)
		configResourceSvc = adminapiservice.NewConfigResource(logger)
		installCtrlPanelSvc = adminapiservice.NewInstallCtrlPanel(logger)
		dexWalletSvc = adminapiservice.NewDexWallet(logger)
		hedgeSvc = adminapiservice.NewHedge(logger)
		lpmonitSvc = adminapiservice.NewLpmonit(logger)
		orderCenterSvc = adminapiservice.NewOrderCenter(logger)
		lpRegisterSvc = adminapiservice.NewLpRegister(logger)
		relayAccountSvc = adminapiservice.NewRelayAccount(logger)
		statusListSvc = adminapiservice.NewStatusList(logger)
		taskManagerSvc = adminapiservice.NewTaskManager(logger)
		tokenManagerSvc = adminapiservice.NewTokenManager(logger)
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		mainLogicEndpoints        *mainlogic.Endpoints
		accountCexEndpoints       *accountcex.Endpoints
		accountDexEndpoints       *accountdex.Endpoints
		ammOrderCenterEndpoints   *ammordercenter.Endpoints
		baseDataEndpoints         *basedata.Endpoints
		bridgeConfigEndpoints     *bridgeconfig.Endpoints
		chainConfigEndpoints      *chainconfig.Endpoints
		configResourceEndpoints   *configresource.Endpoints
		installCtrlPanelEndpoints *installctrlpanel.Endpoints
		dexWalletEndpoints        *dexwallet.Endpoints
		hedgeEndpoints            *hedge.Endpoints
		lpmonitEndpoints          *lpmonit.Endpoints
		orderCenterEndpoints      *ordercenter.Endpoints
		lpRegisterEndpoints       *lpregister.Endpoints
		relayAccountEndpoints     *relayaccount.Endpoints
		statusListEndpoints       *statuslist.Endpoints
		taskManagerEndpoints      *taskmanager.Endpoints
		tokenManagerEndpoints     *tokenmanager.Endpoints
	)
	{
		mainLogicEndpoints = mainlogic.NewEndpoints(mainLogicSvc)
		accountCexEndpoints = accountcex.NewEndpoints(accountCexSvc)
		accountDexEndpoints = accountdex.NewEndpoints(accountDexSvc)
		ammOrderCenterEndpoints = ammordercenter.NewEndpoints(ammOrderCenterSvc)
		baseDataEndpoints = basedata.NewEndpoints(baseDataSvc)
		bridgeConfigEndpoints = bridgeconfig.NewEndpoints(bridgeConfigSvc)
		chainConfigEndpoints = chainconfig.NewEndpoints(chainConfigSvc)
		configResourceEndpoints = configresource.NewEndpoints(configResourceSvc)
		installCtrlPanelEndpoints = installctrlpanel.NewEndpoints(installCtrlPanelSvc)
		dexWalletEndpoints = dexwallet.NewEndpoints(dexWalletSvc)
		hedgeEndpoints = hedge.NewEndpoints(hedgeSvc)
		lpmonitEndpoints = lpmonit.NewEndpoints(lpmonitSvc)
		orderCenterEndpoints = ordercenter.NewEndpoints(orderCenterSvc)
		lpRegisterEndpoints = lpregister.NewEndpoints(lpRegisterSvc)
		relayAccountEndpoints = relayaccount.NewEndpoints(relayAccountSvc)
		statusListEndpoints = statuslist.NewEndpoints(statusListSvc)
		taskManagerEndpoints = taskmanager.NewEndpoints(taskManagerSvc)
		tokenManagerEndpoints = tokenmanager.NewEndpoints(tokenManagerSvc)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Start the servers and send errors (if any) to the error channel.
	switch *hostF {
	case "0.0.0.0":
		{
			addr := "http://0.0.0.0:18006"
			u, err := url.Parse(addr)
			if err != nil {
				logger.Fatalf("invalid URL %#v: %s\n", addr, err)
			}
			if *secureF {
				u.Scheme = "https"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *httpPortF != "" {
				h, _, err := net.SplitHostPort(u.Host)
				if err != nil {
					logger.Fatalf("invalid URL %#v: %s\n", u.Host, err)
				}
				u.Host = net.JoinHostPort(h, *httpPortF)
			} else if u.Port() == "" {
				u.Host = net.JoinHostPort(u.Host, "80")
			}
			handleHTTPServer(ctx, u, mainLogicEndpoints, accountCexEndpoints, accountDexEndpoints, ammOrderCenterEndpoints, baseDataEndpoints, bridgeConfigEndpoints, chainConfigEndpoints, configResourceEndpoints, installCtrlPanelEndpoints, dexWalletEndpoints, hedgeEndpoints, lpmonitEndpoints, orderCenterEndpoints, lpRegisterEndpoints, relayAccountEndpoints, statusListEndpoints, taskManagerEndpoints, tokenManagerEndpoints, &wg, errc, logger, *dbgF)
		}

	default:
		logger.Fatalf("invalid host argument: %q (valid hosts: 0.0.0.0)\n", *hostF)
	}

	// Wait for signal.
	logger.Printf("exiting (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	logger.Println("exited")
}
