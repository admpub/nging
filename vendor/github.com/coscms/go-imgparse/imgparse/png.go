package imgparse

import (
    "image/png"
    "io"
)

func parsePNG(r io.Reader) (int, int, error) {
    cfg, err := png.DecodeConfig(r)
    if err != nil {
        return 0, 0, err
    }

    return cfg.Width, cfg.Height, nil
}
