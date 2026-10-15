package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/co0p/tankism-revival/game"
	"github.com/co0p/tankism-revival/game/scene"
	"github.com/co0p/tankism-revival/lib/fonts"
	"github.com/co0p/tankism-revival/lib/shaders"
	"github.com/co0p/tankism-revival/resources"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowScale = 2
	windowTitle = "Tankism Revival"
)

func main() {
	mapPath := flag.String("map", game.DefaultMapPath, "path to Tiled map JSON (.tmj)")
	flag.Parse()

	// Create font loader
	fontLoader := fonts.NewLoader(resources.FS)

	// Load blur shader for menu background
	blurShader, err := shaders.Load(resources.FS, resources.BlurShaderPath)
	if err != nil {
		log.Fatalf("load blur shader: %v", err)
	}

	// Create scene factories
	sceneFactories := map[scene.SceneID]scene.SceneFactory{
		scene.SceneIDLoading: func(ctx map[string]interface{}) (scene.Scene, error) {
			if ctx == nil {
				ctx = make(map[string]interface{})
			}
			ctx[scene.FontLoaderContextKey] = fontLoader
			return scene.NewLoadingScene(ctx)
		},
		scene.SceneIDMenu: func(ctx map[string]interface{}) (scene.Scene, error) {
			if ctx == nil {
				ctx = make(map[string]interface{})
			}
			ctx[scene.FontLoaderContextKey] = fontLoader
			ctx[scene.BlurShaderContextKey] = blurShader
			return scene.NewMenuScene(ctx)
		},
		scene.SceneIDPlay: func(ctx map[string]interface{}) (scene.Scene, error) {
			// Pass the map path to the play scene
			if ctx == nil {
				ctx = make(map[string]interface{})
			}
			ctx[scene.MapPathContextKey] = *mapPath
			return scene.NewPlayScene(ctx)
		},
		scene.SceneIDGameOver: func(ctx map[string]interface{}) (scene.Scene, error) {
			if ctx == nil {
				ctx = make(map[string]interface{})
			}
			ctx[scene.FontLoaderContextKey] = fontLoader
			ctx[scene.BlurShaderContextKey] = blurShader
			return scene.NewGameOverScene(ctx)
		},
	}

	// Create the scene manager
	manager, err := scene.NewSceneManager(scene.ManagerConfig{
		FadeDurationSeconds: 0.3,
		FadeColor:           color.RGBA{R: 0, G: 0, B: 0, A: 255},
		SceneFactories:      sceneFactories,
		LogicalWidth:        game.LogicalWidth,
		LogicalHeight:       game.LogicalHeight,
	})
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowTitle(windowTitle)
	ebiten.SetWindowSize(game.LogicalWidth*windowScale, game.LogicalHeight*windowScale)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(manager); err != nil {
		// Quit game is a normal exit, not an error
		if err.Error() == "quit game" {
			os.Exit(0)
		}
		log.Fatal(fmt.Errorf("run game: %w", err))
	}
}
