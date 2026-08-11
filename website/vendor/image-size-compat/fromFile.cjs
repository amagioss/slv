'use strict';

const fs = require('node:fs/promises');

const svgRootRe = /<svg\s([^>"']|"[^"]*"|'[^']*')*>/;
const widthRe = /\swidth=(['"])([^%]+?)\1/;
const heightRe = /\sheight=(['"])([^%]+?)\1/;
const viewBoxRe = /\sviewBox=(['"])(.+?)\1/i;
const INCH_CM = 2.54;
const units = {
  in: 96,
  cm: 96 / INCH_CM,
  em: 16,
  ex: 8,
  m: (96 / INCH_CM) * 100,
  mm: 96 / INCH_CM / 10,
  pc: 96 / 72 / 12,
  pt: 96 / 72,
  px: 1,
};
const unitsRe = new RegExp(
  `^([0-9.]+(?:e\\d+)?)(${Object.keys(units).join('|')})?$`,
);

function parseLength(len) {
  const match = unitsRe.exec(len);
  if (!match) {
    return undefined;
  }
  return Math.round(Number(match[1]) * (units[match[2]] || 1));
}

function dimensionsFromSvg(source) {
  const root = source.match(svgRootRe);
  if (!root) {
    return undefined;
  }
  const widthMatch = root[0].match(widthRe);
  const heightMatch = root[0].match(heightRe);
  const viewBoxMatch = root[0].match(viewBoxRe);
  const width = widthMatch && parseLength(widthMatch[2]);
  const height = heightMatch && parseLength(heightMatch[2]);
  if (width && height) {
    return {width, height, type: 'svg'};
  }
  if (viewBoxMatch) {
    const bounds = viewBoxMatch[2].split(' ');
    const vbWidth = parseLength(bounds[2]);
    const vbHeight = parseLength(bounds[3]);
    if (!vbWidth || !vbHeight) {
      return undefined;
    }
    if (width) {
      return {
        width,
        height: Math.floor(width / (vbWidth / vbHeight)),
        type: 'svg',
      };
    }
    if (height) {
      return {
        width: Math.floor(height * (vbWidth / vbHeight)),
        height,
        type: 'svg',
      };
    }
    return {width: vbWidth, height: vbHeight, type: 'svg'};
  }
  return undefined;
}

async function imageSizeFromFile(filePath) {
  const buffer = await fs.readFile(filePath);
  const head = buffer.subarray(0, 1000).toString('utf8');
  if (svgRootRe.test(head)) {
    const svgSize = dimensionsFromSvg(buffer.toString('utf8'));
    if (svgSize) {
      return svgSize;
    }
  }
  const {imageDimensionsFromData} = await import('image-dimensions');
  const size = imageDimensionsFromData(buffer);
  if (!size) {
    throw new TypeError(`Unrecognized image format: ${filePath}`);
  }
  return size;
}

function setConcurrency() {}

module.exports = {imageSizeFromFile, setConcurrency};
