package imgparse

import (
    "errors"
    "io"
)

func ParseRes(r io.Reader, content string) (int, int, error) {
    switch content {
    case "gif":
        return parseGIF(r)
    case "jpeg":
        return parseJPEG(r)
    case "png":
        return parsePNG(r)
    case "webp":
        return parseWebP(r)
    case "webpll":
        return parseWebP(r)
    default:
        return 0, 0, errors.New("Unknown content type.")
    }
}
