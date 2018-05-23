package pbr

import (
	"math"
	"math/rand"
	"time"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

type sampler struct {
	adapt   float64
	bounces int
	direct  int
	branch  int
	camera  *Camera
	scene   *Scene
}

type sample struct {
	row   int
	count int
}

func (s *sampler) start(buffer *rgb.Framebuffer, in <-chan int, done chan<- sample) {
	width := uint(s.camera.Width())
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	go func() {
		for y := range in {
			total := 0
			for x := 0; x < int(width); x++ {
				i := uint(y*int(width) + x)
				count := s.adapted(buffer, uint(i))
				for c := 0; c < count; c++ {
					buffer.Add(i, s.trace(x, y, rnd))
				}
				total += count
			}
			done <- sample{y, total}
		}
	}()
}

func (s *sampler) trace(x, y int, rnd *rand.Rand) (energy rgb.Energy) {
	ray := s.camera.ray(float64(x), float64(y), rnd)
	strength := rgb.Energy{1, 1, 1}
	// lights := s.scene.Lights() TODO: direct lighting

	for i := 0; i < 9; i++ {
		if i > 3 {
			if strength = strength.RandomGain(rnd); strength.Zero() {
				break
			}
		}
		hit := s.scene.Intersect(ray)
		if !hit.Ok {
			energy = energy.Plus(s.scene.EnvAt(ray.Dir).Times(strength))
			break
		}
		point := ray.Moved(hit.Dist)
		normal, mat := hit.Surface.At(point)
		if !mat.Light.Zero() {
			energy = energy.Plus(mat.Light.Times(strength))
			break
		}
		bsdf := mat.BSDF()
		toTangent, fromTangent := tangentMatrix(normal)

		wo := toTangent.MultDir(ray.Dir.Inv())
		wi := bsdf.Sample(wo, rnd)
		weight := wi.Dot(geom.Up) / bsdf.PDF(wi, wo)
		strength = strength.Times(bsdf.Eval(wi, wo)).Scaled(weight)

		ray = geom.NewRay(point, fromTangent.MultDir(wi))
	}
	return energy
}

// TODO: precompute on surfaces?
func tangentMatrix(normal geom.Direction) (to, from *geom.Matrix4) {
	if geom.Vector3(normal).Equals(geom.Vector3(geom.Up)) {
		return geom.Identity(), geom.Identity()
	}
	angle := math.Acos(normal.Dot(geom.Up))
	axis := normal.Cross(geom.Up)
	angleAxis := axis.Scaled(angle)
	m := geom.Rot(angleAxis)
	return m, m.Inverse()
}

func (s *sampler) traceDirect(point geom.Vector3, normal geom.Direction, rnd *rand.Rand) (energy rgb.Energy, coverage float64) {
	for i := 0; i < 1; i++ {
		light := s.scene.Light(rnd)
		ray, solidAngle := light.Box().ShadowRay(point, normal, rnd)
		if solidAngle <= 0 {
			break
		}
		coverage += solidAngle
		hit := s.scene.Intersect(ray)
		if !hit.Ok {
			break
		}
		e := hit.Surface.Material().Emit().Scaled(solidAngle / math.Pi)
		energy = energy.Plus(e)
	}
	return energy, coverage
}

// http://gfx.cs.princeton.edu/pubs/DeCoro_2010_DOR/outliers.pdf
// TODO: backgrounds should be basically completely black on the heatmap
func (s *sampler) adapted(buffer *rgb.Framebuffer, i uint) int {
	if s.adapt == 0 {
		return 1
	}
	count := buffer.Count(i)
	if count < 3 {
		return 1
	}
	needs := 1
	brightness := buffer.Average(i).Mean()
	midtones := (((255 - math.Min(brightness, 255)) / 255) + 3) / 4
	noise := buffer.Noise(i)
	varMean, countMean := buffer.Variance()
	ratio := (noise + 1) / (varMean + 1)
	targetCount := math.Ceil(ratio * countMean * midtones)
	correction := targetCount - count
	limited := math.Max(0, math.Min(s.adapt, correction))
	needs += int(limited)
	return needs
}
