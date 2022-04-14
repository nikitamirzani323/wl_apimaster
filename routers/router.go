package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/nikitamirzani323/wl_api_master/controllers"
	"github.com/nikitamirzani323/wl_api_master/middleware"
)

func Init() *fiber.App {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/"
		},
		Format: "${time} | ${status} | ${latency} | ${ips[0]} | ${method} | ${path} - ${queryParams} ${body}\n",
	}))
	app.Use(recover.New())
	app.Use(compress.New())

	app.Get("/dashboard", monitor.New())
	app.Post("/api/login", controllers.CheckLogin)
	app.Post("/api/generator", controllers.GeneratorPassword)
	api := app.Group("/api/", middleware.JWTProtected())

	api.Post("valid", controllers.Home)
	api.Post("alladmin", controllers.Adminhome)
	api.Post("saveadmin", controllers.AdminSave)
	api.Post("alladminrule", controllers.Adminrulehome)
	api.Post("saveadminrule", controllers.AdminruleSave)
	api.Post("allcurr", controllers.Currhome)
	api.Post("savecurr", controllers.CurrSave)

	api.Post("pasaran", controllers.Pasaranhome)
	api.Post("pasaransave", controllers.Pasaransave)
	api.Post("keluaran", controllers.Keluaranhome)
	api.Post("keluaransave", controllers.Keluaransave)
	api.Post("keluarandelete", controllers.Keluarandelete)

	api.Post("prediksi", controllers.Prediksihome)
	api.Post("prediksisave", controllers.Prediksisave)
	api.Post("prediksidelete", controllers.Prediksidelete)

	api.Post("tafsirmimpi", controllers.Tafsirmimpihome)
	api.Post("tafsirmimpisave", controllers.Tafsirmimpisave)

	api.Post("news", controllers.Newshome)
	api.Post("newssave", controllers.Newssave)
	api.Post("newsdelete", controllers.Newsdelete)
	api.Post("categorynews", controllers.Categoryhome)
	api.Post("categorynewssave", controllers.Categorysave)
	api.Post("categorynewsdelete", controllers.Categorydelete)

	api.Post("movie", controllers.Moviehome)
	api.Post("movienotcdn", controllers.Moviehomenotcdn)
	api.Post("movietrouble", controllers.Movietroublehome)
	api.Post("moviemini", controllers.Movieminihome)
	api.Post("moviesave", controllers.Moviesave)
	api.Post("moviedelete", controllers.Moviedelete)
	api.Post("movieseries", controllers.Moviehomeseries)
	api.Post("movieseriestrouble", controllers.Moviehomeseriestrouble)
	api.Post("movieseriessave", controllers.Movieseriessave)
	api.Post("movieseriesseason", controllers.Seasonhome)
	api.Post("movieseriesseasonsave", controllers.Seasonsave)
	api.Post("movieseriesseasondelete", controllers.Seasondelete)
	api.Post("movieseriesepisode", controllers.Episodehome)
	api.Post("movieseriesepisodesave", controllers.Episodesave)
	api.Post("movieseriesepisodedelete", controllers.Episodedelete)
	api.Post("moviecloudualbum", controllers.Moviecloud)
	api.Post("moviecloudupdate", controllers.Movieupdatecloud)
	api.Post("movieclouddelete", controllers.Moviedeletecloud)
	api.Post("moviecloudupload", controllers.Movieuploadcloud)
	api.Post("genremovie", controllers.Genrehome)
	api.Post("genremoviesave", controllers.Genresave)
	api.Post("genremoviedelete", controllers.Genredelete)

	api.Post("slider", controllers.Sliderhome)
	api.Post("slidersave", controllers.Slidersave)
	api.Post("sliderdelete", controllers.Sliderdelete)

	api.Post("domain", controllers.Domainhome)
	api.Post("domainsave", controllers.DomainSave)

	api.Post("webagen", controllers.Websiteagenhome)
	api.Post("webagensave", controllers.Websiteagensave)
	api.Post("game", controllers.Gamehome)
	api.Post("gamesave", controllers.Gamesave)

	api.Post("cloudflare", controllers.Moviecloud2)
	api.Post("album", controllers.Albumhome)
	api.Post("albumsave", controllers.Albumsave)

	api.Post("crm", controllers.Crmhome)
	api.Post("crmisbtv", controllers.Crmisbtvhome)
	api.Post("crmduniafilm", controllers.Crmduniafilm)
	api.Post("crmsave", controllers.CrmSave)
	api.Post("crmsavesource", controllers.CrmSavesource)
	return app
}
