package main

import (
	"fmt"
	"github.com/nikolausreza131192/pos/entity"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/nikolausreza131192/pos/config"
	"github.com/nikolausreza131192/pos/controllers"
	"github.com/nikolausreza131192/pos/repository"
)

func main() {
	fmt.Println("Starting POS API....")

	// Enable line numbers in logging
	log.SetFlags(log.LstdFlags | log.Llongfile)

	// Init config
	conf := config.InitConfig()

	// Init database
	databases := initDatabase(conf)

	// Init JWT Token Library
	jwtTokenLib := initJWTToken(conf.Auth)

	// Init services
	services := initServices(conf, databases, jwtTokenLib)

	// Init routes
	router := mux.NewRouter()
	initRoutes(router, services)

	fmt.Printf("POS API is running on port %s\n", conf.Main.APIPort)
	srv := &http.Server{
		Addr:    ":" + conf.Main.APIPort,
		Handler: router,
	}
	log.Fatal(srv.ListenAndServe())
}

// Controllers struct define list of all available controllers
type Controllers struct {
	Item controllers.Item
	User controllers.User
	Auth controllers.Auth
}

func initServices(conf config.Config, dbs map[string]*sqlx.DB, jwtTokenLib entity.JWTToken) Controllers {
	fmt.Println("Init services...")
	// Init all repository
	itemRepo := repository.NewItem(repository.ItemRepoParam{
		DB: dbs["stone_work"],
	})
	userRepo := repository.NewUser(repository.UserRepoParam{
		DB: dbs["stone_work"],
	})
	authRepo := repository.NewAuth(repository.AuthRepoParam{
		SecretToken: conf.Auth.SecretToken,
		JWTTokenLib: jwtTokenLib,
	})

	// Init all controllers
	itemController := controllers.NewItem(controllers.ItemControllerParam{
		ItemRP: itemRepo,
	})
	userController := controllers.NewUser(controllers.UserControllerParam{
		UserRP: userRepo,
	})
	loginTime, err := strconv.Atoi(conf.Auth.LoginTime)
	if err != nil {
		log.Fatalln("func initServices fail to parse login time to int", err)
	}
	authController := controllers.NewAuth(controllers.AuthControllerParam{
		AuthRP:    authRepo,
		UserRP:    userRepo,
		LoginTime: loginTime,
	})

	return Controllers{
		Item: itemController,
		User: userController,
		Auth: authController,
	}
}

func initDatabase(conf config.Config) map[string]*sqlx.DB {
	fmt.Println("Init databases...")
	dbs := map[string]*sqlx.DB{}

	for key, db := range conf.Database {
		db, err := sqlx.Open(
			db.Driver,
			db.User+":"+db.Password+"@/"+db.Name+"?parseTime=true",
		)
		if err != nil {
			log.Fatal(err.Error())
		}

		dbs[key] = db
	}
	return dbs
}
