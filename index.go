package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/katallaxie/pkg/logger"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	"github.com/maragudk/gomponents/html"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/spf13/cobra"
	fiber_goth "github.com/zeiss/fiber-goth"
	htmx "github.com/zeiss/fiber-htmx"
)

// Config ...
type Config struct {
	Flags *Flags
}

// Flags ...
type Flags struct {
	Addr string
}

var cfg = &Config{
	Flags: &Flags{},
}

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context())
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfg.Flags.Addr, "addr", ":3000", "addr")

	rootCmd.SilenceUsage = true
}

func run(ctx context.Context) error {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	logger.RedirectStdLog(logger.LogSink)

	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:3000/auth/github/callback"),
	)

	sess := session.Config{
		KeyLookup:      fmt.Sprintf("cookie:%s", gothic.SessionName),
		CookieHTTPOnly: true,
	}
	fiber_goth.ConfigDefault.Session = fiber_goth.NewSessionStore(session.New(sess))

	app := fiber.New()
	app.Static("/static", "./public")

	app.Get("/login/:provider", fiber_goth.NewBeginAuthHandler())
	app.Get("/auth/:provider/callback/", fiber_goth.NewCompleteAuthHandler())
	app.Get("/logout", fiber_goth.NewLogoutHandler())

	app.Get("/", func(ctx *fiber.Ctx) error {
		page := c.HTML5(c.HTML5Props{
			Title:    "index",
			Language: "en",
			Head: []g.Node{
				html.Script(g.Attr("src", "/static/app.js"), g.Attr("type", "application/javascript")),
				html.Link(html.Rel("stylesheet"), html.Href("/static/styles.css")),
			},
			Body: []g.Node{
				html.Button(g.Text("Button"), g.Attr("hx-get", "/api/redirect"), c.Classes{"inline-block cursor-pointer rounded-md bg-gray-800 px-4 py-3 text-center text-sm font-semibold uppercase text-white transition duration-200 ease-in-out hover:bg-gray-900": true}),
			},
		})

		ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return page.Render(ctx)
	})

	app.Get("/api/redirect", htmx.NewHtmxHandler(func(hx *htmx.Htmx) error {
		hx.Redirect("https://google.com")

		return nil
	}))

	err := app.Listen(cfg.Flags.Addr)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
