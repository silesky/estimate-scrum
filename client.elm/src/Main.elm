module Main exposing (..)
effect module WebSocket where { command = MyCmd, subscription = MySub } exposing
  ( send
  , listen
  , keepAlive
  )


import Html exposing (Html, text, div, h1, img)
import Html.Attributes exposing (class, src)



---- MODEL ----

-- create type for 'record'
type alias UserEstimate  = {
    username: String,
    estimate: Int
}

-- create type for array of [0, 3, 5, 8, 13] etc
type alias PointOptions = List Int

type alias Model =
    {
      estimates : List UserEstimate
    }

init : ( Model, Cmd Msg )
init =
    ( {
      estimates = []
    }, Cmd.none )



---- UPDATE ----


type Msg
    = NoOp


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    ( model, Cmd.none )



---- VIEW ----


view : Model -> Html Msg
view model =
    div []
        [ img [ src "/logo.svg" ] []
        , h1 [] [ text "Your Elm App is working!" ]
        ]



---- PROGRAM ----


main : Program Never Model Msg
main =
    Html.program
        { view = view
        , init = init
        , update = update
        , subscriptions = always Sub.none
        }
