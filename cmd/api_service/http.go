package main

import (
	accountcex "admin-panel/gen/account_cex"
	accountdex "admin-panel/gen/account_dex"
	ammordercenter "admin-panel/gen/amm_order_center"
	authenticationlimiter "admin-panel/gen/authentication_limiter"
	basedata "admin-panel/gen/base_data"
	bridgeconfig "admin-panel/gen/bridge_config"
	chainconfig "admin-panel/gen/chain_config"
	configresource "admin-panel/gen/config_resource"
	dexwallet "admin-panel/gen/dex_wallet"
	hedge "admin-panel/gen/hedge"
	accountcexsvr "admin-panel/gen/http/account_cex/server"
	accountdexsvr "admin-panel/gen/http/account_dex/server"
	ammordercentersvr "admin-panel/gen/http/amm_order_center/server"
	authenticationlimitersvr "admin-panel/gen/http/authentication_limiter/server"
	basedatasvr "admin-panel/gen/http/base_data/server"
	bridgeconfigsvr "admin-panel/gen/http/bridge_config/server"
	chainconfigsvr "admin-panel/gen/http/chain_config/server"
	configresourcesvr "admin-panel/gen/http/config_resource/server"
	dexwalletsvr "admin-panel/gen/http/dex_wallet/server"
	hedgesvr "admin-panel/gen/http/hedge/server"
	installctrlpanelsvr "admin-panel/gen/http/install_ctrl_panel/server"
	lpregistersvr "admin-panel/gen/http/lp_register/server"
	lpmonitsvr "admin-panel/gen/http/lpmonit/server"
	mainlogicsvr "admin-panel/gen/http/main_logic/server"
	ordercentersvr "admin-panel/gen/http/order_center/server"
	relayaccountsvr "admin-panel/gen/http/relay_account/server"
	statuslistsvr "admin-panel/gen/http/status_list/server"
	taskmanagersvr "admin-panel/gen/http/task_manager/server"
	tokenmanagersvr "admin-panel/gen/http/token_manager/server"
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
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	goahttp "goa.design/goa/v3/http"
	httpmdlwr "goa.design/goa/v3/http/middleware"
	"goa.design/goa/v3/middleware"
)

// handleHTTPServer starts configures and starts a HTTP server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleHTTPServer(ctx context.Context, u *url.URL, mainLogicEndpoints *mainlogic.Endpoints, accountCexEndpoints *accountcex.Endpoints, accountDexEndpoints *accountdex.Endpoints, ammOrderCenterEndpoints *ammordercenter.Endpoints, authenticationLimiterEndpoints *authenticationlimiter.Endpoints, baseDataEndpoints *basedata.Endpoints, bridgeConfigEndpoints *bridgeconfig.Endpoints, chainConfigEndpoints *chainconfig.Endpoints, configResourceEndpoints *configresource.Endpoints, dexWalletEndpoints *dexwallet.Endpoints, hedgeEndpoints *hedge.Endpoints, installCtrlPanelEndpoints *installctrlpanel.Endpoints, lpmonitEndpoints *lpmonit.Endpoints, orderCenterEndpoints *ordercenter.Endpoints, lpRegisterEndpoints *lpregister.Endpoints, relayAccountEndpoints *relayaccount.Endpoints, statusListEndpoints *statuslist.Endpoints, taskManagerEndpoints *taskmanager.Endpoints, tokenManagerEndpoints *tokenmanager.Endpoints, wg *sync.WaitGroup, errc chan error, logger *log.Logger, debug bool) {

	// Setup goa log adapter.
	var (
		adapter middleware.Logger
	)
	{
		adapter = middleware.NewLogger(logger)
	}

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and configure it to serve
	// HTTP requests to the service endpoints.
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
	}

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.
	var (
		mainLogicServer             *mainlogicsvr.Server
		accountCexServer            *accountcexsvr.Server
		accountDexServer            *accountdexsvr.Server
		ammOrderCenterServer        *ammordercentersvr.Server
		authenticationLimiterServer *authenticationlimitersvr.Server
		baseDataServer              *basedatasvr.Server
		bridgeConfigServer          *bridgeconfigsvr.Server
		chainConfigServer           *chainconfigsvr.Server
		configResourceServer        *configresourcesvr.Server
		dexWalletServer             *dexwalletsvr.Server
		hedgeServer                 *hedgesvr.Server
		installCtrlPanelServer      *installctrlpanelsvr.Server
		lpmonitServer               *lpmonitsvr.Server
		orderCenterServer           *ordercentersvr.Server
		lpRegisterServer            *lpregistersvr.Server
		relayAccountServer          *relayaccountsvr.Server
		statusListServer            *statuslistsvr.Server
		taskManagerServer           *taskmanagersvr.Server
		tokenManagerServer          *tokenmanagersvr.Server
	)
	{
		eh := errorHandler(logger)
		mainLogicServer = mainlogicsvr.New(mainLogicEndpoints, mux, dec, enc, eh, nil)
		accountCexServer = accountcexsvr.New(accountCexEndpoints, mux, dec, enc, eh, nil)
		accountDexServer = accountdexsvr.New(accountDexEndpoints, mux, dec, enc, eh, nil)
		ammOrderCenterServer = ammordercentersvr.New(ammOrderCenterEndpoints, mux, dec, enc, eh, nil)
		authenticationLimiterServer = authenticationlimitersvr.New(authenticationLimiterEndpoints, mux, dec, enc, eh, nil)
		baseDataServer = basedatasvr.New(baseDataEndpoints, mux, dec, enc, eh, nil)
		bridgeConfigServer = bridgeconfigsvr.New(bridgeConfigEndpoints, mux, dec, enc, eh, nil)
		chainConfigServer = chainconfigsvr.New(chainConfigEndpoints, mux, dec, enc, eh, nil)
		configResourceServer = configresourcesvr.New(configResourceEndpoints, mux, dec, enc, eh, nil)
		dexWalletServer = dexwalletsvr.New(dexWalletEndpoints, mux, dec, enc, eh, nil)
		hedgeServer = hedgesvr.New(hedgeEndpoints, mux, dec, enc, eh, nil)
		installCtrlPanelServer = installctrlpanelsvr.New(installCtrlPanelEndpoints, mux, dec, enc, eh, nil)
		lpmonitServer = lpmonitsvr.New(lpmonitEndpoints, mux, dec, enc, eh, nil)
		orderCenterServer = ordercentersvr.New(orderCenterEndpoints, mux, dec, enc, eh, nil)
		lpRegisterServer = lpregistersvr.New(lpRegisterEndpoints, mux, dec, enc, eh, nil)
		relayAccountServer = relayaccountsvr.New(relayAccountEndpoints, mux, dec, enc, eh, nil)
		statusListServer = statuslistsvr.New(statusListEndpoints, mux, dec, enc, eh, nil)
		taskManagerServer = taskmanagersvr.New(taskManagerEndpoints, mux, dec, enc, eh, nil)
		tokenManagerServer = tokenmanagersvr.New(tokenManagerEndpoints, mux, dec, enc, eh, nil)
		if debug {
			servers := goahttp.Servers{
				mainLogicServer,
				accountCexServer,
				accountDexServer,
				ammOrderCenterServer,
				authenticationLimiterServer,
				baseDataServer,
				bridgeConfigServer,
				chainConfigServer,
				configResourceServer,
				dexWalletServer,
				hedgeServer,
				installCtrlPanelServer,
				lpmonitServer,
				orderCenterServer,
				lpRegisterServer,
				relayAccountServer,
				statusListServer,
				taskManagerServer,
				tokenManagerServer,
			}
			servers.Use(httpmdlwr.Debug(mux, os.Stdout))
		}
	}
	// Configure the mux.
	mainlogicsvr.Mount(mux, mainLogicServer)
	accountcexsvr.Mount(mux, accountCexServer)
	accountdexsvr.Mount(mux, accountDexServer)
	ammordercentersvr.Mount(mux, ammOrderCenterServer)
	authenticationlimitersvr.Mount(mux, authenticationLimiterServer)
	basedatasvr.Mount(mux, baseDataServer)
	bridgeconfigsvr.Mount(mux, bridgeConfigServer)
	chainconfigsvr.Mount(mux, chainConfigServer)
	configresourcesvr.Mount(mux, configResourceServer)
	dexwalletsvr.Mount(mux, dexWalletServer)
	hedgesvr.Mount(mux, hedgeServer)
	installctrlpanelsvr.Mount(mux, installCtrlPanelServer)
	lpmonitsvr.Mount(mux, lpmonitServer)
	ordercentersvr.Mount(mux, orderCenterServer)
	lpregistersvr.Mount(mux, lpRegisterServer)
	relayaccountsvr.Mount(mux, relayAccountServer)
	statuslistsvr.Mount(mux, statusListServer)
	taskmanagersvr.Mount(mux, taskManagerServer)
	tokenmanagersvr.Mount(mux, tokenManagerServer)

	// Wrap the multiplexer with additional middlewares. Middlewares mounted
	// here apply to all the service endpoints.
	var handler http.Handler = mux
	{
		handler = httpmdlwr.Log(adapter)(handler)
		handler = httpmdlwr.RequestID()(handler)
	}

	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{Addr: u.Host, Handler: handler, ReadHeaderTimeout: time.Second * 60}
	for _, m := range mainLogicServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range accountCexServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range accountDexServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range ammOrderCenterServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range authenticationLimiterServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range baseDataServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range bridgeConfigServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range chainConfigServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range configResourceServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range dexWalletServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range hedgeServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range installCtrlPanelServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range lpmonitServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range orderCenterServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range lpRegisterServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range relayAccountServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range statusListServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range taskManagerServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range tokenManagerServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}

	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start HTTP server in a separate goroutine.
		go func() {
			logger.Printf("HTTP server listening on %q", u.Host)
			errc <- srv.ListenAndServe()
		}()

		<-ctx.Done()
		logger.Printf("shutting down HTTP server at %q", u.Host)

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			logger.Printf("failed to shutdown: %v", err)
		}
	}()
}

// errorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func errorHandler(logger *log.Logger) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(middleware.RequestIDKey).(string)
		_, _ = w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		logger.Printf("[%s] ERROR: %s", id, err.Error())
	}
}
