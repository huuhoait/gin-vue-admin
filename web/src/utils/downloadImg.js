export const downloadImage = (imgsrc, name) => {
  // Download image url as a file name
  var image = new Image()
  image.setAttribute('crossOrigin', 'anonymous')
  image.onload = function () {
    var canvas = document.createElement('canvas')
    canvas.width = image.width
    canvas.height = image.height
    var context = canvas.getContext('2d')
    context.drawImage(image, 0, 0, image.width, image.height)
    var url = canvas.toDataURL('image/png') // base64

    var a = document.createElement('a') // anchor element
    var event = new MouseEvent('click') // click event
    a.download = name || 'photo' // file name
    a.href = url // href = object url
    a.dispatchEvent(event) // trigger click
  }
  image.src = imgsrc
}
