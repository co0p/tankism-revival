package tiled

import (
	"fmt"
	"math"
	"slices"
)

const (
	GroundLayerName       = "ground"
	PropsLayerName        = "props"
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

func (m *Map) ObstacleAt(cellX, cellY int) (bool, bool) {
	gid, ok := m.TileGIDAt(cellX, cellY)
	if !ok {
		return false, false
	}

	return gid != 0, true
}

func (m *Map) TileGIDAt(cellX, cellY int) (int, bool) {
	return m.TileGIDAtLayer(ObstacleLayerName, cellX, cellY)
}

func (m *Map) TileGIDAtLayer(layerName string, cellX, cellY int) (int, bool) {
	if cellX < 0 || cellY < 0 || cellX >= m.Width || cellY >= m.Height {
		return 0, false
	}

	l := m.layerByName(layerName)
	if l == nil || l.Type != defaultTileLayerType {
		return 0, false
	}

	idx := cellY*m.Width + cellX
	if idx < 0 || idx >= len(l.Data) {
		return 0, false
	}

	return l.Data[idx], true
}

func (m *Map) SetObstacle(cellX, cellY int, blocked bool) error {
	gid := 0
	if blocked {
		gid = 1
	}

	return m.SetTileGID(cellX, cellY, gid)
}

func (m *Map) SetTileGID(cellX, cellY, gid int) error {
	return m.SetTileGIDOnLayer(ObstacleLayerName, cellX, cellY, gid)
}

func (m *Map) SetTileGIDOnLayer(layerName string, cellX, cellY, gid int) error {
	if cellX < 0 || cellY < 0 || cellX >= m.Width || cellY >= m.Height {
		return fmt.Errorf("cell out of bounds")
	}
	if gid < 0 {
		return fmt.Errorf("gid cannot be negative")
	}

	l := m.layerByName(layerName)
	if l == nil || l.Type != defaultTileLayerType {
		return fmt.Errorf("%q tile layer missing", layerName)
	}

	idx := cellY*m.Width + cellX
	if idx < 0 || idx >= len(l.Data) {
		return fmt.Errorf("cell index out of bounds")
	}

	l.Data[idx] = gid
	return nil
}

func (m *Map) ClearObstacles() error {
	return m.ClearTileLayer(ObstacleLayerName)
}

func (m *Map) ClearTileLayer(layerName string) error {
	l := m.layerByName(layerName)
	if l == nil || l.Type != defaultTileLayerType {
		return fmt.Errorf("%q tile layer missing", layerName)
	}

	for i := range l.Data {
		l.Data[i] = 0
	}

	return nil
}

func (m *Map) PropsObjects() []Object {
	l := m.layerByName(PropsLayerName)
	if l == nil || l.Type != defaultObjectLayerType {
		return nil
	}

	objects := make([]Object, 0, len(l.Objects))
	for _, obj := range l.Objects {
		if obj.GID <= 0 {
			continue
		}
		objects = append(objects, obj)
	}

	slices.SortFunc(objects, func(a, b Object) int {
		if a.ID < b.ID {
			return -1
		}
		if a.ID > b.ID {
			return 1
		}
		return 0
	})

	return objects
}

func (m *Map) AddPropsObject(gid int, worldX, worldY, width, height, rotationDegrees float64) error {
	if gid <= 0 {
		return fmt.Errorf("props gid must be > 0")
	}
	if width <= 0 {
		width = float64(m.TileWidth)
	}
	if height <= 0 {
		height = float64(m.TileHeight)
	}

	l := m.layerByName(PropsLayerName)
	if l == nil || l.Type != defaultObjectLayerType {
		return fmt.Errorf("%q object layer missing", PropsLayerName)
	}

	l.Objects = append(l.Objects, Object{
		ID:       m.nextObjectID(),
		Type:     defaultPropObjectType,
		GID:      gid,
		X:        worldX,
		Y:        worldY,
		Width:    width,
		Height:   height,
		Rotation: rotationDegrees,
		Visible:  true,
	})

	return nil
}

func (m *Map) RemoveNearestPropsObject(worldX, worldY, maxDistance float64) (bool, error) {
	l := m.layerByName(PropsLayerName)
	if l == nil || l.Type != defaultObjectLayerType {
		return false, fmt.Errorf("%q object layer missing", PropsLayerName)
	}

	closestIdx := -1
	closestDist := maxDistance

	for i := range l.Objects {
		obj := l.Objects[i]
		if obj.GID <= 0 {
			continue
		}

		distance := math.Hypot(obj.X-worldX, obj.Y-worldY)
		if distance <= closestDist {
			closestDist = distance
			closestIdx = i
		}
	}

	if closestIdx < 0 {
		return false, nil
	}

	l.Objects = slices.Delete(l.Objects, closestIdx, closestIdx+1)
	return true, nil
}

func (m *Map) ClearPropsObjects() error {
	l := m.layerByName(PropsLayerName)
	if l == nil || l.Type != defaultObjectLayerType {
		return fmt.Errorf("%q object layer missing", PropsLayerName)
	}

	l.Objects = l.Objects[:0]
	return nil
}

func (m *Map) Spawn(name string) (Object, bool) {
	l := m.layerByName(SpawnLayerName)
	if l == nil || l.Type != defaultObjectLayerType {
		return Object{}, false
	}

	for _, obj := range l.Objects {
		if obj.Name == name {
			return obj, true
		}
	}

	return Object{}, false
}

func (m *Map) SetSpawn(name string, worldX, worldY float64) error {
	l := m.layerByName(SpawnLayerName)
	if l == nil || l.Type != defaultObjectLayerType {
		return fmt.Errorf("%q layer missing", SpawnLayerName)
	}

	for i := range l.Objects {
		if l.Objects[i].Name != name {
			continue
		}

		l.Objects[i].X = worldX
		l.Objects[i].Y = worldY
		l.Objects[i].Point = true
		l.Objects[i].Visible = true
		return nil
	}

	l.Objects = append(l.Objects, Object{
		ID:      m.nextObjectID(),
		Name:    name,
		Type:    "spawn",
		X:       worldX,
		Y:       worldY,
		Point:   true,
		Visible: true,
	})

	return nil
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

	copyWidth := minInt(oldWidth, newWidth)
	copyHeight := minInt(oldHeight, newHeight)

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

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
