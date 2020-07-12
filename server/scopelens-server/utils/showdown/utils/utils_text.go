package utils

import (
	"github.com/golang/freetype"
	"image"
	"image/draw"
	"io/ioutil"
)

var (
	// Size
	FontSize = 20.0

	// Point sets
	// 1. Info
	PointsInfoX = 88
	PointsInfoY = 114
	PointsInfo = []image.Point{
		{PointsInfoX, PointsInfoY}, {PointsInfoX + OffsetX, PointsInfoY},
		{PointsInfoX, PointsInfoY + OffsetY}, {PointsInfoX + OffsetX, PointsInfoY + OffsetY},
		{PointsInfoX, PointsInfoY + OffsetY*2}, {PointsInfoX + OffsetX, PointsInfoY + OffsetY*2},
	}
	// 3. Move text
	PointsMoveTextX = 400
	PointsMoveTextY = 50
	PointsMoveText = []image.Point{
		{PointsMoveTextX, PointsMoveTextY},
		{PointsMoveTextX + OffsetX, PointsMoveTextY},
		{PointsMoveTextX, PointsMoveTextY + OffsetY},
		{PointsMoveTextX + OffsetX, PointsMoveTextY + OffsetY},
		{PointsMoveTextX, PointsMoveTextY + OffsetY*2},
		{PointsMoveTextX + OffsetX, PointsMoveTextY + OffsetY*2},
	}
)

// Append each pokemon name, ability and item text on the left part of each slot.
func AppendInfo(canvas image.Image, pokemonList *[]Pokemon) (image.Image, error) {
	// Validate Pokemon list
	if _, err := CheckPokemonListValid(pokemonList); err != nil {
		return nil, err
	}

	// Path
	FontPath := SpritePath + "Lato-Bold.ttf"

	// Info line space
	var lineSpace = 34
	
	// Load font file
	fontBytes, err := ioutil.ReadFile(FontPath)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	// Generate output image
	b := canvas.Bounds()
	out := image.NewRGBA(b)
	draw.Draw(out, out.Bounds(), canvas, b.Min, draw.Src)

	// Set FreeType context
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(FontSize)
	c.SetClip(canvas.Bounds())
	c.SetDst(out)
	c.SetSrc(image.White)

	for i := range *pokemonList {
		pt := freetype.Pt(PointsInfo[i].X, PointsInfo[i].Y)
		if _, err = c.DrawString((*pokemonList)[i].Name, pt); err != nil {
			return nil, err
		}
		pt = freetype.Pt(PointsInfo[i].X, PointsInfo[i].Y + lineSpace)
		if _, err = c.DrawString((*pokemonList)[i].Ability, pt); err != nil {
			return nil, err
		}
		pt = freetype.Pt(PointsInfo[i].X, PointsInfo[i].Y + lineSpace * 2)
		if _, err = c.DrawString((*pokemonList)[i].Item, pt); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// Append move text
func AppendMoveText(canvas image.Image, moveText *[]string, slot int) (image.Image, error) {
	// Validate move text slice
	if _, err := CheckMovesListValid(moveText); err != nil {
		return nil, err
	}
	if _, err := CheckSlotNumber(slot); err != nil {
		return nil, err
	}

	// Path
	FontPath := SpritePath + "Lato-Regular.ttf"

	// Move text line space
	var lineSpace = 44
	pt := PointsMoveText[slot]

	// Load font file
	fontBytes, err := ioutil.ReadFile(FontPath)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	// Generate output image
	b := canvas.Bounds()
	out := image.NewRGBA(b)
	draw.Draw(out, out.Bounds(), canvas, b.Min, draw.Src)

	// Set FreeType context
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(FontSize)
	c.SetClip(canvas.Bounds())
	c.SetDst(out)
	c.SetSrc(image.Black)

	for i := range *moveText {
		if _, err = c.DrawString((*moveText)[i], freetype.Pt(pt.X, pt.Y)); err != nil {
			return nil, err
		}
		pt.Y += lineSpace
	}

	return out, nil
}