# pbr: a Physically-Based 3D Renderer in Go

Package pbr implements Physically-Based Rendering with a Monte Carlo path tracer.
[[ Documentation ]](https://godoc.org/github.com/hunterloftis/pbr)
[[ Github ]](https://github.com/hunterloftis/pbr)

![Render](https://user-images.githubusercontent.com/364501/34923521-c39b132c-f96a-11e7-9a27-f79f67268079.png)

```
$ pbr lambo.obj -floor -polar 3.6 -longitude 0.1 -env 293.hdr -rad 500 -lens 60 -fstop 1.4 -to=-0.1,0.5,0.1 -dist 7.5 -focus=-2.2658,0.5542,0
```

This is an unbiased forward path-tracer written in Go and inspired by Disney's [Hyperion video](https://www.disneyanimation.com/technology/innovations/hyperion). It uses a single [BSDF](https://en.wikipedia.org/wiki/Bidirectional_scattering_distribution_function) to simulate realistic materials like plastic, aluminum, glass, rubber, and chalk with intuitive parameters like "color," "gloss," "metalness," and "opacity." It traces rays from physically-based cameras with focal lengths, apertures, and sensors. It has a simple API to create procedural scenes and also includes a robust CLI for rendering photorealistic images from the command line.

### Features

- [Adaptive sampling](#adaptive-sampling)
- [Russian roulette](https://computergraphics.stackexchange.com/questions/2316/is-russian-roulette-really-the-answer)
- Parametric shapes (spheres, cubes, triangles)
- Transformation matrices (translate, rotate, scale)
- Wavefront .obj files (meshes) and .mtl files (materials)
	- [With extended physically-based material properties](http://exocortex.com/blog/extending_wavefront_mtl_to_support_pbr)
- [Physically-based materials](https://www.marmoset.co/posts/basic-theory-of-physically-based-rendering/)
  - Fresnel reflection, transmission, absorption, diffusion
  - Color, refractive indices, gloss, transparency, separate fresnel channels, metals
- Arbitrary light sources ('everything is a light')
- [Environment maps](http://gl.ict.usc.edu/Data/HighResProbes/)
	- [Image-based lighting](https://agraphicsguy.wordpress.com/2016/09/07/image-based-lighting-in-offline-and-real-time-rendering/)
- Physically-based cameras
  - Sensor, aperture, focal length, focus, depth-of-field
- [K-D Tree acceleration](http://slideplayer.com/slide/7653218/)
- [Branched path tracing](https://blender.stackexchange.com/a/3258)
- [Supersampled anti-aliasing](https://en.wikipedia.org/wiki/Supersampling)
- 100% Go with no system dependencies
	- Sequential API, concurrent execution
	- CLI

## Try it

Download:

```
$ go get github.com/hunterloftis/pbr
$ cd $GOPATH/src/github.com/hunterloftis/pbr
$ dep ensure
```

Run:

```
$ cd $GOPATH/src/github.com/hunterloftis/pbr
$ go build ./cmd/pbr
$ pbr fixtures/models/falcon.obj -lat 0.5 -lon 0.5 -complete 5
$ open falcon.png
```

![falcon render](https://user-images.githubusercontent.com/364501/34923876-aaa1f07c-f96d-11e7-8dc2-199373bdb51e.png)

## Adaptive Sampling & Branched Tracing

Adaptive sampling devotes more time to sampling noisy areas than already-resolved ones.
Branched tracing splits primary rays into multiple branches to better sample the most important (first) bounce of each path.
Both of these techniques allow the renderer to spend its Ray-Scene intersection budget more effectively.

Both closeups of the Millennium Falcon were rendered in 10 minutes.
The top image used naive sampling while the bottom used the default adaptive and branching settings:

![falcon adaptive](https://user-images.githubusercontent.com/364501/35202471-f2fd01bc-fef0-11e7-8a2a-e2744228bbc8.png)
![falcon nonadaptive](https://user-images.githubusercontent.com/364501/35202473-f3143e5e-fef0-11e7-8ab5-cfd56a351dd3.png)

This is a heatmap of where the sampler chose to spend more time evaluating noisy pixels:

![falcon heatmap](https://user-images.githubusercontent.com/364501/35202487-00b8ebae-fef1-11e7-8d56-cf059ccd16ab.png)

## Parametric Shapes

## Transformational Matrices

## Wavefront .obj files

## Physically-based materials

## Arbitrary light sources

## Image-based lighting

## Direct lighting

## Thin surfaces

## Physically-based cameras

## Supersampled anti-aliasing

## K-D Tree acceleration

## More examples

I've gitignored the /fixtures directory to keep large binaries out of the repository.
You can download the Makefile fixtures from [Google Drive](https://drive.google.com/drive/folders/1hXQfQ9bZOIt8TvyoaUrRpELMxhKzrOCG?usp=sharing) for a library of models, materials, and HDR environments to play with.

## API

See [GoDoc](https://godoc.org/github.com/hunterloftis/pbr)

## CLI

```
Usage: pbr SCENE [options]

Positional arguments:
  SCENE                  input scene .obj

Options:
  --verbose              verbose output with scene information
  --info, -i             output scene information and exit
  --out OUT, -o OUT      output render .png
  --heat HEAT            output heatmap as .png
  --noise NOISE          output noisemap as .png
  --profile              record performance into profile.pprof
  --width WIDTH          rendering width in pixels [default: 800]
  --height HEIGHT        rendering height in pixels [default: 450]
  --target TARGET        camera target point
  --focus FOCUS          camera focus point (if other than 'target')
  --dist DIST            camera distance from target
  --lat LAT              camera polar angle on target
  --lon LON              camera longitudinal angle on target
  --lens LENS            camera focal length in mm [default: 50]
  --fstop FSTOP          camera f-stop [default: 4]
  --expose EXPOSE        exposure multiplier [default: 1]
  --ambient AMBIENT      the ambient light color [default: &{500 500 500}]
  --env ENV, -e ENV      environment as a panoramic hdr radiosity map (.hdr file)
  --rad RAD              exposure of the hdr (radiosity) environment map [default: 100]
  --floor                create a floor underneath the scene
  --adapt ADAPT          adaptive sampling multiplier [default: 8]
  --bounce BOUNCE, -d BOUNCE
                         number of light bounces (depth) [default: 8]
  --direct DIRECT, -d DIRECT
                         maximum number of direct rays to cast [default: 1]
  --branch BRANCH, -b BRANCH
                         maximum number of branches on first hit [default: 32]
  --complete COMPLETE, -c COMPLETE
                         number of samples-per-pixel at which to exit [default: +Inf]
  --thin                 treat transparent surfaces as having zero thickness
  --help, -h             display this help and exit
  --version              display version and exit
```
