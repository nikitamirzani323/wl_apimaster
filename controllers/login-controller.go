package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/wl_api_master/entities"
	"github.com/nikitamirzani323/wl_api_master/helpers"
	"github.com/nikitamirzani323/wl_api_master/models"
)

const Field_login_redis = "LISTLOGINADMIN_MASTER_WL"

func CheckLogin(c *fiber.Ctx) error {
	render_page := time.Now()
	var errors []*helpers.ErrorResponse
	client := new(entities.Login)
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
	flag_login := false
	ruleadmin := ""
	resultredis, flag := helpers.GetRedis(Field_login_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		login_username, _ := jsonparser.GetString(value, "login_username")
		login_password, _ := jsonparser.GetString(value, "login_password")
		login_idadmin, _ := jsonparser.GetString(value, "login_idadmin")

		if login_username == client.Username {
			hashpass := helpers.HashPasswordMD5(client.Password)
			if hashpass == login_password {
				flag_login = true
				ruleadmin = login_idadmin
			}
		}
	})
	if !flag {
		result, flag_model, rule_model, err := models.Login_Model(client.Username, client.Password)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		if flag_model {
			flag_login = true
			ruleadmin = rule_model
			helpers.SetRedis(Field_login_redis, result, 30*time.Hour)
			log.Println("LIST LOGIN ADMIN MASTER MYSQL")

			models.Update_login(client.Username, client.Ipaddress, client.Timezone)
		} else {
			return c.JSON(fiber.Map{
				"status": fiber.ErrBadRequest,
				"token":  "",
				"time":   time.Since(render_page).String(),
			})
		}

	} else {
		log.Println("LIST LOGIN ADMIN MASTER CACHE")

	}
	temp_token := ""
	if flag_login {
		dataclient := client.Username + "==" + ruleadmin
		dataclient_encr, keymap := helpers.Encryption(dataclient)
		dataclient_encr_final := dataclient_encr + "|" + strconv.Itoa(keymap)
		t, err := helpers.GenerateNewAccessToken(dataclient_encr_final)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		temp_token = t
	}

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"token":  temp_token,
		"time":   time.Since(render_page).String(),
	})
}
func Home(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Home)
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
	client_username, idruleadmin := helpers.Parsing_Decry(temp_decp, "==")
	log.Printf("USERNAME : %s", client_username)
	log.Printf("RULE : %s", idruleadmin)
	log.Printf("PAGE : %s", client.Page)

	ruleadmin := models.Get_AdminRule("ruleadmingroup", idruleadmin)
	flag := models.Get_listitemsearch(ruleadmin, ",", client.Page)
	if !flag {
		c.Status(fiber.StatusForbidden)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Anda tidak bisa akses halaman ini",
			"record":  nil,
		})
	} else {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "ADMIN",
			"record":  nil,
		})
	}
}
func GeneratorPassword(c *fiber.Ctx) error {
	render_page := time.Now()
	client := new(entities.Login)

	pwd := helpers.HashPasswordMD5(client.Password)

	return c.JSON(fiber.Map{
		"status":   fiber.StatusOK,
		"password": pwd,
		"time":     time.Since(render_page).String(),
	})
}
