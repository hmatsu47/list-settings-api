// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// リリース設定の取得（タグ指定分）
	// (GET /tagSettings)
	GetTagSettings(c *gin.Context)
	// リリース設定の取得（URI指定分）
	// (GET /uriSettings)
	GetUriSettings(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetTagSettings operation middleware
func (siw *ServerInterfaceWrapper) GetTagSettings(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetTagSettings(c)
}

// GetUriSettings operation middleware
func (siw *ServerInterfaceWrapper) GetUriSettings(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetUriSettings(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router *gin.Engine, si ServerInterface) *gin.Engine {
	return RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router *gin.Engine, si ServerInterface, options GinServerOptions) *gin.Engine {

	errorHandler := options.ErrorHandler

	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/tagSettings", wrapper.GetTagSettings)

	router.GET(options.BaseURL+"/uriSettings", wrapper.GetUriSettings)

	return router
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7xWX2sbRxD/KmXbx7PvpNjS5Z5aimkFSSiO/RSMWZ3m7la+vT3v7kmnCEEs0SaFQl5K",
	"CyF9SFpa958plEJC3fTDXJXkLV+h7J5k3dWyfXZLweC91c7Mb34z+9sZIpfRmEUQSYGcIeIgYhYJ0B/A",
	"OeObsx214bJIQiTVEsdxSFwsCYvMrmCR2hNuABSr1TscPOSgt82FdzP/VZgbyisajUYG6oBwOYmVE+Sg",
	"bHyYTb7PJsfZ5Ek2mWTj39V6/Cyb/JSNn2eTr7LJr3rxOJvcz8Zfo5GBJPZvg5Qk8q+Ek0ig4iLAW9jX",
	"oQYxIAdhzvFgOfw/s/Ev0z+eTo8fZgdHrw9/nh49Og97wkl17CkNrwB9m5Mq0Lc3W5Vxj4wZCB08L6Yz",
	"rFzKp8rP5EdkoJizGLgkea9REAL7moEZWiE5iXwdkMN+Qjh0kHPn5OCOgSSRoTqZgzhJk7W74EpkoHRF",
	"SBaHxA80naSDHBTAXdxs2nWShCLSzm9gITchBCxUgH9m8vLxvVe/jbPJD/rvOBs/nxN0ViaEYh92E05O",
	"57Ic0lp/z7bdAW7YvT0+yzeHs4v1KY9xqlaogyWsSEJhkey5ngUORZDsufv1RruTF37OWSnratSF3UZI",
	"a/6gNsDdUOO8BencyZIeKHA2ffFxdvDkP2dO9KnVFHbs1ppdWmTu3xLXbzQAvEYUpPt9WiaumHM13mr7",
	"THr1u3ad2c1Uo1SKsuTOFPTjEv3GIWaCSMYHuxGmUJU70m6Kbj31wnXXqmlUEvuipCxV3Fhe0lvv1QLs",
	"BZ313E1RbJYbpXW7B74vkuvtdlKmV1FTjdaIDvpeh6SdNkmbOvJ23jlnadslOA2xkLu8IArniWzpKo0M",
	"FEF6YnyRbbGblLIC7xEXLlXIpM4G9rrr0RT31spkKkKqkQmpaAau3WDr0FdkKi8k8tj8RcKuPgsUk1Cp",
	"KMVSJGvNd321seoyigyUg0YfEs4GWCRv3VRnAiIwUm+dNpMyFo5p+kQGSVuZmXNP6PR7WqhXNv5GvyTq",
	"Dfnr2b3X33736vPD9z5qIQOFxIXZ0zkDcLO1VSWiGRIhV0T+AosVHMdmO2Rtk2IhgZs3Wu9v3Lq9gQp0",
	"lgyQgXrARY61tmqpgyyGCMcEOejaqrVqqbbCMtAdZS4mFf3tgzxfMmctenA0ffjF9MWXb44f5ALx8rP7",
	"06NH0wefvDn+FOmQXE8ILVXGD0BuFeIY5Xmublln9ePJOXPJRKVL4+EklBebl2dGPS0klGI+uFp6uSbd",
	"Uf8XzO8ot+ZifLoqodubrYvY3C4EuQqbS2a8/4fN07nNqUw4KVCp5IBEEniEQ+R4OBSgL79SIuDKYFi4",
	"So5phszFYcCEdK5ZloVGOyeeh1UG4oVOlEo6MoYXzqQL01IKo53lkrbfuU7jNnVtq2OrqffvAAAA///0",
	"309r6wwAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}