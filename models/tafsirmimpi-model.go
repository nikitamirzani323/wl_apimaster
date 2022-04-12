package models

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_api_master/configs"
	"github.com/nikitamirzani323/wl_api_master/db"
	"github.com/nikitamirzani323/wl_api_master/entities"
	"github.com/nikitamirzani323/wl_api_master/helpers"
	"github.com/nleeper/goment"
)

func Fetch_tafsirmimpiHome(search string) (helpers.Response, error) {
	var obj entities.Model_tafsirmimpi
	var arraobj []entities.Model_tafsirmimpi
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "idtafsirmimpi , mimpi, artimimpi, angka2d, angka3d, angka4d,statustafsirmimpi, "
	sql_select += "createtafsirmimpi , createdatetafsirmimpi, updatetafsirmimpi, updatedatetafsirmimpi "
	sql_select += "FROM " + configs.DB_tbl_mst_tafsirmimpi + " "
	if search == "" {
		sql_select += "ORDER BY random() LIMIT 500  "
	} else {
		sql_select += "WHERE LOWER(mimpi) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "OR LOWER(artimimpi) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY mimpi ASC LIMIT 500  "
	}

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idtafsirmimpi_db                                                                               int
			mimpi_db, artimimpi_db, angka2d_db, angka3d_db, angka4d_db, statustafsirmimpi_db               string
			createtafsirmimpi_db, createdatetafsirmimpi_db, updatetafsirmimpi_db, updatedatetafsirmimpi_db string
		)

		err = row.Scan(
			&idtafsirmimpi_db, &mimpi_db, &artimimpi_db, &angka2d_db, &angka3d_db, &angka4d_db, &statustafsirmimpi_db,
			&createtafsirmimpi_db, &createdatetafsirmimpi_db, &updatetafsirmimpi_db, &updatedatetafsirmimpi_db)

		helpers.ErrorCheck(err)
		statuscss := configs.STATUS_CANCEL
		create := ""
		update := ""
		if createtafsirmimpi_db != "" {
			create = createtafsirmimpi_db + ", " + createdatetafsirmimpi_db
		}
		if updatetafsirmimpi_db != "0000-00-00" {
			update = updatetafsirmimpi_db + ", " + updatedatetafsirmimpi_db
		}
		if statustafsirmimpi_db == "Y" {
			statuscss = configs.STATUS_RUNNING
		}

		obj.Tafsirmimpi_id = idtafsirmimpi_db
		obj.Tafsirmimpi_mimpi = mimpi_db
		obj.Tafsirmimpi_artimimpi = artimimpi_db
		obj.Tafsirmimpi_angka2d = angka2d_db
		obj.Tafsirmimpi_angka3d = angka3d_db
		obj.Tafsirmimpi_angka4d = angka4d_db
		obj.Tafsirmimpi_status = statustafsirmimpi_db
		obj.Tafsirmimpi_statuscss = statuscss
		obj.Tafsirmimpi_create = create
		obj.Tafsirmimpi_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_tafsirmimpi(admin, mimpi, artimimpi, angka2d, angka3d, angka4d, status, sData string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	con := db.CreateCon()
	ctx := context.Background()
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_mst_tafsirmimpi + ` (
				idtafsirmimpi , mimpi, artimimpi, angka2d, angka3d, angka4d, statustafsirmimpi, 
				createtafsirmimpi, createdatetafsirmimpi
			) values (
				$1 ,$2, $3, $4, $5, $6, $7,
				$8, $9
			)
		`
		stmt_insert, e_insert := con.PrepareContext(ctx, sql_insert)
		helpers.ErrorCheck(e_insert)
		defer stmt_insert.Close()
		field_column := configs.DB_tbl_mst_tafsirmimpi + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		res_newrecord, e_newrecord := stmt_insert.ExecContext(
			ctx,
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter),
			mimpi, artimimpi, angka2d, angka3d, angka4d, status,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"))
		helpers.ErrorCheck(e_newrecord)
		insert, e := res_newrecord.RowsAffected()
		helpers.ErrorCheck(e)
		if insert > 0 {
			flag = true
			msg = "Succes"
			log.Println("Data Berhasil di save")
		}
	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_tafsirmimpi + `  
				SET mimpi=$1,artimimpi=$2, angka2d=$3, angka3d=$4, angka4d=$5, statustafsirmimpi=$6,
				updatetafsirmimpi=$7, updatedatetafsirmimpi=$8 
				WHERE idtafsirmimpi =$9 
			`
		stmt_record, e := con.PrepareContext(ctx, sql_update)
		helpers.ErrorCheck(e)
		rec_record, e_record := stmt_record.ExecContext(
			ctx,
			mimpi, artimimpi, angka2d, angka3d, angka4d, status,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"),
			idrecord)
		helpers.ErrorCheck(e_record)
		update_record, e_record := rec_record.RowsAffected()
		helpers.ErrorCheck(e_record)

		defer stmt_record.Close()
		if update_record > 0 {
			flag = true
			msg = "Succes"
			log.Printf("Update PASARAN Success : %d\n", idrecord)
		} else {
			log.Println("Update PASARAN failed")
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
