package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
	// Remove newlines and dangerous characters for vCard fields
	return strings.ReplaceAll(strings.ReplaceAll(s, "\n", " "), "\r", " ")
}

// Read and encode the photo as base64
func getBase64Photo(photoPath string) (string, error) {
	if strings.HasPrefix(photoPath, "/") {
		photoPath = "public" + photoPath // Adjust based on your folder structure
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
	// --- Set Gin mode from env or default to release ---
	mode := gin.ReleaseMode
	if val, ok := os.LookupEnv("GIN_MODE"); ok {
		mode = val
	}
	gin.SetMode(mode)

	// --- Detect environment for your own use ---
	env := "production"
	if gin.Mode() == gin.DebugMode {
		env = "development"
	}

	// --- Choose template file (minified in production if available) ---
	templateFile := "public/index.html"
	if env == "production" {
		if _, err := os.Stat("public/index.min.html"); err == nil {
			templateFile = "public/index.min.html"
		}
	}

	// --- Use Gin's HTML template loader ---
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// --- Security headers middleware with Google Fonts CSP ---
	r.Use(func(c *gin.Context) {
		if env == "production" {
			c.Writer.Header().Set("Content-Security-Policy",
				"default-src 'self' https://cdn.tailwindcss.com https://unpkg.com https://cdnjs.cloudflare.com;"+
					" style-src 'self' 'unsafe-inline' https://cdn.tailwindcss.com https://fonts.googleapis.com;"+
					" font-src 'self' https://fonts.gstatic.com;"+
					" script-src 'self' 'unsafe-inline' https://cdn.tailwindcss.com https://unpkg.com https://cdnjs.cloudflare.com 'unsafe-eval';"+
					" img-src * data: blob:;")
		}
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("Referrer-Policy", "no-referrer")
		c.Writer.Header().Set("Cache-Control", "no-store")
		c.Next()
	})

	// --- Load the template for Gin ---
	r.LoadHTMLFiles(templateFile)

	metaDescription := "Award-winning Orlando-based developer with expertise in front-end, back-end, DevOps, and technical leadership. Passionate about generative AI, Unreal Engine, ICVFX, real-time graphics, and video production. Bridging creativity and innovation across industries."
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
		Notes:     "Award-winning Orlando-based developer with expertise in front-end, back-end, DevOps, and technical leadership. Passionate about generative AI, Unreal Engine, ICVFX, real-time graphics, and video production. Bridging creativity and innovation across industries.",
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Card":            data,
			"MetaDescription": metaDescription,
		})
	})

	r.GET("/vcard", func(c *gin.Context) {
		c.Header("Content-Type", "text/vcard")
		c.Header("Content-Disposition", "attachment; filename=\""+data.Name+".vcf\"")
		c.Header("Cache-Control", "no-store")

		// Get photo as base64
		photoBase64, err := getBase64Photo(data.Photo)
		photoSection := ""
		if err != nil {
			log.Printf("Failed to encode photo: %v", err)
		}
		if err == nil && photoBase64 != "" {
			photoType := "JPEG" // Or detect from mime type
			photoSection = "PHOTO;ENCODING=b;TYPE=" + photoType + ":" + photoBase64 + "\n"
		}

		// Custom note handling: add prefix (from query param) with a vCard newline if present
		finalNote := escapeVCardField(data.Notes)
		prefix := c.Query("prefix")
		if len(prefix) > 500 {
			prefix = ""
		}
		if prefix != "" {
			finalNote = escapeVCardField(prefix) + "\\n" + finalNote
		}

		// Format social links for the note field
		extraLinks := ""
		if data.GitHub != "" {
			extraLinks += "GitHub: " + data.GitHub + "\\n"
		}
		if data.LinkedIn != "" {
			extraLinks += "LinkedIn: " + data.LinkedIn + "\\n"
		}
		if data.Instagram != "" {
			extraLinks += "Instagram: " + data.Instagram + "\\n"
		}
		if data.YouTube != "" {
			extraLinks += "YouTube: " + data.YouTube + "\\n"
		}

		noteWithLinks := finalNote
		if extraLinks != "" {
			noteWithLinks += "\\n" + extraLinks
		}

		// Create structured name parts (better compatibility)
		nameParts := strings.Split(data.Name, " ")
		lastName := ""
		firstName := data.Name
		if len(nameParts) > 1 {
			lastName = nameParts[len(nameParts)-1]
			firstName = strings.Join(nameParts[:len(nameParts)-1], " ")
		}

		vcard := "BEGIN:VCARD\n" +
			"VERSION:3.0\n" +
			"FN:" + escapeVCardField(data.Name) + "\n" +
			"N:" + escapeVCardField(lastName) + ";" + escapeVCardField(firstName) + ";;;\n" +
			"ORG:" + escapeVCardField(data.Company) + "\n" +
			"TITLE:" + escapeVCardField(data.Role) + "\n" +
			"EMAIL;TYPE=INTERNET,WORK:" + escapeVCardField(data.Email) + "\n" +
			"TEL;TYPE=WORK,VOICE:" + escapeVCardField(data.Phone) + "\n" +
			"ADR;TYPE=WORK:;;" + escapeVCardField(data.Address) + "\n" +
			"URL;TYPE=WORK:" + escapeVCardField(data.Website) + "\n" +
			photoSection +
			"NOTE:" + noteWithLinks + "\n" +
			"END:VCARD\n"
		c.String(http.StatusOK, vcard)
	})

	// Serve static files
	r.Static("/static", "public/assets")

	log.Println("Server running on :8080")
	srv := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(srv.ListenAndServe())
}
