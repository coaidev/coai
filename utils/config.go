package utils

import (
	"chat/globals"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var configFile = "config/config.yaml"
var configTmpFile = "config/config.tmp.yaml"
var configBackupFile = "config/config.bak.yaml"
var configExampleFile = "config.example.yaml"
var configMutex sync.Mutex

var redirectRoutes = []string{
	"/v1",
	"/mj",
	"/attachments",
}

func SaveConfig(key string, value interface{}) error {
	// save config to file with mutex lock
	configMutex.Lock()
	defer configMutex.Unlock()

	if err := viper.WriteConfigAs(configBackupFile); err != nil {
		return err
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	currentConfig := viper.AllSettings()

	currentConfig[key] = value

	for k, v := range currentConfig {
		viper.Set(k, v)
	}

	if err := viper.WriteConfigAs(configTmpFile); err != nil {
		return err
	}

	if _, err := os.Stat(configFile); err == nil {
		if removeErr := os.Remove(configFile); removeErr != nil {
			return removeErr
		}
	}

	if err := os.Rename(configTmpFile, configFile); err != nil {
		return err
	}

	return nil
}

func ApplySeo(title, icon string) {
	// seo optimization

	if !viper.GetBool("serve_static") {
		return
	}

	content, err := ReadFile("./app/dist/index.html")
	if err != nil {
		globals.Warn(fmt.Sprintf("[service] failed to read index.html: %s", err.Error()))
		return
	}

	if len(title) > 0 {
		content = strings.ReplaceAll(content, "CoAI.Dev", title)
		content = strings.ReplaceAll(content, "chatnio", strings.ToLower(title))
	}

	if len(icon) > 0 {
		content = strings.ReplaceAll(content, "/favicon.ico", icon)
	}

	if err := WriteFile("./app/dist/index.cache.html", content, true); err != nil {
		globals.Warn(fmt.Sprintf("[service] failed to write index.cache.html: %s", err.Error()))
	}

	globals.Info("[service] seo optimization applied to index.cache.html")
}

func ApplyPWAManifest(content string) {
	// pwa manifest rewrite (site.webmanifest -> site.cache.webmanifest)

	if !viper.GetBool("serve_static") {
		return
	}

	if len(content) == 0 {
		// read from site.webmanifest if not provided

		var err error
		content, err = ReadFile("./app/dist/site.webmanifest")
		if err != nil {
			globals.Warn(fmt.Sprintf("[service] failed to read site.webmanifest: %s", err.Error()))
			return
		}
	}

	if err := WriteFile("./app/dist/site.cache.webmanifest", content, true); err != nil {
		globals.Warn(fmt.Sprintf("[service] failed to write site.cache.webmanifest: %s", err.Error()))
	}

	globals.Info("[service] pwa manifest applied to site.cache.webmanifest")
}

func ReadPWAManifest() (content string) {
	// read site.cache.webmanifest content or site.webmanifest if not found

	if !viper.GetBool("serve_static") {
		return
	}

	if text, err := ReadFile("./app/dist/site.cache.webmanifest"); err == nil && len(text) > 0 {
		return text
	}

	if text, err := ReadFile("./app/dist/site.webmanifest"); err != nil {
		globals.Warn(fmt.Sprintf("[service] failed to read site.webmanifest: %s", err.Error()))
	} else {
		content = text
	}

	return
}

func RegisterStaticRoute(engine *gin.Engine) {
	// static files are in ~/app/dist

	if !viper.GetBool("serve_static") {
		engine.NoRoute(func(c *gin.Context) {
			c.JSON(404, gin.H{"status": false, "message": "not found or method not allowed"})
		})
		return
	}

	if !IsFileExist("./app/dist") {
		fmt.Println("[service] app/dist not found, please run `npm run build`")
		return
	}

	ApplySeo(viper.GetString("system.general.title"), viper.GetString("system.general.logo"))
	ApplyPWAManifest(viper.GetString("system.general.pwamanifest"))

	engine.GET("/", func(c *gin.Context) {
		c.File("./app/dist/index.cache.html")
	})

	engine.GET("/site.webmanifest", func(c *gin.Context) {
		c.File("./app/dist/site.cache.webmanifest")
	})

	engine.Use(static.Serve("/", static.LocalFile("./app/dist", true)))
	engine.NoRoute(func(c *gin.Context) {
		c.File("./app/dist/index.cache.html")
	})

	for _, route := range redirectRoutes {
		engine.Any(fmt.Sprintf("%s/*path", route), func(c *gin.Context) {
			c.Request.URL.Path = "/api" + c.Request.URL.Path
			fmt.Println(c.Request.URL.Path)
			engine.HandleContext(c)
		})
	}

	fmt.Println(`[service] start serving static files from ~/app/dist`)
}
