package gamepalette

import "testing"

func TestLegacyGIDsRemainStable(t *testing.T) {
	tests := []struct {
		gid      int
		spriteID string
	}{
		{gid: 1, spriteID: SpriteTileGrass1},
		{gid: 2, spriteID: SpriteTileGrass2},
		{gid: 3, spriteID: SpriteObstacleBarrelGrayTop},
		{gid: 4, spriteID: SpriteObstacleFenceRed},
		{gid: 5, spriteID: SpriteObstacleCrateWood},
		{gid: 6, spriteID: SpriteTankRedSmall},
		{gid: 7, spriteID: SpriteTankBlueSmall},
		{gid: 8, spriteID: SpriteProjectileRed},
		{gid: 9, spriteID: SpriteProjectileBlue},
		{gid: 10, spriteID: SpriteShotRed},
		{gid: 11, spriteID: SpriteShotOrange},
	}

	for _, tc := range tests {
		sprite, ok := SpriteByGID(tc.gid)
		if !ok {
			t.Fatalf("expected gid %d to resolve", tc.gid)
		}
		if sprite.ID != tc.spriteID {
			t.Fatalf("gid %d expected %q, got %q", tc.gid, tc.spriteID, sprite.ID)
		}
	}
}

func TestSpritesUseUniqueGIDs(t *testing.T) {
	seen := make(map[int]string)
	for _, sprite := range Sprites {
		if sprite.GID <= 0 {
			t.Fatalf("sprite %q has invalid gid %d", sprite.ID, sprite.GID)
		}

		if existing, dup := seen[sprite.GID]; dup {
			t.Fatalf("duplicate gid %d for %q and %q", sprite.GID, existing, sprite.ID)
		}
		seen[sprite.GID] = sprite.ID
	}

	if MaxGID() != len(seen) {
		t.Fatalf("expected contiguous gids with MaxGID=%d, got %d unique gids", MaxGID(), len(seen))
	}
}
