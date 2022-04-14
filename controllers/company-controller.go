package controllers

import (
	"log"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/wl_api_master/entities"
	"github.com/nikitamirzani323/wl_api_master/helpers"
	"github.com/nikitamirzani323/wl_api_master/models"
)

const Fieldcompany_home_redis = "LISTCOMPANY_MASTER_WL"

func Companyhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_company)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}

	var obj entities.Model_company
	var arraobj []entities.Model_company
	var obj_listcurr entities.Model_compcurr
	var arraobj_listcurr []entities.Model_compcurr
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcompany_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	listcurr_RD, _, _, _ := jsonparser.Get(jsonredis, "listcurr")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		company_idcompany, _ := jsonparser.GetString(value, "company_idcompany")
		company_startjoin, _ := jsonparser.GetString(value, "company_startjoin")
		company_endjoin, _ := jsonparser.GetString(value, "company_endjoin")
		company_idcurr, _ := jsonparser.GetString(value, "company_idcurr")
		company_nmcompany, _ := jsonparser.GetString(value, "company_nmcompany")
		company_nmowner, _ := jsonparser.GetString(value, "company_nmowner")
		company_phoneowner, _ := jsonparser.GetString(value, "company_phoneowner")
		company_emailowner, _ := jsonparser.GetString(value, "company_emailowner")
		company_urlendpoint, _ := jsonparser.GetString(value, "company_urlendpoint")
		company_status, _ := jsonparser.GetString(value, "company_status")
		company_create, _ := jsonparser.GetString(value, "company_create")
		company_update, _ := jsonparser.GetString(value, "company_update")

		obj.Company_idcomp = company_idcompany
		obj.Company_startjoin = company_startjoin
		obj.Company_endjoin = company_endjoin
		obj.Company_idcurr = company_idcurr
		obj.Company_nmcompany = company_nmcompany
		obj.Company_nmowner = company_nmowner
		obj.Company_phoneowner = company_phoneowner
		obj.Company_emailowner = company_emailowner
		obj.Company_urlendpoint = company_urlendpoint
		obj.Company_status = company_status
		obj.Company_create = company_create
		obj.Company_update = company_update
		arraobj = append(arraobj, obj)
	})
	jsonparser.ArrayEach(listcurr_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		curr_idcurr, _ := jsonparser.GetString(value, "curr_idcurr")

		obj_listcurr.Curr_idcurr = curr_idcurr
		arraobj_listcurr = append(arraobj_listcurr, obj_listcurr)
	})
	if !flag {
		result, err := models.Fetch_companyHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcompany_home_redis, result, 30*time.Minute)
		log.Println("COMPANY MYSQL")
		return c.JSON(result)
	} else {
		log.Println("COMPANY CACHE")
		return c.JSON(fiber.Map{
			"status":       fiber.StatusOK,
			"message":      "Success",
			"record":       arraobj,
			"listcurrency": arraobj_listcurr,
			"time":         time.Since(render_page).String(),
		})
	}
}
func CompanySave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_companysave)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	client_admin, _ := helpers.Parsing_Decry(temp_decp, "==")

	result, err := models.Save_companyHome(
		client_admin,
		client.Company_idcomp, client.Company_idcurr, client.Company_nmcompany,
		client.Company_nmowner, client.Company_phoneowner, client.Company_emailowner, client.Company_urlendpoint,
		client.Company_status, client.Sdata)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_company()
	return c.JSON(result)
}
func _deleteredis_company() {
	val_master := helpers.DeleteRedis(Fieldcompany_home_redis)
	log.Printf("REDIS DELETE MASTER COMPANY : %d", val_master)
}