package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mediocregopher/radix.v2/pool"
)

var redisPool *pool.Pool

func getAppEnvComponentVersion(app string, env string, component string) string {
	version, _ := getRedisClient().Cmd("GET", app+":"+env+":"+component).Str()
	return version
}

func getAppEnv(app string, env string) map[string]string {
	result := map[string]string{}
	components, _ := getRedisClient().Cmd("SMEMBERS", app+":"+env).List()
	for _, component := range components {
		result[component] = getAppEnvComponentVersion(app, env, component)
	}
	return result
}

func getApp(app string) map[string]map[string]string {
	result := map[string]map[string]string{}
	envs, _ := getRedisClient().Cmd("SMEMBERS", app).List()
	for _, env := range envs {
		result[env] = getAppEnv(app, env)
	}
	return result
}

func setComponent(app string, env string, component string, version string) {
	client := getRedisClient()
	client.Cmd("SADD", app+":"+env, component)
	client.Cmd("SET", app+":"+env+":"+component, version)
}

func setEnv(app string, env string, components map[string]string) {
	getRedisClient().Cmd("SADD", app, env)
	for component, version := range components {
		setComponent(app, env, component, version)
	}
}

func setApp(app string, envs map[string]map[string]string) {
	getRedisClient().Cmd("SADD", "apps", app)
	for env, components := range envs {
		setEnv(app, env, components)
	}
}

func sendResult(w http.ResponseWriter, result interface{}) {
	body, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func handleGetAll(w http.ResponseWriter, r *http.Request) {
	state := map[string]map[string]map[string]string{}
	apps, _ := getRedisClient().Cmd("SMEMBERS", "apps").List()
	for _, app := range apps {
		state[app] = getApp(app)
	}
	sendResult(w, state)
}

func getEnvVariable(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
		fmt.Println("INFO: No " + key + " environment variable detected, defaulting to " + value)
	}
	return value
}

func getRedisClient() *pool.Pool {
	if redisPool == nil {
		p, err := pool.New("tcp", getEnvVariable("VMS_REDIS", "127.0.0.1:6379"), 20)
		if err != nil {
			log.Print(err.Error())
		} else {
			redisPool = p
		}
	}
	return redisPool
}

func getListeningPort() string {
	return ":" + getEnvVariable("VMS_PORT", "8080")
}

func main() {
	r := mux.NewRouter()

	r.Methods("GET").Path("/").HandlerFunc(handleGetAll)
	r.Methods("GET").Path("/{app}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		sendResult(w, getApp(params["app"]))
	})
	r.Methods("GET").Path("/{app}/{env}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		sendResult(w, getAppEnv(params["app"], params["env"]))
	})
	r.Methods("GET").Path("/{app}/{env}/{component}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		sendResult(w, getAppEnvComponentVersion(params["app"], params["env"], params["component"]))
	})

	r.Methods("POST").Path("/").Headers("content-type", "application/json").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]map[string]map[string]string
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&payload)
		if err != nil {
			fmt.Println(err)
		}
		for app, envs := range payload {
			setApp(app, envs)
		}
	})
	r.Methods("POST").Path("/{app}").Headers("content-type", "application/json").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]map[string]string
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&payload)
		for env, components := range payload {
			setEnv(mux.Vars(r)["app"], env, components)
		}
	})
	r.Methods("POST").Path("/{app}/{env}").Headers("content-type", "application/json").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		var payload map[string]string
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&payload)
		if err != nil {
			fmt.Println(err)
		}
		for component, version := range payload {
			setComponent(params["app"], params["env"], component, version)
		}
	})
	r.Methods("POST").Path("/{app}/{env}/{component}").Headers("content-type", "application/json").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		var version string
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&version)
		if err != nil {
			fmt.Println(err)
		}
		setComponent(params["app"], params["env"], params["component"], version)
	})

	http.Handle("/", r)

	fmt.Println("INFO: Connected to redis server on address", getRedisClient().Addr)
	port := getListeningPort()
	fmt.Println("INFO: Listening on", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
