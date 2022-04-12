package models

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_api_master/configs"
	"github.com/nikitamirzani323/wl_api_master/db"
	"github.com/nikitamirzani323/wl_api_master/entities"
	"github.com/nikitamirzani323/wl_api_master/helpers"
	"github.com/nleeper/goment"
)

func Fetch_crm(search string, page int) (helpers.Responsemovie, error) {
	var obj entities.Model_crm
	var arraobj []entities.Model_crm
	var res helpers.Responsemovie
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	perpage := 250
	totalrecord := 0
	offset := page
	sql_selectcount := ""
	sql_selectcount += ""
	sql_selectcount += "SELECT "
	sql_selectcount += "COUNT(idusersales) as totalmember  "
	sql_selectcount += "FROM " + configs.DB_tbl_trx_usersales + "  "
	if search != "" {
		sql_selectcount += "WHERE LOWER(phone) LIKE '%" + strings.ToLower(search) + "%' "
		sql_selectcount += "OR LOWER(nama) LIKE '%" + strings.ToLower(search) + "%' "
	}

	row_selectcount := con.QueryRowContext(ctx, sql_selectcount)
	switch e_selectcount := row_selectcount.Scan(&totalrecord); e_selectcount {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e_selectcount)
	}

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "idusersales , phone, nama, "
	sql_select += "source , statususersales,  "
	sql_select += "createusersales, to_char(COALESCE(createdateusersales,NOW()), 'YYYY-MM-DD HH24:MI:SS'), updateusersales, to_char(COALESCE(updatedateusersales,NOW()) , 'YYYY-MM-DD HH24:MI:SS')"
	sql_select += "FROM " + configs.DB_tbl_trx_usersales + "  "
	if search == "" {
		sql_select += "ORDER BY createdateusersales DESC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)
	} else {
		sql_select += "WHERE LOWER(name) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "OR LOWER(phone) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY createdateusersales DESC  LIMIT " + strconv.Itoa(perpage)
	}

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idusersales_db                                                                         int
			phone_db, nama_db, source_db, statususersales_db                                       string
			createusersales_db, createdateusersales_db, updateusersales_db, updatedateusersales_db string
		)

		err = row.Scan(
			&idusersales_db, &phone_db, &nama_db, &source_db, &statususersales_db, &createusersales_db,
			&createdateusersales_db, &updateusersales_db, &updatedateusersales_db)

		helpers.ErrorCheck(err)

		create := ""
		update := ""
		statuscss := ""
		if createusersales_db != "" {
			create = createusersales_db + ", " + createdateusersales_db
		}
		if updateusersales_db != "" {
			update = updateusersales_db + ", " + updatedateusersales_db
		}
		switch statususersales_db {
		case "NEW":
			statuscss = configs.STATUS_NEW
		case "VALID":
			statuscss = configs.STATUS_COMPLETE
		case "INVALID":
			statuscss = configs.STATUS_CANCEL
		}
		obj.Crm_id = idusersales_db
		obj.Crm_phone = phone_db
		obj.Crm_name = nama_db
		obj.Crm_source = source_db
		obj.Crm_status = statususersales_db
		obj.Crm_statuscss = statuscss
		obj.Crm_create = create
		obj.Crm_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Perpage = perpage
	res.Totalrecord = totalrecord
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_crm(admin, phone, nama, status, sData string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_trx_usersales, "phone", phone)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_trx_usersales + ` (
					idusersales , phone, nama, source, statususersales, 
					createusersales, createdateusersales 
				) values (
					$1, $2, $3, $4, $5,
					$6, $7
				)
			`
			field_column := configs.DB_tbl_trx_usersales + tglnow.Format("YYYY")
			idrecord_counter := Get_counter(field_column)
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_usersales, "INSERT",
				tglnow.Format("YY")+strconv.Itoa(idrecord_counter), phone, nama, "", status,
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
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_trx_usersales + `  
				SET statususersales=$1, 
				updateusersales=$2, updatedateusersales=$3 
				WHERE idusersales =$4 
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_trx_usersales, "UPDATE",
			status, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			flag = true
			msg = "Succes"
			log.Println(msg_update)
		} else {
			log.Println(msg_update)
		}
	}

	if flag {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	}

	return res, nil
}
func Save_crmsource(admin, datasource, source, sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		json := []byte(datasource)
		jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			crmisbtv_username, _ := jsonparser.GetString(value, "crmisbtv_username")
			crmisbtv_name, _ := jsonparser.GetString(value, "crmisbtv_name")

			flag_check := CheckDB(configs.DB_tbl_trx_usersales, "phone", crmisbtv_username)
			log.Printf("%s - %s  - %t", crmisbtv_username, crmisbtv_name, flag_check)
			if !flag_check {
				sql_insert := `
					insert into
					` + configs.DB_tbl_trx_usersales + ` (
						idusersales , phone, nama, source, statususersales, 
						createusersales, createdateusersales 
					) values (
						$1, $2, $3, $4, $5,
						$6, $7
					)
				`
				field_column := configs.DB_tbl_trx_usersales + tglnow.Format("YYYY")
				idrecord_counter := Get_counter(field_column)
				flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_usersales, "INSERT",
					tglnow.Format("YY")+strconv.Itoa(idrecord_counter), crmisbtv_username, crmisbtv_name, source, "NEW",
					admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

				if flag_insert {
					flag = true
					msg = "Succes"
					log.Println(msg_insert)
				} else {
					log.Println(msg_insert)
				}
			}
		})
	}

	if flag {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	}

	return res, nil
}
func Fetch_crmisbtv(search string, page int) (helpers.Responsemovie, error) {
	var obj entities.Model_crmisbtv
	var arraobj []entities.Model_crmisbtv
	var res helpers.Responsemovie
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	perpage := 250
	totalrecord := 0
	offset := page
	sql_selectcount := ""
	sql_selectcount += ""
	sql_selectcount += "SELECT "
	sql_selectcount += "COUNT(username) as totalmember  "
	sql_selectcount += "FROM " + configs.DB_tbl_mst_user + "  "
	if search != "" {
		sql_selectcount += "WHERE LOWER(username) LIKE '%" + strings.ToLower(search) + "%' "
		sql_selectcount += "OR LOWER(nmuser) LIKE '%" + strings.ToLower(search) + "%' "
	}

	row_selectcount := con.QueryRowContext(ctx, sql_selectcount)
	switch e_selectcount := row_selectcount.Scan(&totalrecord); e_selectcount {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e_selectcount)
	}

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "username , nmuser, coderef, "
	sql_select += "point_in , point_out, statususer,  "
	sql_select += "to_char(COALESCE(lastlogin,NOW()), 'YYYY-MM-DD HH24:MI:SS'), to_char(COALESCE(createdateuser,NOW()), 'YYYY-MM-DD HH24:MI:SS'),  to_char(COALESCE(updatedateuser,NOW()), 'YYYY-MM-DD HH24:MI:SS') "
	sql_select += "FROM " + configs.DB_tbl_mst_user + "  "
	if search == "" {
		sql_select += "ORDER BY createdateuser DESC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)
	} else {
		sql_select += "WHERE LOWER(username) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "OR LOWER(nmuser) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY createdateuser DESC  LIMIT " + strconv.Itoa(perpage)
	}

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			point_in_db, point_out_db                          int
			username_db, nmuser_db, coderef_db, statususer_db  string
			lastlogin_db, createdateuser_db, updatedateuser_db string
		)

		err = row.Scan(
			&username_db, &nmuser_db, &coderef_db, &point_in_db, &point_out_db, &statususer_db,
			&lastlogin_db, &createdateuser_db, &updatedateuser_db)

		helpers.ErrorCheck(err)

		obj.Crmisbtv_username = username_db
		obj.Crmisbtv_name = nmuser_db
		obj.Crmisbtv_coderef = coderef_db
		obj.Crmisbtv_point = point_in_db - point_out_db
		obj.Crmisbtv_status = statususer_db
		obj.Crmisbtv_lastlogin = lastlogin_db
		obj.Crmisbtv_create = createdateuser_db
		obj.Crmisbtv_update = updatedateuser_db
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Perpage = perpage
	res.Totalrecord = totalrecord
	res.Time = time.Since(start).String()

	return res, nil
}
func Fetch_crmduniafilm(search string, page int) (helpers.Responsemovie, error) {
	var obj entities.Model_crmduniafilm
	var arraobj []entities.Model_crmduniafilm
	var res helpers.Responsemovie
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	perpage := 250
	totalrecord := 0
	offset := page
	sql_selectcount := ""
	sql_selectcount += ""
	sql_selectcount += "SELECT "
	sql_selectcount += "COUNT(username) as totalmember  "
	sql_selectcount += "FROM " + configs.DB_VIEW_MEMBER_DUNIAFILM + "  "
	if search != "" {
		sql_selectcount += "WHERE LOWER(username) LIKE '%" + strings.ToLower(search) + "%' "
		sql_selectcount += "OR LOWER(name) LIKE '%" + strings.ToLower(search) + "%' "
	}

	row_selectcount := con.QueryRowContext(ctx, sql_selectcount)
	switch e_selectcount := row_selectcount.Scan(&totalrecord); e_selectcount {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e_selectcount)
	}

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "username , name "
	sql_select += "FROM " + configs.DB_VIEW_MEMBER_DUNIAFILM + "  "
	if search == "" {
		sql_select += "ORDER BY username ASC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)
	} else {
		sql_select += "WHERE LOWER(username) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "OR LOWER(name) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY username ASC  LIMIT " + strconv.Itoa(perpage)
	}

	log.Println(sql_select)

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			username_db, name_db string
		)

		err = row.Scan(&username_db, &name_db)

		helpers.ErrorCheck(err)

		obj.Crmduniafilm_username = username_db
		obj.Crmduniafilm_name = name_db
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Perpage = perpage
	res.Totalrecord = totalrecord
	res.Time = time.Since(start).String()

	return res, nil
}
