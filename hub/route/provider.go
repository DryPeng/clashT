package route

import (
	"context"
	"net/http"

	C "github.com/DryPeng/clashT/constant"
	"github.com/DryPeng/clashT/constant/provider"
	"github.com/DryPeng/clashT/tunnel"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/samber/lo"
)

func proxyProviderRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getProviders)

	r.Route("/{providerName}", func(r chi.Router) {
		r.Use(parseProviderName, findProviderByName)
		r.Get("/", getProvider)
		r.Put("/", updateProvider)
		r.Get("/healthcheck", healthCheckProvider)
		r.Mount("/", proxyProviderProxyRouter())
	})
	return r
}

func proxyProviderProxyRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/{name}", func(r chi.Router) {
		r.Use(parseProxyName, findProviderProxyByName)
		r.Get("/", getProxy)
		r.Get("/healthcheck", getProxyDelay)
	})
	return r
}

func getProviders(w http.ResponseWriter, r *http.Request) {
	providers := tunnel.Providers()
	render.JSON(w, r, render.M{
		"providers": providers,
	})
}

func getProvider(w http.ResponseWriter, r *http.Request) {
	provider := r.Context().Value(CtxKeyProvider).(provider.ProxyProvider)
	render.JSON(w, r, provider)
}

func updateProvider(w http.ResponseWriter, r *http.Request) {
	provider := r.Context().Value(CtxKeyProvider).(provider.ProxyProvider)
	if err := provider.Update(); err != nil {
		render.Status(r, http.StatusServiceUnavailable)
		render.JSON(w, r, newError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func healthCheckProvider(w http.ResponseWriter, r *http.Request) {
	provider := r.Context().Value(CtxKeyProvider).(provider.ProxyProvider)
	provider.HealthCheck()
	render.NoContent(w, r)
}

func parseProviderName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := getEscapeParam(r, "providerName")
		ctx := context.WithValue(r.Context(), CtxKeyProviderName, name)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func findProviderByName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.Context().Value(CtxKeyProviderName).(string)
		providers := tunnel.Providers()
		provider, exist := providers[name]
		if !exist {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyProvider, provider)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func findProviderProxyByName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			name = r.Context().Value(CtxKeyProxyName).(string)
			pd   = r.Context().Value(CtxKeyProvider).(provider.ProxyProvider)
		)
		proxy, exist := lo.Find(pd.Proxies(), func(proxy C.Proxy) bool {
			return proxy.Name() == name
		})

		if !exist {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyProxy, proxy)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
