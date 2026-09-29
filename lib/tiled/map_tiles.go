package tiled

import "fmt"

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
