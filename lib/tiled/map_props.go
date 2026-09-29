package tiled

import (
	"fmt"
	"math"
	"slices"
)

func (m *Map) PropsObjects() []Object {
	return m.PropsObjectsOnLayer(PropsLayerName)
}

func (m *Map) PropsTallObjects() []Object {
	return m.PropsObjectsOnLayer(PropsTallLayerName)
}

func (m *Map) PropsObjectsOnLayer(layerName string) []Object {
	l := m.layerByName(layerName)
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

func (m *Map) AllPropsObjects() []Object {
	objects := append(m.PropsObjectsOnLayer(PropsLayerName), m.PropsObjectsOnLayer(PropsTallLayerName)...)
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
	return m.AddPropsObjectOnLayer(PropsLayerName, gid, worldX, worldY, width, height, rotationDegrees)
}

func (m *Map) AddPropsObjectOnLayer(layerName string, gid int, worldX, worldY, width, height, rotationDegrees float64) error {
	if gid <= 0 {
		return fmt.Errorf("props gid must be > 0")
	}
	if width <= 0 {
		width = float64(m.TileWidth)
	}
	if height <= 0 {
		height = float64(m.TileHeight)
	}

	l := m.layerByName(layerName)
	if l == nil || l.Type != defaultObjectLayerType {
		return fmt.Errorf("%q object layer missing", layerName)
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
	return m.RemoveNearestPropsObjectOnLayer(PropsLayerName, worldX, worldY, maxDistance)
}

func (m *Map) RemoveNearestPropsObjectOnLayer(layerName string, worldX, worldY, maxDistance float64) (bool, error) {
	l := m.layerByName(layerName)
	if l == nil || l.Type != defaultObjectLayerType {
		return false, fmt.Errorf("%q object layer missing", layerName)
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
	return m.ClearPropsObjectsOnLayer(PropsLayerName)
}

func (m *Map) ClearTallPropsObjects() error {
	return m.ClearPropsObjectsOnLayer(PropsTallLayerName)
}

func (m *Map) ClearPropsObjectsOnLayer(layerName string) error {
	l := m.layerByName(layerName)
	if l == nil || l.Type != defaultObjectLayerType {
		return fmt.Errorf("%q object layer missing", layerName)
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
