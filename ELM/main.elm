module Main exposing (..)

import Browser
import Html exposing (Html, div, textarea, button, text, h1)
import Html.Attributes as Attr
import Html.Events exposing (onInput, onClick, on, keyCode)
import Json.Decode as Decode
import Svg exposing (svg, polyline, circle, g, ellipse, line)
import Svg.Attributes as SvgAttr

-- MODÈLE

type alias Turtle =
    { x : Float
    , y : Float
    , angle : Float 
    , path : List ((Float, Float), (Float, Float), String)
    }

type alias Model =
    { turtle : Turtle
    , commandInput : String
    , strokeColor : String
    }

type Command
    = Forward Float
    | Right Float
    | Left Float
    | CmdCarre Float
    | CmdCercle Float
    | CmdEtoile Float
    | Repeat Int (List Command)
    | Clear
    | ResetAll

init : Model
init =
    { turtle = 
        { x = 300
        , y = 200
        , angle = -90 
        , path = [] 
        }
    , commandInput = ""
    , strokeColor = "#3498db"
    }

-- UPDATE

type Msg
    = UpdateInput String
    | Execute
    | Reset
    | ClearInput
    | ChangeColor String

update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateInput txt ->
            { model | commandInput = txt }

        Reset ->
            init

        ClearInput ->
            {model | commandInput = ""}
            
        ChangeColor newColor ->
            { model | strokeColor = newColor }

        Execute ->
            let
                tokens = 
                    model.commandInput
                        |> String.replace "[" " [ "
                        |> String.replace "]" " ] "
                        |> String.replace "," " "
                        |> String.toLower
                        |> String.words
                
                (commands, _) = parseTokens tokens
            in
            runCommands commands model

parseTokens : List String -> (List Command, List String)
parseTokens tokens =
    case tokens of
        [] -> ( [], [] )
        "]" :: rest -> ( [], rest )

        "fd" :: valStr :: rest ->
            case String.toFloat valStr of
                Just val -> addCommand (Forward val) rest
                Nothing -> addCommand (Forward 50) (valStr :: rest)

        "fd" :: rest ->
            addCommand (Forward 50) rest

        "rt" :: valStr :: rest ->
            case String.toFloat valStr of
                Just val -> addCommand (Right val) rest
                Nothing -> addCommand (Right 90) (valStr :: rest)

        "rt" :: rest ->
            addCommand (Right 90) rest

        "lt" :: valStr :: rest ->
            case String.toFloat valStr of
                Just val -> addCommand (Left val) rest
                Nothing -> addCommand (Left 90) (valStr :: rest)

        "lt" :: rest ->
            addCommand (Left 90) rest

        "carre" :: valStr :: rest -> 
            if valStr == "[" then parseRepeat tokens else addCommand (CmdCarre (Maybe.withDefault 50 (String.toFloat valStr))) rest
        "carre" :: rest -> 
            addCommand (CmdCarre 50) rest

        "cercle" :: valStr :: rest -> 
            if valStr == "[" then parseRepeat tokens else addCommand (CmdCercle (Maybe.withDefault 50 (String.toFloat valStr))) rest
        "cercle" :: rest -> 
            addCommand (CmdCercle 50) rest

        "etoile" :: valStr :: rest -> 
            if valStr == "[" then parseRepeat tokens else addCommand (CmdEtoile (Maybe.withDefault 50 (String.toFloat valStr))) rest
        "etoile" :: rest -> 
            addCommand (CmdEtoile 50) rest

        "clear" :: rest -> 
            addCommand Clear rest

        "reset" :: rest -> 
            addCommand ResetAll rest

        "repeat" :: nStr :: "[" :: rest ->
            let
                n = String.toInt nStr |> Maybe.withDefault 1
                (innerCmds, remainingAfterBracket) = parseTokens rest
                (nextCmds, finalRemaining) = parseTokens remainingAfterBracket
            in
            ( Repeat n innerCmds :: nextCmds, finalRemaining )

        _ :: rest -> 
            parseTokens rest

parseRepeat : List String -> (List Command, List String)
parseRepeat tokens = ([], tokens)

addCommand : Command -> List String -> (List Command, List String)
addCommand cmd rest =
    let 
        (nextCmds, finalRemaining) = parseTokens rest
    in 
    ( cmd :: nextCmds, finalRemaining )

runCommands : List Command -> Model -> Model
runCommands commands model =
    case commands of
        [] -> model
        cmdHead :: cmdTail ->
            runCommands cmdTail (executeSingleCommand cmdHead model)

executeSingleCommand : Command -> Model -> Model
executeSingleCommand cmd model =
    case cmd of
        Forward d -> moveForward d model
        Right d -> rotateTurtle d model
        Left d -> rotateTurtle -d model
        CmdCarre s -> drawSquare s model
        CmdCercle r -> drawCircle r model
        CmdEtoile s -> drawStar s model
        Clear -> 
            let t = model.turtle 
            in { model | turtle = { t | path = [] } }
        ResetAll -> init
        Repeat n subCmds ->
            if n <= 0 then model
            else executeSingleCommand (Repeat (n - 1) subCmds) (runCommands subCmds model)

moveForward : Float -> Model -> Model
moveForward dist model =
    let
        t = model.turtle
        rad = t.angle * (pi / 180)
        newX = t.x + dist * cos rad
        newY = t.y + dist * sin rad
        newPath = t.path ++ [ ((t.x, t.y), (newX, newY), model.strokeColor) ]
    in
    { model | turtle = { t | x = newX, y = newY, path = newPath } }

rotateTurtle : Float -> Model -> Model
rotateTurtle deg model =
    let t = model.turtle in { model | turtle = { t | angle = t.angle + deg } }

drawSquare : Float -> Model -> Model
drawSquare size model =
    model |> moveForward size |> rotateTurtle 90 |> moveForward size |> rotateTurtle 90 |> moveForward size |> rotateTurtle 90 |> moveForward size |> rotateTurtle 90

drawStar : Float -> Model -> Model
drawStar size model =
    let
        orig = model.turtle.angle
        t = model.turtle
        prep = { model | turtle = { t | angle = -72 } }
        drawn = prep |> moveForward size |> rotateTurtle 144 |> moveForward size |> rotateTurtle 144 |> moveForward size |> rotateTurtle 144 |> moveForward size |> rotateTurtle 144 |> moveForward size |> rotateTurtle 144
        finalT = drawn.turtle
    in { drawn | turtle = { finalT | angle = orig } }

drawCircle : Float -> Model -> Model
drawCircle radius model =
    let
        stepSize = (2 * pi * radius) / 36
        drawSteps n m = if n <= 0 then m else drawSteps (n - 1) (m |> moveForward stepSize |> rotateTurtle 10)
    in drawSteps 36 model

-- VUE

view : Model -> Html Msg
view model =
    div [ Attr.style "display" "flex", Attr.style "flex-direction" "column", Attr.style "align-items" "center", Attr.style "justify-content" "center", Attr.style "min-height" "100vh", Attr.style "background-color" "#f0f2f5", Attr.style "font-family" "sans-serif" ]
        [ h1 [ Attr.style "color" "#2c3e50" ] [ text "TcTurtle Elm Edition" ]
        , div [ Attr.style "background" "white", Attr.style "padding" "15px", Attr.style "box-shadow" "0 8px 30px rgba(0,0,0,0.1)", Attr.style "border-radius" "12px" ]
            [ svg [ SvgAttr.width "900", SvgAttr.height "600", SvgAttr.viewBox "0 0 600 400" ]
                ( (List.map viewSegment model.turtle.path) ++ [ viewTurtle model.turtle ] )
            ]
        , div [ Attr.style "margin-top" "20px", Attr.style "width" "450px" ]
            [ div [ Attr.style "margin-bottom" "10px", Attr.style "display" "flex", Attr.style "gap" "10px", Attr.style "align-items" "center", Attr.style "justify-content" "center" ]
                [ text "Couleur du tracé :"
                , div [ Attr.style "display" "flex", Attr.style "gap" "8px" ]
                    (List.map (\c -> div [ onClick (ChangeColor c), Attr.style "background-color" c, Attr.style "width" "24px", Attr.style "height" "24px", Attr.style "border-radius" "50%", Attr.style "cursor" "pointer", Attr.style "border" (if model.strokeColor == c then "2px solid #000" else "1px solid #ddd") ] []) [ "#95a5a6", "#2c3e50", "#3498db", "#e74c3c", "#2ecc71", "#f1c40f" ])
                ]
            , textarea [ onInput UpdateInput, on "keydown" (keyCode |> Decode.andThen (\c -> if c == 13 then Decode.succeed Execute else Decode.fail "")), Attr.value model.commandInput, Attr.style "width" "100%", Attr.style "height" "80px", Attr.style "padding" "10px", Attr.style "border" "1px solid #ccc", Attr.style "border-radius" "8px", Attr.style "font-size" "16px", Attr.style "resize" "none", Attr.style "box-sizing" "border-box" ] []
            , div [ Attr.style "margin-top" "8px", Attr.style "display" "flex", Attr.style "gap" "10px" ]
                [ button (onClick (UpdateInput "repeat 6 [ fd 60, rt 60 ]") :: miniBtnStyles) [ text "Hexagone" ]
                , button (onClick (UpdateInput "repeat 12 [ carre 50, rt 30 ]") :: miniBtnStyles) [ text "Fleurs" ]
                ]
            , div [ Attr.style "display" "flex", Attr.style "gap" "10px", Attr.style "margin-top" "15px" ]
                (List.map (\(l, m, c) -> button (onClick m :: commonBtnStyles c) [ text l ]) [ ("Executer", Execute, "#f3d541ff"), ("Reset", Reset, "#973ce7ff"), ("Effacer texte", ClearInput, "#e74c3c") ])
            ]
        ]

-- Fonction pour dessiner un segment individuel avec sa propre couleur
viewSegment : ((Float, Float), (Float, Float), String) -> Svg.Svg Msg
viewSegment ((x1, y1), (x2, y2), color) =
    line
        [ SvgAttr.x1 (String.fromFloat x1), SvgAttr.y1 (String.fromFloat y1)
        , SvgAttr.x2 (String.fromFloat x2), SvgAttr.y2 (String.fromFloat y2)
        , SvgAttr.stroke color
        , SvgAttr.strokeWidth "3"
        , SvgAttr.strokeLinecap "round"
        ] []

viewTurtle : Turtle -> Svg.Svg Msg
viewTurtle t =
    g [ SvgAttr.transform ("translate(" ++ String.fromFloat t.x ++ "," ++ String.fromFloat t.y ++ ") rotate(" ++ String.fromFloat t.angle ++ ")") ]
        [ circle [ SvgAttr.cx "-8", SvgAttr.cy "-8", SvgAttr.r "3", SvgAttr.fill "#1b5e20" ] []
        , circle [ SvgAttr.cx "8", SvgAttr.cy "-8", SvgAttr.r "3", SvgAttr.fill "#1b5e20" ] []
        , circle [ SvgAttr.cx "-8", SvgAttr.cy "8", SvgAttr.r "3", SvgAttr.fill "#1b5e20" ] []
        , circle [ SvgAttr.cx "8", SvgAttr.cy "8", SvgAttr.r "3", SvgAttr.fill "#1b5e20" ] []
        , ellipse [ SvgAttr.cx "0", SvgAttr.cy "0", SvgAttr.rx "13", SvgAttr.ry "11", SvgAttr.fill "#4caf50" ] []
        , circle [ SvgAttr.cx "16", SvgAttr.cy "0", SvgAttr.r "5", SvgAttr.fill "#2e7d32" ] []
        ]

miniBtnStyles = [ Attr.style "font-size" "11px", Attr.style "padding" "5px 10px", Attr.style "cursor" "pointer", Attr.style "border-radius" "4px", Attr.style "background" "#fff", Attr.style "border" "1px solid #ccc" ]
commonBtnStyles color = [ Attr.style "background" color, Attr.style "color" "white", Attr.style "border" "none", Attr.style "padding" "12px", Attr.style "border-radius" "6px", Attr.style "cursor" "pointer", Attr.style "flex" "1", Attr.style "font-weight" "bold" ]

main = Browser.sandbox { init = init, update = update, view = view }
