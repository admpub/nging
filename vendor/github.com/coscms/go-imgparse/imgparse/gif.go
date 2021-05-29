package imgparse

import (
    "image/gif"
    "io"
)

func parseGIF(r io.Reader) (int, int, error) {
    cfg, err := gif.DecodeConfig(r)
    if err != nil {
        return 0, 0, err
    }

    return cfg.Width, cfg.Height, nil
}
