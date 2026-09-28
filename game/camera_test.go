package game

import "testing"

func TestCameraFollowTargetCentersWhenNotNearEdges(t *testing.T) {
	worldMap := Map{Width: 2000, Height: 1500}
	camera := NewCamera(400, 300)

	camera.FollowTarget(worldMap, 900, 700)

	if camera.X != 700 || camera.Y != 550 {
		t.Fatalf("expected centered camera (700,550), got (%f,%f)", camera.X, camera.Y)
	}
}

func TestCameraFollowTargetClampsAtEdges(t *testing.T) {
	worldMap := Map{Width: 1000, Height: 800}
	camera := NewCamera(400, 300)

	camera.FollowTarget(worldMap, 50, 40)
	if camera.X != 0 || camera.Y != 0 {
		t.Fatalf("expected top-left clamp (0,0), got (%f,%f)", camera.X, camera.Y)
	}

	camera.FollowTarget(worldMap, 980, 790)
	if camera.X != 600 || camera.Y != 500 {
		t.Fatalf("expected bottom-right clamp (600,500), got (%f,%f)", camera.X, camera.Y)
	}
}

func TestCameraWorldToScreen(t *testing.T) {
	camera := NewCamera(400, 300)
	camera.X = 100
	camera.Y = 50

	sx, sy := camera.WorldToScreen(130, 120)
	if sx != 30 || sy != 70 {
		t.Fatalf("expected screen point (30,70), got (%f,%f)", sx, sy)
	}
}
