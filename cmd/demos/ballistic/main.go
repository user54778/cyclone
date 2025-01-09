package main

import (
	"flag"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/user54778/cyclone/internal/math64"
	"github.com/user54778/cyclone/internal/physics"
)

// shotType represents the type of ballistic being shot.
type shotType int

const (
	Unused shotType = iota
	Pistol
	Artillery
	Fireball
	Laser
)

// ammoRound is a type to represent a single ammunition round record.
type AmmoRound struct {
	particle  physics.Particle // Every weapon fires a particle.
	shotType  shotType         // Different bullet types per weapon
	startTime int
}

func convertToRayVec3(v math64.Vector3) rl.Vector3 {
	return rl.Vector3{X: float32(v.X), Y: float32(v.Y), Z: float32(v.Z)}
}

// Render draws the ammo round.
func (r *AmmoRound) Render() {
	position := r.particle.Position
	rlPosition := convertToRayVec3(position)

	var color rl.Color
	switch r.shotType {
	case Pistol:
		color = rl.Black
	case Artillery:
		color = rl.Brown
	case Fireball:
		color = rl.Red
	case Laser:
		color = rl.Yellow
	}

	// Draw a solid black sphere
	rl.DrawSphereEx(rlPosition, 1, 5, 4, color)

	// Draw a flattened sphere representing the particle's shadow.
	gray := rl.NewColor(50, 50, 50, 128)
	shadowPosition := rl.Vector3{X: rlPosition.X, Y: 0, Z: rlPosition.Z}
	rl.DrawSphereWires(shadowPosition, 1, 5, 4, gray)
}

// BallisticDemo holds state of the Weapon ammo and ammo type.
type BallisticDemo struct {
	ammo            []AmmoRound
	currentShotType shotType
}

func NewBallisticDemo(ammoRounds int) BallisticDemo {
	return BallisticDemo{
		ammo:            make([]AmmoRound, ammoRounds),
		currentShotType: Pistol,
	}
}

// fire is a function that deals with the particle specifics of the ballistics.
func (demo *BallisticDemo) Fire() {
	for i := range demo.ammo {
		shot := &demo.ammo[i]
		// Get the first available round
		if shot.shotType == Unused {
			// Set properties of the particle.
			// Note here that mass of the particle should be exaggerated to more than real life.
			switch demo.currentShotType {
			case Pistol:
				shot.particle.SetMass(2.0)                                     // 2.0kg
				shot.particle.Velocity = math64.NewVector3(0.0, 0.0, 35.0)     // 35 m/s
				shot.particle.Acceleration = math64.NewVector3(0.0, -1.0, 0.0) // Effect of gravity
				shot.particle.Damping = 0.99                                   // No friction
			case Artillery:
				shot.particle.SetMass(200.0)                                // 200.0kg
				shot.particle.Velocity = math64.NewVector3(0.0, 30.0, 40.0) // 50 m/s
				shot.particle.Acceleration = math64.NewVector3(0.0, -20.0, 0.0)
				shot.particle.Damping = 0.99
			case Fireball:
				shot.particle.SetMass(1.0)                                    // 1.0kg - mostly blast damage
				shot.particle.Velocity = math64.NewVector3(0.0, 0.0, 10.0)    // 5 m/s
				shot.particle.Acceleration = math64.NewVector3(0.0, 0.6, 0.0) // Floats up
				shot.particle.Damping = 0.9
			case Laser:
				shot.particle.SetMass(0.1)                                    // 0.1kg; almost no mass. This is the kind of laser as seen in movies, not a realistic one
				shot.particle.Velocity = math64.NewVector3(0.0, 0.0, 100.0)   // 100 m/s
				shot.particle.Acceleration = math64.NewVector3(0.0, 0.0, 0.0) // No effect of gravity
				shot.particle.Damping = 0.99
			}

			// Data common to all particle types
			shot.particle.Position = math64.NewVector3(0.0, 1.5, 0.0)
			shot.startTime = int(rl.GetTime() * 1000) // In ms
			shot.shotType = demo.currentShotType

			rl.TraceLog(rl.LogInfo, "Particle Type and Properties: %#v", shot)
			// Exit after firing once
			return
		}
	}
}

// update is a function that is used to update the particle positions.
// This is where the integrator is used.
func (demo *BallisticDemo) Update() {
	duration := rl.GetFrameTime() // Last frame's duration in seconds
	if duration <= 0.0 {
		return
	}

	for i := range demo.ammo {
		shot := &demo.ammo[i]
		if shot.shotType != Unused {
			shot.particle.Integrate(float64(duration))

			// Special logic for fireball since it's onscreen longer.
			if shot.shotType == Fireball {
				if shot.particle.Position.Z > 50.0 || shot.particle.Position.Y > 20.0 {
					shot.shotType = Unused
					continue
				}
			}

			// Bounds checks
			// 1) Particle hasn't fallen below ground
			// 2) Particle's lifetime < 5s
			// 3) Particle hasn't moved past visible play area
			if shot.particle.Position.Y < 0.0 || shot.startTime+5000 < int(rl.GetTime()) || shot.particle.Position.Z > 200.0 {
				shot.shotType = Unused
			}
			rl.TraceLog(rl.LogInfo, "Updated particle: %#v", shot)
		}
	}
}

// switchWeapon switchs the users weapon type based on the numbers 1-4.
func (demo *BallisticDemo) switchWeapon() {
	switch {
	case rl.IsKeyPressed(rl.KeyOne), rl.IsKeyPressed(rl.KeyKp1):
		demo.currentShotType = Pistol
		rl.TraceLog(rl.LogInfo, "Current Ammo Type: Pistol")
	case rl.IsKeyPressed(rl.KeyTwo), rl.IsKeyPressed(rl.KeyKp2):
		demo.currentShotType = Artillery
		rl.TraceLog(rl.LogInfo, "Current Ammo Type: Artillery")
	case rl.IsKeyPressed(rl.KeyThree), rl.IsKeyPressed(rl.KeyKp3):
		demo.currentShotType = Fireball
		rl.TraceLog(rl.LogInfo, "Current Ammo Type: Fireball")
	case rl.IsKeyPressed(rl.KeyFour), rl.IsKeyPressed(rl.KeyKp4):
		demo.currentShotType = Laser
		rl.TraceLog(rl.LogInfo, "Current Ammo Type: Laser")
	}
}

// mouse checks if the user is clicking the left mouse button, and will fire if so.
func (demo *BallisticDemo) mouse() {
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		demo.Fire()
		rl.TraceLog(rl.LogInfo, "Fire!")
	}
}

func main() {
	var maxRounds int
	flag.IntVar(&maxRounds, "rounds", 16, "max amount of bullet rounds that can be on screen")

	flag.Parse()

	demo := NewBallisticDemo(maxRounds)

	rl.InitWindow(1280, 720, "ballistic")
	defer rl.CloseWindow()

	camera := &rl.Camera{}
	camera.Position = rl.NewVector3(-25.0, 8.0, 5.0)
	camera.Target = rl.NewVector3(0.0, 5.0, 22.0) // Camera looking at point
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)      // Where it rotates over the y-unit vector
	camera.Fovy = 45.0                            // How close I am
	// Camera Perspective set to projection by default.

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Input
		demo.switchWeapon()
		demo.mouse()

		// Game logic
		demo.Update()

		// Rendering
		rl.BeginDrawing()
		rl.ClearBackground(rl.LightGray)

		rl.BeginMode3D(*camera)

		/*
			// Draw a sphere at the firing point.
			rl.DrawSphereEx(rl.NewVector3(0.0, 1.5, 0.0), 0.1, 5, 5, rl.Gray)
			// Add a shadow projected onto the ground.
			rl.DrawSphereEx(rl.NewVector3(0.0, -1.5, 0.0), 0.1, 5, 5, rl.Black)
		*/

		// Draw scale lines
		for i := 0; i < 200; i++ {
			rl.DrawLine3D(rl.NewVector3(-5.0, 0.0, float32(i)), rl.NewVector3(5.0, 0.0, float32(i)), rl.Gray)
		}

		// Render each particle in turn
		for _, round := range demo.ammo {
			if round.shotType != Unused {
				// rl.TraceLog(rl.LogInfo, "Rendering particle")
				round.Render()
			}
		}
		// rl.TraceLog(rl.LogInfo, "Demo object: %#v", demo)

		rl.EndMode3D()

		// Back in 2D mode. Draw the FPS in the top-left.
		rl.DrawFPS(10, 10) // NOTE: Draw FPS after finishing 3D rendering since FPS will show performance of entire frame,
		// including both 3D rendering and other logic.

		rl.EndDrawing()
	}
}

// display is used to draw the main textures of the simulation, including the starting
// sphere at the firing point with its shadow, rendering the particle and the name of the shot type.
/*
func display(camera rl.Camera) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.LightGray)

	rl.BeginMode3D(camera)

	// Draw a sphere at the firing point.
	rl.DrawSphereEx(rl.NewVector3(0.0, 1.5, 0.0), 0.1, 5, 5, rl.Gray)
	// Add a shadow projected onto the ground.
	rl.DrawSphereEx(rl.NewVector3(0.0, -1.5, 0.0), 0.1, 5, 5, rl.Black)

	// Draw scale lines
	for i := 0; i < 200; i++ {
		rl.DrawLine3D(rl.NewVector3(-5.0, 0.0, float32(i)), rl.NewVector3(5.0, 0.0, float32(i)), rl.Gray)
	}

	rl.EndMode3D()

	// Back in 2D mode. Draw the FPS in the top-left.
	rl.DrawFPS(10, 10) // NOTE: Draw FPS after finishing 3D rendering since FPS will show performance of entire frame,
	// including both 3D rendering and other logic.

	rl.EndDrawing()
}
*/
/*
	*
	// Camera3D type, defines a camera position/orientation in 3d space
	type Camera3D struct {
		// Camera position
		Position Vector3
		// Camera target it looks-at
		Target Vector3
		// Camera up vector (rotation over its axis)
		Up Vector3
		// Camera field-of-view apperture in Y (degrees) in perspective, used as near plane width in orthographic
		Fovy float32
		// Camera type, controlling projection type, either CameraPerspective or CameraOrthographic.
		Projection CameraProjection
	}
	camera := &rl.Camera{}
	camera.Position = rl.NewVector3(0.0, 10.0, 10.0) // Camera Position
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)     // Camera looking at point
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0) // Where it rotates over the y-unit vector
	camera.Fovy = 45.0 // How close I am
*/
/*
	rl.DrawCube(rl.NewVector3(-4.0, 0.0, 2.0), 2.0, 5.0, 2.0, rl.Red)
	rl.DrawCubeWires(rl.NewVector3(-4.0, 0.0, 2.0), 2.0, 5.0, 2.0, rl.Gold)
	rl.DrawCubeWires(rl.NewVector3(-4.0, 0.0, -2.0), 3.0, 6.0, 2.0, rl.Maroon)

	rl.DrawSphere(rl.NewVector3(-1.0, 0.0, -2.0), 1.0, rl.Green)
	rl.DrawSphereWires(rl.NewVector3(1.0, 0.0, 2.0), 2.0, 16, 16, rl.Lime)

	rl.DrawCylinder(rl.NewVector3(4.0, 0.0, -2.0), 1.0, 2.0, 3.0, 4, rl.SkyBlue)
	rl.DrawCylinderWires(rl.NewVector3(4.0, 0.0, -2.0), 1.0, 2.0, 3.0, 4, rl.DarkBlue)
	rl.DrawCylinderWires(rl.NewVector3(4.5, -1.0, 2.0), 1.0, 1.0, 2.0, 6, rl.Brown)

	rl.DrawCylinder(rl.NewVector3(1.0, 0.0, -4.0), 0.0, 1.5, 3.0, 8, rl.Gold)
	rl.DrawCylinderWires(rl.NewVector3(1.0, 0.0, -4.0), 0.0, 1.5, 3.0, 8, rl.Pink)

	rl.DrawGrid(10, 1.0) // Draw a grid
*/
