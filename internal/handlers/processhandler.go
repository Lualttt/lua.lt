package handlers

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func Process(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32000)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		slog.Error("process page failed to parse form", "error", err)
		return
	}

	blackData := r.PostFormValue("black")
	redData := r.PostFormValue("red")

	if blackData == "" ||
		redData == "" ||
		!strings.HasPrefix(blackData, "data:image/png;base64,") ||
		!strings.HasPrefix(redData, "data:image/png;base64,") {

		http.Error(w, "wrong form data", http.StatusBadRequest)
		slog.Error("process page wrong form data")
		return
	}

	blackBase64 := strings.TrimPrefix(blackData, "data:image/png;base64,")
	redBase64 := strings.TrimPrefix(redData, "data:image/png;base64,")

	blackImage, blackErr := decodeImage(blackBase64)
	redImage, redErr := decodeImage(redBase64)
	if blackErr != nil || redErr != nil {
		http.Error(w, "Internal server error!", http.StatusInternalServerError)
		slog.Error("process page failed to decode image", "blackError", blackErr, "redError", redErr)
		return
	}

	whiteBackground := image.NewRGBA(image.Rect(0, 0, 264, 176))
	draw.Draw(whiteBackground, whiteBackground.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	blackImage = cropImage(blackImage, 264, 176)
	redImage = cropImage(redImage, 264, 176)

	blackImage = alphaComposite(whiteBackground, blackImage)
	redImage = alphaComposite(whiteBackground, redImage)

	redImage = convertToBlack(redImage)

	blackBuffer, blackErr := encodeImageToBuffer(blackImage)
	redBuffer, redErr := encodeImageToBuffer(redImage)
	if blackErr != nil || redErr != nil {
		http.Error(w, "Internal server error!", http.StatusInternalServerError)
		slog.Error("process page failed to encode image", "blackError", blackErr, "redError", redErr)
		return
	}

	err = sendImage(blackBuffer, redBuffer)
	if err != nil {
		http.Error(w, "Internal server error!", http.StatusInternalServerError)
		slog.Error("process page error sending image", "error", err)
		return
	}
}

func decodeImage(base64Str string) (image.Image, error) {
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	return png.Decode(bytes.NewReader(data))
}

func cropImage(img image.Image, width int, height int) image.Image {
	cropped := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(cropped, cropped.Bounds(), img, image.Point{}, draw.Over)
	return cropped
}

func alphaComposite(dst, src image.Image) *image.RGBA {
	dstCopy := image.NewRGBA(dst.Bounds())
	draw.Draw(dstCopy, dstCopy.Bounds(), dst, image.Point{}, draw.Over)
	draw.Draw(dstCopy, dstCopy.Bounds(), src, image.Point{}, draw.Over)
	return dstCopy
}

func convertToBlack(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	gray := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, g, _, a := img.At(x, y).RGBA()
			gray.Set(x, y, color.RGBA{uint8(g), uint8(g), uint8(g), uint8(a)})
		}
	}
	return gray
}

func encodeImageToBuffer(img image.Image) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := png.Encode(buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func sendImage(blackBuffer []byte, redBuffer []byte) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	blackPart, err := writer.CreateFormFile("black", "black")
	if err != nil {
		return err
	}
	_, err = blackPart.Write(blackBuffer)
	if err != nil {
		return err
	}

	redPart, err := writer.CreateFormFile("red", "red")
	if err != nil {
		return err
	}
	_, err = redPart.Write(redBuffer)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", os.Getenv("LUALT_DRAWING_ADDRESS"), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
