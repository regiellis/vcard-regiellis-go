package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type CardData struct {
	Name      string
	Role      string
	Email     string
	Company   string
	Phone     string
	Address   string
	Photo     string
	Website   string
	GitHub    string
	LinkedIn  string
	Instagram string
	YouTube   string
	Notes     string
}

func escapeVCardField(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\n", " "), "\r", " ")
}

func getBase64Photo(photoPath string) (string, error) {
	if strings.HasPrefix(photoPath, "/") {
		photoPath = photoPath
	}
	fileData, err := os.ReadFile(photoPath)
	if err != nil {
		return "", err
	}
	mimeType := http.DetectContentType(fileData)
	if !strings.HasPrefix(mimeType, "image/") {
		return "", fmt.Errorf("not an image: %s", mimeType)
	}
	return base64.StdEncoding.EncodeToString(fileData), nil
}

func main() {
	// Force Gin release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Use Gin's default router
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Security headers middleware (CSP temporarily disabled)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Referrer-Policy", "no-referrer")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Next()
	})

	r.LoadHTMLFiles("public/index.html")

	metaDescription := "Award-winning Orlando-based Full stack developer, DevOps, and technical leadership. Passionate about generative AI, Unreal Engine, ICVFX, real-time graphics, and video production. Bridging creativity and innovation across industries."
	data := CardData{
		Name:      "Regi Ellis",
		Role:      "Web Veteran & DevOps Tooling | Python, Go, AI Integration | Unreal Engine Enthusiast",
		Email:     "regi@bynine.io",
		Company:   "Playlogic LLC",
		Phone:     "+1.407.437.0818",
		Address:   "1540 International Pkwy, Suite 200, Lake Mary, FL 32746",
		Photo:     "/static/images/regiellis_profile.jpg",
		Website:   "https://regiellis.com",
		GitHub:    "https://github.com/regiellis",
		LinkedIn:  "https://www.linkedin.com/in/rellis/",
		Instagram: "https://www.instagram.com/_its.just.regi_/",
		YouTube:   "https://www.youtube.com/@_its.just.regi_/videos",
		Notes:     "Award-winning Orlando-based Full stack developer, DevOps, and technical leadership. Passionate about generative AI, Unreal Engine, ICVFX, real-time graphics, and video production. Bridging creativity and innovation across industries.",
	}

	r.GET("/", func(c *gin.Context) {
		prefix := c.Query("prefix")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Card":            data,
			"MetaDescription": metaDescription,
			"Prefix":          prefix,
		})
	})

	r.GET("/vcard", func(c *gin.Context) {
		c.Header("X-Robots-Tag", "noindex, nofollow")
		c.String(http.StatusNotFound, "Not Found")
	})
	r.Static("/static", "public/assets")

	// Use PORT from env (Leapcell sets this), default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
