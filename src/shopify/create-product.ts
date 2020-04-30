import { ProductConfig } from '../files/product-config';
import { Product } from '../files/product';
import Shopify, { IProduct } from 'shopify-api-node';
import { maskString } from '../utils/mask-string';
import mustache from 'mustache';
import { shopify } from './shopify';
import { readImages } from '../files/read-image';

// Override Mustache's escape function to not HTML-escape variables.
mustache.escape = (text: string): string => text;

/**
 * Returns the specifications of the provided product in an HTML table.
 */
function createSpecificationsTable(p: Product, c: ProductConfig): string {
  const rows = c.specifications.reduce<string[]>((prev, curr) => {
    let label = curr.label;
    if (curr.units) {
      label = `${label} [${curr.units}]`
    }
    prev.push(`
    <tr>
      <th>${label}</th>
      <td>${p[curr.key]}</td>
    </tr>
    `);
    return prev;
  }, []);

  return `<table>${rows.join('')}</table>`
}

/**
 * Returns the title for the provided product.
 */
function createTitle(p: Product, c: ProductConfig): string {
  return mustache.render(c.title, p);
}

/**
 * Returns the image filenames based on the configuration.
 */
function imageFilenames(p: Product, c: ProductConfig): string[] {
  const imageConfig = c.image;
  if (imageConfig && p[imageConfig.key]) {  
    const filenames = p[imageConfig.key].split('|').map(fn => fn.trim());
    return filenames.map(filename => {
      // If the image configuration also defines a `charIndices` and a
      // `filenamePattern` property, then mask the filename value using those
      // properties.
      if (imageConfig.charIndices && imageConfig.filenamePattern) {
        filename = maskString(
          filename,
          imageConfig.charIndices,
          imageConfig.filenamePattern
        ) + '.jpg';
      }
      return `${imageConfig.dir}/${filename}`;
    });
  }
  return [];
}

/**
 * Create a product in the Shopify store.
 */
export async function createProduct(p: Product, c: ProductConfig): Promise<IProduct> {
  const title = createTitle(p, c);
  const table = createSpecificationsTable(p, c);
  const images = await readImages(imageFilenames(p, c));

  /* eslint-disable @typescript-eslint/camelcase */
  const product: NewProduct = {
    title: title,
    body_html: table,
    vendor: p.vendor,
    product_type: c.type,
    tags: [p.vendor, p.name].join(', '),
    variants: [{
      inventory_management: 'shopify',
      weight: parseFloat(p.weight),
      weight_unit: p.weight_unit || 'kg'
    }],
    images: images.map(image => {
      return {
        attachment: image.base64,
        filename: image.filename
      };
    })
  };
  /* eslint-enable @typescript-eslint/camelcase */
  return shopify.product.create(product);
}

type RecursivePartial<T> = {
  [P in keyof T]?: RecursivePartial<T[P]>;
};

interface Image {
  base64: string;
  filename: string;
}

type NewProduct = RecursivePartial<Shopify.IProduct & {
  images: {
    attachment: string;
    filename: string;
  }[];
}>;