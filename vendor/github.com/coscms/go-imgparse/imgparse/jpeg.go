package imgparse

import (
    "image/jpeg"
    "io"
)

func parseJPEG(r io.Reader) (int, int, error) {
    cfg, err := jpeg.DecodeConfig(r)
    if err != nil {
        return 0, 0, err
    }

    return cfg.Width, cfg.Height, nil
}
