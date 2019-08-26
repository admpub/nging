//css: npm install clean-css
var cleanCSS = require('clean-css');
var process = require('process');
var fs = require('fs')

function cssMinifier(flieIn, fileOut) {
   var flieIn = Array.isArray(flieIn) ? flieIn : [flieIn];
   var origCode, finalCode = '';
   var clean = new cleanCSS({})
   for (var i = 0; i < flieIn.length; i++) {
      origCode = fs.readFileSync(flieIn[i], 'utf8');
      finalCode += clean.minify(origCode).styles;
   }
   fs.writeFileSync(fileOut, finalCode, 'utf8');
}
var arguments = process.argv.splice(2); //node ./minifyCSS.js src.js min.js
cssMinifier(arguments[0], arguments[1]);  //单个文件压缩
//cssMinifier(['./file-src/index_20120913.css','./file-src/indexw_20120913.css'], './file-smin/index.css');