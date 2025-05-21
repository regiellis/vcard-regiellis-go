# Regi Ellis Smart Card Site

![Regi Ellis](https://regiellis.com/static/images/og-image.png)

[![Go](https://img.shields.io/badge/Go-1.22-blue?logo=go)](https://go.dev/)
[![TailwindCSS](https://img.shields.io/badge/TailwindCSS-3.x-38bdf8?logo=tailwindcss)](https://tailwindcss.com/)
[![Alpine.js](https://img.shields.io/badge/Alpine.js-3.x-77c1d2?logo=alpine.js)](https://alpinejs.dev/)
[![Gin](https://img.shields.io/badge/Gin-Framework-00b386?logo=go)](https://gin-gonic.com/)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-yellow.svg)](https://www.gnu.org/licenses/gpl-3.0)

---

Personal project to create a digital business card. The site auto-downloads a vCard file and has a dynamic note prefix feature to use with a NFC card for differnt use cases.

## Features

- **Digital vCard Download:**
  - Automatically prompts users to download a vCard (.vcf) file with dynamic parameters (e.g., custom note prefix via URL).
  - vCard includes contact info, social links, and a profile photo.
- **Animated Background:**
  - Subtle, animated canvas background for a modern look.
- **Branding Shapes:**
  - Interactive, device-orientation-reactive branding shapes for visual flair.
- **Responsive Design:**
  - Mobile-first, works well on all devices.
- **NFC-Ready:**
  - Supports dynamic note prefixing for NFC use cases (e.g., `?prefix=Hello+from+NFC`).

## Usage

1. **Run the server:**

   ```sh
   go run server.go
   # The site will be available at http://localhost:8080
   ```

2. **Download vCard:**
   - Visit `/` and the vCard will auto-download.
   - To add a custom note prefix (e.g., for NFC), use:

     `http://localhost:8080/?prefix=Your+Custom+Message`

   - The prefix will appear as the first line in the vCard note field.

## Project Structure

- `server.go` — Go backend serving the site and vCard endpoint
- `public/index.html` — Main HTML, TailwindCSS, Alpine.js, and JS for UI/animation
- `public/assets/` — Images and styles
- `.gitignore` — Clean repo, ignores binaries, temp, and system files

## Customization

- Edit `CardData` in `server.go` for your own info.
- Update images in `public/assets/images/`.
- Tweak styles in `public/assets/styles/www.css` or Tailwind classes in `index.html`.

## License

Personal project by Regi Ellis. Not intended for commercial use.
