import { ProductConfig } from '../files/product-config';
import { Product } from '../files/product';
import Shopify, { IProduct } from 'shopify-api-node';
import { maskString } from '../utils/mask-string';
import { existsSync, readFileSync } from 'fs-extra';
import mustache from 'mustache';
import { shopify } from './shopify';

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
 * Returns the base64 encoded image of the provided product.
 */
function readImage(p: Product, c: ProductConfig): Image | undefined {
  const image = c.image;
  if (!(image && p[image.key])) { return; }

  // By default, assign the value in the column defined by `image.key` as the
  // image filename.
  let filename = p[image.key];
  // If the image configuration also defines a `charIndices` and a
  // `filenamePattern` property, then mask the filename value using those
  // properties.
  if (image.charIndices && image.filenamePattern) {
    filename = maskString(
      filename,
      image.charIndices,
      image.filenamePattern
    ) + '.jpg';
  }
  
  const filePath = `${image.dir}/${filename}`;
  if (!existsSync(filePath)) { return; }

  return {
    filename: filename,
    base64: readFileSync(filePath).toString('base64')
  };
}

/**
 * Create a product in the Shopify store.
 */
export function createProduct(p: Product, c: ProductConfig): Promise<IProduct> {
  const title = createTitle(p, c);
  const table = createSpecificationsTable(p, c);
  const image = readImage(p, c);

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
      weight_unit: p.weightUnit || 'kg'
    }],
    images: [
      ...image ? [{
        attachment: image.base64,
        filename: image.filename
      }] : []
    ]
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