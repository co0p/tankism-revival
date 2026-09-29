package tiled

import (
	"fmt"

	"github.com/co0p/tankism-revival/lib/mathutil"
)

const (
	GroundLayerName       = "ground"
	PropsLayerName        = "props"
	PropsTallLayerName    = "props_tall"
	ObstacleLayerName     = "obstacles"
	SpawnLayerName        = "spawns"
	PlayerSpawnObjectName = "player_spawn"
	EnemySpawnObjectName  = "enemy_spawn"

	defaultMapType         = "map"
	defaultMapVersion      = "1.10"
	defaultTiledVersion    = "1.11.0"
	defaultOrientation     = "orthogonal"
	defaultRenderOrder     = "right-down"
	defaultTileLayerType   = "tilelayer"
	defaultObjectLayerType = "objectgroup"
	defaultObjectDrawOrder = "topdown"

	defaultPropObjectType = "prop_sprite"
)

type Map struct {
	CompressionLevel int       `json:"compressionlevel,omitempty"`
	Height           int       `json:"height"`
	Width            int       `json:"width"`
	Infinite         bool      `json:"infinite"`
	Layers           []Layer   `json:"layers"`
	NextLayerID      int       `json:"nextlayerid"`
	NextObjectID     int       `json:"nextobjectid"`
	Orientation      string    `json:"orientation"`
	RenderOrder      string    `json:"renderorder"`
	TiledVersion     string    `json:"tiledversion"`
	TileHeight       int       `json:"tileheight"`
	Tilesets         []Tileset `json:"tilesets,omitempty"`
	TileWidth        int       `json:"tilewidth"`
	Type             string    `json:"type"`
	Version          string    `json:"version"`
}

type Layer struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Width     int      `json:"width,omitempty"`
	Height    int      `json:"height,omitempty"`
	Data      []int    `json:"data,omitempty"`
	Visible   bool     `json:"visible"`
	Opacity   float64  `json:"opacity"`
	X         int      `json:"x"`
	Y         int      `json:"y"`
	DrawOrder string   `json:"draworder,omitempty"`
	Objects   []Object `json:"objects,omitempty"`
}

type Object struct {
	ID       int     `json:"id"`
	Name     string  `json:"name,omitempty"`
	Type     string  `json:"type,omitempty"`
	GID      int     `json:"gid,omitempty"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Width    float64 `json:"width,omitempty"`
	Height   float64 `json:"height,omitempty"`
	Rotation float64 `json:"rotation,omitempty"`
	Point    bool    `json:"point,omitempty"`
	Visible  bool    `json:"visible"`
}

type Tileset struct {
	FirstGID    int           `json:"firstgid"`
	Name        string        `json:"name,omitempty"`
	TileWidth   int           `json:"tilewidth,omitempty"`
	TileHeight  int           `json:"tileheight,omitempty"`
	TileCount   int           `json:"tilecount,omitempty"`
	Columns     int           `json:"columns"`
	Image       string        `json:"image,omitempty"`
	ImageWidth  int           `json:"imagewidth,omitempty"`
	ImageHeight int           `json:"imageheight,omitempty"`
	Source      string        `json:"source,omitempty"`
	Tiles       []TilesetTile `json:"tiles,omitempty"`
}

type TilesetTile struct {
	ID          int    `json:"id"`
	Image       string `json:"image"`
	ImageWidth  int    `json:"imagewidth"`
	ImageHeight int    `json:"imageheight"`
}

func NewOrthogonalMap(width, height, tileWidth, tileHeight int) Map {
	m := Map{
		CompressionLevel: 0,
		Height:           height,
		Width:            width,
		Infinite:         false,
		Layers:           nil,
		NextLayerID:      1,
		NextObjectID:     1,
		Orientation:      defaultOrientation,
		RenderOrder:      defaultRenderOrder,
		TiledVersion:     defaultTiledVersion,
		TileHeight:       tileHeight,
		TileWidth:        tileWidth,
		Type:             defaultMapType,
		Version:          defaultMapVersion,
		Tilesets:         nil,
	}

	_ = m.EnsureEditorLayers()
	return m
}

func (m *Map) EnsureEditorLayers() error {
	if m.Width <= 0 || m.Height <= 0 {
		return fmt.Errorf("map width and height must be > 0")
	}
	if m.TileWidth <= 0 || m.TileHeight <= 0 {
		return fmt.Errorf("tile width and height must be > 0")
	}

	if m.Type == "" {
		m.Type = defaultMapType
	}
	if m.Version == "" {
		m.Version = defaultMapVersion
	}
	if m.TiledVersion == "" {
		m.TiledVersion = defaultTiledVersion
	}
	if m.Orientation == "" {
		m.Orientation = defaultOrientation
	}
	if m.RenderOrder == "" {
		m.RenderOrder = defaultRenderOrder
	}

	if err := m.ensureTileLayer(GroundLayerName); err != nil {
		return err
	}
	if err := m.ensureTileLayer(ObstacleLayerName); err != nil {
		return err
	}
	if err := m.ensurePropsLayer(); err != nil {
		return err
	}
	if err := m.ensureTallPropsLayer(); err != nil {
		return err
	}
	if err := m.ensureSpawnLayer(); err != nil {
		return err
	}

	if _, ok := m.Spawn(PlayerSpawnObjectName); !ok {
		m.SetSpawn(PlayerSpawnObjectName, float64(m.TileWidth)/2, float64(m.TileHeight)/2)
	}
	if _, ok := m.Spawn(EnemySpawnObjectName); !ok {
		m.SetSpawn(EnemySpawnObjectName, float64(m.TileWidth)*2.5, float64(m.TileHeight)*2.5)
	}

	return nil
}

func (m *Map) Resize(width, height int) error {
	if width <= 0 || height <= 0 {
		return fmt.Errorf("map width and height must be > 0")
	}

	m.Width = width
	m.Height = height

	return m.EnsureEditorLayers()
}

func (m *Map) ensureTileLayer(layerName string) error {
	tileLayer := m.layerByName(layerName)
	if tileLayer == nil {
		m.Layers = append(m.Layers, Layer{
			ID:      m.nextLayerID(),
			Name:    layerName,
			Type:    defaultTileLayerType,
			Width:   m.Width,
			Height:  m.Height,
			Data:    make([]int, m.Width*m.Height),
			Visible: true,
			Opacity: 1,
			X:       0,
			Y:       0,
		})
		tileLayer = m.layerByName(layerName)
	}

	if tileLayer.Type == "" {
		tileLayer.Type = defaultTileLayerType
	}
	if tileLayer.Type != defaultTileLayerType {
		return fmt.Errorf("layer %q must be type %q", layerName, defaultTileLayerType)
	}

	oldWidth := tileLayer.Width
	oldHeight := tileLayer.Height
	if oldWidth <= 0 {
		oldWidth = m.Width
	}
	if oldHeight <= 0 {
		oldHeight = m.Height
	}

	tileLayer.Data = resizeTileData(tileLayer.Data, oldWidth, oldHeight, m.Width, m.Height)
	tileLayer.Width = m.Width
	tileLayer.Height = m.Height
	tileLayer.Visible = true
	tileLayer.Opacity = 1

	return nil
}

func (m *Map) ensurePropsLayer() error {
	propsLayer := m.layerByName(PropsLayerName)
	if propsLayer == nil {
		m.Layers = append(m.Layers, Layer{
			ID:        m.nextLayerID(),
			Name:      PropsLayerName,
			Type:      defaultObjectLayerType,
			Visible:   true,
			Opacity:   1,
			X:         0,
			Y:         0,
			DrawOrder: defaultObjectDrawOrder,
		})
		return nil
	}

	if propsLayer.Type == defaultTileLayerType {
		convertedObjects := make([]Object, 0)
		for row := 0; row < propsLayer.Height; row++ {
			for col := 0; col < propsLayer.Width; col++ {
				idx := row*propsLayer.Width + col
				if idx < 0 || idx >= len(propsLayer.Data) {
					continue
				}
				gid := propsLayer.Data[idx]
				if gid == 0 {
					continue
				}

				convertedObjects = append(convertedObjects, Object{
					ID:      m.nextObjectID(),
					Type:    defaultPropObjectType,
					GID:     gid,
					X:       (float64(col) + 0.5) * float64(m.TileWidth),
					Y:       (float64(row) + 0.5) * float64(m.TileHeight),
					Width:   float64(m.TileWidth),
					Height:  float64(m.TileHeight),
					Visible: true,
				})
			}
		}

		propsLayer.Type = defaultObjectLayerType
		propsLayer.Width = 0
		propsLayer.Height = 0
		propsLayer.Data = nil
		propsLayer.Objects = convertedObjects
	}

	if propsLayer.Type != defaultObjectLayerType {
		return fmt.Errorf("layer %q must be type %q", PropsLayerName, defaultObjectLayerType)
	}

	propsLayer.Visible = true
	propsLayer.Opacity = 1
	if propsLayer.DrawOrder == "" {
		propsLayer.DrawOrder = defaultObjectDrawOrder
	}

	return nil
}

func (m *Map) ensureTallPropsLayer() error {
	tallPropsLayer := m.layerByName(PropsTallLayerName)
	if tallPropsLayer == nil {
		m.Layers = append(m.Layers, Layer{
			ID:        m.nextLayerID(),
			Name:      PropsTallLayerName,
			Type:      defaultObjectLayerType,
			Visible:   true,
			Opacity:   1,
			X:         0,
			Y:         0,
			DrawOrder: defaultObjectDrawOrder,
		})
		return nil
	}

	if tallPropsLayer.Type == "" {
		tallPropsLayer.Type = defaultObjectLayerType
	}
	if tallPropsLayer.Type != defaultObjectLayerType {
		return fmt.Errorf("layer %q must be type %q", PropsTallLayerName, defaultObjectLayerType)
	}

	tallPropsLayer.Visible = true
	tallPropsLayer.Opacity = 1
	if tallPropsLayer.DrawOrder == "" {
		tallPropsLayer.DrawOrder = defaultObjectDrawOrder
	}

	return nil
}

func (m *Map) ensureSpawnLayer() error {
	spawnLayer := m.layerByName(SpawnLayerName)
	if spawnLayer == nil {
		m.Layers = append(m.Layers, Layer{
			ID:        m.nextLayerID(),
			Name:      SpawnLayerName,
			Type:      defaultObjectLayerType,
			Visible:   true,
			Opacity:   1,
			X:         0,
			Y:         0,
			DrawOrder: defaultObjectDrawOrder,
		})
		return nil
	}

	if spawnLayer.Type == "" {
		spawnLayer.Type = defaultObjectLayerType
	}
	if spawnLayer.Type != defaultObjectLayerType {
		return fmt.Errorf("layer %q must be type %q", SpawnLayerName, defaultObjectLayerType)
	}

	spawnLayer.Visible = true
	spawnLayer.Opacity = 1
	if spawnLayer.DrawOrder == "" {
		spawnLayer.DrawOrder = defaultObjectDrawOrder
	}

	return nil
}

func (m *Map) layerByName(name string) *Layer {
	for i := range m.Layers {
		if m.Layers[i].Name == name {
			return &m.Layers[i]
		}
	}

	return nil
}

func (m *Map) nextLayerID() int {
	id := m.NextLayerID
	if id <= 0 {
		id = 1
	}
	m.NextLayerID = id + 1
	return id
}

func (m *Map) nextObjectID() int {
	id := m.NextObjectID
	if id <= 0 {
		id = 1
	}
	m.NextObjectID = id + 1
	return id
}

func resizeTileData(oldData []int, oldWidth, oldHeight, newWidth, newHeight int) []int {
	resized := make([]int, newWidth*newHeight)
	if oldWidth <= 0 || oldHeight <= 0 {
		return resized
	}

	copyWidth := mathutil.MinInt(oldWidth, newWidth)
	copyHeight := mathutil.MinInt(oldHeight, newHeight)

	for row := 0; row < copyHeight; row++ {
		oldRowStart := row * oldWidth
		newRowStart := row * newWidth
		for col := 0; col < copyWidth; col++ {
			oldIdx := oldRowStart + col
			newIdx := newRowStart + col
			if oldIdx < 0 || oldIdx >= len(oldData) || newIdx < 0 || newIdx >= len(resized) {
				continue
			}
			resized[newIdx] = oldData[oldIdx]
		}
	}

	return resized
}
