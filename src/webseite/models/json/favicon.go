package json

import "webseite/websocket"

type ServerFavicon struct {
    Id   int32
    Icon string
}

func (s *ServerFavicon) Send(c *websocket.Connection) {
    fav := JSONResponse{
        Ident: "favicon",
        Value: s,
    }

    fav.Send(c)
}