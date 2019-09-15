const TextToSVG = require('text-to-svg');
const lg = TextToSVG.loadSync('./LucidaGrande.ttf');

const attributes = {fill: 'magenta', stroke: 'lightcyan'};
const options = {x: 0, y: 0, fontSize: 112, fontFamily: 'Lucida Grande', anchor: 'top', attributes: attributes};
 
console.log(lg.getSVG('✓', options));
