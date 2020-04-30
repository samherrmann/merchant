import { existsSync } from 'fs-extra';
import sharp from 'sharp';
import { basename } from 'path';

/**
 * Returns the base64 encoded image.
 */
export async function readImage(filePath: string, maxWidth = 4000, maxHeight = 4000): Promise<Image | undefined> {
  // Check if file exists.
  if (!existsSync(filePath)) { return; }

  // Read image file and resize if it's larger than maximum allowed size.
  const fileBuffer = await sharp(filePath).resize({
    withoutEnlargement: true,
    fit: 'inside',
    width: maxWidth,
    height: maxHeight
  }).toBuffer();

  return {
    filename: basename(filePath),
    base64: fileBuffer.toString('base64')
  };
}

/**
 * Returns the base64 encoded images.
 */
export async function readImages(filePaths: string[], maxWidth = 4000, maxHeight = 4000): Promise<Image[]> {
  return Promise.all(
    filePaths.map(path => readImage(path, maxWidth, maxHeight))
  ).then(images => images.filter((image): image is Image => !!image));
}

export interface Image {
  base64: string;
  filename: string;
}