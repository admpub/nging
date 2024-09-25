/*
npm install imagemin imagemin-pngquant imagemin-webp -g
npm install imagemin-mozjpeg -g //[or] npm install imagemin-jpegtran -g
*/

//node ./minifyIMG.js /home/images/ /home/compressed/ 
var arguments = process.argv.splice(2); 
if (arguments.length == 1) {
   var srcFile = arguments[0];
   var pos = srcFile.lastIndexOf('.');
   arguments.push(srcFile.substring(0, pos)+'.min'+srcFile.substring(pos))
}
const imageDir = arguments[0];
const output = arguments[1];

const imagemin = require('imagemin');
const jpgImages = imageDir+'*.jpg';
const pngImages = imageDir+'*.png';

// jpeg
const imageminMozjpeg = require('imagemin-mozjpeg');
const imageminJpegtran = require('imagemin-jpegtran');
const optimiseJPEGImages = () =>
  imagemin([jpgImages], {
	destination: output,
	plugins: [
		//imageminJpegtran(),
    imageminMozjpeg({quality: 70}),
	]
});
optimiseJPEGImages().catch(error => console.log(error));

// png
const imageminPngquant = require('imagemin-pngquant');
const optimisePNGImages = () =>
  imagemin([pngImages], {
	  destination: output,
    plugins: [
      imageminPngquant({ quality: [0.6, 0.8] })
    ],
});
optimisePNGImages().catch(error => console.log(error));


// webp
const imageminWebp = require('imagemin-webp');
const convertPNGToWebp = () =>
  imagemin([pngImages], {
	  destination: output,
    plugins: [
      imageminWebp({quality: 85}),
    ],
});
const convertJPGToWebp = () =>
  imagemin([jpgImages], {
	destination: output,
	plugins: [
    imageminWebp({quality: 75}),
	]
});
convertPNGToWebp().catch(error => console.log(error));
convertJPGToWebp().catch(error => console.log(error));

// html code
/*
<picture>
    <source srcset="sample_image.webp" type="image/webp">
    <source srcset="sample_image.jpg" type="image/jpg">
    <img src="sample_image.jpg" alt="">
</picture>
*/