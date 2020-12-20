package main
import "github.com/gen2brain/raylib-go/raylib" // This will handle all that graphics mumbo-jumbo

// TYPE DEFINITIONS ------------------------------------------------------

type kind struct {
	transparent bool
	texture     rl.Texture2D
}

type block struct {
	__empty  bool
	ignore   bool
	kind     string
	opacity  uint16
	position rl.Vector3
}

// WORLD/ENGINE DEFINITION -----------------------------------------------

type world struct {
	blocks   []*block
	kinds  	 map[string]*kind
}

func newWorld(objs_amnt uint) world {
	var blockslice = make([]*block, objs_amnt)
	
	for i := range blockslice {			// I'll fill the entire slice with "placeholder" objects
		blockslice[i] = new(block)
		blockslice[i].__empty = true 	// And then mark them as "empty" so then they arent processed
	}

	return world{blockslice, make(map[string]*kind)}
}

func (self *world) NewBlock(x, y, z float32, kind string, opacity uint16) (correct bool) {
	for i := range self.blocks {
		if self.blocks[i].__empty {
			self.blocks[i] = &block{false, false, kind, opacity, rl.Vector3{x, y, z}};

			correct = true
			break
		}
	}
	return
}

func (self *world) NewKind(name string, transparent bool, texture rl.Texture2D) {
	self.kinds[name] = &kind{transparent, texture}
}

func (self *world) Update(delta float32) {}

func (self *world) Draw() {
	for _, block := range self.blocks {
		if block.__empty {
			break
		}
		
		if !block.ignore {
			kind := *self.kinds[block.kind]
			rl.DrawCubeTexture(kind.texture, block.position, 1, 1, 1, rl.White)
		}
	}
}

// MAIN LOOP DEFINITION --------------------------------------------------

func main() {
	rl.InitWindow(800, 450, "Voxelia")
	rl.SetTargetFPS(60)

	var WorldInstance = newWorld(100 * 100 * 100)
	WorldInstance.NewKind("Arrow", false, rl.LoadTexture("assets/test_arrow.png"))
	WorldInstance.NewBlock(0, 0, 0, "Arrow", 255)
	WorldInstance.NewBlock(0, 1, 0, "Arrow", 255)

	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(5.0, 4.0, 0.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Type = rl.CameraPerspective

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		WorldInstance.Update(rl.GetFrameTime())

		rl.BeginMode3D(camera)
		WorldInstance.Draw()
		rl.EndMode3D()

		rl.DrawFPS(5, 5)
		
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
