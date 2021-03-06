package gole

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/iamshubha/golang-postgresql/pkg/model"
	"github.com/iamshubha/golang-postgresql/pkg/util"
)

func StartWorking(w http.ResponseWriter, r *http.Request) {
	golemodel := model.GoleCreateMOdel{}
	db := util.GetDB()
	defer db.Close()
	err := json.NewDecoder(r.Body).Decode(&golemodel)
	if err != nil {
		fmt.Println(err)
		return
	}
	if golemodel.Userid == 0 || golemodel.Workon == "" {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Please send correct parameaters",
		})
		return
	}
	sqlQuery := `
	INSERT INTO goletable (userid, workon, starttime)
	VALUES ($1,$2,$3)  
	RETURNING id;
	`
	_, err = db.Exec(sqlQuery, golemodel.Userid, golemodel.Workon, time.Now())
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Fail",
		})
		log.Print(err)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success",
	})
}

func StopWorking(w http.ResponseWriter, r *http.Request) {

	golemodel := model.IdAndUserid{}
	db := util.GetDB()
	defer db.Close()
	err := json.NewDecoder(r.Body).Decode(&golemodel)
	if err != nil {
		fmt.Println(err)
		return
	}
	if golemodel.UserId == 0 || golemodel.Id == 0 {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Please send correct parameaters",
		})
		return
	}
	sqlQuery := `
	UPDATE goletable SET stoptime = $3 
	WHERE userid = $1 AND id = $2;
	`
	_, err = db.Exec(sqlQuery, golemodel.UserId, golemodel.Id, time.Now())
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Fail",
		})
		log.Print(err)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success",
	})
}

func GetWorkiteams(w http.ResponseWriter, r *http.Request) {

	urlData := mux.Vars(r)
	id, ok := urlData["id"]
	if !ok {
		log.Println(ok)
	}
	db := util.GetDB()
	defer db.Close()
	sqlQuery := `
	SELECT * FROM goletable WHERE userid =$1;
	`
	dataRows, err := db.Query(sqlQuery, id)
	if err != nil {
		log.Println(err)
	}
	defer dataRows.Close()
	data := make([]model.GoleDetails, 0)
	for dataRows.Next() {
		goleDetailsModel := model.GoleDetails{}
		dataRows.Scan(&goleDetailsModel.Id, &goleDetailsModel.Userid, &goleDetailsModel.Workon, &goleDetailsModel.Starttime, &goleDetailsModel.Stoptime, &goleDetailsModel.Total)
		data = append(data, goleDetailsModel)
	}
	log.Println(data)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "data found",
		"data":    data,
	})
}

func DeleteWorking(w http.ResponseWriter, r *http.Request) {
	userdetails := model.IdAndUserid{}
	err := json.NewDecoder(r.Body).Decode(&userdetails)
	if err != nil {
		log.Println(err)
		return
	}

	sqlQuery := `
	DELETE FROM goletable
	WHERE  userid = $1 AND id = $2;
	`

	db := util.GetDB()
	defer db.Close()
	_, err = db.Exec(sqlQuery, userdetails.UserId, userdetails.Id)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Success",
	})
}
