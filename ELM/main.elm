module Main exposing (..)

import Browser
import Html exposing (Html, div, textarea, button, text, h1)
import Html.Attributes as Attr
import Html.Events exposing (onInput, onClick)
import Svg exposing (svg, polyline, circle, g, ellipse)
import Svg.Attributes as SvgAttr

-- MODÈLE

type alias Turtle =
    { x : Float
    , y : Float
    , angle : Float 
    , path : List (Float, Float)
    }

type alias Model =
    { turtle : Turtle
    , commandInput : String
    }

init : Model
init =
    { turtle = 
        { x = 300
        , y = 200
        , angle = -90 
        , path = [ (300, 200) ] 
        }
    , commandInput = ""
    }

-- UPDATE

type Msg
    = UpdateInput String
    | Execute
    | Reset

update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateInput txt ->
            { model | commandInput = txt }

        Reset ->
            init

        Execute ->
            let
                words = String.words (String.toLower model.commandInput)
            in
            case words of

                [ "fd"] ->
                    moveForward 50 model
                
                [ "rt" ] ->
                    rotateTurtle 90 model

                [ "lt" ] ->
                    rotateTurtle -90 model

                [ "fd", valStr ] ->
                    moveForward (Maybe.withDefault 50 (String.toFloat valStr)) model
                
                [ "rt", valStr ] ->
                    rotateTurtle (Maybe.withDefault 90 (String.toFloat valStr)) model

                [ "lt", valStr ] ->
                    rotateTurtle (-(Maybe.withDefault 90 (String.toFloat valStr))) model

                [ "cercle", valStr] -> 
                        drawCircle (Maybe.withDefault 50 (String.toFloat valStr)) model

                [ "carre", valStr] ->
                    drawSquare (Maybe.withDefault 50 (String.toFloat valStr)) model

                [ "etoile", valStr ] ->
                    drawStar (Maybe.withDefault 50 (String.toFloat valStr)) model
                
                [ "reset" ] ->
                    init

                [ "clear" ] ->
                    let t = model.turtle in
                    { model | turtle = { t | path = [ (t.x, t.y) ] } }

                _ ->
                    model

moveForward : Float -> Model -> Model
moveForward dist model =
    let
        t = model.turtle
        rad = t.angle * (pi / 180)
        newX = t.x + dist * cos rad
        newY = t.y + dist * sin rad
    in
    { model | turtle = { t | x = newX, y = newY, path = t.path ++ [ (newX, newY) ] } }

rotateTurtle : Float -> Model -> Model
rotateTurtle deg model =
    let
        t = model.turtle
    in
    { model | turtle = { t | angle = t.angle + deg } }

drawSquare : Float -> Model -> Model
drawSquare size model =
    model
        |> moveForward size |> rotateTurtle 90
        |> moveForward size |> rotateTurtle 90
        |> moveForward size |> rotateTurtle 90
        |> moveForward size |> rotateTurtle 90

drawStar : Float -> Model -> Model
drawStar size model =
    let
        originalAngle = model.turtle.angle

        preparedModel = 
            let t = model.turtle 
            in { model | turtle = { t | angle = -72 } }
        
        drawnModel =
            preparedModel
                |> moveForward size |> rotateTurtle 144
                |> moveForward size |> rotateTurtle 144
                |> moveForward size |> rotateTurtle 144
                |> moveForward size |> rotateTurtle 144
                |> moveForward size |> rotateTurtle 144
    in
    let finalTurtle = drawnModel.turtle
    in { drawnModel | turtle = { finalTurtle | angle = originalAngle } }

drawCircle : Float -> Model -> Model
drawCircle radius model =
    let
        stepSize = (2 * pi * radius) / 36
        
        drawSteps n currentModel =
            if n <= 0 then
                currentModel
            else
                drawSteps (n - 1) (currentModel |> moveForward stepSize |> rotateTurtle 10)
    in
    drawSteps 360 model


-- VUE

view : Model -> Html Msg
view model =
    div 
        [ Attr.style "display" "flex"
        , Attr.style "flex-direction" "column"
        , Attr.style "align-items" "center"
        , Attr.style "justify-content" "center"
        , Attr.style "min-height" "100vh"
        , Attr.style "background-color" "#f0f2f5"
        , Attr.style "font-family" "sans-serif"
        ]
        [ h1 [ Attr.style "color" "#2c3e50" ] [ text "TcTurtle Elm Edition" ]
        
        , div 
            [ Attr.style "background" "white"
            , Attr.style "padding" "15px"
            , Attr.style "box-shadow" "0 8px 30px rgba(0,0,0,0.1)"
            , Attr.style "border-radius" "12px"
            ]
            [ svg 
                [ SvgAttr.width "900", SvgAttr.height "600", SvgAttr.viewBox "0 0 600 400" ]
                [ polyline 
                    [ SvgAttr.points (pointsToString model.turtle.path)
                    , SvgAttr.fill "none"
                    , SvgAttr.stroke "#3498db"
                    , SvgAttr.strokeWidth "3"
                    , SvgAttr.strokeLinecap "round"
                    ] []
                , viewTurtle model.turtle
                ]
            ]
        
        , div [ Attr.style "margin-top" "20px", Attr.style "width" "450px" ]
            [ textarea 
                [ onInput UpdateInput
                , Attr.value model.commandInput
                , Attr.placeholder "Veuillez entrer une commande (ex : fd 100, rt 90, carre 100...)"
                , Attr.style "width" "100%"
                , Attr.style "height" "60px"
                , Attr.style "padding" "10px"
                , Attr.style "border" "1px solid #ccc"
                , Attr.style "border-radius" "8px"
                , Attr.style "box-sizing" "border-box"
                , Attr.style "font-size" "16px"
                , Attr.style "resize" "none"
                ] []
            
            , div [ Attr.style "display" "flex", Attr.style "gap" "10px", Attr.style "margin-top" "10px" ]
                (List.map (\(label, msg, color) -> 
                    button (onClick msg :: commonBtnStyles color) [ text label ]
                ) 
                [ ("Executer", Execute, "#f3d541ff")
                , ("Reset", Reset, "#973ce7ff")
                ])
            ]
        ]

viewTurtle : Turtle -> Svg.Svg Msg
viewTurtle t =
    g [ SvgAttr.transform (turtleTransform t) ]
        [
          circle [ SvgAttr.cx "-8", SvgAttr.cy "-8", SvgAttr.r "3", SvgAttr.fill "#1b5e20" ] []
        , circle [ SvgAttr.cx "8", SvgAttr.cy "-8", SvgAttr.r "3", SvgAttr.fill "#1b5e20" ] []
        , circle [ SvgAttr.cx "-8", SvgAttr.cy "8", SvgAttr.r "3", SvgAttr.fill "#1b5e20" ] []
        , circle [ SvgAttr.cx "8", SvgAttr.cy "8", SvgAttr.r "3", SvgAttr.fill "#1b5e20" ] []
        , ellipse [ SvgAttr.cx "0", SvgAttr.cy "0", SvgAttr.rx "13", SvgAttr.ry "11", SvgAttr.fill "#4caf50" ] []
        , circle [ SvgAttr.cx "16", SvgAttr.cy "0", SvgAttr.r "5", SvgAttr.fill "#2e7d32" ] []
        ]

commonBtnStyles : String -> List (Html.Attribute Msg)
commonBtnStyles color =
    [ Attr.style "background" color
    , Attr.style "color" "white"
    , Attr.style "border" "none"
    , Attr.style "padding" "10px"
    , Attr.style "border-radius" "6px"
    , Attr.style "cursor" "pointer"
    , Attr.style "flex" "1"
    , Attr.style "font-weight" "bold"
    ]

pointsToString : List (Float, Float) -> String
pointsToString path =
    path
        |> List.map (\( x, y ) -> String.fromFloat x ++ "," ++ String.fromFloat y)
        |> String.join " "

turtleTransform : Turtle -> String
turtleTransform t =
    "translate(" ++ String.fromFloat t.x ++ "," ++ String.fromFloat t.y ++ ") rotate(" ++ String.fromFloat t.angle ++ ")"

-- MAIN

main =
    Browser.sandbox { init = init, update = update, view = view }
