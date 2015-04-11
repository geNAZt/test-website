package json

import "webseite/websocket"

type ServerFavicon struct {
    Id   int32
    Icon string
}

func (s *ServerFavicon) Send(c *websocket.Connection) bool {
    fav := JSONResponse{
        Ident: "favicon",
        Value: s,
    }

    return fav.Send(c)
}