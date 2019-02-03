package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/balance-transfer-go/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"

	"github.com/mux"
	"github.com/op/go-logging"
)

var hfc utils.FabricSetup
var log = logging.MustGetLogger("Main")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)
var goPath string

/**
Basic Flow:
1) Prepare client context
2) Create resource managememt client
3) Create new channel
4) Peer(s) join channel
5) Install chaincode onto peer(s) filesystem
6) Instantiate chaincode on channel
7) Query peer for channels, installed/instantiated chaincodes etc.
**/

func main() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.DEBUG, "")
	logging.SetBackend(backendLeveled, backendFormatter)

	log.Debug("<<<<<<<<<<<<< STARTING MAIN >>>>>>>>>>>>>")
	goPath = "/home/anshu/go"
	hfc = utils.FabricSetup{
		AdminUser:         "Admin",
		OrdererOrgName:    "ordererorg",
		ConfigFileName:    "config.yaml",
		Secret:            []byte("thisismysecret"),
		IdentityTypeUser:  "user",
		RegistrarUsername: "admin",
		RegistrarPassword: "adminpw",
		ChannelID:         "mychannel",
	}
	hfc.Init()
	r := mux.NewRouter()
	r.HandleFunc("/users", login).Methods("POST")
	r.Handle("/channels", authMiddleware(http.HandlerFunc(channelCC))).Methods("POST")
	r.Handle("/channels/{channelName}", authMiddleware(http.HandlerFunc(joinCC))).Methods("POST")
	r.Handle("/chaincodes", authMiddleware(http.HandlerFunc(installCC))).Methods("POST")
	r.Handle("/channels/{channelName}/instantiate", authMiddleware(http.HandlerFunc(instantiateCC))).Methods("POST")
	r.Handle("/channels/{channelName}/chaincodes/{chaincodeName}", authMiddleware(http.HandlerFunc(invokeCC))).Methods("POST")
	r.Handle("/channels/{channelName}/chaincodes/{chaincodeName}", authMiddleware(http.HandlerFunc(queryCC))).Methods("GET")
	r.Handle("/channels/{channelName}/blocks", authMiddleware(http.HandlerFunc(queryBlkHeight))).Methods("GET")
	r.Handle("/channels/{channelName}/blockByID", authMiddleware(http.HandlerFunc(queryBlkByID))).Methods("GET")
	http.ListenAndServe(":4000", handlers.LoggingHandler(os.Stdout, r))
}

func queryBlkByID(w http.ResponseWriter, r *http.Request) {
	log.Debug("==================== Q U E R Y   B L O C K  B Y  I D ====================")
	type queryBlkCC struct {
		Peer    string
		BlockID string
	}
	type response struct {
		res []byte
	}
	decoder := json.NewDecoder(r.Body)
	body := queryBlkCC{}
	decoder.Decode(&body)

	vars := mux.Vars(r)
	username := r.Header.Get("username")
	orgName := r.Header.Get("orgName")

	blkID, _ := strconv.ParseUint(body.BlockID, 10, 64)
	res := hfc.QueryBlockByID(vars["channelName"], username, orgName, body.Peer, blkID)
	w.Write(res)
}

func queryBlkHeight(w http.ResponseWriter, r *http.Request) {
	log.Debug("==================== Q U E R Y   B L O C K S ====================")
	type queryBlkCC struct {
		Peer string
	}
	type response struct {
		res []byte
	}
	decoder := json.NewDecoder(r.Body)
	body := queryBlkCC{}
	decoder.Decode(&body)

	vars := mux.Vars(r)
	username := r.Header.Get("username")
	orgName := r.Header.Get("orgName")

	res := hfc.QueryBlockHeight(vars["channelName"], username, orgName, body.Peer)
	w.Write(res)

}

func instantiateCC(w http.ResponseWriter, r *http.Request) {
	log.Debug("==================== I N S T A N T I A T I N G  C H A I N C O D E ====================")
	type instantiateCCBody struct {
		Peers            string
		ChainCodeName    string
		Args             []string
		ChainCodeType    string
		ChainCodeVersion string
		ChainCodePath    string
	}
	type response struct {
		res []byte
	}
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	orgName := r.Header.Get("orgName")
	username := r.Header.Get("username")
	body := instantiateCCBody{}
	decoder.Decode(&body)

	res := hfc.InstantiateChainCode(username, orgName, vars["channelName"], body.ChainCodeName, body.ChainCodePath, body.ChainCodeVersion, utils.GetArgs(body.Args))
	w.Write(res)
}

func installCC(w http.ResponseWriter, r *http.Request) {
	log.Debugf("$$$$$$$$$$$$$$$$$$ I N S T A L L I N G  C H A I N C O D E  $$$$$$$$$$$$$$$$$$")
	type installCCBody struct {
		Peers            string
		ChainCodeName    string
		ChainCodePath    string
		ChainCodeType    string
		ChainCodeVersion string
	}
	type response struct {
		res []byte
	}
	decoder := json.NewDecoder(r.Body)
	orgName := r.Header.Get("orgName")
	body := installCCBody{}
	decoder.Decode(&body)

	res := hfc.InstallChainCode(body.Peers, body.ChainCodeName, body.ChainCodePath, body.ChainCodeType, body.ChainCodeVersion, orgName, goPath)
	w.Write(res)
}

func joinCC(w http.ResponseWriter, r *http.Request) {
	type joinCCBody struct {
		Peers string
	}
	type response struct {
		res []byte
	}
	vars := mux.Vars(r)
	orgName := r.Header.Get("orgName")
	decoder := json.NewDecoder(r.Body)
	body := joinCCBody{}
	decoder.Decode(&body)

	res := hfc.JoinChannel(vars["channelName"], body.Peers, orgName)
	w.Write(res)
}

func channelCC(w http.ResponseWriter, r *http.Request) {
	log.Debug("Trying to create channel.................")

	type channelCreateBody struct {
		ChannelName       string
		ChannelConfigPath string
	}
	type response struct {
		res []byte
	}
	username := r.Header.Get("username")
	orgName := r.Header.Get("orgName")
	//log.Infof("Username is %s and Orgname is %s", username, orgName)
	decoder := json.NewDecoder(r.Body)
	body := channelCreateBody{}
	decoder.Decode(&body)

	res := hfc.CreateChannel(username, orgName, body.ChannelName, body.ChannelConfigPath)
	w.Write(res)
}

func queryCC(w http.ResponseWriter, r *http.Request) {
	log.Debugf("--------------- Q U E R Y  C H A I N C O D E ---------------")

	type queryCcBody struct {
		Fcn  string
		Args []string
		Peer string
	}
	type response struct {
		res []byte
	}
	vars := mux.Vars(r)
	username := r.Header.Get("username")
	orgName := r.Header.Get("orgName")
	decoder := json.NewDecoder(r.Body)
	body := queryCcBody{}
	decoder.Decode(&body)

	res := hfc.QueryChainCode(vars["channelName"], vars["chaincodeName"], username, orgName, body.Fcn, utils.GetArgs(body.Args), body.Peer)
	w.Write(res)
}

func invokeCC(w http.ResponseWriter, r *http.Request) {
	log.Debug("<<<<<<<<<<<<< INVOKING CHAINCODE >>>>>>>>>>>>>")
	type invokeCCBody struct {
		Peers []string
		Fcn   string
		Args  []string
	}
	type response struct {
		TxID string
	}
	vars := mux.Vars(r)
	username := r.Header.Get("username")
	orgName := r.Header.Get("orgName")
	decoder := json.NewDecoder(r.Body)
	body := invokeCCBody{}
	decoder.Decode(&body)

	//	client := utils.GetClient(hfc.Sdk, vars["channelName"], username, orgName)
	res := hfc.InvokeChainCode(vars["channelName"], vars["chaincodeName"], username, orgName, body.Fcn, utils.GetArgs(body.Args), body.Peers)
	w.Write(res)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("authorization")
		if tokenString != "" {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return hfc.Secret, nil
			})
			if err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					r.Header.Add("username", claims["username"].(string))
					r.Header.Add("orgName", claims["orgName"].(string))
					next.ServeHTTP(w, r)
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	log.Infof("==================== LOGIN ==================")
	// define response
	type response struct {
		Success bool
		Message string
		Token   string
	}

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	username := r.Form.Get("username")
	orgName := r.Form.Get("orgName")
	secret := r.Form.Get("secret")
	if username != "" && orgName != "" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": username,
			"orgName":  orgName,
			"exp":      time.Now().Unix() + 360000,
		})
		tokenString, err := token.SignedString(hfc.Secret)
		log.Infof(tokenString, err)
		message, success := utils.GetRegisteredUser(strings.ToLower(username), strings.ToLower(orgName), secret, hfc.IdentityTypeUser, hfc.Sdk)
		res := response{
			Success: success,
			Message: message,
		}

		if success {
			res.Token = tokenString
		}
		out, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
