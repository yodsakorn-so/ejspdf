# ejspdf

![Go](https://github.com/yodsakorn-so/ejspdf/actions/workflows/go.yml/badge.svg)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.22+-blue)

Render **EJS templates to PDF** using **Go + Chrome (chromedp)**.

`ejspdf` เหมาะสำหรับงานสร้างเอกสาร เช่น
- ใบรับรองแพทย์
- ใบแจ้งหนี้ (Invoice)
- ใบเสร็จ
- เอกสารราชการ / เอกสารภาษาไทย

โดยใช้ HTML + CSS เต็มรูปแบบ และแปลงเป็น PDF ด้วย Chrome

---

## ✨ Features

- ✅ Render EJS template (Node.js style)
- ✅ Full HTML / CSS support
- ✅ A4 PDF output
- ✅ รองรับภาษาไทย (Thai font friendly)
- ✅ ใช้ Chrome / Chromium จริงในการ render
- ✅ ใช้ได้กับ Windows / macOS / Linux

---

## 📦 Installation

```bash
go get github.com/yodsakorn-so/ejspdf
