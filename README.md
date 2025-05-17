## md2img

Markdown to Image.

### Docker Image

[y1jiong/md2img](https://hub.docker.com/r/y1jiong/md2img)

### Usage

```text
POST /markdown?width=0&mobile=false&html=false
POST /html?width=0&mobile=false
POST /url?width=0&mobile=false
```

Post `markdown` / `html` / `url` to get image.

### Custom template.html

You can customize the template by modifying the [template.html](internal/service/markdown/template.html) file in the
repository. You can change the styles, fonts, and other elements to suit your needs. Finally, place template.html in the
same directory as the binary file.
