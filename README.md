# Inertia.js Go (GoFiber) Adapter

Inertia-Fiber uses the GoFiber view engine as default.

> The Inertia.js server-side adapter for Go. Visit [inertiajs.com](https://inertiajs.com) to learn more.
> *Note: All features have been tested in the example project. If, however, you do find a bug, please report it in the issues.*

This code is taken from [ztcollazo](https://github.com/ztcollazo/fiber_inertia) with some improvement and additions.

> [!WARNING]
> This code does not yet support `TypeScript` and `ReactJS`, only supports `VueJS` and `Vite`.

## Usage

```golang
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/theartefak/inertia-fiber"
)

func main() {
    engine := inertia.New()
    app := fiber.New(fiber.Config{
        Views: engine,
    })

    app.Use(engine.Middleware())

    app.Get("/", func(c *fiber.Ctx) error {
        return c.Render("Index", fiber.Map{
            "greeting": "Hello World",
        })
    }).Name("Home")

    app.Listen("127.0.0.1:8000")
}
```

##### And then in the `app.html`

```html
<!doctype html>
<html lang="en">
    <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />

        <title inertia>Artefak</title>
        {{ .Vite }}
        {{ .Ziggy }}
    </head>
    <body>
        {{ .Inertia }}
    </body>
</html>
```

##### Example with config

```golang
...
//go:embed *
var fs embed.FS

func main() {
    engine := inertia.New(inertia.Config{
        FS         : http.FS(fs),         // The file system to use for loading templates and assets.
        AssetsPath : "./resources/js",    // The path to the assets directory.
        Template   : "Index",             // The name of the template to use.
    })
...
```

### Default Config

- **Root** : `./resources/views`
- **FS** : None (You can use embed.FS and then call http.FS on it.)
- **AssetsPath** : `./resources/js`
- **Template** : `app` 

##### Extra functions :

- Share : Shares a prop for every request.
- AddProp : Add a prop from middleware for the next request.
- AddParam : Share a param with the root template.

## Credits

- [**GoFiber**](https://github.com/gofiber/fiber)
- [**ztcollazo**](https://github.com/ztcollazo/fiber_inertia)
- [**theArtechnology**](https://github.com/theArtechnology/fiber-inertia)

and many more.

## License

This open-sourced software is licensed under the [MIT license](LICENSE).
