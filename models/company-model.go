package models

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_api_master/configs"
	"github.com/nikitamirzani323/wl_api_master/db"
	"github.com/nikitamirzani323/wl_api_master/entities"
	"github.com/nikitamirzani323/wl_api_master/helpers"
	"github.com/nleeper/goment"
)

func Fetch_companyHome() (helpers.ResponseCompany, error) {
	var obj entities.Model_company
	var arraobj []entities.Model_company
	var res helpers.ResponseCompany
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idcompany , to_char(COALESCE(startjoincompany,now()), 'YYYY-MM-DD HH24:MI:SS'), to_char(COALESCE(endjoincompany,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			idcurr , nmcompany, nmowner, phoneowner, emailowner, companyurl, statuscompany, 
			createcompany, to_char(COALESCE(createdatecompany,now()), 'YYYY-MM-DD HH24:MI:SS') as createdatecompany, 
			updatecompany, to_char(COALESCE(updatedatecompany,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedatecompany 
			FROM ` + configs.DB_tbl_mst_company + ` 
			ORDER BY createcompany DESC 
		`

	row, err := con.QueryContext(ctx, sql_select)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcompany_db, startjoincompany_db, endjoincompany_db                                               string
			idcurr_db, nmcompany_db, nmowner_db, phoneowner_db, emailowner_db, companyurl_db, statuscompany_db string
			createcompany_db, createdatecompany_db, updatecompany_db, updatedatecompany_db                     string
		)

		err = row.Scan(
			&idcompany_db, &startjoincompany_db, &endjoincompany_db,
			&idcurr_db, &nmcompany_db, &nmowner_db, &phoneowner_db, &emailowner_db, &companyurl_db, &statuscompany_db,
			&createcompany_db, &createdatecompany_db, &updatecompany_db, &updatedatecompany_db)

		helpers.ErrorCheck(err)
		if startjoincompany_db == "0000-00-00 00:00:00" {
			startjoincompany_db = ""
		}
		if endjoincompany_db == "0000-00-00 00:00:00" {
			endjoincompany_db = ""
		}
		create := createcompany_db + ", " + createdatecompany_db
		update := ""
		if updatecompany_db != "" {
			update = updatecompany_db + ", " + updatedatecompany_db
		}
		obj.Company_idcomp = idcompany_db
		obj.Company_startjoin = startjoincompany_db
		obj.Company_endjoin = endjoincompany_db
		obj.Company_idcurr = idcurr_db
		obj.Company_nmcompany = nmcompany_db
		obj.Company_nmowner = nmowner_db
		obj.Company_phoneowner = phoneowner_db
		obj.Company_emailowner = emailowner_db
		obj.Company_urlendpoint = companyurl_db
		obj.Company_status = statuscompany_db
		obj.Company_create = create
		obj.Company_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	var objCurr entities.Model_compcurr
	var arraobjCurr []entities.Model_compcurr
	sql_listcurr := `SELECT 
		idcurr 	
		FROM ` + configs.DB_tbl_mst_currency + ` 
	`
	row_listcurr, err_listcurr := con.QueryContext(ctx, sql_listcurr)
	helpers.ErrorCheck(err_listcurr)
	for row_listcurr.Next() {
		var (
			idcurr_db string
		)

		err = row_listcurr.Scan(&idcurr_db)

		helpers.ErrorCheck(err)

		objCurr.Curr_idcurr = idcurr_db
		arraobjCurr = append(arraobjCurr, objCurr)
		msg = "Success"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Listcurr = arraobjCurr
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_companyHome(
	admin, idcompany, idcurr, nmcompany,
	nmowner, phoneowner, emailowner, companyurl,
	status, sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_company, "idcompany", idcompany)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_company + ` (
					idcompany , startjoincompany, idcurr, nmcompany, nmowner, 
					phoneowner, emailowner, companyurl, statuscompany
					createcompany, createdatecompany
				) values (
					$1, $2, $3, $4, $5,  
					$6, $7, $8, $9, 
					$10, $11 
				)
			`
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_company, "INSERT",
				idcompany, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idcurr, nmcompany, nmowner, phoneowner, emailowner, companyurl, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				flag = true
				msg = "Succes"
				log.Println(msg_insert)
			} else {
				log.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_mst_company + `   
				SET nmcompany=$1, nmowner=$2, phoneowner=$3, emailowner=$4,  
				companyurl =$5, statuscompany=$6,  
				updatecompany=$7, updatedatecompany=$8 
				WHERE idcompany =$9 
			`
		flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_mst_company, "UPDATE",
			nmcompany, nmowner, phoneowner, emailowner, companyurl, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idcompany)

		if flag_update {
			flag = true
			msg = "Succes"
			log.Println(msg_update)
		} else {
			log.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
