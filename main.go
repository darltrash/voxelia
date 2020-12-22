package main
import "github.com/gen2brain/raylib-go/raylib" // This will handle all that graphics mumbo-jumbo
//import "https://github.com/yuin/gopher-lua" // For scripting support and thingies :)
import "strconv"

//import "fmt"

// TYPE DEFINITIONS ----------------------------------------------------------

type kind struct {
	transparent bool 
	modeltype 	uint8
	texturepath string
	texture     rl.Texture2D
}

type block struct {
	ignore   bool
	kind     uint8
	opacity  uint8
	position rl.Vector3
}

type world struct {
	blocks   []block
	kinds  	 []kind
	color    rl.Color
}

// METHOD DEFINITIONS --------------------------------------------------------

func NewWorld(block_amnt int) world {
	return world{make([]block, 100), make([]kind, 128), rl.Color{255, 255, 255, 0}}
}

func (self *world) NewBlock(x, y, z float32, kind int, opacity uint16) int {
	self.blocks = append(self.blocks, block{false, uint8(kind), uint8(opacity), rl.Vector3{x, y, z}});
	return len(self.blocks)
}

func (self *world) NewKind(texture string, modeltype uint8, transparent bool) int {
	self.kinds = append(self.kinds, kind{transparent, modeltype, texture, rl.LoadTexture("assets/" + texture)})
	return len(self.kinds)-1
}

func (self *world) Draw(Canvas rl.RenderTexture2D, camera rl.Camera3D) {
	rl.BeginTextureMode(Canvas)
		rl.ClearBackground(rl.Blank)
		rl.BeginMode3D(camera)
			for _, block := range self.blocks {
				if !block.ignore {
					kind := self.kinds[block.kind]
					self.color.A = block.opacity

					switch modeltype := kind.modeltype; modeltype {
						case 0: // Box
							rl.DrawCubeTexture(kind.texture, block.position, 1, 1, 1, self.color)
						case 1: // Billboard
							rl.DrawBillboard(camera, kind.texture, block.position, 1.0, self.color)
						case 2: // Plane (Currently not working!)
							rl.DrawPlane(block.position, rl.Vector2{1, 1}, self.color)
					}
				}
			}
		rl.EndMode3D()
	rl.EndTextureMode()
}

// SETUP STUFF ----------------------------------------------------------------

func float2string(input float32) string {
	return strconv.FormatFloat(float64(input), 'f', 6, 64)
}

func SetShaderUniform(shader rl.Shader, uniform string, data []float32, kind int32) {
	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, uniform), data, kind)
}

func main() {
	rl.InitWindow(800, 450, "Voxelia")
	rl.SetTargetFPS(600)
	rl.SetExitKey(0);

	var Shader = rl.LoadShader("", "assets/main.frag")	
	var BlockRender = rl.LoadRenderTexture(800, 450)
	
	var WorldInstance = NewWorld(100)
	var (
		GRASSBLOCK = WorldInstance.NewKind("grass_block.png", 0, false)
		GRASSBILL  = WorldInstance.NewKind("grass_bill.png", 1, true)
	)

	// GENERATE BLOCKS AND THINGS! ---------------------------------------------

	for x := float32(-5.); x < 6; x++ {
		for y := float32(-5.); y < 6; y++ {
		    WorldInstance.NewBlock(x, 0, y, GRASSBLOCK, 255)
		    if rl.GetRandomValue(1, 3)==1 {
		    	WorldInstance.NewBlock(x, 1, y, GRASSBILL, 255)
		    }
		}
	}

	var distance = float32(2.0)
	
	var camera = rl.Camera3D{
		Position: rl.NewVector3(10.0, 4.0, 0.),
		Target:   rl.NewVector3(0.0, 0.0, 0.0),
		Up:       rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:     45.0,
		Type:     rl.CameraPerspective,
	}

	var timer = float32(0)

	// MAIN LOOP DEFINITION ----------------------------------------------------

	for !rl.WindowShouldClose() {
		timer += rl.GetFrameTime() * 50
		camera.Position.X = camera.Target.X + distance * 5. 
		camera.Position.Y = distance * 3.

		// DRAWING/PROCESSING --------------------------------------------------
	
		rl.BeginDrawing()
			rl.ClearBackground(rl.RayWhite)
			rl.DrawRectangleGradientV(0, 0, 800, 450, rl.SkyBlue, rl.Blue);

			WorldInstance.Draw(BlockRender, camera)

			rl.BeginShaderMode(Shader)
				rl.DrawTextureRec(BlockRender.Texture, rl.NewRectangle(0, 0, float32(BlockRender.Texture.Width), float32(-BlockRender.Texture.Height)), rl.NewVector2(0, 0), rl.White)
			rl.EndShaderMode()

			rl.DrawText(float2string(camera.Position.X), 10, 10, 10, rl.Black)
			rl.DrawText(float2string(camera.Position.Y), 10, 20, 10, rl.Black)
			rl.DrawText(float2string(camera.Position.Z), 10, 30, 10, rl.Black)
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
