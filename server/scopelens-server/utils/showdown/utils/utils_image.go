package utils

import (
	"image"
	"image/draw"
	"image/png"
	"os"
)

var (
	// Point sets
	// 1. Pokemon
	PointsPmX = 95
	PointsPmY = 30
	PointsPm  = []image.Point{
		{PointsPmX, PointsPmY}, {PointsPmX + OffsetX, PointsPmY},
		{PointsPmX, PointsPmY + OffsetY}, {PointsPmX + OffsetX, PointsPmY + OffsetY},
		{PointsPmX, PointsPmY + OffsetY*2}, {PointsPmX + OffsetX, PointsPmY + OffsetY*2},
	}
	// 2. Items
	PointsItemX = 150
	PointsItemY = 50
	PointsItem = []image.Point{
		{PointsItemX, PointsItemY}, {PointsItemX + OffsetX, PointsItemY},
		{PointsItemX, PointsItemY + OffsetY}, {PointsItemX + OffsetX, PointsItemY + OffsetY},
		{PointsItemX, PointsItemY + OffsetY*2}, {PointsItemX + OffsetX, PointsItemY + OffsetY*2},
	}
	// 3. Move Icons
	PointsMoveIconsX = 358
	PointsMoveIconsY = 23
	PointsMoveIcons = []image.Point{
		{PointsMoveIconsX, PointsMoveIconsY},
		{PointsMoveIconsX + OffsetX, PointsMoveIconsY},
		{PointsMoveIconsX, PointsMoveIconsY + OffsetY},
		{PointsMoveIconsX + OffsetX, PointsMoveIconsY + OffsetY},
		{PointsMoveIconsX, PointsMoveIconsY + OffsetY*2},
		{PointsMoveIconsX + OffsetX, PointsMoveIconsY + OffsetY*2},
	}
	// 4. Types
	PointsTypesX = 200
	PointsTypesY = 24
	PointsTypes = []image.Point{
		{PointsTypesX, PointsTypesY},
		{PointsTypesX + OffsetX, PointsTypesY},
		{PointsTypesX, PointsTypesY + OffsetY},
		{PointsTypesX + OffsetX, PointsTypesY + OffsetY},
		{PointsTypesX, PointsTypesY + OffsetY*2},
		{PointsTypesX + OffsetX, PointsTypesY + OffsetY*2},
	}
)

// Save image.Image object to path `dst`.
func SaveImage(loadedImage image.Image, dst string) error {
	// Save an image
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	// Encode to `PNG` with `BestSpeed` level, then save to file
	enc := png.Encoder{CompressionLevel: png.BestSpeed}
	err = enc.Encode(f, loadedImage)
	if err != nil {
		return err
	}
	return nil
}

// Core append method. Call to return a concatenated image.
func AppendImage(canvas image.Image, targetPath string, offset image.Point) (*image.RGBA, error) {
	// Open target image
	targetFile, err := os.Open(targetPath)
	if err != nil {
		return nil, err
	}
	defer targetFile.Close()
	target, err := png.Decode(targetFile)
	if err != nil {
		return nil, err
	}

	// Append target image on canvas
	b := canvas.Bounds()
	out := image.NewRGBA(b)
	draw.Draw(out, b, canvas, image.Point{}, draw.Src)
	draw.Draw(out, target.Bounds().Add(offset), target, image.Point{}, draw.Over)

	return out, nil
}

// Append all pokemon sprites.
func AppendPokemon(canvas image.Image, pokemonList *[]Pokemon) (image.Image, error) {
	// Validate Pokemon list
	if _, err := CheckPokemonListValid(pokemonList); err != nil {
		return nil, err
	}

	// Append loop
	for i := range *pokemonList {
		var err error
		targetPath := SpritePath + "2d/" + String2Filename((*pokemonList)[i].Name) + ".png"

		// When file not found error occurred, use unknown.png instead.
		if _, err := os.Stat(targetPath); err != nil {
			targetPath = SpritePath + "2d/unknown.png"
		}

		// Append pokemon sprites at specified points
		canvas, err = AppendImage(canvas, targetPath, PointsPm[i])
		if err != nil {
			return nil, err
		}
	}
	return canvas, nil
}

// Append all items sprites.
func AppendItems(canvas image.Image, pokemonList *[]Pokemon) (image.Image, error) {
	// Validate Pokemon list
	if _, err := CheckPokemonListValid(pokemonList); err != nil {
		return nil, err
	}

	// Append loop
	for i := range *pokemonList {
		var err error
		targetPath := SpritePath + "items/" + String2Filename((*pokemonList)[i].Item) + ".png"

		// When file not found error occurred, use unknown.png instead.
		if _, err := os.Stat(targetPath); err != nil {
			targetPath = SpritePath + "items/unknown.png"
		}

		// Append pokemon sprites at specified points
		canvas, err = AppendImage(canvas, targetPath, PointsItem[i])
		if err != nil {
			return nil, err
		}
	}
	return canvas, nil
}

// Append pokemon type(s)
func AppendTypes(canvas image.Image, pokemonList *[]Pokemon) (image.Image, error) {
	// Validate Pokemon list
	if _, err := CheckPokemonListValid(pokemonList); err != nil {
		return nil, err
	}

	lineSpace := 36

	// Append a type bar func
	appendAType := func(canvas image.Image, t string, pt image.Point) (image.Image, error) {
		var err error
		targetPath := SpritePath + "types/" + String2Filename(t) + ".png"

		// When file not found error occurred, use unknown.png instead.
		if _, err := os.Stat(targetPath); err != nil {
			targetPath = SpritePath + "types/unknown.png"
		}

		// Append pokemon sprites at specified points
		canvas, err = AppendImage(canvas, targetPath, pt)
		if err != nil {
			return nil, err
		}
		return canvas, nil
	}

	// Append loop
	for i := range *pokemonList {
		var err error
		types := (*pokemonList)[i].Type
		// Validate types slice
		if _, err := CheckTypesValid(&types); err != nil {
			return nil, err
		}

		// Append a type bar
		for _, t := range types {
			canvas, err = appendAType(canvas, t, PointsTypes[i])
			if err != nil {
				return nil, err
			}
			PointsTypes[i].Y += lineSpace
		}

	}
	return canvas, nil
}

// Append move icons
func AppendMoveIcons(canvas image.Image, moveTypes *[]string, slot int) (image.Image, error) {
	// Validate move types slice
	if _, err := CheckMovesListValid(moveTypes); err != nil {
		return nil, err
	}
	if _, err := CheckSlotNumber(slot); err != nil {
		return nil, err
	}

	// Coordinates
	lineSpace := 44
	pt := PointsMoveIcons[slot]

	// Append loop
	for i := range *moveTypes {
		var err error
		targetPath := SpritePath + "moves/" + String2Filename((*moveTypes)[i]) + ".png"

		// When file not found error occurred, use unknown.png instead.
		if _, err := os.Stat(targetPath); err != nil {
			targetPath = SpritePath + "moves/unknown.png"
		}

		// Append move icons at specified points
		canvas, err = AppendImage(canvas, targetPath, pt)
		if err != nil {
			return nil, err
		}
		pt.Y += lineSpace
	}
	return canvas, nil
}

// Append both move type icons and move name text on canvas
// by calling AppendMoveIcons() and AppendMoveText().
func AppendMoves(canvas image.Image, pokemonList *[]Pokemon) (image.Image, error) {
	// Validate Pokemon list
	if _, err := CheckPokemonListValid(pokemonList); err != nil {
		return nil, err
	}

	// Append loop
	for i := range *pokemonList{
		var err error
		// Move names a.k.a. move text
		moves := (*pokemonList)[i].Moves
		// Move types
		moveTypes := make([]string, len(moves))
		for i := range moves {
			if tp, ok := Move2Type[String2Filename(moves[i])]; ok {
				moveTypes[i] = tp
			} else {
				moveTypes[i] = "Unknown"
			}
		}

		// Append move icons
		canvas, err = AppendMoveIcons(canvas, &moveTypes, i)
		if err != nil {
			return nil, err
		}

		// Append move text
		canvas, err = AppendMoveText(canvas, &moves, i)
		if err != nil {
			return nil, err
		}
	}

	return canvas, nil
}